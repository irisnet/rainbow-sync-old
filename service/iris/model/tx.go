package iris

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/irisnet/irishub/app/v1/stake"
	"github.com/irisnet/irishub/app/v1/distribution"
	"github.com/irisnet/irishub/app/v1/gov"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/slashing"
	dtypes "github.com/irisnet/irishub/app/v1/distribution/types"
	dtags "github.com/irisnet/irishub/app/v1/distribution/tags"
	"github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/app/v1/rand"
)

type IrisTx struct {
	Time      time.Time         `json:"time" bson:"time"`
	Height    int64             `json:"height" bson:"height"`
	TxHash    string            `json:"tx_hash" bson:"tx_hash"`
	From      string            `json:"from" bson:"from"`
	To        string            `json:"to" bson:"to"`
	Initiator string            `json:"initiator" bson:"initiator"`
	Amount    []*Coin           `json:"amount" bson:"amount"`
	Type      string            `json:"type" bson:"type"`
	Fee       *Fee              `json:"fee" bson:"fee"`
	ActualFee *ActualFee        `json:"actual_fee" bson:"actual_fee"`
	Memo      string            `json:"memo" bson:"memo"`
	Status    string            `json:"status" bson:"status"`
	Code      uint32            `json:"code" bson:"code"`
	Tags      map[string]string `json:"tags" bson:"tags"`
	Msgs      []DocTxMsg        `bson:"msgs"`
}

type DocTxMsg struct {
	Type string `bson:"type"`
	Msg  Msg    `bson:"msg"`
}

type Msg interface {
	Type() string
	BuildMsg(msg interface{})
}

const (
	CollectionNameIrisTx = "sync_iris_tx"
)

func (d IrisTx) Name() string {
	return CollectionNameIrisTx
}

func (d IrisTx) PkKvPair() map[string]interface{} {
	return bson.M{}
}

type (
	MsgTransfer = bank.MsgSend
	MsgBurn = bank.MsgBurn
	MsgSetMemoRegexp = bank.MsgSetMemoRegexp

	MsgStakeCreate = stake.MsgCreateValidator
	MsgStakeEdit = stake.MsgEditValidator
	MsgStakeDelegate = stake.MsgDelegate
	MsgStakeBeginUnbonding = stake.MsgBeginUnbonding
	MsgBeginRedelegate = stake.MsgBeginRedelegate
	MsgUnjail = slashing.MsgUnjail
	MsgSetWithdrawAddress = dtypes.MsgSetWithdrawAddress
	MsgWithdrawDelegatorReward = distribution.MsgWithdrawDelegatorReward
	MsgWithdrawDelegatorRewardsAll = distribution.MsgWithdrawDelegatorRewardsAll
	MsgWithdrawValidatorRewardsAll = distribution.MsgWithdrawValidatorRewardsAll

	MsgDeposit = gov.MsgDeposit
	MsgSubmitProposal = gov.MsgSubmitProposal
	MsgSubmitSoftwareUpgradeProposal = gov.MsgSubmitSoftwareUpgradeProposal
	MsgSubmitTaxUsageProposal = gov.MsgSubmitCommunityTaxUsageProposal
	MsgSubmitTokenAdditionProposal = gov.MsgSubmitTokenAdditionProposal
	MsgVote = gov.MsgVote

	MsgRequestRand = rand.MsgRequestRand

	AssetIssueToken = asset.MsgIssueToken
	AssetEditToken = asset.MsgEditToken
	AssetMintToken = asset.MsgMintToken
	AssetTransferTokenOwner = asset.MsgTransferTokenOwner
	AssetCreateGateway = asset.MsgCreateGateway
	AssetEditGateWay = asset.MsgEditGateway
	AssetTransferGatewayOwner = asset.MsgTransferGatewayOwner

	SdkCoins = types.Coins
	KVPair = types.KVPair
)

var (
	TagDistributionReward = dtags.Reward
)

type Coin struct {
	Denom  string  `bson:"denom" json:"denom"`
	Amount float64 `bson:"amount" json:"amount"`
}

type Coins []*Coin

type Fee struct {
	Amount Coins `bson:"amount" json:"amount"`
	Gas    int64 `bson:"gas" json:"gas"`
}

type ActualFee struct {
	Denom  string  `json:"denom"`
	Amount float64 `json:"amount"`
}
