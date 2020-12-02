package record

import (
	"github.com/cosmos/cosmos-sdk/types"
	. "github.com/irisnet/rainbow-sync/msgs"
)

func HandleTxMsg(v types.Msg) (MsgDocInfo, bool) {
	var (
		msgDocInfo MsgDocInfo
	)
	ok := true
	switch v.Type() {
	case new(MsgRecordCreate).Type():
		docMsg := DocMsgRecordCreate{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	default:
		ok = false
	}
	return msgDocInfo, ok
}
