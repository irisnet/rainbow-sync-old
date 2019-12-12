package msg

import (
	"github.com/irisnet/rainbow-sync/constant"
	imodel "github.com/irisnet/rainbow-sync/model"
)

type DocTxMsgRequestRand struct {
	Consumer      string `bson:"consumer"`       // request address
	BlockInterval uint64 `bson:"block-interval"` // block interval after which the requested random number will be generated
}

func (doctx *DocTxMsgRequestRand) Type() string {
	return constant.Iris_TxTypeRequestRand
}

func (doctx *DocTxMsgRequestRand) BuildMsg(txMsg interface{}) {
	msg := txMsg.(imodel.MsgRequestRand)
	doctx.Consumer = msg.Consumer.String()
	doctx.BlockInterval = msg.BlockInterval
}
