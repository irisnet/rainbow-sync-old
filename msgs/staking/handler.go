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
	case new(MsgStakeBeginUnbonding).Type():
		docMsg := DocTxMsgBeginUnbonding{}
		return docMsg.HandleTxMsg(msg), ok
	case new(MsgStakeDelegate).Type():
		docMsg := DocTxMsgDelegate{}
		return docMsg.HandleTxMsg(msg), ok
	case new(MsgStakeEdit).Type():
		docMsg := DocMsgEditValidator{}
		return docMsg.HandleTxMsg(msg), ok
	case new(MsgStakeCreate).Type():
		docMsg := DocTxMsgCreateValidator{}
		return docMsg.HandleTxMsg(msg), ok
	default:
		ok = false
	}
	return MsgDocInfo{}, ok
}
