package service

import (
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/utils"
	"strings"
)

type (
	DocMsgStartRequestContext struct {
		RequestContextId string `bson:"request_context_id" yaml:"request_context_id"`
		Consumer         string `bson:"consumer" yaml:"consumer"`
	}
)

func (m *DocMsgStartRequestContext) GetType() string {
	return MsgTypeStartRequestContext
}

func (m *DocMsgStartRequestContext) BuildMsg(v interface{}) {
	msg := v.(*MsgStartRequestContext)

	m.RequestContextId = strings.ToUpper(msg.RequestContextId)
	m.Consumer = msg.Consumer
}

func (m *DocMsgStartRequestContext) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgStartRequestContext
	)

	utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(v), &msg)
	addrs = append(addrs, msg.Consumer)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
