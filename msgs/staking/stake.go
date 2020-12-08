package staking

import (
	stake "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/irisnet/rainbow-sync/model"
	. "github.com/irisnet/rainbow-sync/msgs"
)

// MsgDelegate - struct for bonding transactions
type DocTxMsgBeginRedelegate struct {
	DelegatorAddress    string `bson:"delegator_address"`
	ValidatorSrcAddress string `bson:"validator_src_address"`
	ValidatorDstAddress string `bson:"validator_dst_address"`
	Amount              string `bson:"amount"`
}

type Description struct {
	Moniker         string `bson:"moniker"`
	Identity        string `bson:"identity"`
	Website         string `bson:"website"`
	SecurityContact string `bson:"security_contact"`
	Details         string `bson:"details"`
}

type CommissionRates struct {
	Rate          string `bson:"rate"`            // the commission rate charged to delegators
	MaxRate       string `bson:"max_rate"`        // maximum commission rate which validator can ever charge
	MaxChangeRate string `bson:"max_change_rate"` // maximum daily increase of the validator commission
}

func (doctx *DocTxMsgBeginRedelegate) GetType() string {
	return MsgTypeBeginRedelegate
}

func (doctx *DocTxMsgBeginRedelegate) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgBeginRedelegate)
	doctx.DelegatorAddress = msg.DelegatorAddress
	doctx.ValidatorSrcAddress = msg.ValidatorSrcAddress
	doctx.ValidatorDstAddress = msg.ValidatorDstAddress
	doctx.Amount = msg.Amount.String()
}
func (m *DocTxMsgBeginRedelegate) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgBeginRedelegate
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.DelegatorAddress, msg.ValidatorDstAddress, msg.ValidatorSrcAddress)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}

// MsgBeginUnbonding - struct for unbonding transactions
type DocTxMsgBeginUnbonding struct {
	DelegatorAddress string `bson:"delegator_address"`
	ValidatorAddress string `bson:"validator_address"`
	Amount           string `bson:"amount"`
}

func (doctx *DocTxMsgBeginUnbonding) GetType() string {
	return MsgTypeStakeBeginUnbonding
}

func (doctx *DocTxMsgBeginUnbonding) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgUndelegate)
	doctx.ValidatorAddress = msg.ValidatorAddress
	doctx.DelegatorAddress = msg.DelegatorAddress
	doctx.Amount = msg.Amount.String()
}
func (m *DocTxMsgBeginUnbonding) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgUndelegate
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.DelegatorAddress, msg.ValidatorAddress)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}

// MsgDelegate - struct for bonding transactions
type DocTxMsgDelegate struct {
	DelegatorAddress string     `bson:"delegator_address"`
	ValidatorAddress string     `bson:"validator_address"`
	Amount           model.Coin `bson:"amount"`
}

func (doctx *DocTxMsgDelegate) GetType() string {
	return MsgTypeStakeDelegate
}

func (doctx *DocTxMsgDelegate) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgDelegate)
	doctx.ValidatorAddress = msg.ValidatorAddress
	doctx.DelegatorAddress = msg.DelegatorAddress
	doctx.Amount = model.BuildDocCoin(msg.Amount)
}
func (m *DocTxMsgDelegate) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgDelegate
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.DelegatorAddress, msg.ValidatorAddress)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}

// MsgEditValidator - struct for editing a validator
type DocMsgEditValidator struct {
	Description       Description `bson:"description"`
	ValidatorAddress  string      `bson:"validator_address"`
	CommissionRate    string      `bson:"commission_rate"`
	MinSelfDelegation string      `bson:"min_self_delegation"`
}

func (doctx *DocMsgEditValidator) GetType() string {
	return MsgTypeStakeEditValidator
}

func (doctx *DocMsgEditValidator) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgEditValidator)
	doctx.ValidatorAddress = msg.ValidatorAddress
	commissionRate := msg.CommissionRate
	if commissionRate == nil {
		doctx.CommissionRate = ""
	} else {
		doctx.CommissionRate = commissionRate.String()
	}
	doctx.Description = loadDescription(msg.Description)
	doctx.MinSelfDelegation = msg.MinSelfDelegation.String()
}
func (m *DocMsgEditValidator) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgEditValidator
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.ValidatorAddress)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}

// MsgCreateValidator defines an SDK message for creating a new validator.
type DocTxMsgCreateValidator struct {
	Description       Description     `bson:"description"`
	Commission        CommissionRates `bson:"commission"`
	MinSelfDelegation string          `bson:"min_self_delegation"`
	DelegatorAddress  string          `bson:"delegator_address"`
	ValidatorAddress  string          `bson:"validator_address"`
	Pubkey            string          `bson:"pubkey"`
	Value             model.Coin      `bson:"value"`
}

func (doctx *DocTxMsgCreateValidator) GetType() string {
	return MsgTypeStakeCreateValidator
}

func (doctx *DocTxMsgCreateValidator) BuildMsg(txMsg interface{}) {
	msg := txMsg.(*MsgCreateValidator)
	doctx.ValidatorAddress = msg.ValidatorAddress
	doctx.Pubkey = msg.Pubkey
	doctx.DelegatorAddress = msg.DelegatorAddress
	doctx.MinSelfDelegation = msg.MinSelfDelegation.String()
	doctx.Commission = CommissionRates{
		Rate:          msg.Commission.Rate.String(),
		MaxChangeRate: msg.Commission.MaxChangeRate.String(),
		MaxRate:       msg.Commission.MaxRate.String(),
	}
	doctx.Description = loadDescription(msg.Description)
	doctx.Value = model.BuildDocCoin(msg.Value)
}
func (m *DocTxMsgCreateValidator) HandleTxMsg(v SdkMsg) MsgDocInfo {

	var (
		addrs []string
		msg   MsgCreateValidator
	)

	ConvertMsg(v, &msg)
	addrs = append(addrs, msg.DelegatorAddress, msg.ValidatorAddress)
	handler := func() (Msg, []string) {
		return m, addrs
	}

	return CreateMsgDocInfo(v, handler)
}

func loadDescription(description stake.Description) Description {
	return Description{
		Moniker:         description.Moniker,
		Details:         description.Details,
		Identity:        description.Identity,
		Website:         description.Website,
		SecurityContact: description.SecurityContact,
	}
}
