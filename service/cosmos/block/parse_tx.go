package block

import (
	model "github.com/irisnet/rainbow-sync/service/cosmos/db"
	cmodel "github.com/irisnet/rainbow-sync/service/cosmos/model"
	"github.com/irisnet/rainbow-sync/service/cosmos/helper"
	"github.com/irisnet/rainbow-sync/service/cosmos/logger"
	"github.com/irisnet/rainbow-sync/service/cosmos/constant"
	"github.com/tendermint/tendermint/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	"fmt"
	"time"
	"gopkg.in/mgo.v2/txn"
	"gopkg.in/mgo.v2/bson"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cutils "github.com/irisnet/rainbow-sync/service/cosmos/utils"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"encoding/json"
)

const (
	COSMOS = "Cosmos"
)

type (
	MsgTransfer = bank.MsgSend
	MsgMultiSend = bank.MsgMultiSend

	MsgUnjail = slashing.MsgUnjail
	MsgSetWithdrawAddress = dtypes.MsgSetWithdrawAddress
	MsgWithdrawDelegatorReward = dtypes.MsgWithdrawDelegatorReward
	MsgWithdrawValidatorCommission = dtypes.MsgWithdrawValidatorCommission

	MsgDeposit = gov.MsgDeposit
	MsgSubmitProposal = gov.MsgSubmitProposal
	MsgVote = gov.MsgVote
	Proposal = gov.Proposal

	MsgVerifyInvariant = crisis.MsgVerifyInvariant

	MsgDelegate = stypes.MsgDelegate
	MsgUndelegate = stypes.MsgUndelegate
	MsgBeginRedelegate = stypes.MsgBeginRedelegate
	MsgCreateValidator = stypes.MsgCreateValidator
	MsgEditValidator = stypes.MsgEditValidator
)

type Cosmos_Block struct{}

func (cosmos *Cosmos_Block) Name() string {
	return COSMOS
}

func (cosmos *Cosmos_Block) SaveDocsWithTxn(blockDoc *cmodel.Block, cosmosTxs []cmodel.CosmosTx, taskDoc cmodel.SyncCosmosTask) error {
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

func (cosmos *Cosmos_Block) ParseBlock(b int64, client *cosmoshelper.CosmosClient) (resBlock *cmodel.Block, cosmosTxs []cmodel.CosmosTx, resErr error) {

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
func (cosmos *Cosmos_Block) ParseCosmosTxs(b int64, client *cosmoshelper.CosmosClient) ([]cmodel.CosmosTx, error) {
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

func (cosmos *Cosmos_Block) ParseCosmosTxModel(txBytes types.Tx, block *types.Block) []cmodel.CosmosTx {
	var (
		authTx     auth.StdTx
		methodName = "parseCosmosTxModel"
		txdetail   cmodel.CosmosTx
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
	msgStat, err := parseRawlog(result.Log)
	if err != nil {
		logger.Error("get parseRawlog err", logger.String("method", methodName),
			logger.String("err", err.Error()),
			logger.String("Chain Block", cosmos.Name()))
	}

	fee := cutils.BuildFee(authTx.Fee)
	txdetail.TxHash = cutils.BuildHex(txBytes.Hash())
	txdetail.Height = block.Height
	txdetail.Memo = authTx.Memo
	txdetail.Fee = &fee
	txdetail.Time = block.Time
	txdetail.Status = status
	txdetail.Code = result.Code

	length_msgStat := len(msgStat)

	msgs := authTx.GetMsgs()
	len_msgs := len(msgs)
	if len_msgs <= 0 {
		logger.Error("can't get msgs", logger.String("method", methodName),
			logger.String("Chain Block", cosmos.Name()))
		return nil
	}
	txs := make([]cmodel.CosmosTx, 0, len_msgs)
	for i, msg := range msgs {
		txdetail.Initiator = ""
		txdetail.From = ""
		txdetail.To = ""
		txdetail.Amount = nil
		txdetail.Type = ""
		txdetail.Events = parseEvents(result)
		if length_msgStat > i {
			txdetail.Status = msgStat[i]
		}
		switch msg.(type) {
		case MsgDelegate:
			msg := msg.(MsgDelegate)
			txdetail.Initiator = msg.DelegatorAddress.String()
			txdetail.From = msg.DelegatorAddress.String()
			txdetail.To = msg.ValidatorAddress.String()
			txdetail.Amount = cutils.ParseCoins(sdk.Coins{msg.Amount})
			txdetail.Type = constant.Cosmos_TxTypeStakeDelegate
			txs = append(txs, txdetail)

		case MsgUndelegate:
			msg := msg.(MsgUndelegate)
			txdetail.Initiator = msg.DelegatorAddress.String()
			txdetail.From = msg.DelegatorAddress.String()
			txdetail.To = msg.ValidatorAddress.String()
			txdetail.Amount = cutils.ParseCoins(sdk.Coins{msg.Amount})
			txdetail.Type = constant.Cosmos_TxTypeStakeUnDelegate
			txs = append(txs, txdetail)

		case MsgEditValidator:
			msg := msg.(MsgEditValidator)
			txdetail.Initiator = msg.ValidatorAddress.String()
			txdetail.From = msg.ValidatorAddress.String()
			txdetail.To = ""
			txdetail.Amount = []*cmodel.Coin{}
			txdetail.Type = constant.Cosmos_TxTypeStakeEditValidator
			txs = append(txs, txdetail)

		case MsgCreateValidator:
			msg := msg.(MsgCreateValidator)
			txdetail.Initiator = msg.DelegatorAddress.String()
			txdetail.From = msg.DelegatorAddress.String()
			txdetail.To = msg.ValidatorAddress.String()
			txdetail.Amount = cutils.ParseCoins(sdk.Coins{msg.Value})
			txdetail.Type = constant.Cosmos_TxTypeStakeCreateValidator
			txs = append(txs, txdetail)

		case MsgBeginRedelegate:
			msg := msg.(MsgBeginRedelegate)
			txdetail.Initiator = msg.DelegatorAddress.String()
			txdetail.From = msg.ValidatorSrcAddress.String()
			txdetail.To = msg.ValidatorDstAddress.String()
			txdetail.Amount = cutils.ParseCoins(sdk.Coins{msg.Amount})
			txdetail.Type = constant.Cosmos_TxTypeBeginRedelegate
			txs = append(txs, txdetail)

		case MsgTransfer:
			msg := msg.(MsgTransfer)
			txdetail.Initiator = msg.FromAddress.String()
			txdetail.From = msg.FromAddress.String()
			txdetail.To = msg.ToAddress.String()
			txdetail.Amount = cutils.ParseCoins(msg.Amount)
			txdetail.Type = constant.Cosmos_TxTypeTransfer
			txs = append(txs, txdetail)

		case MsgMultiSend:
			msg := msg.(MsgMultiSend)
			txdetail.Initiator = msg.Inputs[0].Address.String()
			txdetail.From = msg.Inputs[0].Address.String()
			txdetail.To = msg.Outputs[0].Address.String()
			txdetail.Amount = cutils.ParseCoins(msg.Inputs[0].Coins)
			txdetail.Type = constant.Cosmos_TxTypeMultiSend
			txs = append(txs, txdetail)

		case MsgVerifyInvariant:
			msg := msg.(MsgVerifyInvariant)
			txdetail.Initiator = msg.Sender.String()
			txdetail.From = msg.Sender.String()
			txdetail.To = ""
			txdetail.Amount = []*cmodel.Coin{}
			txdetail.Type = constant.Cosmos_TxTypeVerifyInvariant
			txs = append(txs, txdetail)

		case MsgUnjail:
			msg := msg.(MsgUnjail)
			txdetail.Initiator = msg.ValidatorAddr.String()
			txdetail.From = msg.ValidatorAddr.String()
			txdetail.Type = constant.Cosmos_TxTypeUnjail
			txs = append(txs, txdetail)
		case MsgSetWithdrawAddress:
			msg := msg.(MsgSetWithdrawAddress)
			txdetail.Initiator = msg.DelegatorAddress.String()
			txdetail.From = msg.DelegatorAddress.String()
			txdetail.To = msg.WithdrawAddress.String()
			txdetail.Type = constant.Cosmos_TxTypeSetWithdrawAddress
			txs = append(txs, txdetail)

		case MsgWithdrawDelegatorReward:
			msg := msg.(MsgWithdrawDelegatorReward)
			txdetail.Initiator = msg.DelegatorAddress.String()
			txdetail.From = msg.DelegatorAddress.String()
			txdetail.To = msg.ValidatorAddress.String()
			txdetail.Amount = []*cmodel.Coin{}
			coin := parseRewards(txdetail.Events)
			if coin != nil {
				txdetail.Amount = []*cmodel.Coin{coin}
			}
			txdetail.Type = constant.Cosmos_TxTypeWithdrawDelegatorReward
			txs = append(txs, txdetail)

		case MsgWithdrawValidatorCommission:
			msg := msg.(MsgWithdrawValidatorCommission)
			txdetail.Initiator = msg.ValidatorAddress.String()
			txdetail.From = msg.ValidatorAddress.String()
			txdetail.Type = constant.Cosmos_TxTypeWithdrawDelegatorRewardsAll
			txs = append(txs, txdetail)

		case MsgSubmitProposal:
			msg := msg.(MsgSubmitProposal)
			txdetail.Initiator = msg.Proposer.String()
			txdetail.From = msg.Proposer.String()
			txdetail.To = ""
			txdetail.Amount = cutils.ParseCoins(msg.InitialDeposit)
			txdetail.Type = constant.Cosmos_TxTypeSubmitProposal
			txs = append(txs, txdetail)

		case MsgDeposit:
			msg := msg.(MsgDeposit)
			txdetail.Initiator = msg.Depositor.String()
			txdetail.From = msg.Depositor.String()
			txdetail.Amount = cutils.ParseCoins(msg.Amount)
			txdetail.Type = constant.Cosmos_TxTypeDeposit
			txs = append(txs, txdetail)
		case MsgVote:
			msg := msg.(MsgVote)
			txdetail.Initiator = msg.Voter.String()
			txdetail.From = msg.Voter.String()
			txdetail.Amount = []*cmodel.Coin{}
			txdetail.Type = constant.Cosmos_TxTypeVote
			txs = append(txs, txdetail)

		default:
			logger.Warn("unknown msg type")
		}
	}

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

func parseRewards(events []cmodel.Event) (rewards *cmodel.Coin) {

	var totalrewards cmodel.Coin
	for _, val := range events {

		if val.Type == constant.Cosmos_TxEventWithdrawRewards {
			amount := cutils.ParseRewards(val.Attributes["amount"])
			if amount != nil {
				if totalrewards.Denom == "" {
					totalrewards.Denom = amount.Denom
				}
				totalrewards.Amount += amount.Amount
			}
		}
	}

	if totalrewards.Amount > 0 {
		rewards = &totalrewards
	}
	return
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
