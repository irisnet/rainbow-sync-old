package cron

import (
	"github.com/irisnet/rainbow-sync/block"
	"github.com/irisnet/rainbow-sync/db"
	"github.com/irisnet/rainbow-sync/lib/pool"
	"github.com/irisnet/rainbow-sync/logger"
	"github.com/irisnet/rainbow-sync/model"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
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
				total, err := GetErrTxsByPage(skip, 20)
				if err != nil {
					logger.Error("Get ErrTxs ByPage have error", logger.String("err", err.Error()))
				}
				if total < 20 {
					runValue = false
					logger.Debug("finish Get ErrTxs ByPage.")
				} else {
					skip = skip + total
					logger.Debug("continue Get ErrTxs ByPage", logger.Int("skip", skip))
				}
			}

			logger.Debug("finish repair  err tx.")
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

func GetErrTxsByPage(skip, limit int) (int, error) {

	res, err := new(model.ErrTx).Find(skip, limit)
	if err != nil {
		return 0, err
	}

	if len(res) > 0 {
		doWork(res)
	}

	return len(res), nil
}
func doWork(failtxs []model.ErrTx) {
	client := pool.GetClient()
	defer func() {
		client.Release()
	}()

	for _, val := range failtxs {
		txs, msgs, err := block.ParseTxs(val.Height, client)
		if err != nil {
			continue
		}
		if err := RepareTxs(txs, msgs, val); err != nil {
			continue
		}
		val.Repair = 1
		if err := db.Update(&val); err != nil {
			logger.Error(err.Error(),
				logger.Int64("height", val.Height),
				logger.String("txhash", val.TxHash),
			)
		}
	}

}

func RepareTxs(txs []*model.Tx, txMsgs []model.TxMsg, failtx model.ErrTx) error {
	var (
		ops []txn.Op
	)
	txAndMsgNum := len(txs) + len(txMsgs)
	if txAndMsgNum > 0 {
		ops = make([]txn.Op, 0, txAndMsgNum)
		for _, v := range txs {
			if failtx.TxHash != v.TxHash {
				continue
			}
			op := txn.Op{
				C:      model.CollectionNameIrisTx,
				Id:     bson.NewObjectId(),
				Insert: v,
			}
			ops = append(ops, op)
		}

		for _, v := range txMsgs {
			if failtx.TxHash != v.TxHash {
				continue
			}
			op := txn.Op{
				C:      model.CollectionNameIrisTxMsg,
				Id:     bson.NewObjectId(),
				Insert: v,
			}
			ops = append(ops, op)
		}
	}
	if len(ops) <= 0 {
		return nil
	}
	return db.Txn(ops)
}
