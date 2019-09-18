package cosmos

import (
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"strings"
	"github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	cmodel "github.com/irisnet/rainbow-sync/service/cosmos/model"
	"strconv"
	"github.com/irisnet/rainbow-sync/service/cosmos/logger"
)

var (
	cdc *codec.Codec
)

// 初始化账户地址前缀
func init() {
	cdc = app.MakeCodec()
}

func GetCodec() *codec.Codec {
	return cdc
}

func BuildHex(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}

func BuildFee(fee auth.StdFee) cmodel.Fee {
	return cmodel.Fee{
		Amount: ParseCoins(fee.Amount),
		Gas:    int64(fee.Gas),
	}
}

func ParseCoins(coinsStr sdk.Coins) (coins []*cmodel.Coin) {

	coins = make([]*cmodel.Coin, 0, len(coinsStr))
	for _, coinStr := range coinsStr {
		coin := ParseCoin(coinStr)
		coins = append(coins, &coin)
	}
	return coins
}

func ParseCoin(sdkcoin sdk.Coin) (coin cmodel.Coin) {
	amount, err := strconv.ParseInt(sdkcoin.Amount.String(), 10, 64)
	if err != nil {
		logger.Error("ParseCoin have error", logger.String("error", err.Error()))
	}
	return cmodel.Coin{
		Denom:  sdkcoin.Denom,
		Amount: amount,
	}

}
