package msg

import (
	"github.com/irisnet/rainbow-sync/service/cosmos/constant"
	imodel "github.com/irisnet/rainbow-sync/service/cosmos/model"
)

type (
	DocTxMsgTransfer struct {
		FromAddress string    `bson:"from_address"`
		ToAddress   string    `bson:"to_address" `
		Amount      []CoinStr `bson:"amount"`
	}
)

func (m *DocTxMsgTransfer) Type() string {
	return constant.TxMsgTypeTransfer
}

func (m *DocTxMsgTransfer) BuildMsg(txMsg interface{}) {
	msg := txMsg.(imodel.MsgTransfer)

	var amount []CoinStr

	for _, v := range msg.Amount {
		coinStr := CoinStr{
			Denom:  v.Denom,
			Amount: v.Amount.String(),
		}
		amount = append(amount, coinStr)
	}

	m.FromAddress = msg.FromAddress.String()
	m.ToAddress = msg.ToAddress.String()
	m.Amount = amount
}
