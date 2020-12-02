package distribution

import (
	"github.com/irisnet/rainbow-sync/model"
	. "github.com/irisnet/rainbow-sync/msgs"
)

type DocTxMsgSetWithdrawAddress struct {
	DelegatorAddress string `bson:"delegator_address"`
	WithdrawAddress  string `bson:"withdraw_address"`
}

func (doctx *DocTxMsgSetWithdrawAddress) GetType() string {
	return MsgTypeSetWithdrawAddress
}

func (doctx *DocTxMsgSetWithdrawAddress) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgStakeSetWithdrawAddress)
	doctx.DelegatorAddress = msg.DelegatorAddress
	doctx.WithdrawAddress = msg.WithdrawAddress
}

func (m *DocTxMsgSetWithdrawAddress) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgStakeSetWithdrawAddress
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.DelegatorAddress, msg.WithdrawAddress)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}

// msg struct for delegation withdraw from a single validator
type DocTxMsgWithdrawDelegatorReward struct {
	DelegatorAddress string `bson:"delegator_address"`
	ValidatorAddress string `bson:"validator_address"`
}

func (doctx *DocTxMsgWithdrawDelegatorReward) GetType() string {
	return MsgTypeWithdrawDelegatorReward
}

func (doctx *DocTxMsgWithdrawDelegatorReward) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgWithdrawDelegatorReward)
	doctx.DelegatorAddress = msg.DelegatorAddress
	doctx.ValidatorAddress = msg.ValidatorAddress
}
func (m *DocTxMsgWithdrawDelegatorReward) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgWithdrawDelegatorReward
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.DelegatorAddress, msg.ValidatorAddress)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}

// msg struct for delegation withdraw for all of the delegator's delegations
type DocTxMsgFundCommunityPool struct {
	Amount    []model.Coin `bson:"amount"`
	Depositor string       `bson:"depositor"`
}

func (doctx *DocTxMsgFundCommunityPool) GetType() string {
	return MsgTypeMsgFundCommunityPool
}

func (doctx *DocTxMsgFundCommunityPool) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgFundCommunityPool)
	doctx.Depositor = msg.Depositor
	doctx.Amount = model.BuildDocCoins(msg.Amount)
}
func (m *DocTxMsgFundCommunityPool) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgFundCommunityPool
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.Depositor)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}

// msg struct for validator withdraw
type DocTxMsgWithdrawValidatorCommission struct {
	ValidatorAddress string `bson:"validator_address"`
}

func (doctx *DocTxMsgWithdrawValidatorCommission) GetType() string {
	return MsgTypeMsgWithdrawValidatorCommission
}

func (doctx *DocTxMsgWithdrawValidatorCommission) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgWithdrawValidatorCommission)
	doctx.ValidatorAddress = msg.ValidatorAddress
}

func (m *DocTxMsgWithdrawValidatorCommission) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgWithdrawValidatorCommission
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.ValidatorAddress)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}
