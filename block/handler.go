package block

import (
	"github.com/cosmos/cosmos-sdk/types"
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/msgs/bank"
	"github.com/irisnet/rainbow-sync/msgs/coinswap"
	"github.com/irisnet/rainbow-sync/msgs/crisis"
	"github.com/irisnet/rainbow-sync/msgs/distribution"
	"github.com/irisnet/rainbow-sync/msgs/evidence"
	"github.com/irisnet/rainbow-sync/msgs/gov"
	"github.com/irisnet/rainbow-sync/msgs/htlc"
	"github.com/irisnet/rainbow-sync/msgs/ibc"
	"github.com/irisnet/rainbow-sync/msgs/nft"
	"github.com/irisnet/rainbow-sync/msgs/record"
	"github.com/irisnet/rainbow-sync/msgs/service"
	"github.com/irisnet/rainbow-sync/msgs/slashing"
	"github.com/irisnet/rainbow-sync/msgs/staking"
	"github.com/irisnet/rainbow-sync/msgs/token"
)

func HandleTxMsg(v types.Msg) MsgDocInfo {
	if BankDocInfo, ok := bank.HandleTxMsg(v); ok {
		return BankDocInfo
	}
	if IServiceDocInfo, ok := service.HandleTxMsg(v); ok {
		return IServiceDocInfo
	}
	if NftDocInfo, ok := nft.HandleTxMsg(v); ok {
		return NftDocInfo
	}
	if RecordDocInfo, ok := record.HandleTxMsg(v); ok {
		return RecordDocInfo
	}
	if TokenDocInfo, ok := token.HandleTxMsg(v); ok {
		return TokenDocInfo
	}
	if CoinswapDocInfo, ok := coinswap.HandleTxMsg(v); ok {
		return CoinswapDocInfo
	}
	if CrisisDocInfo, ok := crisis.HandleTxMsg(v); ok {
		return CrisisDocInfo
	}
	if DistrubutionDocInfo, ok := distribution.HandleTxMsg(v); ok {
		return DistrubutionDocInfo
	}
	if SlashingDocInfo, ok := slashing.HandleTxMsg(v); ok {
		return SlashingDocInfo
	}
	if EvidenceDocInfo, ok := evidence.HandleTxMsg(v); ok {
		return EvidenceDocInfo
	}
	if HtlcDocInfo, ok := htlc.HandleTxMsg(v); ok {
		return HtlcDocInfo
	}
	if StakingDocInfo, ok := staking.HandleTxMsg(v); ok {
		return StakingDocInfo
	}
	if GovDocInfo, ok := gov.HandleTxMsg(v); ok {
		return GovDocInfo
	}
	if IbcDocInfo, ok := ibc.HandleTxMsg(v); ok {
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
