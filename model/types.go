package model

import (
	"github.com/cosmos/cosmos-sdk/x/bank"
	ibcBank "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/types"
	ibcChannel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/irisnet/rainbow-sync/db"
	"github.com/irismod/coinswap"
)

var (
	Collections = []db.Docs{
		new(SyncZoneTask),
		new(ZoneTx),
		new(Block),
	}
)

func CheckIndex() {
	for _, v := range Collections {
		v.EnsureIndexes()
	}
}

type (
	MsgTransfer = bank.MsgSend
	MsgMultiSend = bank.MsgMultiSend

	IBCBankMsgTransfer = ibcBank.MsgTransfer
	IBCPacket = ibcChannel.MsgPacket
	IBCTimeout = ibcChannel.MsgTimeout
	MsgAddLiquidity = coinswap.MsgAddLiquidity
	MsgRemoveLiquidity = coinswap.MsgRemoveLiquidity
	MsgSwapOrder = coinswap.MsgSwapOrder
)
