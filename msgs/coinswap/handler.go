package coinswap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/irisnet/rainbow-sync/msgs"
)

func HandleTxMsg(msg sdk.Msg) (MsgDocInfo, bool) {
	ok := true
	switch msg.Type() {
	case new(MsgAddLiquidity).Type():
		docMsg := DocTxMsgAddLiquidity{}
		return docMsg.HandleTxMsg(msg), ok
	case new(MsgRemoveLiquidity).Type():
		docMsg := DocTxMsgRemoveLiquidity{}
		return docMsg.HandleTxMsg(msg), ok
	case new(MsgSwapOrder).Type():
		docMsg := DocTxMsgSwapOrder{}
		return docMsg.HandleTxMsg(msg), ok
	default:
		ok = false
	}
	return MsgDocInfo{}, ok
}
