package block

import (
	"github.com/irisnet/rainbow-sync/service/iris/logger"
	imodel "github.com/irisnet/rainbow-sync/service/iris/model"
	"github.com/irisnet/rainbow-sync/service/iris/utils"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/tendermint/tendermint/types"
	"github.com/irisnet/rainbow-sync/service/iris/helper"
	"github.com/irisnet/rainbow-sync/service/iris/constant"
)

// parse iris txs  from block result txs
func (iris *Iris_Block) ParseIrisTxs(b int64, client *helper.Client) ([]*imodel.IrisTx, error) {
	resblock, err := client.Block(&b)
	if err != nil {
		logger.Warn("get block result err, now try again", logger.String("err", err.Error()),
			logger.String("Chain Block", iris.Name()), logger.Any("height", b))
		// there is possible parse block fail when in iterator
		var err2 error
		client2 := helper.GetClient()
		resblock, err2 = client2.Block(&b)
		client2.Release()
		if err2 != nil {
			return nil, err2
		}
	}

	irisTxs := make([]*imodel.IrisTx, 0, len(resblock.Block.Txs))
	for _, tx := range resblock.Block.Txs {
		iristx := iris.ParseIrisTxModel(tx, resblock.Block)
		irisTxs = append(irisTxs, &iristx)
	}

	return irisTxs, nil
}

// parse iris tx from iris block result tx
func (iris *Iris_Block) ParseIrisTxModel(txBytes types.Tx, block *types.Block) imodel.IrisTx {

	var (
		authTx     auth.StdTx
		methodName = "ParseTx"
		docTx      imodel.IrisTx
		actualFee  *imodel.ActualFee
		docTxMsgs  []imodel.DocTxMsg
	)

	cdc := utils.GetCodec()

	err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &authTx)
	if err != nil {
		logger.Error(err.Error())
		return docTx
	}

	height := block.Height
	time := block.Time
	txHash := utils.BuildHex(txBytes.Hash())
	fee := utils.BuildFee(authTx.Fee)
	memo := authTx.Memo

	// get tx status, gasUsed, gasPrice and actualFee from tx result
	status, result, err := utils.QueryTxResult(txBytes.Hash())
	if err != nil {
		logger.Error("get txResult err", logger.String("method", methodName), logger.String("err", err.Error()))
	}
	gasUsed := Min(result.GasUsed, fee.Gas)
	if len(fee.Amount) > 0 {
		gasPrice := fee.Amount[0].Amount / float64(fee.Gas)
		actualFee = &imodel.ActualFee{
			Denom:  fee.Amount[0].Denom,
			Amount: float64(gasUsed) * gasPrice,
		}
	} else {
		actualFee = &imodel.ActualFee{}
	}
	msgs := authTx.GetMsgs()
	if len(msgs) <= 0 {
		logger.Error("can't get msgs", logger.String("method", methodName))
		return docTx
	}
	msg := msgs[0]

	docTx = imodel.IrisTx{
		Height:    height,
		Time:      time,
		TxHash:    txHash,
		Fee:       fee,
		ActualFee: actualFee,
		Memo:      memo,
		Status:    status,
		Code:      result.Code,
		Tags:      parseTags(result),
	}
	switch msg.(type) {
	case imodel.MsgTransfer:
		msg := msg.(imodel.MsgTransfer)

		docTx.From = msg.Inputs[0].Address.String()
		docTx.To = msg.Outputs[0].Address.String()
		docTx.Initiator = msg.Inputs[0].Address.String()
		docTx.Amount = utils.ParseCoins(msg.Inputs[0].Coins.String())
		docTx.Type = constant.Iris_TxTypeTransfer
	case imodel.MsgBurn:
		msg := msg.(imodel.MsgBurn)
		docTx.From = msg.Owner.String()
		docTx.To = ""
		docTx.Initiator = msg.Owner.String()
		docTx.Amount = utils.ParseCoins(msg.Coins.String())
		docTx.Type = constant.Iris_TxTypeBurn
	case imodel.MsgSetMemoRegexp:
		msg := msg.(imodel.MsgSetMemoRegexp)
		docTx.From = msg.Owner.String()
		docTx.To = ""
		docTx.Initiator = msg.Owner.String()
		docTx.Amount = []*imodel.Coin{}
		docTx.Type = constant.Iris_TxTypeSetMemoRegexp
		txMsg := imodel.DocTxMsgSetMemoRegexp{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, imodel.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx

	case imodel.MsgStakeCreate:
		msg := msg.(imodel.MsgStakeCreate)
		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()
		docTx.Initiator = msg.DelegatorAddr.String()
		docTx.Amount = []*imodel.Coin{utils.ParseCoin(msg.Delegation.String())}
		docTx.Type = constant.Iris_TxTypeStakeCreateValidator

	case imodel.MsgStakeEdit:
		msg := msg.(imodel.MsgStakeEdit)

		docTx.From = msg.ValidatorAddr.String()
		docTx.To = ""
		docTx.Initiator = msg.ValidatorAddr.String()
		docTx.Amount = []*imodel.Coin{}
		docTx.Type = constant.Iris_TxTypeStakeEditValidator

	case imodel.MsgStakeDelegate:
		msg := msg.(imodel.MsgStakeDelegate)

		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()
		docTx.Initiator = msg.DelegatorAddr.String()
		docTx.Amount = []*imodel.Coin{utils.ParseCoin(msg.Delegation.String())}
		docTx.Type = constant.Iris_TxTypeStakeDelegate

	case imodel.MsgStakeBeginUnbonding:
		msg := msg.(imodel.MsgStakeBeginUnbonding)

		shares := utils.ParseFloat(msg.SharesAmount.String())
		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()
		docTx.Initiator = msg.DelegatorAddr.String()
		coin := imodel.Coin{
			Amount: shares,
		}
		docTx.Amount = []*imodel.Coin{&coin}
		docTx.Type = constant.Iris_TxTypeStakeBeginUnbonding
	case imodel.MsgBeginRedelegate:
		msg := msg.(imodel.MsgBeginRedelegate)

		shares := utils.ParseFloat(msg.SharesAmount.String())
		docTx.From = msg.ValidatorSrcAddr.String()
		docTx.To = msg.ValidatorDstAddr.String()
		docTx.Initiator = msg.DelegatorAddr.String()
		coin := imodel.Coin{
			Amount: shares,
		}
		docTx.Amount = []*imodel.Coin{&coin}
		docTx.Type = constant.Iris_TxTypeBeginRedelegate
	case imodel.MsgUnjail:
		msg := msg.(imodel.MsgUnjail)

		docTx.From = msg.ValidatorAddr.String()
		docTx.Initiator = msg.ValidatorAddr.String()
		docTx.Type = constant.Iris_TxTypeUnjail
	case imodel.MsgSetWithdrawAddress:
		msg := msg.(imodel.MsgSetWithdrawAddress)

		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.WithdrawAddr.String()
		docTx.Initiator = msg.DelegatorAddr.String()
		docTx.Type = constant.Iris_TxTypeSetWithdrawAddress
	case imodel.MsgWithdrawDelegatorReward:
		msg := msg.(imodel.MsgWithdrawDelegatorReward)

		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()
		docTx.Initiator = msg.DelegatorAddr.String()
		docTx.Type = constant.Iris_TxTypeWithdrawDelegatorReward

		for _, tag := range result.Tags {
			key := string(tag.Key)
			if key == imodel.TagDistributionReward {
				reward := string(tag.Value)
				docTx.Amount = utils.ParseCoins(reward)
				break
			}
		}
	case imodel.MsgWithdrawDelegatorRewardsAll:
		msg := msg.(imodel.MsgWithdrawDelegatorRewardsAll)

		docTx.From = msg.DelegatorAddr.String()
		docTx.Initiator = msg.DelegatorAddr.String()
		docTx.Type = constant.Iris_TxTypeWithdrawDelegatorRewardsAll
		for _, tag := range result.Tags {
			key := string(tag.Key)
			if key == imodel.TagDistributionReward {
				reward := string(tag.Value)
				docTx.Amount = utils.ParseCoins(reward)
				break
			}
		}
	case imodel.MsgWithdrawValidatorRewardsAll:
		msg := msg.(imodel.MsgWithdrawValidatorRewardsAll)

		docTx.From = msg.ValidatorAddr.String()
		docTx.Initiator = msg.ValidatorAddr.String()
		docTx.Type = constant.Iris_TxTypeWithdrawValidatorRewardsAll
		for _, tag := range result.Tags {
			key := string(tag.Key)
			if key == imodel.TagDistributionReward {
				reward := string(tag.Value)
				docTx.Amount = utils.ParseCoins(reward)
				break
			}
		}
	case imodel.MsgSubmitProposal:
		msg := msg.(imodel.MsgSubmitProposal)

		docTx.From = msg.Proposer.String()
		docTx.To = ""
		docTx.Initiator = msg.Proposer.String()
		docTx.Amount = utils.ParseCoins(msg.InitialDeposit.String())
		docTx.Type = constant.Iris_TxTypeSubmitProposal

	case imodel.MsgSubmitSoftwareUpgradeProposal:
		msg := msg.(imodel.MsgSubmitSoftwareUpgradeProposal)

		docTx.From = msg.Proposer.String()
		docTx.To = ""
		docTx.Initiator = msg.Proposer.String()
		docTx.Amount = utils.ParseCoins(msg.InitialDeposit.String())
		docTx.Type = constant.Iris_TxTypeSubmitProposal

	case imodel.MsgSubmitTaxUsageProposal:
		msg := msg.(imodel.MsgSubmitTaxUsageProposal)

		docTx.From = msg.Proposer.String()
		docTx.To = ""
		docTx.Initiator = msg.Proposer.String()
		docTx.Amount = utils.ParseCoins(msg.InitialDeposit.String())
		docTx.Type = constant.Iris_TxTypeSubmitProposal

	case imodel.MsgSubmitTokenAdditionProposal:
		msg := msg.(imodel.MsgSubmitTokenAdditionProposal)

		docTx.From = msg.Proposer.String()
		docTx.To = ""
		docTx.Initiator = msg.Proposer.String()
		docTx.Amount = utils.ParseCoins(msg.InitialDeposit.String())
		docTx.Type = constant.Iris_TxTypeSubmitProposal
		txMsg := imodel.DocTxMsgSubmitTokenAdditionProposal{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, imodel.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx

	case imodel.MsgDeposit:
		msg := msg.(imodel.MsgDeposit)

		docTx.From = msg.Depositor.String()
		docTx.Initiator = msg.Depositor.String()
		docTx.Amount = utils.ParseCoins(msg.Amount.String())
		docTx.Type = constant.Iris_TxTypeDeposit
	case imodel.MsgVote:
		msg := msg.(imodel.MsgVote)

		docTx.From = msg.Voter.String()
		docTx.Initiator = msg.Voter.String()
		docTx.Amount = []*imodel.Coin{}
		docTx.Type = constant.Iris_TxTypeVote

	case imodel.MsgRequestRand:
		msg := msg.(imodel.MsgRequestRand)

		docTx.From = msg.Consumer.String()
		docTx.Initiator = msg.Consumer.String()
		docTx.Amount = []*imodel.Coin{}
		docTx.Type = constant.Iris_TxTypeRequestRand
		txMsg := imodel.DocTxMsgRequestRand{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, imodel.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

		return docTx

	case imodel.AssetIssueToken:
		msg := msg.(imodel.AssetIssueToken)

		docTx.From = msg.Owner.String()
		docTx.Type = constant.TxTypeAssetIssueToken
		txMsg := imodel.DocTxMsgIssueToken{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, imodel.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

		return docTx
	case imodel.AssetEditToken:
		msg := msg.(imodel.AssetEditToken)

		docTx.From = msg.Owner.String()
		docTx.Type = constant.TxTypeAssetEditToken
		txMsg := imodel.DocTxMsgEditToken{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, imodel.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

		return docTx
	case imodel.AssetMintToken:
		msg := msg.(imodel.AssetMintToken)

		docTx.From = msg.Owner.String()
		docTx.To = msg.To.String()
		docTx.Type = constant.TxTypeAssetMintToken
		txMsg := imodel.DocTxMsgMintToken{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, imodel.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

		return docTx
	case imodel.AssetTransferTokenOwner:
		msg := msg.(imodel.AssetTransferTokenOwner)

		docTx.From = msg.SrcOwner.String()
		docTx.To = msg.DstOwner.String()
		docTx.Type = constant.TxTypeAssetTransferTokenOwner
		txMsg := imodel.DocTxMsgTransferTokenOwner{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, imodel.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

		return docTx
	case imodel.AssetCreateGateway:
		msg := msg.(imodel.AssetCreateGateway)

		docTx.From = msg.Owner.String()
		docTx.Type = constant.TxTypeAssetCreateGateway
		txMsg := imodel.DocTxMsgCreateGateway{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, imodel.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

		return docTx
	case imodel.AssetEditGateWay:
		msg := msg.(imodel.AssetEditGateWay)

		docTx.From = msg.Owner.String()
		docTx.Type = constant.TxTypeAssetEditGateway
		txMsg := imodel.DocTxMsgEditGateway{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, imodel.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})

		return docTx
	case imodel.AssetTransferGatewayOwner:
		msg := msg.(imodel.AssetTransferGatewayOwner)

		docTx.From = msg.Owner.String()
		docTx.To = msg.To.String()
		docTx.Type = constant.TxTypeAssetTransferGatewayOwner
		txMsg := imodel.DocTxMsgTransferGatewayOwner{}
		txMsg.BuildMsg(msg)
		docTx.Msgs = append(docTxMsgs, imodel.DocTxMsg{
			Type: txMsg.Type(),
			Msg:  &txMsg,
		})
		return docTx
	default:
		logger.Warn("unknown msg type")
	}

	return docTx

}

func parseTags(result abci.ResponseDeliverTx) map[string]string {
	tags := make(map[string]string, 0)
	for _, tag := range result.Tags {
		key := string(tag.Key)
		value := string(tag.Value)
		tags[key] = value
	}
	return tags
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
