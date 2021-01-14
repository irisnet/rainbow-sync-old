package slashing

import (
	. "github.com/irisnet/rainbow-sync/msgs"
)

// MsgUnjail - struct for unjailing jailed validator
type DocTxMsgUnjail struct {
	ValidatorAddr string `bson:"address"` // address of the validator operator
}

func (doctx *DocTxMsgUnjail) GetType() string {
	return MsgTypeUnjail
}

func (doctx *DocTxMsgUnjail) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgUnjail)
	doctx.ValidatorAddr = msg.ValidatorAddr
}
func (m *DocTxMsgUnjail) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgUnjail
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.ValidatorAddr)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
