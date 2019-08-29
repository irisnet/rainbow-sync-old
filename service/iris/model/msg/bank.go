package msg

import (
	"github.com/irisnet/rainbow-sync/service/iris/constant"
	imodel "github.com/irisnet/rainbow-sync/service/iris/model"
)

type DocTxMsgSetMemoRegexp struct {
	Owner      string `bson:"owner"`
	MemoRegexp string `bson:"memo_regexp"`
}

func (doctx *DocTxMsgSetMemoRegexp) Type() string {
	return constant.Iris_TxTypeSetMemoRegexp
}

func (doctx *DocTxMsgSetMemoRegexp) BuildMsg(txMsg interface{}) {
	msg := txMsg.(imodel.MsgSetMemoRegexp)
	doctx.MemoRegexp = msg.MemoRegexp
	doctx.Owner = msg.Owner.String()
}
