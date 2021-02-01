package gov

import (
	"github.com/irisnet/rainbow-sync/lib/cdc"
	"github.com/irisnet/rainbow-sync/model"
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/utils"
)

type DocTxMsgSubmitProposal struct {
	Proposer       string       `bson:"proposer"`        //  Address of the proposer
	InitialDeposit []model.Coin `bson:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
	Content        interface{}  `bson:"content"`
}

func (doctx *DocTxMsgSubmitProposal) GetType() string {
	return MsgTypeSubmitProposal
}

func (doctx *DocTxMsgSubmitProposal) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgSubmitProposal)

	doctx.Content = CovertContent(msg.GetContent())
	doctx.Proposer = msg.Proposer
	doctx.InitialDeposit = model.BuildDocCoins(msg.InitialDeposit)
}

func CovertContent(content GovContent) interface{} {
	switch content.ProposalType() {
	case ProposalTypeCancelSoftwareUpgrade:
		var data ContentCancelSoftwareUpgradeProposal
		utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(content), &data)
		return data
	case ProposalTypeSoftwareUpgrade:
		var data ContentSoftwareUpgradeProposal
		utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(content), &data)
		return data
	case ProposalTypeCommunityPoolSpend:
		var data ContentCommunityPoolSpendProposal
		utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(content), &data)
		return data
	case ProposalTypeClientUpdate:
		var data ContentClientUpdateProposal
		utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(content), &data)
		return data
	case ProposalTypeText:
		var data ContentTextProposal
		utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(content), &data)
		return data
	case ProposalTypeParameterChange:
		var data ContentParameterChangeProposal
		utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(content), &data)
		return data
	}
	return content
}

func (m *DocTxMsgSubmitProposal) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgSubmitProposal
	)

	data, _ := cdc.GetMarshaler().MarshalJSON(v)
	cdc.GetMarshaler().UnmarshalJSON(data, &msg)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}

// MsgVote
type DocTxMsgVote struct {
	ProposalID uint64 `bson:"proposal_id"` // ID of the proposal
	Voter      string `bson:"voter"`       //  address of the voter
	Option     int32  `bson:"option"`      //  option from OptionSet chosen by the voter
}

func (doctx *DocTxMsgVote) GetType() string {
	return MsgTypeVote
}

func (doctx *DocTxMsgVote) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgVote)
	doctx.Voter = msg.Voter
	doctx.Option = int32(msg.Option)
	doctx.ProposalID = msg.ProposalId
}

func (m *DocTxMsgVote) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgVote
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.Voter)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}

// MsgDeposit
type DocTxMsgDeposit struct {
	ProposalID uint64       `bson:"proposal_id"` // ID of the proposal
	Depositor  string       `bson:"depositor"`   // Address of the depositor
	Amount     []model.Coin `bson:"amount"`      // Coins to add to the proposal's deposit
}

func (doctx *DocTxMsgDeposit) GetType() string {
	return MsgTypeDeposit
}

func (doctx *DocTxMsgDeposit) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgDeposit)
	doctx.Depositor = msg.Depositor
	doctx.Amount = model.BuildDocCoins(msg.Amount)
	doctx.ProposalID = msg.ProposalId
}

func (m *DocTxMsgDeposit) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgDeposit
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.Depositor)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
