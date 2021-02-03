package msgparser

import "github.com/kaifei-bianjie/msg-parser"

var (
	MsgClient msg_parser.MsgClient
)

func init() {
	MsgClient = msg_parser.NewMsgClient()
}
