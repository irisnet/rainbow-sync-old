package msg

import (
	"github.com/irisnet/rainbow-sync/service/cosmos/constant"
	imodel "github.com/irisnet/rainbow-sync/service/cosmos/model"
)

type (
	DocTxMsgTransfer struct {
		FromAddress string    `json:"from_address"`
		ToAddress   string    `json:"to_address" `
		Amount      []CoinStr `json:"amount"`
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
