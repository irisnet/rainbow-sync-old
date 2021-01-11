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
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/irisnet/irismod/modules/coinswap"
	"github.com/irisnet/irismod/modules/htlc"
	"github.com/irisnet/irismod/modules/nft"
	"github.com/irisnet/irismod/modules/record"
	"github.com/irisnet/irismod/modules/service"
	"github.com/irisnet/irismod/modules/token"
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
		sdkparams.AppModuleBasic{},
		upgrade.AppModuleBasic{},
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
	var cdc = codec.NewLegacyAmino()

	interfaceRegistry := ctypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	encodecfg = params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             cdc,
	}
	std.RegisterLegacyAminoCodec(encodecfg.Amino)
	std.RegisterInterfaces(encodecfg.InterfaceRegistry)
	moduleBasics.RegisterLegacyAminoCodec(encodecfg.Amino)
	moduleBasics.RegisterInterfaces(encodecfg.InterfaceRegistry)
}
func GetTxDecoder() sdk.TxDecoder {
	return encodecfg.TxConfig.TxDecoder()
}
