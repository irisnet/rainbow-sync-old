package bank

import (
	"github.com/irisnet/rainbow-sync/model"
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/utils"
)

type (
	DocMsgSend struct {
		FromAddress string       `bson:"from_address"`
		ToAddress   string       `bson:"to_address"`
		Amount      []model.Coin `bson:"amount"`
	}
)

func (m *DocMsgSend) GetType() string {
	return MsgTypeSend
}

func (m *DocMsgSend) BuildMsg(v interface{}) {
	msg := v.(*MsgSend)
	m.FromAddress = msg.FromAddress
	m.ToAddress = msg.ToAddress
	m.Amount = model.BuildDocCoins(msg.Amount)
}

func (m *DocMsgSend) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgSend
	)

	utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(v), &msg)

	addrs = append(addrs, msg.FromAddress, msg.ToAddress)

	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
