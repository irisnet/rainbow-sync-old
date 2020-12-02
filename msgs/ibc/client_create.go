package ibc

import (
	. "github.com/irisnet/rainbow-sync/msgs"
)

// MsgCreateClient defines a message to create an IBC client
type DocMsgCreateClient struct {
	ClientID       string `bson:"client_id" yaml:"client_id"`
	ClientState    string `bson:"client_state"`
	ConsensusState string `bson:"consensus_state"`
	Signer         string `bson:"signer" yaml:"signer"`
}

func (m *DocMsgCreateClient) GetType() string {
	return MsgTypeCreateClient
}

func (m *DocMsgCreateClient) BuildMsg(v interface{}) {
	msg := v.(*MsgCreateClient)

	m.ClientID = msg.ClientId
	m.Signer = msg.Signer
	if msg.ClientState != nil {
		m.ClientState = msg.ClientState.String()
	}
	if msg.ConsensusState != nil {
		m.ConsensusState = msg.ConsensusState.String()
	}
}

func (m *DocMsgCreateClient) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
	)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
