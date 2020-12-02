package block

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/irisnet/rainbow-sync/constant"
	"github.com/irisnet/rainbow-sync/db"
	"github.com/irisnet/rainbow-sync/lib/cdc"
	"github.com/irisnet/rainbow-sync/lib/pool"
	"github.com/irisnet/rainbow-sync/logger"
	"github.com/irisnet/rainbow-sync/model"
	"github.com/irisnet/rainbow-sync/utils"
	"github.com/tendermint/tendermint/types"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"time"
)

func SaveDocsWithTxn(blockDoc *model.Block, irisTxs []*model.Tx, taskDoc model.SyncTask) error {
	var (
		ops, irisTxsOps []txn.Op
	)

	if blockDoc.Height == 0 {
		return fmt.Errorf("invalid block, height equal 0")
	}

	blockOp := txn.Op{
		C:      model.CollectionNameBlock,
		Id:     bson.NewObjectId(),
		Insert: blockDoc,
	}

	length_txs := len(irisTxs)
	if length_txs > 0 {
		irisTxsOps = make([]txn.Op, 0, length_txs)
		for _, v := range irisTxs {
			op := txn.Op{
				C:      model.CollectionNameIrisTx,
				Id:     bson.NewObjectId(),
				Insert: v,
			}
			irisTxsOps = append(irisTxsOps, op)
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

	ops = make([]txn.Op, 0, length_txs+2)
	ops = append(append(ops, blockOp, updateOp), irisTxsOps...)

	if len(ops) > 0 {
		err := db.Txn(ops)
		if err != nil {
			return err
		}
	}

	return nil
}

func ParseBlock(b int64, client *pool.Client) (*model.Block, []*model.Tx, error) {

	defer func() {
		if err := recover(); err != nil {
			logger.Error("parse  block fail", logger.Int64("height", b),
				logger.Any("err", err))
		}
	}()

	resBlock := &model.Block{
		Height:     b,
		CreateTime: time.Now().Unix(),
	}

	txs, err := ParseTxs(b, client)
	if err != nil {
		return nil, nil, err
	}

	return resBlock, txs, nil
}

func ParseTxs(b int64, client *pool.Client) ([]*model.Tx, error) {
	resblock, err := client.Block(&b)
	if err != nil {
		logger.Warn("get block result err, now try again", logger.String("err", err.Error()),
			logger.Any("height", b))
		// there is possible parse block fail when in iterator
		var err2 error
		client2 := pool.GetClient()
		resblock, err2 = client2.Block(&b)
		client2.Release()
		if err2 != nil {
			return nil, err2
		}
	}
	txs := make([]*model.Tx, 0, len(resblock.Block.Txs))
	for _, tx := range resblock.Block.Txs {
		tx := ParseTx(tx, resblock.Block, client)
		if tx.Height > 0 {
			txs = append(txs, &tx)
		}
	}
	return txs, nil
}

// parse iris tx from iris block result tx
func ParseTx(txBytes types.Tx, block *types.Block, client *pool.Client) model.Tx {

	var (
		docTxMsgs  []model.DocTxMsg
		methodName = "ParseTx"
		docTx      model.Tx
		actualFee  *model.ActualFee
	)
	Tx, err := cdc.GetTxDecoder()(txBytes)
	if err != nil {
		logger.Error(err.Error())
		return docTx
	}
	authTx := Tx.(signing.Tx)
	fee := BuildFee(authTx.GetFee(), authTx.GetGas())
	memo := authTx.GetMemo()
	height := block.Height
	txHash := utils.BuildHex(txBytes.Hash())
	res, err := client.Tx(txBytes.Hash(), false)
	if err != nil {
		logger.Warn("QueryTxResult have error, now try again", logger.String("err", err.Error()))
		time.Sleep(time.Duration(1) * time.Second)
		var err1 error
		client2 := pool.GetClient()
		res, err1 = client2.Tx(txBytes.Hash(), false)
		client2.Release()
		if err1 != nil {
			logger.Error("get txResult err", logger.String("method", methodName), logger.String("err", err1.Error()))
		}
	}

	gasUsed := utils.Min(res.TxResult.GasUsed, fee.Gas)
	if len(fee.Amount) > 0 {
		gasPrice := utils.ParseFloat(fee.Amount[0].Amount) / float64(fee.Gas)
		actualFee = &model.ActualFee{
			Denom:  fee.Amount[0].Denom,
			Amount: fmt.Sprint(float64(gasUsed) * gasPrice),
		}
	} else {
		actualFee = &model.ActualFee{}
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
	docTx.Status = constant.TxStatusSuccess
	if res.TxResult.Code != 0 {
		docTx.Status = constant.TxStatusFail
		docTx.Log = res.TxResult.Log

	}

	msgs := authTx.GetMsgs()
	if len(msgs) == 0 {
		return docTx
	}
	for i, v := range msgs {
		msgDocInfo := HandleTxMsg(v)
		if len(msgDocInfo.Addrs) == 0 {
			continue
		}
		if i == 0 {
			docTx.Type = msgDocInfo.DocTxMsg.Type
		}

		docTx.Signers = append(docTx.Signers, removeDuplicatesFromSlice(msgDocInfo.Signers)...)
		docTx.Addrs = append(docTx.Addrs, removeDuplicatesFromSlice(msgDocInfo.Addrs)...)
		docTxMsgs = append(docTxMsgs, msgDocInfo.DocTxMsg)
		docTx.Types = append(docTx.Types, msgDocInfo.DocTxMsg.Type)
	}
	docTx.Addrs = removeDuplicatesFromSlice(docTx.Addrs)
	docTx.Types = removeDuplicatesFromSlice(docTx.Types)
	docTx.Signers = removeDuplicatesFromSlice(docTx.Signers)
	docTx.Msgs = docTxMsgs

	// don't save txs which have not parsed
	if docTx.TxHash == "" {
		return model.Tx{}
	}

	return docTx

}

func BuildFee(fee sdk.Coins, gas uint64) *model.Fee {
	return &model.Fee{
		Amount: model.BuildDocCoins(fee),
		Gas:    int64(gas),
	}
}
