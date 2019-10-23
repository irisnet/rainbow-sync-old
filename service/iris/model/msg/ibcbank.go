package msg

import (
	"github.com/irisnet/rainbow-sync/service/iris/constant"
	imodel "github.com/irisnet/rainbow-sync/service/iris/model"
)

type (
	DocTxMsgIBCBankTransfer struct {
		SrcPort      string `bson:"src_port"`
		SrcChannel   string `bson:"src_channel"`
		Denomination string `bson:"denomination"`
		Amount       string `bson:"amount"`
		Sender       string `bson:"sender"`
		Receiver     string `bson:"receiver"`
		Source       bool   `bson:"source"`
	}

	DocTxMsgIBCBankReceivePacket struct {
		Packet Packet `bson:"packet"`
		Height uint64 `bson:"height"`
		Signer string `bson:"signer"`
	}

	Packet struct {
		Msequence           uint64 `bson:"m_sequence"`
		Mtimeout            uint64 `bson:"m_timeout"`
		MsourcePort         string `bson:"m_source_port"`
		MsourceChannel      string `bson:"m_source_channel"`
		MdestinationPort    string `bson:"m_destination_port"`
		MdestinationChannel string `bson:"m_destination_channel"`
		Mdata               string `bson:"m_data"`
	}
)

func (m *DocTxMsgIBCBankTransfer) Type() string {
	return constant.TxMsgTypeIBCBankTransfer
}

func (m *DocTxMsgIBCBankTransfer) BuildMsg(txMsg interface{}) {
	msg := txMsg.(imodel.IBCBankMsgTransfer)

	m.SrcPort = msg.SrcPort
	m.SrcChannel = msg.SrcChannel
	m.Denomination = msg.Denomination
	m.Amount = msg.Amount.String()
	m.Sender = msg.Sender
	m.Receiver = msg.Receiver
	m.Source = msg.Source
}

func (m *DocTxMsgIBCBankReceivePacket) Type() string {
	return constant.TxMsgTypeIBCBankRecvTransferPacket
}

func (m *DocTxMsgIBCBankReceivePacket) BuildMsg(txMsg interface{}) {
	msg := txMsg.(imodel.IBCBankMsgReceivePacket)

	packet := Packet{
		Msequence:           msg.Packet.Sequence(),
		Mtimeout:            msg.Packet.TimeoutHeight(),
		MsourcePort:         msg.Packet.SourcePort(),
		MsourceChannel:      msg.Packet.SourceChannel(),
		MdestinationPort:    msg.Packet.DestPort(),
		MdestinationChannel: msg.Packet.DestChannel(),
		Mdata:               string(msg.Packet.Data()),
	}

	m.Packet = packet
	m.Signer = msg.Signer.String()
	m.Height = msg.Height
}
