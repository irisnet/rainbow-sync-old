package block

import (
	"github.com/irisnet/rainbow-sync/conf"
	"github.com/kaifei-bianjie/msg-parser/codec"
)

var (

	// Bech32ChainPrefix defines the prefix of this chain
	Bech32ChainPrefix = "i"

	// PrefixAcc is the prefix for account
	PrefixAcc = "a"

	// PrefixValidator is the prefix for validator keys
	PrefixValidator = "v"

	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus = "c"

	// PrefixPublic is the prefix for public
	PrefixPublic = "p"

	// PrefixAddress is the prefix for address
	PrefixAddress = "a"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = conf.SvrConf.Bech32ChainPrefix + PrefixAcc + PrefixAddress
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = conf.SvrConf.Bech32ChainPrefix + PrefixAcc + PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = conf.SvrConf.Bech32ChainPrefix + PrefixValidator + PrefixAddress
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = conf.SvrConf.Bech32ChainPrefix + PrefixValidator + PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = conf.SvrConf.Bech32ChainPrefix + PrefixConsensus + PrefixAddress
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = conf.SvrConf.Bech32ChainPrefix + PrefixConsensus + PrefixPublic

	TxStatusSuccess = "success"
	TxStatusFail    = "fail"
)

func init() {
	codec.SetBech32Prefix(Bech32PrefixAccAddr, Bech32PrefixAccPub, Bech32PrefixValAddr, Bech32PrefixValPub, Bech32PrefixConsAddr, Bech32PrefixConsPub)
}
