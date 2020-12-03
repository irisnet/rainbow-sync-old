package nft

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
	case new(MsgNFTMint).Type():
		docMsg := DocMsgNFTMint{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgNFTEdit).Type():
		docMsg := DocMsgNFTEdit{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgNFTTransfer).Type():
		docMsg := DocMsgNFTTransfer{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgNFTBurn).Type():
		docMsg := DocMsgNFTBurn{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgIssueDenom).Type():
		docMsg := DocMsgIssueDenom{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	default:
		ok = false
	}
	return msgDocInfo, ok
}
