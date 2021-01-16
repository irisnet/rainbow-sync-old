package token

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
	case new(MsgMintToken).Type():
		docMsg := DocMsgMintToken{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgBurnToken).Type():
		docMsg := DocMsgBurnToken{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgEditToken).Type():
		docMsg := DocMsgEditToken{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgIssueToken).Type():
		docMsg := DocMsgIssueToken{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgTransferTokenOwner).Type():
		docMsg := DocMsgTransferTokenOwner{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	default:
		ok = false
	}
	return msgDocInfo, ok
}
