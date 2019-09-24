package msg

import (
	"github.com/irisnet/rainbow-sync/service/iris/constant"
	imodel "github.com/irisnet/rainbow-sync/service/iris/model"
	"github.com/irisnet/irishub/app/v1/gov"
	"strconv"
)

type Param struct {
	Subspace string `bson:"subspace"`
	Key      string `bson:"key"`
	Value    string `bson:"value"`
}

type Params []Param

type DocTxMsgSubmitProposal struct {
	Title          string       `bson:"title"`          //  Title of the proposal
	Description    string       `bson:"description"`    //  Description of the proposal
	Proposer       string       `bson:"proposer"`       //  Address of the proposer
	InitialDeposit imodel.Coins `bson:"initialDeposit"` //  Initial deposit paid by sender. Must be strictly positive.
	ProposalType   string       `bson:"proposalType"`   //  Initial deposit paid by sender. Must be strictly positive.
	Params         Params       `bson:"params"`
}

func (doctx *DocTxMsgSubmitProposal) Type() string {
	return constant.Iris_TxTypeSubmitProposal
}

func (doctx *DocTxMsgSubmitProposal) BuildMsg(txMsg interface{}) {
	msg := txMsg.(imodel.MsgSubmitProposal)
	doctx.Title = msg.Title
	doctx.Description = msg.Description
	doctx.ProposalType = msg.ProposalType.String()
	doctx.Proposer = msg.Proposer.String()
	doctx.Params = loadParams(msg.Params)
	doctx.InitialDeposit = loadInitialDeposit(msg.InitialDeposit)
}

type DocTxMsgSubmitSoftwareUpgradeProposal struct {
	DocTxMsgSubmitProposal
	Version      uint64 `bson:"version"`
	Software     string `bson:"software"`
	SwitchHeight uint64 `bson:"switch_height"`
	Threshold    string `bson:"threshold"`
}

func (doctx *DocTxMsgSubmitSoftwareUpgradeProposal) Type() string {
	return constant.Iris_TxTypeSubmitProposal
}

func (doctx *DocTxMsgSubmitSoftwareUpgradeProposal) BuildMsg(txMsg interface{}) {
	msg := txMsg.(imodel.MsgSubmitSoftwareUpgradeProposal)
	doctx.Title = msg.Title
	doctx.Description = msg.Description
	doctx.ProposalType = msg.ProposalType.String()
	doctx.Proposer = msg.Proposer.String()
	doctx.Params = loadParams(msg.Params)
	doctx.InitialDeposit = loadInitialDeposit(msg.InitialDeposit)
	doctx.Version = msg.Version
	doctx.Software = msg.Software
	doctx.SwitchHeight = msg.SwitchHeight
	doctx.Threshold = msg.Threshold.String()
}

type DocTxMsgSubmitCommunityTaxUsageProposal struct {
	DocTxMsgSubmitProposal
	Usage       string `bson:"usage"`
	DestAddress string `bson:"dest_address"`
	Percent     string `bson:"percent"`
}

func (doctx *DocTxMsgSubmitCommunityTaxUsageProposal) Type() string {
	return constant.Iris_TxTypeSubmitProposal
}

func (doctx *DocTxMsgSubmitCommunityTaxUsageProposal) BuildMsg(txMsg interface{}) {
	msg := txMsg.(imodel.MsgSubmitTaxUsageProposal)
	doctx.Title = msg.Title
	doctx.Description = msg.Description
	doctx.ProposalType = msg.ProposalType.String()
	doctx.Proposer = msg.Proposer.String()
	doctx.Params = loadParams(msg.Params)
	doctx.InitialDeposit = loadInitialDeposit(msg.InitialDeposit)
	doctx.Usage = msg.Usage.String()
	doctx.DestAddress = msg.DestAddress.String()
	doctx.Percent = msg.Percent.String()
}

type DocTxMsgSubmitTokenAdditionProposal struct {
	DocTxMsgSubmitProposal
	Symbol          string `bson:"symbol"`
	CanonicalSymbol string `bson:"canonical_symbol"`
	Name            string `bson:"name"`
	Decimal         uint8  `bson:"decimal"`
	MinUnitAlias    string `bson:"min_unit_alias"`
	InitialSupply   uint64 `bson:"initial_supply"`
}

func (doctx *DocTxMsgSubmitTokenAdditionProposal) Type() string {
	return constant.Iris_TxTypeSubmitTokenAdditionProposal
}

func (doctx *DocTxMsgSubmitTokenAdditionProposal) BuildMsg(txMsg interface{}) {
	msg := txMsg.(imodel.MsgSubmitTokenAdditionProposal)
	doctx.Title = msg.Title
	doctx.Description = msg.Description
	doctx.ProposalType = msg.ProposalType.String()
	doctx.Proposer = msg.Proposer.String()
	doctx.Params = loadParams(msg.Params)
	doctx.InitialDeposit = loadInitialDeposit(msg.InitialDeposit)
	doctx.Symbol = msg.Symbol
	doctx.MinUnitAlias = msg.MinUnitAlias
	doctx.CanonicalSymbol = msg.CanonicalSymbol
	doctx.Name = msg.Name
	doctx.Decimal = msg.Decimal
	doctx.InitialSupply = msg.InitialSupply
}

func loadParams(params []gov.Param) (result []Param) {
	for _, val := range params {
		result = append(result, Param{Subspace: val.Subspace, Value: val.Value, Key: val.Key})
	}
	return
}

func loadInitialDeposit(coins imodel.SdkCoins) (result imodel.Coins) {
	for _, val := range coins {
		amt, _ := strconv.ParseFloat(val.Amount.String(), 64)
		result = append(result, &imodel.Coin{Amount: amt, Denom: val.Denom})
	}
	return
}

// MsgVote
type DocTxMsgVote struct {
	ProposalID uint64 `bson:"proposal_id"` // ID of the proposal
	Voter      string `bson:"voter"`       //  address of the voter
	Option     string `bson:"option"`      //  option from OptionSet chosen by the voter
}

func (doctx *DocTxMsgVote) Type() string {
	return constant.Iris_TxTypeVote
}

func (doctx *DocTxMsgVote) BuildMsg(txMsg interface{}) {
	msg := txMsg.(imodel.MsgVote)
	doctx.Voter = msg.Voter.String()
	doctx.Option = msg.Option.String()
	doctx.ProposalID = msg.ProposalID
}
