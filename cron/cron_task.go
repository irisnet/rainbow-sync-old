package cron

import (
	"github.com/irisnet/rainbow-sync/block"
	"github.com/irisnet/rainbow-sync/db"
	"github.com/irisnet/rainbow-sync/lib/pool"
	"github.com/irisnet/rainbow-sync/logger"
	"github.com/irisnet/rainbow-sync/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"os/signal"
	"time"
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

	var res []model.Tx
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

func doWork(iristxs []model.Tx) {
	client := pool.GetClient()
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
	}

}

func ParseUnknownTxs(height int64, client *pool.Client) (resTxs []*model.Tx, err error) {
	resTxs, err = block.ParseTxs(height, client)
	if err != nil {
		logger.Error("Parse block txs fail", logger.Int64("block", height),
			logger.String("err", err.Error()))
	}
	return
}

func UpdateUnknowTxs(iristx []*model.Tx) error {

	update_fn := func(tx *model.Tx) error {
		fn := func(c *mgo.Collection) error {
			return c.Update(bson.M{"tx_hash": tx.TxHash},
				bson.M{"$set": bson.M{"actual_fee": tx.ActualFee, "status": tx.Status, "events": tx.Events}})
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
