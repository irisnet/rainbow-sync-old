package nft

import (
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/utils"
	"strings"
)

type (
	DocMsgNFTEdit struct {
		Sender string `bson:"sender"`
		Id     string `bson:"id"`
		Denom  string `bson:"denom"`
		URI    string `bson:"uri"`
		Data   string `bson:"data"`
		Name   string `bson:"name"`
	}
)

func (m *DocMsgNFTEdit) GetType() string {
	return MsgTypeNFTEdit
}

func (m *DocMsgNFTEdit) BuildMsg(v interface{}) {
	msg := v.(*MsgNFTEdit)

	m.Sender = msg.Sender
	m.Id = strings.ToLower(msg.Id)
	m.Denom = strings.ToLower(msg.DenomId)
	m.URI = msg.URI
	m.Data = msg.Data
	m.Name = msg.Name
}

func (m *DocMsgNFTEdit) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgNFTEdit
	)

	utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(v), &msg)
	addrs = append(addrs, msg.Sender)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
