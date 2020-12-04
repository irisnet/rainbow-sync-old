package token

import (
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/utils"
)

type DocMsgTransferTokenOwner struct {
	SrcOwner string `bson:"src_owner"`
	DstOwner string `bson:"dst_owner"`
	Symbol   string `bson:"symbol"`
}

func (m *DocMsgTransferTokenOwner) GetType() string {
	return MsgTypeTransferTokenOwner
}

func (m *DocMsgTransferTokenOwner) BuildMsg(v interface{}) {
	msg := v.(*MsgTransferTokenOwner)

	m.Symbol = msg.Symbol
	m.SrcOwner = msg.SrcOwner
	m.DstOwner = msg.DstOwner
}

func (m *DocMsgTransferTokenOwner) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgTransferTokenOwner
	)

	utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(v), &msg)
	addrs = append(addrs, msg.SrcOwner, msg.DstOwner)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
