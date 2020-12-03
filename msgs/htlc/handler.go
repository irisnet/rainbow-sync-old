package htlc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/irisnet/rainbow-sync/msgs"
)

func HandleTxMsg(msg sdk.Msg) (MsgDocInfo, bool) {
	ok := true
	switch msg.Type() {
	case new(MsgClaimHTLC).Type():
		docMsg := DocTxMsgClaimHTLC{}
		return docMsg.HandleTxMsg(msg), ok
	case new(MsgCreateHTLC).Type():
		docMsg := DocTxMsgCreateHTLC{}
		return docMsg.HandleTxMsg(msg), ok
	case new(MsgRefundHTLC).Type():
		docMsg := DocTxMsgRefundHTLC{}
		return docMsg.HandleTxMsg(msg), ok
	default:
		ok = false
	}
	return MsgDocInfo{}, ok
}
