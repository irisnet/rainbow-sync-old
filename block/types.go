package block

import m "github.com/kaifei-bianjie/msg-parser/modules"

type CustomMsgDocInfo struct {
	m.MsgDocInfo
	Denoms []string
}
