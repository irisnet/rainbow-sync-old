package token

import (
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/utils"
)

type DocMsgMintToken struct {
	Symbol string `bson:"symbol"`
	Amount uint64 `bson:"amount"`
	To     string `bson:"to"`
	Owner  string `bson:"owner"`
}

func (m *DocMsgMintToken) GetType() string {
	return MsgTypeMintToken
}

func (m *DocMsgMintToken) BuildMsg(v interface{}) {
	msg := v.(*MsgMintToken)

	m.Symbol = msg.Symbol
	m.Amount = msg.Amount
	m.To = msg.To
	m.Owner = msg.Owner
}

func (m *DocMsgMintToken) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgMintToken
	)

	utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(v), &msg)
	addrs = append(addrs, msg.Owner, msg.To)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
