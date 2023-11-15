package msgparser

import (
	"github.com/irisnet/rainbow-sync/lib/logger"
	"github.com/irisnet/rainbow-sync/utils"
	. "github.com/kaifei-bianjie/common-parser/modules"
	"github.com/kaifei-bianjie/common-parser/types"
	cosmosmod_parser "github.com/kaifei-bianjie/cosmosmod-parser"
	. "github.com/kaifei-bianjie/cosmosmod-parser/modules"
	"github.com/kaifei-bianjie/cosmosmod-parser/modules/bank"
	"github.com/kaifei-bianjie/cosmosmod-parser/modules/distribution"
	"github.com/kaifei-bianjie/cosmosmod-parser/modules/ibc"
	"github.com/kaifei-bianjie/cosmosmod-parser/modules/staking"
	irismod_parser "github.com/kaifei-bianjie/irismod-parser"
	. "github.com/kaifei-bianjie/irismod-parser/modules"
	"github.com/kaifei-bianjie/irismod-parser/modules/coinswap"
	"strings"
)

type MsgParser interface {
	HandleTxMsg(v types.SdkMsg) CustomMsgDocInfo
}

var (
	irisModClient   irismod_parser.MsgClient
	cosmosModClient cosmosmod_parser.MsgClient
)

func NewMsgParser(router Router) MsgParser {
	return &msgParser{
		rh: router,
	}
}

type msgParser struct {
	rh Router
}

// Handler returns the MsgServiceHandler for a given msg or nil if not found.
func (parser msgParser) getModule(v types.SdkMsg) string {
	var (
		route string
	)

	data := types.MsgTypeURL(v)
	if strings.HasPrefix(data, "/ibc.core.") {
		route = IbcRouteKey
	} else if strings.HasPrefix(data, "/ibc.applications.") {
		route = IbcTransferRouteKey
	} else if strings.HasPrefix(data, "/cosmos.bank.") {
		route = BankRouteKey
	} else if strings.HasPrefix(data, "/cosmos.crisis.") {
		route = CrisisRouteKey
	} else if strings.HasPrefix(data, "/cosmos.distribution.") {
		route = DistributionRouteKey
	} else if strings.HasPrefix(data, "/cosmos.slashing.") {
		route = SlashingRouteKey
	} else if strings.HasPrefix(data, "/cosmos.evidence.") {
		route = EvidenceRouteKey
	} else if strings.HasPrefix(data, "/cosmos.staking.") {
		route = StakingRouteKey
	} else if strings.HasPrefix(data, "/cosmos.gov.") {
		route = GovRouteKey
		//} else if strings.HasPrefix(data, "/tibc.core.") {
		//	route = TIbcRouteKey
		//} else if strings.HasPrefix(data, "/tibc.apps.") {
		//	route = TIbcTransferRouteKey
	} else if strings.HasPrefix(data, "/irismod.nft.") {
		route = NftRouteKey
		//} else if strings.HasPrefix(data, "/irismod.farm.") {
		//	route = FarmRouteKey
	} else if strings.HasPrefix(data, "/irismod.coinswap.") {
		route = CoinswapRouteKey
	} else if strings.HasPrefix(data, "/irismod.token.") {
		route = TokenRouteKey
	} else if strings.HasPrefix(data, "/irismod.record.") {
		route = RecordRouteKey
	} else if strings.HasPrefix(data, "/irismod.service.") {
		route = ServiceRouteKey
	} else if strings.HasPrefix(data, "/irismod.htlc.") {
		route = HtlcRouteKey
	} else if strings.HasPrefix(data, "/irismod.random.") {
		route = RandomRouteKey
	} else if strings.HasPrefix(data, "/irismod.oracle.") {
		route = OracleRouteKey
	} else {
		route = data
	}
	return route
}

func (parser *msgParser) HandleTxMsg(v types.SdkMsg) CustomMsgDocInfo {
	module := parser.getModule(v)
	handleFunc, err := parser.rh.GetRoute(module)
	if err != nil {
		logger.Error(err.Error(),
			logger.String("route", module),
			logger.String("type", module))
		return CustomMsgDocInfo{}
	}
	return handleFunc(v)
}

func init() {
	irisModClient = irismod_parser.NewMsgClient()
	cosmosModClient = cosmosmod_parser.NewMsgClient()
}

func handleBank(v types.SdkMsg) CustomMsgDocInfo {
	var (
		msgDoc CustomMsgDocInfo
		denoms []string
	)
	bankDocInfo, _ := cosmosModClient.Bank.HandleTxMsg(v)
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
	docInfo, _ := cosmosModClient.Crisis.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}
func handleDistribution(v types.SdkMsg) CustomMsgDocInfo {
	var (
		msgDoc CustomMsgDocInfo
		denoms []string
	)
	distrubutionDocInfo, _ := cosmosModClient.Distribution.HandleTxMsg(v)
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
	docInfo, _ := cosmosModClient.Slashing.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}
func handleStaking(v types.SdkMsg) CustomMsgDocInfo {
	var (
		msgDoc CustomMsgDocInfo
		denoms []string
	)
	stakingDocInfo, _ := cosmosModClient.Staking.HandleTxMsg(v)
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
	docInfo, _ := cosmosModClient.Evidence.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}
func handleGov(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := cosmosModClient.Gov.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}
func handleIbc(v types.SdkMsg) CustomMsgDocInfo {
	var (
		msgDoc CustomMsgDocInfo
		denoms []string
	)
	ibcDocInfo, _ := cosmosModClient.Ibc.HandleTxMsg(v)
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
	case MsgTypeAcknowledgement:
		break
	default:
		// clear msgDoc info for skip no use ibc tx msg
		msgDoc = CustomMsgDocInfo{}
		logger.Warn("skip no use ibc tx",
			logger.String("type", ibcDocInfo.DocTxMsg.Type))
	}
	msgDoc.Denoms = utils.RemoveDuplicatesFromSlice(denoms)
	return msgDoc
}

func handleNft(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := irisModClient.Nft.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleService(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := irisModClient.Service.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleToken(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := irisModClient.Token.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleOracle(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := irisModClient.Oracle.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleRecord(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := irisModClient.Record.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleRandom(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := irisModClient.Random.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleHtlc(v types.SdkMsg) CustomMsgDocInfo {
	docInfo, _ := irisModClient.Htlc.HandleTxMsg(v)
	var msgDoc CustomMsgDocInfo
	msgDoc.MsgDocInfo = docInfo
	return msgDoc
}

func handleCoinswap(v types.SdkMsg) CustomMsgDocInfo {
	var (
		msgDoc CustomMsgDocInfo
		denoms []string
	)
	coinswapDocInfo, _ := irisModClient.Coinswap.HandleTxMsg(v)
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
