package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/rainbow-sync/model"
	"github.com/irisnet/rainbow-sync/utils"
)

func CreateMsgDocInfo(msg sdk.Msg, handler func() (Msg, []string)) MsgDocInfo {
	var (
		docTxMsg model.DocTxMsg
		signers  []string
		addrs    []string
	)

	m, addrcollections := handler()

	m.BuildMsg(msg)
	docTxMsg = model.DocTxMsg{
		Type: m.GetType(),
		Msg:  m,
	}

	signers = BuildDocSigners(msg.GetSigners())
	addrs = append(addrs, signers...)
	addrs = append(addrs, addrcollections...)

	return MsgDocInfo{
		DocTxMsg: docTxMsg,
		Signers:  signers,
		Addrs:    addrs,
	}
}
func BuildDocSigners(signers []sdk.AccAddress) []string {
	var (
		//firstSigner string
		allSigners []string
	)
	if len(signers) == 0 {
		return allSigners
	}
	for _, v := range signers {
		allSigners = append(allSigners, v.String())
	}

	return allSigners
}
func ConvertMsg(v interface{}, msg interface{}) {
	utils.UnMarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(v), &msg)
}
