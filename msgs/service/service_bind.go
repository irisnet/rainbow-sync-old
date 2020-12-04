package service

import (
	. "github.com/irisnet/rainbow-sync/msgs"
	"github.com/irisnet/rainbow-sync/utils"
)

type (
	DocMsgBindService struct {
		ServiceName string `bson:"service_name"`
		Provider    string `bson:"provider"`
		Deposit     Coins  `bson:"deposit"`
		Pricing     string `bson:"pricing"`
		QoS         uint64 `bson:"qos"`
		Owner       string `bson:"owner"`
	}
)

func (m *DocMsgBindService) GetType() string {
	return MsgTypeBindService
}

func (m *DocMsgBindService) BuildMsg(v interface{}) {
	msg := v.(*MsgBindService)

	var coins Coins
	for _, one := range msg.Deposit {
		coins = append(coins, &Coin{Denom: one.Denom, Amount: one.Amount.String()})
	}
	m.ServiceName = msg.ServiceName
	m.Provider = msg.Provider
	m.Deposit = coins
	m.Pricing = msg.Pricing
	m.QoS = msg.QoS
	m.Owner = msg.Owner
}

func (m *DocMsgBindService) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgBindService
	)

	utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(v), &msg)

	addrs = append(addrs, msg.Owner, msg.Provider)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
