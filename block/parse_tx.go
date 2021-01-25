package block

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/irisnet/rainbow-sync/db"
	"github.com/irisnet/rainbow-sync/lib/cdc"
	"github.com/irisnet/rainbow-sync/lib/pool"
	"github.com/irisnet/rainbow-sync/logger"
	"github.com/irisnet/rainbow-sync/model"
	"github.com/irisnet/rainbow-sync/utils"
	aTypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"time"
)

func SaveDocsWithTxn(blockDoc *model.Block, txs []*model.Tx, txMsgs []model.TxMsg, taskDoc model.SyncTask) error {
	var (
		ops, insertOps []txn.Op
	)

	if blockDoc.Height == 0 {
		return fmt.Errorf("invalid block, height equal 0")
	}

	blockOp := txn.Op{
		C:      model.CollectionNameBlock,
		Id:     bson.NewObjectId(),
		Insert: blockDoc,
	}

	txAndMsgNum := len(txs) + len(txMsgs)
	if txAndMsgNum > 0 {
		insertOps = make([]txn.Op, 0, txAndMsgNum)
		for _, v := range txs {
			op := txn.Op{
				C:      model.CollectionNameIrisTx,
				Id:     bson.NewObjectId(),
				Insert: v,
			}
			insertOps = append(insertOps, op)
		}

		for _, v := range txMsgs {
			op := txn.Op{
				C:      model.CollectionNameIrisTxMsg,
				Id:     bson.NewObjectId(),
				Insert: v,
			}
			insertOps = append(insertOps, op)
		}
	}

	updateOp := txn.Op{
		C:      model.CollectionNameSyncTask,
		Id:     taskDoc.ID,
		Assert: txn.DocExists,
		Update: bson.M{
			"$set": bson.M{
				"current_height":   taskDoc.CurrentHeight,
				"status":           taskDoc.Status,
				"last_update_time": taskDoc.LastUpdateTime,
			},
		},
	}

	ops = make([]txn.Op, 0, txAndMsgNum+2)
	ops = append(append(ops, blockOp, updateOp), insertOps...)

	if len(ops) > 0 {
		err := db.Txn(ops)
		if err != nil {
			return err
		}
	}

	return nil
}

func ParseBlock(b int64, client *pool.Client) (resBlock *model.Block, resTxs []*model.Tx, resTxMsgs []model.TxMsg, resErr error) {

	defer func() {
		if err := recover(); err != nil {
			logger.Error("parse  block fail", logger.Int64("height", b),
				logger.Any("err", err))
			resErr = fmt.Errorf("%v", err)
		}
	}()

	resBlock = &model.Block{
		Height:     b,
		CreateTime: time.Now().Unix(),
	}

	txs, msgs, err := ParseTxs(b, client)
	if err != nil {
		resErr = err
		return
	}

	resTxs = txs
	resTxMsgs = msgs

	return
}

func ParseTxs(b int64, client *pool.Client) ([]*model.Tx, []model.TxMsg, error) {
	ctx := context.Background()
	resblock, err := client.Block(ctx, &b)
	if err != nil {
		time.Sleep(1 * time.Second)
		// there is possible parse block fail when in iterator
		var err2 error
		client2 := pool.GetClient()
		resblock, err2 = client2.Block(ctx, &b)
		client2.Release()
		if err2 != nil {
			var docFailTx model.ErrTx
			docFailTx.Height = b
			docFailTx.Log = fmt.Sprintf("parse block err:%v", err2.Error())
			if err := docFailTx.Save(); err != nil && err.Error() != db.ExistError {
				logger.Error("save error block failed", logger.String("err", err.Error()))
			}
			return nil, nil, err2
		}
	}
	txs := make([]*model.Tx, 0, len(resblock.Block.Txs))
	var docMsgs []model.TxMsg
	for _, tx := range resblock.Block.Txs {
		tx, msgs := ParseTx(tx, resblock.Block, client)
		if tx.Height > 0 {
			txs = append(txs, &tx)
			docMsgs = append(docMsgs, msgs...)
		}
	}
	return txs, docMsgs, nil
}

// parse iris tx from iris block result tx
func ParseTx(txBytes types.Tx, block *types.Block, client *pool.Client) (model.Tx, []model.TxMsg) {

	var (
		docMsgs   []model.TxMsg
		docTxMsgs []model.DocTxMsg
		docTx     model.Tx
		docFailTx model.ErrTx
		actualFee model.Coin
	)
	txHash := utils.BuildHex(txBytes.Hash())
	Tx, err := cdc.GetTxDecoder()(txBytes)
	if err != nil {
		docFailTx.Height = block.Height
		docFailTx.TxHash = txHash
		docFailTx.Log = fmt.Sprintf("TxDecoder have error:%v", err.Error())
		if err := docFailTx.Save(); err != nil && err.Error() != db.ExistError {
			logger.Error("save TxDecoder txs failed",
				logger.Int64("height", block.Height),
				logger.String("txhash", txHash),
				logger.String("err", err.Error()))
		}
		return docTx, docMsgs
	}
	authTx := Tx.(signing.Tx)
	fee := BuildFee(authTx.GetFee(), authTx.GetGas())
	memo := authTx.GetMemo()
	height := block.Height

	ctx := context.Background()
	res, err := client.Tx(ctx, txBytes.Hash(), false)
	if err != nil {
		time.Sleep(1 * time.Second)
		var err1 error
		client2 := pool.GetClient()
		res, err1 = client2.Tx(ctx, txBytes.Hash(), false)
		client2.Release()
		if err1 != nil {
			docFailTx.Height = height
			docFailTx.TxHash = txHash
			docFailTx.Log = err1.Error()
			if err := docFailTx.Save(); err != nil && err.Error() != db.ExistError {
				logger.Error("save txResult  txs failed",
					logger.Int64("height", block.Height),
					logger.String("txhash", txHash),
					logger.String("err", err.Error()))
			}
			return docTx, docMsgs
		}
	}

	if len(fee.Amount) > 0 {
		actualFee = fee.Amount[0]
	}

	docTx = model.Tx{
		Height:    height,
		Time:      block.Time.Unix(),
		TxHash:    txHash,
		Fee:       fee,
		ActualFee: actualFee,
		Memo:      memo,
		TxIndex:   res.Index,
	}
	docTx.Status = TxStatusSuccess
	if res.TxResult.Code != 0 {
		docTx.Status = TxStatusFail
		docTx.Log = res.TxResult.Log

	}
	docTx.Events = parseEvents(res.TxResult.Events)
	eventsIndexMap := make(map[int]model.MsgEvent)
	if res.TxResult.Code == 0 {
		eventsIndexMap = splitEvents(res.TxResult.Log)
	}

	msgs := authTx.GetMsgs()
	if len(msgs) == 0 {
		return docTx, docMsgs
	}
	for i, v := range msgs {
		msgDocInfo := HandleTxMsg(v)
		if len(msgDocInfo.Addrs) == 0 {
			continue
		}

		docTx.Signers = append(docTx.Signers, removeDuplicatesFromSlice(msgDocInfo.Signers)...)
		docTx.Addrs = append(docTx.Addrs, removeDuplicatesFromSlice(msgDocInfo.Addrs)...)
		docTxMsgs = append(docTxMsgs, msgDocInfo.DocTxMsg)
		docTx.Types = append(docTx.Types, msgDocInfo.DocTxMsg.Type)

		docMsg := model.TxMsg{
			Time:      docTx.Time,
			TxFee:     docTx.ActualFee,
			Height:    docTx.Height,
			TxHash:    docTx.TxHash,
			Type:      msgDocInfo.DocTxMsg.Type,
			MsgIndex:  i,
			TxIndex:   res.Index,
			TxStatus:  docTx.Status,
			TxMemo:    memo,
			TxLog:     docTx.Log,
			GasUsed:   res.TxResult.GasUsed,
			GasWanted: res.TxResult.GasWanted,
		}
		docMsg.Msg = msgDocInfo.DocTxMsg
		if val, ok := eventsIndexMap[i]; ok {
			docMsg.Events = val.Events
		}
		docMsg.Addrs = removeDuplicatesFromSlice(msgDocInfo.Addrs)
		docMsg.Signers = removeDuplicatesFromSlice(msgDocInfo.Signers)
		docMsgs = append(docMsgs, docMsg)

	}
	docTx.Addrs = removeDuplicatesFromSlice(docTx.Addrs)
	docTx.Types = removeDuplicatesFromSlice(docTx.Types)
	docTx.Signers = removeDuplicatesFromSlice(docTx.Signers)
	docTx.Msgs = docTxMsgs

	// don't save txs which have not parsed
	if docTx.TxHash == "" {
		return model.Tx{}, docMsgs
	}

	for i, _ := range docMsgs {
		docMsgs[i].TxAddrs = docTx.Addrs
		docMsgs[i].TxSigners = docTx.Signers
	}
	return docTx, docMsgs

}

func BuildFee(fee sdk.Coins, gas uint64) *model.Fee {
	return &model.Fee{
		Amount: model.BuildDocCoins(fee),
		Gas:    int64(gas),
	}
}

func parseEvents(events []aTypes.Event) []model.Event {
	var eventDocs []model.Event
	if len(events) > 0 {
		for _, e := range events {
			var kvPairDocs []model.KvPair
			if len(e.Attributes) > 0 {
				for _, v := range e.Attributes {
					kvPairDocs = append(kvPairDocs, model.KvPair{
						Key:   string(v.Key),
						Value: string(v.Value),
					})
				}
			}
			eventDocs = append(eventDocs, model.Event{
				Type:       e.Type,
				Attributes: kvPairDocs,
			})
		}
	}

	return eventDocs
}

func splitEvents(log string) map[int]model.MsgEvent {
	var eventDocs []model.MsgEvent
	if log != "" {
		utils.UnMarshalJsonIgnoreErr(log, &eventDocs)

	}

	msgIndexMap := make(map[int]model.MsgEvent, len(eventDocs))
	for _, val := range eventDocs {
		msgIndexMap[val.MsgIndex] = val
	}
	return msgIndexMap
}
