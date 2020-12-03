package service

import (
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/utils"
)

type (
	DocMsgWithdrawEarnedFees struct {
		Owner    string
		Provider string
	}
)

func (m *DocMsgWithdrawEarnedFees) GetType() string {
	return MsgTypeWithdrawEarnedFees
}

func (m *DocMsgWithdrawEarnedFees) BuildMsg(v interface{}) {
	msg := v.(*MsgWithdrawEarnedFees)

	m.Owner = msg.Owner.String()
	m.Provider = msg.Provider.String()
}

func (m *DocMsgWithdrawEarnedFees) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgWithdrawEarnedFees
	)

	utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(v), &msg)

	addrs = append(addrs, msg.Owner.String(), msg.Provider.String())
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
