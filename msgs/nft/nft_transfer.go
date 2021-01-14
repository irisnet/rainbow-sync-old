package nft

import (
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/utils"
	"strings"
)

type (
	DocMsgNFTTransfer struct {
		Sender    string `bson:"sender"`
		Recipient string `bson:"recipient"`
		URI       string `bson:"uri"`
		Name      string `bson:"name"`
		Denom     string `bson:"denom"`
		Id        string `bson:"id"`
		Data      string `bson:"data"`
	}
)

func (m *DocMsgNFTTransfer) GetType() string {
	return MsgTypeNFTTransfer
}

func (m *DocMsgNFTTransfer) BuildMsg(v interface{}) {
	msg := v.(*MsgNFTTransfer)

	m.Sender = msg.Sender
	m.Recipient = msg.Recipient
	m.Id = strings.ToLower(msg.Id)
	m.Denom = strings.ToLower(msg.DenomId)
	m.URI = msg.URI
	m.Data = msg.Data
	m.Name = msg.Name
}

func (m *DocMsgNFTTransfer) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgNFTTransfer
	)

	utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(v), &msg)
	addrs = append(addrs, msg.Sender, msg.Recipient)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
