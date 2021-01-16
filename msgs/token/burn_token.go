package token

import (
	. "github.com/irisnet/rainbow-sync/msgs"
)

// MsgBurnToken defines an SDK message for burning some tokens.
type DocMsgBurnToken struct {
	Symbol string `bson:"symbol"`
	Amount uint64 `bson:"amount"`
	Sender string `bson:"sender"`
}

func (m *DocMsgBurnToken) GetType() string {
	return MsgTypeBurnToken
}

func (m *DocMsgBurnToken) BuildMsg(v interface{}) {
	msg := v.(*MsgBurnToken)

	m.Symbol = msg.Symbol
	m.Amount = msg.Amount
	m.Sender = msg.Sender
}

func (m *DocMsgBurnToken) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgBurnToken
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.Sender)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
