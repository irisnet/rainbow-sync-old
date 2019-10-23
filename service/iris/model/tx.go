package iris

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type IrisTx struct {
	Time          time.Time  `bson:"time"`
	Height        int64      `bson:"height"`
	TxHash        string     `bson:"tx_hash"`
	From          string     `bson:"from"`
	To            string     `bson:"to"`
	Initiator     string     `bson:"initiator"`
	Amount        []*Coin    `bson:"amount"`
	Type          string     `bson:"type"`
	Fee           *Fee       `bson:"fee"`
	Memo          string     `bson:"memo"`
	Status        string     `bson:"status"`
	Code          uint32     `bson:"code"`
	Log           string     `bson:"log"`
	Events        []Event    `bson:"events"`
	Msgs          []DocTxMsg `bson:"msgs"`
	IBCPacketHash string     `bson:"ibc_packet_hash"`
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

type Coin struct {
	Denom  string  `bson:"denom"  json:"denom"`
	Amount float64 `bson:"amount"  json:"amount"`
}

type Coins []*Coin

type Fee struct {
	Amount []*Coin `bson:"amount" json:"amount"`
	Gas    int64   `bson:"gas" json:"gas"`
}
