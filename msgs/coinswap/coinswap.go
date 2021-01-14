package coinswap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/rainbow-sync/model"
	. "github.com/irisnet/rainbow-sync/msgs"
)

type DocTxMsgAddLiquidity struct {
	MaxToken     model.Coin `bson:"max_token"`      // coin to be deposited as liquidity with an upper bound for its amount
	ExactIrisAmt string     `bson:"exact_iris_amt"` // exact amount of native asset being add to the liquidity pool
	MinLiquidity string     `bson:"min_liquidity"`  // lower bound UNI sender is willing to accept for deposited coins
	Deadline     int64      `bson:"deadline"`
	Sender       string     `bson:"sender"`
}

func (doctx *DocTxMsgAddLiquidity) GetType() string {
	return MsgTypeAddLiquidity
}

func (doctx *DocTxMsgAddLiquidity) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgAddLiquidity)
	doctx.Sender = msg.Sender
	doctx.MinLiquidity = msg.MinLiquidity.String()
	doctx.ExactIrisAmt = msg.ExactStandardAmt.String()
	doctx.Deadline = msg.Deadline
	doctx.MaxToken = model.BuildDocCoin(msg.MaxToken)
}

func (m *DocTxMsgAddLiquidity) HandleTxMsg(v sdk.Msg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgAddLiquidity
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.Sender)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}

type DocTxMsgRemoveLiquidity struct {
	MinToken          string     `bson:"min_token"`          // coin to be withdrawn with a lower bound for its amount
	WithdrawLiquidity model.Coin `bson:"withdraw_liquidity"` // amount of UNI to be burned to withdraw liquidity from a reserve pool
	MinIrisAmt        string     `bson:"min_iris_amt"`       // minimum amount of the native asset the sender is willing to accept
	Deadline          int64      `bson:"deadline"`
	Sender            string     `bson:"sender"`
}

func (doctx *DocTxMsgRemoveLiquidity) GetType() string {
	return MsgTypeRemoveLiquidity
}

func (doctx *DocTxMsgRemoveLiquidity) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgRemoveLiquidity)
	doctx.Sender = msg.Sender
	doctx.MinIrisAmt = msg.MinStandardAmt.String()
	doctx.MinToken = msg.MinToken.String()
	doctx.Deadline = msg.Deadline
	doctx.WithdrawLiquidity = model.BuildDocCoin(msg.WithdrawLiquidity)
}
func (m *DocTxMsgRemoveLiquidity) HandleTxMsg(v sdk.Msg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgRemoveLiquidity
	)
	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.Sender)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}

type DocTxMsgSwapOrder struct {
	Input      Input  `bson:"input"`        // the amount the sender is trading
	Output     Output `bson:"output"`       // the amount the sender is receiving
	Deadline   int64  `bson:"deadline"`     // deadline for the transaction to still be considered valid
	IsBuyOrder bool   `bson:"is_buy_order"` // boolean indicating whether the order should be treated as a buy or sell
}

type Input struct {
	Address string     `bson:"address"`
	Coin    model.Coin `bson:"coin"`
}

type Output struct {
	Address string     `bson:"address"`
	Coin    model.Coin `bson:"coin"`
}

func (doctx *DocTxMsgSwapOrder) GetType() string {
	return MsgTypeSwapOrder
}

func (doctx *DocTxMsgSwapOrder) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgSwapOrder)
	doctx.Deadline = msg.Deadline
	doctx.IsBuyOrder = msg.IsBuyOrder
	doctx.Input = Input{
		Address: msg.Input.Address,
		Coin:    model.BuildDocCoin(msg.Input.Coin),
	}
	doctx.Output = Output{
		Address: msg.Output.Address,
		Coin:    model.BuildDocCoin(msg.Output.Coin),
	}
}
func (m *DocTxMsgSwapOrder) HandleTxMsg(v sdk.Msg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgSwapOrder
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.Output.Address, msg.Input.Address)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
