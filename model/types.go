package model

import (
	"github.com/cosmos/cosmos-sdk/x/bank"
	ibcBank "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/types"
	ibcChannel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
)

type (
	MsgTransfer = bank.MsgSend
	MsgMultiSend = bank.MsgMultiSend

	IBCBankMsgTransfer = ibcBank.MsgTransfer
	IBCPacket = ibcChannel.MsgPacket
	IBCTimeout = ibcChannel.MsgTimeout
)
