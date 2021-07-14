package msgparser

import (
	"github.com/irisnet/rainbow-sync/lib/logger"
	"github.com/irisnet/rainbow-sync/utils"
	"github.com/kaifei-bianjie/msg-parser"
	. "github.com/kaifei-bianjie/msg-parser/modules"
	"github.com/kaifei-bianjie/msg-parser/modules/bank"
	"github.com/kaifei-bianjie/msg-parser/modules/coinswap"
	"github.com/kaifei-bianjie/msg-parser/modules/distribution"
	"github.com/kaifei-bianjie/msg-parser/modules/ibc"
	"github.com/kaifei-bianjie/msg-parser/modules/staking"
	"github.com/kaifei-bianjie/msg-parser/types"
	"strings"
)

type MsgParser interface {
	HandleTxMsg(v types.SdkMsg) CustomMsgDocInfo
}

var (
	_client msg_parser.MsgClient
)

func NewMsgParser(router Router) MsgParser {
	return &msgParser{
		rh: router,
	}
}

type msgParser struct {
	rh Router
}

func (parser *msgParser) HandleTxMsg(v types.SdkMsg) CustomMsgDocInfo {
	handleFunc, err := parser.rh.GetRoute(v.Route())
	if err != nil {
		logger.Error(err.Error(),
			logger.String("route", v.Route()),
			logger.String("type", v.Type()))
		return CustomMsgDocInfo{}
	}
	return handleFunc(v)
}

func init() {
	_client = msg_parser.NewMsgClient()
}

func handleBank(v types.SdkMsg) CustomMsgDocInfo {
	var (
		msgDoc CustomMsgDocInfo
		denoms []string
	)
	bankDocInfo, _ := _client.Bank.HandleTxMsg(v)
	msgDoc.MsgDocInfo = bankDocInfo
	switch bankDocInfo.DocTxMsg.Type {
	case MsgTypeSend:
		doc := bankDocInfo.DocTxMsg.Msg.(*bank.DocMsgSend)
		denoms = parseDenoms(doc.Amount)
	case MsgTypeMultiSend:
		doc := bankDocInfo.DocTxMsg.Msg.(*bank.DocMsgMultiSend)
		if len(doc.Inputs) > 0 {
			for _, v := range doc.Inputs {
				denoms = append(denoms, parseDenoms(v.Coins)...)
			}
		}
		if len(doc.Outputs) > 0 {
			for _, v := range doc.Outputs {
				denoms = append(denoms, parseDenoms(v.Coins)...)
			}
		}
	}
	msgDoc.Denoms = utils.RemoveDuplicatesFromSlice(denoms)
	return msgDoc
}
func handleCrisis(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := _client.Crisis.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}
func handleDistribution(v types.SdkMsg) CustomMsgDocInfo {
	var (
		msgDoc CustomMsgDocInfo
		denoms []string
	)
	distrubutionDocInfo, _ := _client.Distribution.HandleTxMsg(v)
	msgDoc.MsgDocInfo = distrubutionDocInfo
	switch distrubutionDocInfo.DocTxMsg.Type {
	case MsgTypeMsgFundCommunityPool:
		doc := distrubutionDocInfo.DocTxMsg.Msg.(*distribution.DocTxMsgFundCommunityPool)
		denoms = append(denoms, parseDenoms(doc.Amount)...)
	case MsgTypeWithdrawDelegatorReward:
	case MsgTypeMsgWithdrawValidatorCommission:
		break
	}
	msgDoc.Denoms = utils.RemoveDuplicatesFromSlice(denoms)
	return msgDoc
}
func handleSlashing(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := _client.Slashing.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}
func handleStaking(v types.SdkMsg) CustomMsgDocInfo {
	var (
		msgDoc CustomMsgDocInfo
		denoms []string
	)
	stakingDocInfo, _ := _client.Staking.HandleTxMsg(v)
	msgDoc.MsgDocInfo = stakingDocInfo
	switch stakingDocInfo.DocTxMsg.Type {
	case MsgTypeStakeDelegate:
		doc := stakingDocInfo.DocTxMsg.Msg.(*staking.DocTxMsgDelegate)
		denoms = append(denoms, parseDenoms(convertCoins([]Coin{doc.Amount}))...)
	case MsgTypeStakeBeginUnbonding:
		doc := stakingDocInfo.DocTxMsg.Msg.(*staking.DocTxMsgBeginUnbonding)
		denoms = append(denoms, parseDenoms([]types.Coin{doc.Amount})...)
	case MsgTypeBeginRedelegate:
		doc := stakingDocInfo.DocTxMsg.Msg.(*staking.DocTxMsgBeginRedelegate)
		denoms = append(denoms, parseDenoms([]types.Coin{doc.Amount})...)
	}
	msgDoc.Denoms = utils.RemoveDuplicatesFromSlice(denoms)
	return msgDoc
}
func handleEvidence(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := _client.Evidence.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}
func handleGov(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := _client.Gov.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}
func handleIbc(v types.SdkMsg) CustomMsgDocInfo {
	var (
		msgDoc CustomMsgDocInfo
		denoms []string
	)
	ibcDocInfo, _ := _client.Ibc.HandleTxMsg(v)
	msgDoc.MsgDocInfo = ibcDocInfo
	switch ibcDocInfo.DocTxMsg.Type {
	case MsgTypeIBCTransfer:
		doc := ibcDocInfo.DocTxMsg.Msg.(*ibc.DocMsgTransfer)
		denoms = append(denoms, doc.Token.Denom)
	case MsgTypeTimeout:
		doc := ibcDocInfo.DocTxMsg.Msg.(*ibc.DocMsgTimeout)
		denom := doc.Packet.Data.Denom
		if strings.Contains(denom, "/") {
			denom = ibc.GetIbcPacketDenom(doc.Packet, doc.Packet.Data.Denom)
		}
		denoms = append(denoms, denom)
	case MsgTypeRecvPacket:
		doc := ibcDocInfo.DocTxMsg.Msg.(*ibc.DocMsgRecvPacket)
		denom := doc.Packet.Data.Denom
		if strings.Contains(denom, "/") {
			denom = ibc.GetIbcPacketDenom(doc.Packet, doc.Packet.Data.Denom)
		}
		denoms = append(denoms, denom)
	}
	msgDoc.Denoms = utils.RemoveDuplicatesFromSlice(denoms)
	return msgDoc
}

func handleNft(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := _client.Nft.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleService(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := _client.Service.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleToken(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := _client.Token.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleOracle(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := _client.Oracle.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleRecord(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := _client.Record.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleRandom(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := _client.Random.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleHtlc(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := _client.Htlc.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleCoinswap(v types.SdkMsg) CustomMsgDocInfo {
	var (
		msgDoc CustomMsgDocInfo
		denoms []string
	)
	coinswapDocInfo, _ := _client.Coinswap.HandleTxMsg(v)
	msgDoc.MsgDocInfo = coinswapDocInfo
	switch coinswapDocInfo.DocTxMsg.Type {
	case MsgTypeSwapOrder:
		doc := coinswapDocInfo.DocTxMsg.Msg.(*coinswap.DocTxMsgSwapOrder)
		denoms = append(denoms, parseDenoms([]types.Coin{doc.Input.Coin})...)
		denoms = append(denoms, parseDenoms([]types.Coin{doc.Output.Coin})...)
	case MsgTypeAddLiquidity:
		doc := coinswapDocInfo.DocTxMsg.Msg.(*coinswap.DocTxMsgAddLiquidity)
		denoms = append(denoms, parseDenoms([]types.Coin{doc.MaxToken})...)
	case MsgTypeRemoveLiquidity:
		doc := coinswapDocInfo.DocTxMsg.Msg.(*coinswap.DocTxMsgRemoveLiquidity)
		denoms = append(denoms, parseDenoms([]types.Coin{doc.WithdrawLiquidity})...)
	}
	msgDoc.Denoms = utils.RemoveDuplicatesFromSlice(denoms)
	return msgDoc
}

func parseDenoms(coins []types.Coin) []string {
	if len(coins) == 0 {
		return nil
	}
	var denoms []string
	for _, v := range coins {
		denoms = append(denoms, v.Denom)
	}

	return denoms
}

// convert coins defined in modules to coins defined in types
func convertCoins(mCoins []Coin) types.Coins {
	var coins types.Coins
	if len(mCoins) == 0 {
		return coins
	}
	for _, v := range mCoins {
		coins = append(coins, types.Coin{
			Denom:  v.Denom,
			Amount: v.Amount,
		})
	}

	return coins
}
