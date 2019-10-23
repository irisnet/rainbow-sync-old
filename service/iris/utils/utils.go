package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	ibcBank "github.com/cosmos/cosmos-sdk/x/ibc/mock/bank"
	"github.com/irisnet/rainbow-sync/service/iris/conf"
	"github.com/irisnet/rainbow-sync/service/iris/constant"
	"github.com/irisnet/rainbow-sync/service/iris/helper"
	"github.com/irisnet/rainbow-sync/service/iris/logger"
	imodel "github.com/irisnet/rainbow-sync/service/iris/model"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
	"strings"
	"time"
)

var (
	cdc          *codec.Codec
	ModuleBasics = module.NewBasicManager(
		bank.AppModuleBasic{},
		auth.AppModuleBasic{},
		ibc.AppModule{},
		ibcBank.AppModule{},
	)
)

// 初始化账户地址前缀
func init() {
	cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
	config := sdk.GetConfig()
	conf.SetNetworkType(conf.IrisNetwork)
	irisConfig := conf.GetConfig()
	config.SetBech32PrefixForAccount(irisConfig.GetBech32AccountAddrPrefix(), irisConfig.GetBech32AccountPubPrefix())
	config.SetBech32PrefixForValidator(irisConfig.GetBech32ValidatorAddrPrefix(), irisConfig.GetBech32ValidatorPubPrefix())
	config.SetBech32PrefixForConsensusNode(irisConfig.GetBech32ConsensusAddrPrefix(), irisConfig.GetBech32ConsensusPubPrefix())
	config.Seal()

}

func GetCodec() *codec.Codec {
	return cdc.Seal()
}

func BuildHex(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}

func ParseCoins(coinsStr sdk.Coins) (coins []*imodel.Coin) {

	coins = make([]*imodel.Coin, 0, len(coinsStr))
	for _, coinStr := range coinsStr {
		coin := ParseCoin(coinStr)
		coins = append(coins, &coin)
	}
	return coins
}

func ParseCoin(sdkcoin sdk.Coin) (coin imodel.Coin) {
	amount, err := strconv.ParseFloat(sdkcoin.Amount.String(), 64)
	if err != nil {
		logger.Error("ParseCoin have error", logger.String("error", err.Error()))
	}
	return imodel.Coin{
		Denom:  sdkcoin.Denom,
		Amount: amount,
	}

}

func getPrecision(amount string) string {
	length := len(amount)
	if length > 15 {
		nums := strings.Split(amount, ".")
		if len(nums) > 2 {
			return amount
		}

		if len_num0 := len(nums[0]); len_num0 > 15 {
			amount = string([]byte(nums[0])[:15])
			for i := 1; i <= len_num0-15; i++ {
				amount += "0"
			}
		} else {
			//leng_num1 := len(nums[1])
			leng_append := 16 - len_num0
			amount = nums[0] + "." + string([]byte(nums[1])[:leng_append])
			//for i := 1; i <= leng_num1-leng_append; i++ {
			//	amount += "0"
			//}
		}
	}
	return amount
}

func BuildFee(fee auth.StdFee) imodel.Fee {
	return imodel.Fee{
		Amount: ParseCoins(fee.Amount),
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
		logger.Warn("QueryTxResult have error, now try again", logger.String("err", err.Error()))
		time.Sleep(time.Duration(1) * time.Second)
		var err1 error
		client2 := helper.GetClient()
		res, err1 = client2.Tx(txHash, false)
		client2.Release()
		if err1 != nil {
			return "unknown", resDeliverTx, err1
		}
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

func Md5Encrypt(data []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func Base64Decode(str string) ([]byte, error) {
	enc := base64.Encoding{}
	data, err := enc.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return data, nil
}
