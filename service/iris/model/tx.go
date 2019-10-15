package iris

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/types"
)

type IrisTx struct {
	Time      time.Time  `json:"time" bson:"time"`
	Height    int64      `json:"height" bson:"height"`
	TxHash    string     `json:"tx_hash" bson:"tx_hash"`
	From      string     `json:"from" bson:"from"`
	To        string     `json:"to" bson:"to"`
	Initiator string     `json:"initiator" bson:"initiator"`
	Amount    []*Coin    `json:"amount" bson:"amount"`
	Type      string     `json:"type" bson:"type"`
	Fee       *Fee       `json:"fee" bson:"fee"`
	Memo      string     `json:"memo" bson:"memo"`
	Status    string     `json:"status" bson:"status"`
	Code      uint32     `json:"code" bson:"code"`
	Log       string     `json:"log" bson:"log"`
	Events    []Event    `bson:"events"`
	Msgs      []DocTxMsg `bson:"msgs"`
}

type DocTxMsg struct {
	Type string `bson:"type"`
	Msg  Msg    `bson:"msg"`
}

type Msg interface {
	Type() string
	BuildMsg(msg interface{})
}

type Event struct {
	Type       string            `bson:"type" json:"type"`
	Attributes map[string]string `bson:"attributes" json:"attributes"`
}
const (
	CollectionNameIrisTx = "sync_iris_tx"
)

func (d IrisTx) Name() string {
	return CollectionNameIrisTx
}

func (d IrisTx) PkKvPair() map[string]interface{} {
	return bson.M{}
}

type (
	MsgTransfer = bank.MsgSend

	SdkCoins = types.Coins
	KVPair = types.KVPair
)

type Coin struct {
	Denom  string `bson:"denom"  json:"denom"`
	Amount int64  `bson:"amount"  json:"amount"`
}

type Coins []*Coin

type Fee struct {
	Amount []*Coin `bson:"amount" json:"amount"`
	Gas    int64   `bson:"gas" json:"gas"`
}
