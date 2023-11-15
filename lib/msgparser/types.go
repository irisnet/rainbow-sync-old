package msgparser

import m "github.com/kaifei-bianjie/common-parser/modules"

const (
	BankRouteKey         string = "bank"
	StakingRouteKey      string = "staking"
	DistributionRouteKey string = "distribution"
	CrisisRouteKey       string = "crisis"
	EvidenceRouteKey     string = "evidence"
	GovRouteKey          string = "gov"
	SlashingRouteKey     string = "slashing"
	IbcRouteKey          string = "ibc"
	IbcTransferRouteKey  string = "transfer"
	NftRouteKey          string = "nft"
	ServiceRouteKey      string = "service"
	TokenRouteKey        string = "token"
	HtlcRouteKey         string = "htlc"
	CoinswapRouteKey     string = "coinswap"
	RandomRouteKey       string = "random"
	OracleRouteKey       string = "oracle"
	RecordRouteKey       string = "record"
)

type CustomMsgDocInfo struct {
	m.MsgDocInfo
	Denoms []string
}
