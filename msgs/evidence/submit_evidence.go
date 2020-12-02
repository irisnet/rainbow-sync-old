package evidence

import (
	"encoding/json"
	. "github.com/irisnet/rainbow-sync/msgs"
)

// MsgSubmitEvidence defines an sdk.Msg type that supports submitting arbitrary
// Evidence.
type DocMsgSubmitEvidence struct {
	Submitter string `bson:"submitter"`
	Evidence  string `bson:"evidence"`
}

func (m *DocMsgSubmitEvidence) GetType() string {
	return MsgTypeSubmitEvidence
}

func (m *DocMsgSubmitEvidence) BuildMsg(v interface{}) {
	msg := v.(*MsgSubmitEvidence)
	m.Submitter = msg.Submitter
	evidence, _ := json.Marshal(msg.Evidence)
	m.Evidence = string(evidence)

}

func (m *DocMsgSubmitEvidence) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgSubmitEvidence
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.Submitter)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
