package slashing

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/irisnet/rainbow-sync/msgs"
)

func HandleTxMsg(msg sdk.Msg) (MsgDocInfo, bool) {
	ok := true
	switch msg.Type() {
	case new(MsgUnjail).Type():
		docMsg := DocTxMsgUnjail{}
		return docMsg.HandleTxMsg(msg), ok

	default:
		ok = false
	}
	return MsgDocInfo{}, ok
}
