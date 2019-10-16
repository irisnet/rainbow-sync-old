package cosmos

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"encoding/hex"
	"strings"
	cmodel "github.com/irisnet/rainbow-sync/service/cosmos/model"
	"strconv"
	"github.com/irisnet/rainbow-sync/service/cosmos/logger"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/supply"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"regexp"
	"fmt"
)

var (
	cdc          *codec.Codec
	ModuleBasics = module.NewBasicManager(
		bank.AppModuleBasic{},
		auth.AppModuleBasic{},
		staking.AppModuleBasic{},
		gov.AppModuleBasic{},
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
	)
)

// 初始化账户地址前缀
func init() {
	cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	dtypes.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
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

func ParseRewards(coinStr string) (coin *cmodel.Coin) {
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

	amt, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		logger.Error("Convert str to int failed", logger.Any("amount", amount))
		return coin
	}

	return &cmodel.Coin{
		Denom:  denom,
		Amount: amt,
	}
}