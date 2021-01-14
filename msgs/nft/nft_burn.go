package nft

import (
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/utils"
	"strings"
)

type (
	DocMsgNFTBurn struct {
		Sender string `bson:"sender"`
		Id     string `bson:"id"`
		Denom  string `bson:"denom"`
	}
)

func (m *DocMsgNFTBurn) GetType() string {
	return MsgTypeNFTBurn
}

func (m *DocMsgNFTBurn) BuildMsg(v interface{}) {
	msg := v.(*MsgNFTBurn)

	m.Sender = msg.Sender
	m.Id = strings.ToLower(msg.Id)
	m.Denom = strings.ToLower(msg.DenomId)
}

func (m *DocMsgNFTBurn) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgNFTBurn
	)

	utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(v), &msg)
	addrs = append(addrs, msg.Sender)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
