package block

import (
	"github.com/irisnet/rainbow-sync/lib/msgparser"
	. "github.com/kaifei-bianjie/msg-parser/modules"
	"github.com/kaifei-bianjie/msg-parser/types"
)

func HandleTxMsg(v types.SdkMsg) MsgDocInfo {
	if BankDocInfo, ok := msgparser.MsgClient.Bank.HandleTxMsg(v); ok {
		return BankDocInfo
	}
	if IServiceDocInfo, ok := msgparser.MsgClient.Service.HandleTxMsg(v); ok {
		return IServiceDocInfo
	}
	if NftDocInfo, ok := msgparser.MsgClient.Nft.HandleTxMsg(v); ok {
		return NftDocInfo
	}
	if RecordDocInfo, ok := msgparser.MsgClient.Record.HandleTxMsg(v); ok {
		return RecordDocInfo
	}
	if TokenDocInfo, ok := msgparser.MsgClient.Token.HandleTxMsg(v); ok {
		return TokenDocInfo
	}
	if CoinswapDocInfo, ok := msgparser.MsgClient.Coinswap.HandleTxMsg(v); ok {
		return CoinswapDocInfo
	}
	if CrisisDocInfo, ok := msgparser.MsgClient.Crisis.HandleTxMsg(v); ok {
		return CrisisDocInfo
	}
	if DistrubutionDocInfo, ok := msgparser.MsgClient.Distribution.HandleTxMsg(v); ok {
		return DistrubutionDocInfo
	}
	if SlashingDocInfo, ok := msgparser.MsgClient.Slashing.HandleTxMsg(v); ok {
		return SlashingDocInfo
	}
	if EvidenceDocInfo, ok := msgparser.MsgClient.Evidence.HandleTxMsg(v); ok {
		return EvidenceDocInfo
	}
	if HtlcDocInfo, ok := msgparser.MsgClient.Htlc.HandleTxMsg(v); ok {
		return HtlcDocInfo
	}
	if StakingDocInfo, ok := msgparser.MsgClient.Staking.HandleTxMsg(v); ok {
		return StakingDocInfo
	}
	if GovDocInfo, ok := msgparser.MsgClient.Gov.HandleTxMsg(v); ok {
		return GovDocInfo
	}
	if IbcDocInfo, ok := msgparser.MsgClient.Ibc.HandleTxMsg(v); ok {
		return IbcDocInfo
	}
	return MsgDocInfo{}
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
