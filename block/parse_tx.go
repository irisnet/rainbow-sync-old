package block

import (
	"fmt"
	"github.com/irisnet/rainbow-sync/db"
	"github.com/irisnet/rainbow-sync/lib/logger"
	"github.com/irisnet/rainbow-sync/lib/msgparser"
	"github.com/irisnet/rainbow-sync/lib/pool"
	"github.com/irisnet/rainbow-sync/model"
	"github.com/irisnet/rainbow-sync/utils"
	"github.com/kaifei-bianjie/common-parser/codec"
	msgsdktypes "github.com/kaifei-bianjie/common-parser/types"
	. "github.com/kaifei-bianjie/cosmosmod-parser/modules"
	"github.com/kaifei-bianjie/cosmosmod-parser/modules/ibc"
	aTypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"time"
)

var _parser msgparser.MsgParser

func init() {
	router := msgparser.RegisteRouter()
	//if conf.Server.OnlySupportModule != "" {
	//	modules := strings.Split(conf.Server.OnlySupportModule, ",")
	//	msgRoute := msgparser.NewRouter()
	//	for _, one := range modules {
	//		fn, exist := msgparser.RouteHandlerMap[one]
	//		if !exist {
	//			logger.Fatal("no support module: " + one)
	//		}
	//		msgRoute = msgRoute.AddRoute(one, fn)
	//		if one == msgparser.IbcRouteKey {
	//			msgRoute = msgRoute.AddRoute(msgparser.IbcTransferRouteKey, msgparser.RouteHandlerMap[one])
	//		}
	//	}
	//	if msgRoute.GetRoutesLen() > 0 {
	//		router = msgRoute
	//	}
	//
	//}
	_parser = msgparser.NewMsgParser(router)
}
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

func ParseBlock(b int64, client *pool.Client) (*model.Block, []*model.Tx, []model.TxMsg, error) {

	defer func() {
		if err := recover(); err != nil {
			logger.Error("parse  block fail", logger.Int64("height", b),
				logger.Any("err", err))
		}
	}()
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
			return nil, nil, nil, utils.ConvertErr(b, "", "ParseBlock", err2)
		}
	}
	blockDoc := model.Block{
		Height:     b,
		CreateTime: time.Now().Unix(),
	}

	blockResults, err := client.BlockResults(ctx, &b)

	if err != nil {
		time.Sleep(1 * time.Second)
		blockResults, err = client.BlockResults(ctx, &b)
		if err != nil {
			return &blockDoc, nil, nil, utils.ConvertErr(b, "", "ParseBlockResult", err)
		}
	}

	if len(resblock.Block.Txs) != len(blockResults.TxsResults) {
		return nil, nil, nil, utils.ConvertErr(b, "", "block.Txs length not equal blockResult", nil)
	}

	txs := make([]*model.Tx, 0, len(resblock.Block.Txs))
	var docMsgs []model.TxMsg
	for index, tx := range resblock.Block.Txs {
		txResult := blockResults.TxsResults[index]
		tx, msgs, err := ParseTx(tx, txResult, resblock.Block, index)
		if err != nil {
			return &blockDoc, txs, docMsgs, err
		}
		if tx.Height > 0 {
			txs = append(txs, &tx)
			docMsgs = append(docMsgs, msgs...)
		}
	}
	return &blockDoc, txs, docMsgs, nil
}

// parse iris tx from iris block result tx
func ParseTx(txBytes types.Tx, txResult *aTypes.ResponseDeliverTx, block *types.Block, index int) (model.Tx, []model.TxMsg, error) {

	var (
		docMsgs   []model.TxMsg
		docTxMsgs []msgsdktypes.TxMsg
		docTx     model.Tx
		actualFee msgsdktypes.Coin
	)
	height := block.Height
	txHash := utils.BuildHex(txBytes.Hash())
	authTx, err := codec.GetSigningTx(txBytes)
	if err != nil {
		logger.Warn(err.Error(),
			logger.String("errTag", "TxDecoder"),
			logger.String("txhash", txHash),
			logger.Int64("height", block.Height))
		return docTx, docMsgs, nil
	}
	fee := msgsdktypes.BuildFee(authTx.GetFee(), authTx.GetGas())
	memo := authTx.GetMemo()

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
		TxIndex:   uint32(index),
		TxId:      buildTxId(height, uint32(index)),
	}
	docTx.Status = utils.TxStatusSuccess
	if txResult.Code != 0 {
		docTx.Status = utils.TxStatusFail
		docTx.Log = txResult.Log

	}
	docTx.Events = parseEvents(txResult.Events)
	eventsIndexMap := make(map[int]model.MsgEvent)
	if txResult.Code == 0 {
		eventsIndexMap = splitEvents(txResult.Log)
	}

	msgs := authTx.GetMsgs()
	if len(msgs) == 0 {
		return docTx, docMsgs, nil
	}
	for i, v := range msgs {
		msgDocInfo := _parser.HandleTxMsg(v)
		if len(msgDocInfo.Addrs) == 0 {
			continue
		}

		switch msgDocInfo.DocTxMsg.Type {
		case MsgTypeIBCTransfer:
			if ibcTranferMsg, ok := msgDocInfo.DocTxMsg.Msg.(*ibc.DocMsgTransfer); ok {
				if val, exist := eventsIndexMap[i]; exist {
					ibcTranferMsg.PacketId = buildPacketId(val.Events)
					msgDocInfo.DocTxMsg.Msg = ibcTranferMsg
				}

			} else {
				logger.Warn("ibc transfer handler packet_id failed", logger.String("errTag", "TxMsg"),
					logger.String("txhash", txHash),
					logger.Int("msg_index", i),
					logger.Int64("height", height))
			}
		}

		docTx.Signers = append(docTx.Signers, utils.RemoveDuplicatesFromSlice(msgDocInfo.Signers)...)
		docTx.Addrs = append(docTx.Addrs, utils.RemoveDuplicatesFromSlice(msgDocInfo.Addrs)...)
		docTxMsgs = append(docTxMsgs, msgDocInfo.DocTxMsg)
		docTx.Types = append(docTx.Types, msgDocInfo.DocTxMsg.Type)

		docMsg := model.TxMsg{
			Time:      docTx.Time,
			TxFee:     docTx.ActualFee,
			Height:    docTx.Height,
			TxHash:    docTx.TxHash,
			Type:      msgDocInfo.DocTxMsg.Type,
			MsgIndex:  i,
			TxIndex:   uint32(index),
			TxStatus:  docTx.Status,
			TxMemo:    memo,
			TxLog:     docTx.Log,
			GasUsed:   txResult.GasUsed,
			GasWanted: txResult.GasWanted,
		}
		docMsg.Msg = msgDocInfo.DocTxMsg
		if val, ok := eventsIndexMap[i]; ok {
			docMsg.Events = val.Events
		}
		docMsg.Addrs = utils.RemoveDuplicatesFromSlice(msgDocInfo.Addrs)
		docMsg.Signers = utils.RemoveDuplicatesFromSlice(msgDocInfo.Signers)
		docMsg.Denoms = msgDocInfo.Denoms
		docMsgs = append(docMsgs, docMsg)

	}
	docTx.Addrs = utils.RemoveDuplicatesFromSlice(docTx.Addrs)
	docTx.Types = utils.RemoveDuplicatesFromSlice(docTx.Types)
	docTx.Signers = utils.RemoveDuplicatesFromSlice(docTx.Signers)
	docTx.Msgs = docTxMsgs

	// don't save txs which have not parsed
	if len(docTx.Addrs) == 0 {
		logger.Warn(utils.NoSupportMsgTypeTag,
			logger.String("errTag", "TxMsg"),
			logger.String("txhash", txHash),
			logger.Int64("height", height))
		return docTx, docMsgs, nil
	}

	for i, _ := range docMsgs {
		docMsgs[i].TxAddrs = docTx.Addrs
		docMsgs[i].TxSigners = docTx.Signers
	}
	return docTx, docMsgs, nil

}

// unique index: (height,tx_index)
// txIndex: max value is 9999
// return height*10000+tx_index
func buildTxId(height int64, txIndex uint32) uint64 {
	if txIndex > 9999 {
		logger.Warn("build TxId failed for only support txIndex max value is 9999",
			logger.Int64("height", height),
			logger.Uint32("tx_index", txIndex))
		return uint64(height*10000 + 9999)
	}
	return uint64(height*10000) + uint64(txIndex)
}

func buildPacketId(events []model.Event) string {
	if len(events) > 0 {
		var mapKeyValue map[string]string
		for _, e := range events {
			if len(e.Attributes) > 0 && e.Type == utils.IbcTransferEventTypeSendPacket {
				mapKeyValue = make(map[string]string, len(e.Attributes))
				for _, v := range e.Attributes {
					mapKeyValue[v.Key] = v.Value
				}
				break
			}
		}

		if len(mapKeyValue) > 0 {
			scPort := mapKeyValue[utils.IbcTransferEventAttriKeyPacketScPort]
			scChannel := mapKeyValue[utils.IbcTransferEventAttriKeyPacketScChannel]
			dcPort := mapKeyValue[utils.IbcTransferEventAttriKeyPacketDcPort]
			dcChannel := mapKeyValue[utils.IbcTransferEventAttriKeyPacketDcChannels]
			sequence := mapKeyValue[utils.IbcTransferEventAttriKeyPacketSequence]
			return fmt.Sprintf("%v%v%v%v%v", scPort, scChannel, dcPort, dcChannel, sequence)
		}
	}
	return ""
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
