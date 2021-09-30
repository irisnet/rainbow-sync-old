package msgparser

import m "github.com/kaifei-bianjie/msg-parser/modules"

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
	FarmRouteKey         string = "farm"
	TIbcTransferRouteKey string = "NFT"
	TIbcRouteKey         string = "tibc"
)

var RouteHandlerMap = map[string]Handler{
	BankRouteKey:         handleBank,
	StakingRouteKey:      handleStaking,
	DistributionRouteKey: handleDistribution,
	CrisisRouteKey:       handleCrisis,
	EvidenceRouteKey:     handleEvidence,
	GovRouteKey:          handleGov,
	SlashingRouteKey:     handleSlashing,
	IbcRouteKey:          handleIbc,
	IbcTransferRouteKey:  handleIbc,
	NftRouteKey:          handleNft,
	ServiceRouteKey:      handleService,
	TokenRouteKey:        handleToken,
	HtlcRouteKey:         handleHtlc,
	CoinswapRouteKey:     handleCoinswap,
	RandomRouteKey:       handleRandom,
	OracleRouteKey:       handleOracle,
	RecordRouteKey:       handleRecord,
	FarmRouteKey:         handleFarm,
	TIbcTransferRouteKey: handleTIbc,
	TIbcRouteKey:         handleTIbc,
}

type CustomMsgDocInfo struct {
	m.MsgDocInfo
	Denoms []string
}
