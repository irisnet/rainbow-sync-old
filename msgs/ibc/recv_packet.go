package ibc

import (
	"github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	. "github.com/irisnet/rainbow-sync/msgs"
)

type DocMsgRecvPacket struct {
	Packet          Packet `bson:"packet"`
	ProofCommitment string `bson:"proof_commitment"`
	ProofHeight     Height `bson:"proof_height"`
	Signer          string `bson:"signer"`
}

type Height struct {
	// the epoch that the client is currently on
	RevisionNumber uint64 `bson:"revision_number"`
	// the height within the given epoch
	RevisionHeight uint64 `bson:"revision_height"`
}

type Packet struct {
	// number corresponds to the order of sends and receives, where a Packet
	// with an earlier sequence number must be sent and received before a Packet
	// with a later sequence number.
	Sequence uint64 `bson:"sequence"`
	// identifies the port on the sending chain.
	SourcePort string `bson:"source_port"`
	// identifies the channel end on the sending chain.
	SourceChannel string `bson:"source_channel"`
	// identifies the port on the receiving chain.
	DestinationPort string `bson:"destination_port"`
	// identifies the channel end on the receiving chain.
	DestinationChannel string `bson:"destination_channel"`
	// actual opaque bytes transferred directly to the application module
	Data string `bson:"data"`
	// block height after which the packet times out
	TimeoutHeight Height `bson:"timeout_height"`
	// block timestamp (in nanoseconds) after which the packet times out
	TimeoutTimestamp uint64 `bson:"timeout_timestamp"`
}

func (m *DocMsgRecvPacket) GetType() string {
	return MsgTypeRecvPacket
}

func (m *DocMsgRecvPacket) BuildMsg(v interface{}) {
	msg := v.(*MsgRecvPacket)

	m.ProofCommitment = string(msg.ProofCommitment)
	m.ProofHeight = Height{RevisionHeight: msg.ProofHeight.RevisionHeight, RevisionNumber: msg.ProofHeight.RevisionNumber}
	m.Signer = msg.Signer

	m.Packet = DecodeToIBCRecord(msg.Packet)
}
func DecodeToIBCRecord(packet types.Packet) Packet {
	return Packet{
		Sequence:           packet.Sequence,
		SourcePort:         packet.SourcePort,
		SourceChannel:      packet.SourceChannel,
		DestinationChannel: packet.DestinationChannel,
		DestinationPort:    packet.DestinationPort,
		Data:               string(packet.Data),
		TimeoutHeight:      Height{RevisionHeight: packet.TimeoutHeight.RevisionHeight, RevisionNumber: packet.TimeoutHeight.RevisionNumber},
		TimeoutTimestamp:   packet.TimeoutTimestamp,
	}
}

func (m *DocMsgRecvPacket) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
	)

	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
