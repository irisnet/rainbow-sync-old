package msg

import (
	"github.com/irisnet/rainbow-sync/model"
	"encoding/json"
	"github.com/irisnet/rainbow-sync/constant"
)

// Packet defines a type that carries data across different chains through IBC
type Packet struct {
	Sequence           uint64     `bson:"sequence" json:"sequence"`            // number corresponds to the order of sends and receives, where a Packet with an earlier sequence number must be sent and received before a Packet with a later sequence number.
	SourcePort         string     `bson:"source_port" json:"source_port"`         // identifies the port on the sending chain.
	SourceChannel      string     `bson:"source_channel" json:"source_channel"`      // identifies the channel end on the sending chain.
	DestinationPort    string     `bson:"destination_port" json:"destination_port"`    // identifies the port on the receiving chain.
	DestinationChannel string     `bson:"destination_channel" json:"destination_channel"` // identifies the channel end on the receiving chain.
	TimeoutHeight      uint64     `bson:"timeout_height" json:"timeout_height"`      // block height after which the packet times out
	Data               SendPacket `bson:"data" json:"data"`                 // opaque value which can be defined by the application logic of the associated modules.
}

type SendPacket struct {
	Type  string `bson:"type" json:"type"`
	Value Data   `bson:"value" json:"value"`
}

type Data struct {
	Amount   []CoinStr `bson:"amount" json:"amount"`
	Receiver string    `bson:"receiver" json:"receiver"`
	Sender   string    `bson:"sender" json:"sender"`
	//Source   bool      `bson:"source" json:"source"`
	//Timeout  string    `bson:"timeout" json:"timeout"`
}

// MsgPacket receives incoming IBC packet
type DocTxMsgIBCMsgPacket struct {
	Packet Packet `bson:"packet"`
	//Proof       string `bson:"proof"`
	ProofHeight uint64 `bson:"proof_height" `
	Signer      string `bson:"signer"`
}

func (m *DocTxMsgIBCMsgPacket) Type() string {
	return constant.TxMsgTypeIBCBankMsgPacket
}

func (m *DocTxMsgIBCMsgPacket) BuildMsg(txMsg interface{}) {
	msg := txMsg.(model.IBCPacket)

	var sendpacket SendPacket
	json.Unmarshal(msg.GetData(), &sendpacket)

	packet := Packet{
		Sequence:           msg.Packet.GetSequence(),
		TimeoutHeight:      msg.Packet.GetTimeoutHeight(),
		SourcePort:         msg.Packet.GetSourcePort(),
		SourceChannel:      msg.Packet.GetSourceChannel(),
		DestinationPort:    msg.Packet.GetDestPort(),
		DestinationChannel: msg.Packet.GetDestChannel(),
		Data:               sendpacket,
	}

	m.Packet = packet
	m.Signer = msg.Signer.String()
	m.ProofHeight = msg.ProofHeight
	//proofdata, _ := json.Marshal(msg.Proof)
	//m.Proof = string(proofdata)
}
