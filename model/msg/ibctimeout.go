package msg

import (
	"github.com/irisnet/rainbow-sync/constant"
	"github.com/irisnet/rainbow-sync/model"
	"encoding/json"
)

// MsgTimeout receives timed-out packet
type DocTxMsgIBCTimeout struct {
	Packet                       `bson:"packet"`
	NextSequenceRecv uint64      `bson:"next_sequence_recv"`
	Proof            interface{} `bson:"proof"`
	ProofHeight      uint64      `bson:"proof_height"`
	Signer           string      `bson:"signer"`
}

func (m *DocTxMsgIBCTimeout) Type() string {
	return constant.TxMsgTypeIBCBankMsgTimeout
}

func (m *DocTxMsgIBCTimeout) BuildMsg(txMsg interface{}) {
	msg := txMsg.(model.IBCTimeout)

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
	proofdata, _ := json.Marshal(msg.Proof)
	m.Proof = string(proofdata)
	m.NextSequenceRecv = msg.NextSequenceRecv
	m.ProofHeight = msg.ProofHeight
	m.Signer = msg.Signer.String()
}
