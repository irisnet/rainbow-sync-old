package block

import (
	"github.com/irisnet/rainbow-sync/lib/msgparser"
	. "github.com/kaifei-bianjie/msg-parser/modules"
	"github.com/kaifei-bianjie/msg-parser/modules/bank"
	"github.com/kaifei-bianjie/msg-parser/modules/coinswap"
	"github.com/kaifei-bianjie/msg-parser/modules/distribution"
	"github.com/kaifei-bianjie/msg-parser/modules/ibc"
	"github.com/kaifei-bianjie/msg-parser/modules/staking"
	"github.com/kaifei-bianjie/msg-parser/types"
	"strings"
)

func HandleTxMsg(v types.SdkMsg) CustomMsgDocInfo {
	var (
		msgDoc CustomMsgDocInfo
		denoms []string
	)
	if bankDocInfo, ok := msgparser.MsgClient.Bank.HandleTxMsg(v); ok {
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
		msgDoc.Denoms = removeDuplicatesFromSlice(denoms)
		return msgDoc
	}
	if iServiceDocInfo, ok := msgparser.MsgClient.Service.HandleTxMsg(v); ok {
		msgDoc.MsgDocInfo = iServiceDocInfo
		return msgDoc
	}
	if nftDocInfo, ok := msgparser.MsgClient.Nft.HandleTxMsg(v); ok {
		msgDoc.MsgDocInfo = nftDocInfo
		return msgDoc
	}
	if recordDocInfo, ok := msgparser.MsgClient.Record.HandleTxMsg(v); ok {
		msgDoc.MsgDocInfo = recordDocInfo
		return msgDoc
	}
	if tokenDocInfo, ok := msgparser.MsgClient.Token.HandleTxMsg(v); ok {
		msgDoc.MsgDocInfo = tokenDocInfo
		return msgDoc
	}
	if coinswapDocInfo, ok := msgparser.MsgClient.Coinswap.HandleTxMsg(v); ok {
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
		msgDoc.Denoms = removeDuplicatesFromSlice(denoms)
		return msgDoc
	}
	if crisisDocInfo, ok := msgparser.MsgClient.Crisis.HandleTxMsg(v); ok {
		msgDoc.MsgDocInfo = crisisDocInfo
		return msgDoc
	}
	if distrubutionDocInfo, ok := msgparser.MsgClient.Distribution.HandleTxMsg(v); ok {
		msgDoc.MsgDocInfo = distrubutionDocInfo
		switch distrubutionDocInfo.DocTxMsg.Type {
		case MsgTypeMsgFundCommunityPool:
			doc := distrubutionDocInfo.DocTxMsg.Msg.(*distribution.DocTxMsgFundCommunityPool)
			denoms = append(denoms, parseDenoms(doc.Amount)...)
		case MsgTypeWithdrawDelegatorReward:
		case MsgTypeMsgWithdrawValidatorCommission:
			break
		}
		msgDoc.Denoms = removeDuplicatesFromSlice(denoms)
		return msgDoc
	}
	if slashingDocInfo, ok := msgparser.MsgClient.Slashing.HandleTxMsg(v); ok {
		msgDoc.MsgDocInfo = slashingDocInfo
		return msgDoc
	}
	if evidenceDocInfo, ok := msgparser.MsgClient.Evidence.HandleTxMsg(v); ok {
		msgDoc.MsgDocInfo = evidenceDocInfo
		return msgDoc
	}
	if htlcDocInfo, ok := msgparser.MsgClient.Htlc.HandleTxMsg(v); ok {
		msgDoc.MsgDocInfo = htlcDocInfo
		return msgDoc
	}
	if stakingDocInfo, ok := msgparser.MsgClient.Staking.HandleTxMsg(v); ok {
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
		msgDoc.Denoms = removeDuplicatesFromSlice(denoms)
		return msgDoc
	}
	if govDocInfo, ok := msgparser.MsgClient.Gov.HandleTxMsg(v); ok {
		msgDoc.MsgDocInfo = govDocInfo
		return msgDoc
	}
	if ibcDocInfo, ok := msgparser.MsgClient.Ibc.HandleTxMsg(v); ok {
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
		msgDoc.Denoms = removeDuplicatesFromSlice(denoms)
		return msgDoc
	}
	return msgDoc
}

func removeDuplicatesFromSlice(data []string) (result []string) {
	tempSet := make(map[string]string, len(data))
	for _, val := range data {
		if _, ok := tempSet[val]; ok || val == "" {
			continue
		}
		tempSet[val] = val
	}
	for one := range tempSet {
		result = append(result, one)
	}
	return
}
