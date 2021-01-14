package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/irisnet/rainbow-sync/msgs"
)

func HandleTxMsg(msg sdk.Msg) (MsgDocInfo, bool) {
	ok := true
	switch msg.Type() {
	case new(MsgBeginRedelegate).Type():
		docMsg := DocTxMsgBeginRedelegate{}
		return docMsg.HandleTxMsg(msg), ok
	case new(MsgUndelegate).Type():
		docMsg := DocTxMsgBeginUnbonding{}
		return docMsg.HandleTxMsg(msg), ok
	case new(MsgDelegate).Type():
		docMsg := DocTxMsgDelegate{}
		return docMsg.HandleTxMsg(msg), ok
	case new(MsgEditValidator).Type():
		docMsg := DocMsgEditValidator{}
		return docMsg.HandleTxMsg(msg), ok
	case new(MsgCreateValidator).Type():
		docMsg := DocTxMsgCreateValidator{}
		return docMsg.HandleTxMsg(msg), ok
	default:
		ok = false
	}
	return MsgDocInfo{}, ok
}
