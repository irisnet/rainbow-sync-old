package msgsdk

import "github.com/weichang-bianjie/msg-sdk"

var (
	MsgClient msg_sdk.MsgClient
)

func init() {
	MsgClient = msg_sdk.NewMsgClient()
}
