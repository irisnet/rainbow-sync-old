package ibc

import (
	. "github.com/irisnet/rainbow-sync/msgs"
)

// MsgUpdateClient defines a message to update an IBC client
type DocMsgUpdateClient struct {
	ClientID string `bson:"client_id" yaml:"client_id"`
	Header   string `bson:"header" yaml:"header"`
	Signer   string `bson:"signer" yaml:"signer"`
}

func (m *DocMsgUpdateClient) GetType() string {
	return MsgTypeUpdateClient
}

func (m *DocMsgUpdateClient) BuildMsg(v interface{}) {
	msg := v.(*MsgUpdateClient)

	m.ClientID = msg.ClientId
	m.Signer = msg.Signer
	if msg.Header != nil {
		m.Header = msg.Header.String()
	}
}

func (m *DocMsgUpdateClient) HandleTxMsg(v SdkMsg) MsgDocInfo {
	var (
		addrs []string
	)

	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
