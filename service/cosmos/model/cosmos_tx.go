package cosmos

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type (
	CosmosTx struct {
		Time      time.Time `bson:"time"`
		Height    int64     `bson:"height"`
		TxHash    string    `bson:"tx_hash"`
		From      string    `bson:"from"`
		To        string    `bson:"to"`
		Initiator string    `bson:"initiator"`
		Amount    []*Coin   `bson:"amount"`
		Type      string    `bson:"type"`
		Fee       *Fee      `bson:"fee"`
		Memo      string    `bson:"memo"`
		Status    string    `bson:"status"`
		Code      uint32    `bson:"code"`
		Events    []Event   `bson:"events"`
		//Msg    Msg       `bson:"msg"`
	}
)

const (
	CollectionNameCosmosTx = "sync_cosmos_tx"
)

func (d CosmosTx) Name() string {
	return CollectionNameCosmosTx
}

func (d CosmosTx) PkKvPair() map[string]interface{} {
	return bson.M{}
}

type Coin struct {
	Denom  string `bson:"denom"  json:"denom"`
	Amount int64  `bson:"amount"  json:"amount"`
}

type Coins []*Coin

type Fee struct {
	Amount []*Coin `bson:"amount" json:"amount"`
	Gas    int64   `bson:"gas" json:"gas"`
}

type Tag map[string]string

type RawLog struct {
	MsgIndex int    `json:"msg_index"`
	Success  bool   `json:"success"`
	Log      string `json:"log"`
}

type Event struct {
	Type       string            `bson:"type" json:"type"`
	Attributes map[string]string `bson:"attributes" json:"attributes"`
}
