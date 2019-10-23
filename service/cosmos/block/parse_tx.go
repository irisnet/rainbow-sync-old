package block

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/irisnet/rainbow-sync/service/cosmos/constant"
	model "github.com/irisnet/rainbow-sync/service/cosmos/db"
	"github.com/irisnet/rainbow-sync/service/cosmos/helper"
	"github.com/irisnet/rainbow-sync/service/cosmos/logger"
	cmodel "github.com/irisnet/rainbow-sync/service/cosmos/model"
	imsg "github.com/irisnet/rainbow-sync/service/cosmos/model/msg"
	cutils "github.com/irisnet/rainbow-sync/service/cosmos/utils"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"strconv"
	"time"
)

const (
	COSMOS = "Cosmos"
)

type CosmosBlock struct{}

func (cosmos *CosmosBlock) Name() string {
	return COSMOS
}

func (cosmos *CosmosBlock) SaveDocsWithTxn(blockDoc *cmodel.Block, cosmosTxs []cmodel.CosmosTx, taskDoc cmodel.SyncCosmosTask) error {
	var (
		ops, cosmosTxsOps []txn.Op
	)

	if blockDoc.Height == 0 {
		return fmt.Errorf("invalid block, height equal 0")
	}

	blockOp := txn.Op{
		C:      cmodel.CollectionNameBlock,
		Id:     bson.NewObjectId(),
		Insert: blockDoc,
	}

	if length := len(cosmosTxs); length > 0 {

		cosmosTxsOps = make([]txn.Op, 0, length)
		for _, v := range cosmosTxs {
			op := txn.Op{
				C:      cmodel.CollectionNameCosmosTx,
				Id:     bson.NewObjectId(),
				Insert: v,
			}
			cosmosTxsOps = append(cosmosTxsOps, op)
		}
	}

	updateOp := txn.Op{
		C:      cmodel.CollectionNameSyncCosmosTask,
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

	ops = make([]txn.Op, 0, len(cosmosTxsOps)+2)
	ops = append(append(ops, blockOp, updateOp), cosmosTxsOps...)

	if len(ops) > 0 {
		err := model.Txn(ops)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cosmos *CosmosBlock) ParseBlock(b int64, client *cosmoshelper.CosmosClient) (resBlock *cmodel.Block, cosmosTxs []cmodel.CosmosTx, resErr error) {

	defer func() {
		if err := recover(); err != nil {
			logger.Error("parse cosmos block fail", logger.Int64("height", b),
				logger.Any("err", err), logger.String("Chain Block", cosmos.Name()))

			resBlock = &cmodel.Block{}
			cosmosTxs = nil
			resErr = fmt.Errorf("%v", err)
		}
	}()

	cosmosTxsdata, err := cosmos.ParseCosmosTxs(b, client)
	if err != nil {
		logger.Error("parse cosmos asset error", logger.String("error", err.Error()),
			logger.String("Chain Block", cosmos.Name()))
	}

	resBlock = &cmodel.Block{
		Height:     b,
		CreateTime: time.Now().Unix(),
	}
	cosmosTxs = cosmosTxsdata
	resErr = err
	return
}

// parse cosmos txs  from block result txs
func (cosmos *CosmosBlock) ParseCosmosTxs(b int64, client *cosmoshelper.CosmosClient) ([]cmodel.CosmosTx, error) {
	resblock, err := client.Block(&b)
	if err != nil {
		logger.Warn("get block result err, now try again", logger.String("err", err.Error()),
			logger.String("Chain Block", cosmos.Name()))
		// there is possible parse block fail when in iterator
		var err2 error
		client2 := cosmoshelper.GetCosmosClient()
		resblock, err2 = client2.Block(&b)
		client2.Release()
		if err2 != nil {
			return nil, err2
		}
	}

	//fmt.Printf("======>>resblock.Block.Txs:%+v\n",resblock.Block.Txs)
	//fmt.Println("length:",len(resblock.Block.Txs))

	cosmosTxs := make([]cmodel.CosmosTx, 0, len(resblock.Block.Txs))
	for _, tx := range resblock.Block.Txs {
		if cosmostx := cosmos.ParseCosmosTxModel(tx, resblock.Block); len(cosmostx) > 0 {
			cosmosTxs = append(cosmosTxs, cosmostx...)
		}
	}

	return cosmosTxs, nil
}

func (cosmos *CosmosBlock) ParseCosmosTxModel(txBytes types.Tx, block *types.Block) []cmodel.CosmosTx {
	var (
		authTx     auth.StdTx
		methodName = "parseCosmosTxModel"
		txdetail   cmodel.CosmosTx
		docTxMsgs  []cmodel.DocTxMsg
	)

	cdc := cutils.GetCodec()
	err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &authTx)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	status, result, err := QueryTxResult(txBytes.Hash())
	if err != nil {
		logger.Error("get txResult err", logger.String("method", methodName),
			logger.String("err", err.Error()),
			logger.String("Chain Block", cosmos.Name()))
	}
	//msgStat, err := parseRawlog(result.Log)
	//if err != nil {
	//	logger.Error("get parseRawlog err", logger.String("method", methodName),
	//		logger.String("err", err.Error()),
	//		logger.String("Chain Block", cosmos.Name()))
	//}

	fee := cutils.BuildFee(authTx.Fee)
	txdetail.TxHash = cutils.BuildHex(txBytes.Hash())
	txdetail.Height = block.Height
	txdetail.Memo = authTx.Memo
	txdetail.Fee = &fee
	txdetail.Time = block.Time
	txdetail.Status = status
	txdetail.Code = result.Code
	txdetail.Events = parseEvents(result)

	//length_msgStat := len(msgStat)

	msgs := authTx.GetMsgs()
	lenMsgs := len(msgs)
	if lenMsgs <= 0 {
		logger.Error("can't get msgs", logger.String("method", methodName),
			logger.String("Chain Block", cosmos.Name()))
		return nil
	}
	txs := make([]cmodel.CosmosTx, 0, lenMsgs)
	for _, msg := range msgs {
		txdetail.Initiator = ""
		txdetail.From = ""
		txdetail.To = ""
		txdetail.Amount = nil
		txdetail.Type = ""
		//if length_msgStat > i {
		//	txdetail.Status = msgStat[i]
		//}
		switch msg.(type) {

		case cmodel.MsgTransfer:
			msg := msg.(cmodel.MsgTransfer)
			txdetail.Initiator = msg.FromAddress.String()
			txdetail.From = msg.FromAddress.String()
			txdetail.To = msg.ToAddress.String()
			txdetail.Amount = cutils.ParseCoins(msg.Amount)
			txdetail.Type = constant.TxTypeTransfer

			docTxMsg := imsg.DocTxMsgTransfer{}
			docTxMsg.BuildMsg(msg)
			txdetail.Msgs = append(docTxMsgs, cmodel.DocTxMsg{
				Type: docTxMsg.Type(),
				Msg:  &docTxMsg,
			})
			break

		case cmodel.IBCBankMsgTransfer:
			msg := msg.(cmodel.IBCBankMsgTransfer)
			txdetail.Initiator = msg.Sender
			txdetail.From = txdetail.Initiator
			txdetail.To = msg.Receiver
			txdetail.Amount = buildCoins(msg.Denomination, msg.Amount.String())
			txdetail.Type = constant.TxTypeIBCBankTransfer
			txdetail.IBCPacketHash = buildIBCPacketHashByEvents(txdetail.Events)
			txMsg := imsg.DocTxMsgIBCBankTransfer{}
			txMsg.BuildMsg(msg)
			txdetail.Msgs = append(docTxMsgs, cmodel.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})
			break
		case cmodel.IBCBankMsgReceivePacket:
			msg := msg.(cmodel.IBCBankMsgReceivePacket)
			txdetail.Initiator = msg.Signer.String()
			txdetail.Type = constant.TxTypeIBCBankRecvTransferPacket

			if transPacketData, err := buildIBCPacketData(msg.Packet.Data()); err != nil {
				logger.Error("build ibc packet data fail", logger.String("packetData", string(msg.Packet.Data())),
					logger.String("err", err.Error()))
			} else {
				txdetail.From = transPacketData.Sender
				txdetail.To = transPacketData.Receiver
				txdetail.Amount = buildCoins(transPacketData.Denomination, transPacketData.Amount)
			}

			if hash, err := buildIBCPacketHashByPacket(msg.Packet.(cmodel.IBCPacket)); err != nil {
				logger.Error("build ibc packet hash fail", logger.String("err", err.Error()))
			} else {
				txdetail.IBCPacketHash = hash
			}

			txMsg := imsg.DocTxMsgIBCBankReceivePacket{}
			txMsg.BuildMsg(msg)
			txdetail.Msgs = append(docTxMsgs, cmodel.DocTxMsg{
				Type: txMsg.Type(),
				Msg:  &txMsg,
			})
			break
		default:
			logger.Warn("unknown msg type")
		}
	}
	txs = append(txs, txdetail)

	return txs
}

// get tx status and log by query txHash
func QueryTxResult(txHash []byte) (string, *abci.ResponseDeliverTx, error) {
	status := constant.TxStatusSuccess

	client := cosmoshelper.GetCosmosClient()
	defer client.Release()

	res, err := client.Tx(txHash, false)
	if err != nil {
		return "unknown", nil, err
	}
	result := res.TxResult
	if result.Code != 0 {
		status = constant.TxStatusFail
	}

	return status, &result, nil
}

func parseEvents(result *abci.ResponseDeliverTx) []cmodel.Event {

	var events []cmodel.Event
	for _, val := range result.GetEvents() {
		one := cmodel.Event{
			Type: val.Type,
		}
		one.Attributes = make(map[string]string, len(val.Attributes))
		for _, attr := range val.Attributes {
			one.Attributes[string(attr.Key)] = string(attr.Value)
		}
		events = append(events, one)
	}

	return events
}

func buildCoins(denom string, amountStr string) []*cmodel.Coin {
	var coins []*cmodel.Coin
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		logger.Error("convert str to float64 fail", logger.String("amountStr", amountStr),
			logger.String("err", err.Error()))
		amount = 0
	}
	coin := cmodel.Coin{
		Denom:  denom,
		Amount: amount,
	}
	return append(coins, &coin)
}

func buildIBCPacketHashByEvents(events []cmodel.Event) string {
	var packetStr string
	if len(events) == 0 {
		return ""
	}

	for _, e := range events {
		if e.Type == constant.EventTypeSendPacket {
			for k, v := range e.Attributes {
				if k == constant.EventAttributesKeyPacket {
					packetStr = v
					break
				}
			}
		}
	}

	if packetStr == "" {
		return ""
	}

	return cutils.Md5Encrypt([]byte(packetStr))
}

func buildIBCPacketHashByPacket(packet cmodel.IBCPacket) (string, error) {
	data, err := packet.MarshalJSON()
	if err != nil {
		return "", err
	}
	return cutils.Md5Encrypt(data), nil
}

func buildIBCPacketData(packetData []byte) (cmodel.IBCTransferPacketDataValue, error) {
	var transferPacketData cmodel.IBCTransferPacketData
	err := json.Unmarshal(packetData, &transferPacketData)
	if err != nil {
		return transferPacketData.Value, err
	}

	return transferPacketData.Value, nil
}

func parseRawlog(rawlog string) (map[int]string, error) {

	var Stats []cmodel.RawLog
	if err := json.Unmarshal([]byte(rawlog), &Stats); err != nil {
		return nil, err
	}

	msgStat := make(map[int]string, len(Stats))
	for _, stat := range Stats {
		if stat.Success {
			msgStat[stat.MsgIndex] = constant.TxStatusSuccess
		} else {
			msgStat[stat.MsgIndex] = constant.TxStatusFail
		}

	}
	return msgStat, nil
}
