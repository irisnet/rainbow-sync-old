package msg

import (
	"github.com/irisnet/rainbow-sync/constant"
	"github.com/irisnet/rainbow-sync/model"
	"github.com/irisnet/rainbow-sync/utils"
)

type DocTxMsgIBCBankTransfer struct {
	SourcePort    string      `bson:"source_port"`    // the port on which the packet will be sent
	SourceChannel string      `bson:"source_channel"` // the channel by which the packet will be sent
	DestHeight    uint64      `bson:"dest_height"`    // the current height of the destination chain
	Amount        model.Coins `bson:"amount"`         // the tokens to be transferred
	Sender        string      `bson:"sender"`         // the sender address
	Receiver      string      `bson:"receiver"`       // the recipient address on the destination chain
	//Source        bool        `bson:"source"`         // indicates if the sending chain is the source chain of the tokens to be transferred
}

func (m *DocTxMsgIBCBankTransfer) Type() string {
	return constant.TxMsgTypeIBCBankTransfer
}

func (m *DocTxMsgIBCBankTransfer) BuildMsg(txMsg interface{}) {
	msg := txMsg.(model.IBCBankMsgTransfer)

	m.SourcePort = msg.SourcePort
	m.SourceChannel = msg.SourceChannel
	m.DestHeight = msg.DestHeight
	m.Amount = utils.ParseCoins(msg.Amount)
	m.Sender = msg.Sender.String()
	m.Receiver = msg.Receiver
	//m.Source = msg.Source
}
