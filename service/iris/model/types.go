package iris

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	ibcBank "github.com/cosmos/cosmos-sdk/x/ibc/mock/bank"
)

type (
	MsgTransfer             = bank.MsgSend
	IBCBankMsgTransfer      = ibcBank.MsgTransfer
	IBCBankMsgReceivePacket = ibcBank.MsgRecvTransferPacket
	IBCPacket               = ibcBank.Packet

	SdkCoins = types.Coins
	KVPair   = types.KVPair

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
