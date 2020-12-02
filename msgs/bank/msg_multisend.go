package bank

import (
	"github.com/irisnet/rainbow-sync/model"
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/utils"
)

type (
	DocMsgMultiSend struct {
		Inputs   []Item   `bson:"inputs"`
		Outputs  []Item   `bson:"outputs"`
		TempData []string `bson:"-"`
	}
	Item struct {
		Address string       `bson:"address"`
		Coins   []model.Coin `bson:"coins"`
	}
)

func (m *DocMsgMultiSend) GetType() string {
	return MsgTypeMultiSend
}

func (m *DocMsgMultiSend) BuildMsg(v interface{}) {
	msg := v.(*MsgMultiSend)
	for _, one := range msg.Inputs {
		m.Inputs = append(m.Inputs, Item{Address: one.Address, Coins: model.BuildDocCoins(one.Coins)})
		m.TempData = append(m.TempData, one.Address)
	}
	for _, one := range msg.Outputs {
		m.Outputs = append(m.Outputs, Item{Address: one.Address, Coins: model.BuildDocCoins(one.Coins)})
		m.TempData = append(m.TempData, one.Address)
	}

}

func (m *DocMsgMultiSend) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgMultiSend
	)

	utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(v), &msg)

	for _, one := range msg.Inputs {
		addrs = append(addrs, one.Address)
	}
	for _, one := range msg.Outputs {
		addrs = append(addrs, one.Address)
	}

	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
