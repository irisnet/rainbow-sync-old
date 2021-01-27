package block

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/rainbow-sync/lib/msgsdk"
	. "github.com/weichang-bianjie/msg-sdk/modules"
)

func HandleTxMsg(v types.Msg) MsgDocInfo {
	if BankDocInfo, ok := msgsdk.MsgClient.Bank.HandleTxMsg(v); ok {
		return BankDocInfo
	}
	if IServiceDocInfo, ok := msgsdk.MsgClient.Service.HandleTxMsg(v); ok {
		return IServiceDocInfo
	}
	if NftDocInfo, ok := msgsdk.MsgClient.Nft.HandleTxMsg(v); ok {
		return NftDocInfo
	}
	if RecordDocInfo, ok := msgsdk.MsgClient.Record.HandleTxMsg(v); ok {
		return RecordDocInfo
	}
	if TokenDocInfo, ok := msgsdk.MsgClient.Token.HandleTxMsg(v); ok {
		return TokenDocInfo
	}
	if CoinswapDocInfo, ok := msgsdk.MsgClient.Coinswap.HandleTxMsg(v); ok {
		return CoinswapDocInfo
	}
	if CrisisDocInfo, ok := msgsdk.MsgClient.Crisis.HandleTxMsg(v); ok {
		return CrisisDocInfo
	}
	if DistrubutionDocInfo, ok := msgsdk.MsgClient.Distribution.HandleTxMsg(v); ok {
		return DistrubutionDocInfo
	}
	if SlashingDocInfo, ok := msgsdk.MsgClient.Slashing.HandleTxMsg(v); ok {
		return SlashingDocInfo
	}
	if EvidenceDocInfo, ok := msgsdk.MsgClient.Evidence.HandleTxMsg(v); ok {
		return EvidenceDocInfo
	}
	if HtlcDocInfo, ok := msgsdk.MsgClient.Htlc.HandleTxMsg(v); ok {
		return HtlcDocInfo
	}
	if StakingDocInfo, ok := msgsdk.MsgClient.Staking.HandleTxMsg(v); ok {
		return StakingDocInfo
	}
	if GovDocInfo, ok := msgsdk.MsgClient.Gov.HandleTxMsg(v); ok {
		return GovDocInfo
	}
	if IbcDocInfo, ok := msgsdk.MsgClient.Ibc.HandleTxMsg(v); ok {
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
