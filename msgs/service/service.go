package service

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
	case new(MsgDefineService).Type():
		docMsg := DocMsgDefineService{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgBindService).Type():
		docMsg := DocMsgBindService{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgCallService).Type():
		docMsg := DocMsgCallService{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgRespondService).Type():
		docMsg := DocMsgServiceResponse{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgUpdateServiceBinding).Type():
		docMsg := DocMsgUpdateServiceBinding{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgSetWithdrawAddress).Type():
		docMsg := DocMsgSetWithdrawAddress{}
		ConvertMsg(v, &docMsg)
		if docMsg.WithdrawAddress == "" {
			ok = false
			return msgDocInfo, ok
		}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgDisableServiceBinding).Type():
		docMsg := DocMsgDisableServiceBinding{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgEnableServiceBinding).Type():
		docMsg := DocMsgEnableServiceBinding{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgRefundServiceDeposit).Type():
		docMsg := DocMsgRefundServiceDeposit{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgUpdateRequestContext).Type():
		docMsg := DocMsgUpdateRequestContext{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgPauseRequestContext).Type():
		docMsg := DocMsgPauseRequestContext{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgStartRequestContext).Type():
		docMsg := DocMsgStartRequestContext{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgKillRequestContext).Type():
		docMsg := DocMsgKillRequestContext{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	case new(MsgWithdrawEarnedFees).Type():
		docMsg := DocMsgWithdrawEarnedFees{}
		msgDocInfo = docMsg.HandleTxMsg(v)
		break
	default:
		ok = false
	}
	return msgDocInfo, ok
}
