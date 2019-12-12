package cron

import (
	"time"
	"os"
	"os/signal"
	"github.com/irisnet/rainbow-sync/logger"
	"github.com/irisnet/rainbow-sync/db"
	model "github.com/irisnet/rainbow-sync/model"
	"github.com/irisnet/rainbow-sync/block"
	"github.com/irisnet/rainbow-sync/helper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"fmt"
)

type CronService struct{}

func (s *CronService) StartCronService() {
	fn := func() {
		logger.Debug("Start  CronService ...")
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		stop := make(chan os.Signal)
		signal.Notify(stop, os.Interrupt)

		fn_update := func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("CronService have error", logger.Any("err", r))
				}
			}()

			runValue := true
			skip := 0
			for runValue {
				total, err := GetUnknownTxsByPage(skip, 20)
				if err != nil {
					logger.Error("GetUnknownTxsByPage have error", logger.String("err", err.Error()))
				}
				if total < 20 {
					runValue = false
					logger.Debug("finish GetUnknownTxsByPage.")
				} else {
					skip = skip + total
					logger.Debug("continue GetUnknownTxsByPage", logger.Int("skip", skip))
				}
			}

			logger.Debug("finish update  txs.")
		}
		fn_update()
		for {
			select {
			case <-ticker.C:
				fn_update()
			case <-stop:
				close(stop)
				logger.Debug(" CronService Quit...")
				return
			}

		}
	}

	go fn()
}

func GetUnknownTxsByPage(skip, limit int) (int, error) {

	var res []model.IrisTx
	q := bson.M{"status": "unknown"}
	sorts := []string{"-height"}

	fn := func(c *mgo.Collection) error {
		return c.Find(q).Sort(sorts...).Skip(skip).Limit(limit).All(&res)
	}

	if err := db.ExecCollection(model.CollectionNameIrisTx, fn); err != nil {
		return 0, err
	}

	if len(res) > 0 {
		doWork(res)
	}

	return len(res), nil
}

func GetCoinFlowByHash(txhash string) ([]model.IrisAssetDetail, error) {
	var res []model.IrisAssetDetail
	q := bson.M{"tx_hash": txhash}
	sorts := []string{"-height"}

	fn := func(c *mgo.Collection) error {
		return c.Find(q).Sort(sorts...).All(&res)
	}

	if err := db.ExecCollection(model.CollectionNameAssetDetail, fn); err != nil {
		return nil, err
	}

	return res, nil
}

func doWork(iristxs []model.IrisTx) {
	client := helper.GetClient()
	defer func() {
		client.Release()
	}()

	for _, val := range iristxs {
		txs, err := ParseUnknownTxs(val.Height, client)
		if err != nil {
			continue
		}
		if err := UpdateUnknowTxs(txs); err != nil {
			logger.Warn("UpdateUnknowTxs have error", logger.String("error", err.Error()))
		}
		if err := UpdateCoinFlow(val.TxHash, val.Height, client); err != nil {
			logger.Warn("UpdateCoinFlow have error", logger.String("error", err.Error()))
		}
	}

}

func ParseUnknownTxs(height int64, client *helper.Client) (resIrisTxs []*model.IrisTx, err error) {
	var irisBlock block.Iris_Block
	resIrisTxs, err = irisBlock.ParseIrisTxs(height, client)
	if err != nil {
		logger.Error("Parse block txs fail", logger.Int64("block", height),
			logger.String("err", err.Error()))
	}
	return
}

func ParseCoinflows(height int64, client *helper.Client) (coinflows []*model.IrisAssetDetail, err error) {
	var irisBlock block.Iris_Block
	coinflows, err = irisBlock.ParseIrisAssetDetail(height, client)
	if err != nil {
		logger.Error("Parse block coinflow fail", logger.Int64("block", height),
			logger.String("err", err.Error()))
	}
	return
}

func UpdateUnknowTxs(iristx []*model.IrisTx) error {

	update_fn := func(tx *model.IrisTx) error {
		fn := func(c *mgo.Collection) error {
			return c.Update(bson.M{"tx_hash": tx.TxHash},
				bson.M{"$set": bson.M{"actual_fee": tx.ActualFee, "status": tx.Status, "tags": tx.Tags}})
		}

		if err := db.ExecCollection(model.CollectionNameIrisTx, fn); err != nil {
			return err
		}
		return nil
	}

	for _, dbval := range iristx {
		update_fn(dbval)
	}

	return nil
}

func UpdateCoinFlow(txhash string, height int64, client *helper.Client) error {

	coinflows, err := GetCoinFlowByHash(txhash)
	if err != nil {
		return err
	}
	var ops []txn.Op

	if len(coinflows) > 0 {
		return fmt.Errorf("coinflow not need to update")
	}
	assetdetail, err := ParseCoinflows(height, client)
	for _, dbval := range assetdetail {
		ops = append(ops, txn.Op{
			C:      model.CollectionNameAssetDetail,
			Id:     bson.NewObjectId(),
			Insert: dbval,
		})

	}

	if len(ops) > 0 {
		err := db.Txn(ops)
		if err != nil {
			return err
		}
	}

	return nil
}
