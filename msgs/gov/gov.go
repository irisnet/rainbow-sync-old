package gov

import (
	"github.com/irisnet/rainbow-sync/model"
	. "github.com/irisnet/rainbow-sync/msgs"
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

	doctx.Content = msg.GetContent()
	doctx.Proposer = msg.Proposer
	doctx.InitialDeposit = model.BuildDocCoins(msg.InitialDeposit)
}

func (m *DocTxMsgSubmitProposal) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		//msg   MsgSubmitProposal
	)

	//ConvertMsg(v, &msg)
	//addrs = append(addrs, msg.Proposer)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}

// MsgVote
type DocTxMsgVote struct {
	ProposalID uint64 `bson:"proposal_id"` // ID of the proposal
	Voter      string `bson:"voter"`       //  address of the voter
	Option     string `bson:"option"`      //  option from OptionSet chosen by the voter
}

func (doctx *DocTxMsgVote) GetType() string {
	return MsgTypeVote
}

func (doctx *DocTxMsgVote) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgVote)
	doctx.Voter = msg.Voter
	doctx.Option = msg.Option.String()
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
