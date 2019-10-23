package cosmos

import (
	"github.com/cosmos/cosmos-sdk/x/bank"
	ibcBank "github.com/cosmos/cosmos-sdk/x/ibc/mock/bank"
)

type (
	MsgTransfer  = bank.MsgSend
	MsgMultiSend = bank.MsgMultiSend

	IBCBankMsgTransfer      = ibcBank.MsgTransfer
	IBCBankMsgReceivePacket = ibcBank.MsgRecvTransferPacket
	IBCPacket               = ibcBank.Packet

	IBCTransferPacketData struct {
		Type  string                     `json:"type"`
		Value IBCTransferPacketDataValue `json:"value"`
	}

	IBCTransferPacketDataValue struct {
		Denomination string `json:"denomination"`
		Amount       string `json:"amount"`
		Sender       string `json:"sender"`
		Receiver     string `json:"receiver"`
		Source       bool   `json:"source"`
	}
)
