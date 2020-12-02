package cdc

import (
	"github.com/cosmos/cosmos-sdk/codec"
	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/irismod/coinswap"
	"github.com/irismod/htlc"
	"github.com/irismod/nft"
	"github.com/irismod/record"
	"github.com/irismod/service"
	"github.com/irismod/token"
)

var (
	encodecfg    params.EncodingConfig
	moduleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		service.AppModuleBasic{},
		nft.AppModuleBasic{},
		htlc.AppModuleBasic{},
		coinswap.AppModuleBasic{},
		record.AppModuleBasic{},
		token.AppModuleBasic{},
		gov.AppModuleBasic{},
		staking.AppModuleBasic{},
		distribution.AppModuleBasic{},
		slashing.AppModuleBasic{},
		evidence.AppModuleBasic{},
		crisis.AppModuleBasic{},
		htlc.AppModuleBasic{},
		coinswap.AppModuleBasic{},
	)
)

// 初始化账户地址前缀
func init() {
	var cdc = codec.New()

	interfaceRegistry := ctypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, std.DefaultPublicKeyCodec{}, tx.DefaultSignModes)

	encodecfg = params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             cdc,
	}
	std.RegisterCodec(encodecfg.Amino)
	std.RegisterInterfaces(encodecfg.InterfaceRegistry)
	moduleBasics.RegisterCodec(encodecfg.Amino)
	moduleBasics.RegisterInterfaces(encodecfg.InterfaceRegistry)
}
func GetTxDecoder() sdk.TxDecoder {
	return encodecfg.TxConfig.TxDecoder()
}
