package utils

import (
	"strings"
	"encoding/hex"
	"strconv"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	imodel "github.com/irisnet/rainbow-sync/service/iris/model"
	"github.com/irisnet/rainbow-sync/service/iris/helper"
	"github.com/irisnet/rainbow-sync/service/iris/constant"
	"github.com/irisnet/rainbow-sync/service/iris/logger"
	"fmt"
	"regexp"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/types"
	"github.com/irisnet/rainbow-sync/service/iris/conf"
)

var (
	cdc *codec.Codec
)

// 初始化账户地址前缀
func init() {
	if conf.IrisNetwork == types.Mainnet {
		types.SetNetworkType(types.Mainnet)
	}
	cdc = app.MakeLatestCodec()
}

func GetCodec() *codec.Codec {
	return cdc
}

func BuildHex(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}

func ParseCoins(coinsStr string) (coins imodel.Coins) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		coin := ParseCoin(coinStr)
		coins = append(coins, coin)
	}
	return coins
}

func ParseCoin(coinStr string) (coin *imodel.Coin) {
	var (
		reDnm  = `[A-Za-z0-9]{2,}\S*`
		reAmt  = `[0-9]+[.]?[0-9]*`
		reSpc  = `[[:space:]]*`
		reCoin = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmt, reSpc, reDnm))
	)

	coinStr = strings.TrimSpace(coinStr)

	matches := reCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		logger.Error("invalid coin expression", logger.Any("coin", coinStr))
		return coin
	}
	denom, amount := matches[2], matches[1]

	amt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		logger.Error("Convert str to int failed", logger.Any("amount", amount))
		return coin
	}

	return &imodel.Coin{
		Denom:  denom,
		Amount: amt,
	}
}

func BuildFee(fee auth.StdFee) *imodel.Fee {
	return &imodel.Fee{
		Amount: ParseCoins(fee.Amount.String()),
		Gas:    int64(fee.Gas),
	}
}

// get tx status and log by query txHash
func QueryTxResult(txHash []byte) (string, abci.ResponseDeliverTx, error) {
	var resDeliverTx abci.ResponseDeliverTx
	status := constant.TxStatusSuccess

	client := helper.GetClient()
	defer client.Release()

	res, err := client.Tx(txHash, false)
	if err != nil {
		return "unknown", resDeliverTx, err
	}
	result := res.TxResult
	if result.Code != 0 {
		status = constant.TxStatusFail
	}

	return status, result, nil
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func ParseFloat(s string, bit ...int) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logger.Error("common.ParseFloat error", logger.String("value", s))
		return 0
	}

	if len(bit) > 0 {
		return RoundFloat(f, bit[0])
	}
	return f
}

func RoundFloat(num float64, bit int) (i float64) {
	format := "%" + fmt.Sprintf("0.%d", bit) + "f"
	s := fmt.Sprintf(format, num)
	i, err := strconv.ParseFloat(s, 0)
	if err != nil {
		logger.Error("common.RoundFloat error", logger.String("format", format))
		return 0
	}
	return i
}
