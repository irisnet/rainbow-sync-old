package iris

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/modules/slashing"
	dtypes "github.com/irisnet/irishub/modules/distribution/types"
	dtags "github.com/irisnet/irishub/modules/distribution/tags"
	"github.com/irisnet/irishub/types"
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
	//Msg    Msg       `bson:"msg"`
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
	StakeValidator = stake.Validator
	Delegation = stake.Delegation
	UnbondingDelegation = stake.UnbondingDelegation

	MsgDeposit = gov.MsgDeposit
	MsgSubmitProposal = gov.MsgSubmitProposal
	MsgSubmitSoftwareUpgradeProposal = gov.MsgSubmitSoftwareUpgradeProposal
	MsgSubmitTaxUsageProposal = gov.MsgSubmitTxTaxUsageProposal
	MsgVote = gov.MsgVote
	Proposal = gov.Proposal
	SdkVote = gov.Vote

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
