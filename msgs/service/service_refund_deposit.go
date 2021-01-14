package service

import (
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/utils"
)

type (
	DocMsgRefundServiceDeposit struct {
		ServiceName string `bson:"service_name" yaml:"service_name"`
		Provider    string `bson:"provider" yaml:"provider"`
		Owner       string `bson:"owner" yaml:"owner"`
	}
)

func (m *DocMsgRefundServiceDeposit) GetType() string {
	return MsgTypeRefundServiceDeposit
}

func (m *DocMsgRefundServiceDeposit) BuildMsg(v interface{}) {
	msg := v.(*MsgRefundServiceDeposit)

	m.ServiceName = msg.ServiceName
	m.Provider = msg.Provider
	m.Owner = msg.Owner
}

func (m *DocMsgRefundServiceDeposit) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgRefundServiceDeposit
	)

	utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(v), &msg)
	addrs = append(addrs, msg.Owner, msg.Provider)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
