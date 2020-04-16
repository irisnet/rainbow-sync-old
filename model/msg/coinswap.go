package msg

import (
	"github.com/irisnet/rainbow-sync/constant"
	"github.com/irisnet/rainbow-sync/utils"
	types "github.com/irisnet/rainbow-sync/model"
)

// MsgAddLiquidity - struct for adding liquidity to a reserve pool
type DocTxMsgAddLiquidity struct {
	MaxToken         *types.Coin `bson:"max_token"`          // coin to be deposited as liquidity with an upper bound for its amount
	ExactStandardAmt string      `bson:"exact_standard_amt"` // exact amount of native asset being add to the liquidity pool
	MinLiquidity     string      `bson:"min_liquidity"`      // lower bound UNI sender is willing to accept for deposited coins
	Deadline         int64       `bson:"deadline"`           // deadline of tx
	Sender           string      `bson:"sender"`             // msg sender
}

func (doctx *DocTxMsgAddLiquidity) Type() string {
	return constant.TxTypeAddLiquidity
}

func (doctx *DocTxMsgAddLiquidity) BuildMsg(txMsg interface{}) {
	msg := txMsg.(types.MsgAddLiquidity)
	doctx.Sender = msg.Sender.String()
	doctx.MinLiquidity = msg.MinLiquidity.String()
	doctx.ExactStandardAmt = msg.ExactStandardAmt.String()
	doctx.Deadline = msg.Deadline
	doctx.MaxToken = utils.ParseRewards(msg.MaxToken.String())
}

// MsgRemoveLiquidity - struct for removing liquidity from a reserve pool
type DocTxMsgRemoveLiquidity struct {
	MinToken          string      `bson:"min_token"`          // coin to be withdrawn with a lower bound for its amount
	WithdrawLiquidity *types.Coin `bson:"withdraw_liquidity"` // amount of UNI to be burned to withdraw liquidity from a reserve pool
	MinStandardAmt    string      `bson:"min_standard_amt"`   // minimum amount of the native asset the sender is willing to accept
	Deadline          int64       `bson:"deadline"`           // deadline of tx
	Sender            string      `bson:"sender"`             // msg sender
}

func (doctx *DocTxMsgRemoveLiquidity) Type() string {
	return constant.TxTypeRemoveLiquidity
}

func (doctx *DocTxMsgRemoveLiquidity) BuildMsg(txMsg interface{}) {
	msg := txMsg.(types.MsgRemoveLiquidity)
	doctx.Sender = msg.Sender.String()
	doctx.MinStandardAmt = msg.MinStandardAmt.String()
	doctx.MinToken = msg.MinToken.String()
	doctx.Deadline = msg.Deadline
	doctx.WithdrawLiquidity = utils.ParseRewards(msg.WithdrawLiquidity.String())
}

type DocTxMsgSwapOrder struct {
	Input      Input  `bson:"input"`        // the amount the sender is trading
	Output     Output `bson:"output"`       // the amount the sender is receiving
	Deadline   int64  `bson:"deadline"`     // deadline for the transaction to still be considered valid
	IsBuyOrder bool   `bson:"is_buy_order"` // boolean indicating whether the order should be treated as a buy or sell
}

type Input struct {
	Address string      `bson:"address"`
	Coin    *types.Coin `bson:"coin"`
}

type Output struct {
	Address string      `bson:"address"`
	Coin    *types.Coin `bson:"coin"`
}

func (doctx *DocTxMsgSwapOrder) Type() string {
	return constant.TxTypeSwapOrder
}

func (doctx *DocTxMsgSwapOrder) BuildMsg(txMsg interface{}) {
	msg := txMsg.(types.MsgSwapOrder)
	doctx.Deadline = msg.Deadline
	doctx.IsBuyOrder = msg.IsBuyOrder
	doctx.Input = Input{
		Address: msg.Input.Address.String(),
		Coin:    utils.ParseRewards(msg.Input.Coin.String()),
	}
	doctx.Output = Output{
		Address: msg.Output.Address.String(),
		Coin:    utils.ParseRewards(msg.Output.Coin.String()),
	}
}
