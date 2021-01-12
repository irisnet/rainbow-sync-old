package nft

import (
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/utils"
	"strings"
)

type DocMsgIssueDenom struct {
	Id     string `bson:"id"`
	Name   string `bson:"name"`
	Sender string `bson:"sender"`
	Schema string `bson:"schema"`
}

func (m *DocMsgIssueDenom) GetType() string {
	return MsgTypeIssueDenom
}

func (m *DocMsgIssueDenom) BuildMsg(v interface{}) {
	msg := v.(*MsgIssueDenom)

	m.Sender = msg.Sender
	m.Schema = msg.Schema
	m.Id = strings.ToLower(msg.Id)
	m.Name = msg.Name
}

func (m *DocMsgIssueDenom) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgIssueDenom
	)

	utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(v), &msg)
	addrs = append(addrs, msg.Sender)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
