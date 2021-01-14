package ibc

import (
	"github.com/irisnet/rainbow-sync/model"
	. "github.com/irisnet/rainbow-sync/msgs"
)

type DocMsgTransfer struct {
	// the port on which the packet will be sent
	SourcePort string `bson:"source_port"`
	// the channel by which the packet will be sent
	SourceChannel string `bson:"source_channel"`
	// the tokens to be transferred
	Token model.Coin `bson:"token"`
	// the sender address
	Sender string `bson:"sender"`
	// the recipient address on the destination chain
	Receiver string `bson:"receiver"`
	// Timeout height relative to the current block height.
	// The timeout is disabled when set to 0.
	TimeoutHeight Height `bson:"timeout_height"`
	// Timeout timestamp (in nanoseconds) relative to the current block timestamp.
	// The timeout is disabled when set to 0.
	TimeoutTimestamp uint64 `bson:"timeout_timestamp"`
}

func (m *DocMsgTransfer) GetType() string {
	return MsgTypeIbcTransfer
}

func (m *DocMsgTransfer) BuildMsg(v interface{}) {
	msg := v.(*MsgTransfer)
	m.Token = model.BuildDocCoin(msg.Token)
	m.SourcePort = msg.SourcePort
	m.SourceChannel = msg.SourceChannel
	m.Sender = msg.Sender
	m.Receiver = msg.Receiver
	m.TimeoutHeight = Height{RevisionNumber: msg.TimeoutHeight.RevisionNumber, RevisionHeight: msg.TimeoutHeight.RevisionHeight}
	m.TimeoutTimestamp = msg.TimeoutTimestamp

}

func (m *DocMsgTransfer) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
		msg   MsgTransfer
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.Sender, msg.Receiver)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
