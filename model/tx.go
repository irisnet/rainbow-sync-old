package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/rainbow-sync/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Tx struct {
	Time      int64       `bson:"time"`
	Height    int64       `bson:"height"`
	TxHash    string      `bson:"tx_hash"`
	Type      string      `bson:"type"`
	Fee       *Fee        `bson:"fee"`
	ActualFee *ActualFee  `bson:"actual_fee"`
	Memo      string      `bson:"memo"`
	Status    string      `bson:"status"`
	Log       string      `bson:"log"`
	Types     []string    `bson:"types"`
	Events    []Event     `bson:"events"`
	Msgs      []DocTxMsg  `bson:"msgs"`
	Signers   []string    `bson:"signers"`
	Addrs     []string    `bson:"addrs"`
	TxIndex   uint32      `bson:"tx_index"`
	Ext       interface{} `bson:"ext"`
}

type (
	Event struct {
		Type       string   `bson:"type"`
		Attributes []KvPair `bson:"attributes"`
	}

	KvPair struct {
		Key   string `bson:"key"`
		Value string `bson:"value"`
	}
)

type DocTxMsg struct {
	Type string `bson:"type"`
	Msg  Msg    `bson:"msg"`
}

type Msg interface {
	GetType() string
	BuildMsg(msg interface{})
}

const (
	CollectionNameIrisTx = "sync_iris_tx"
)

func (d Tx) Name() string {
	return CollectionNameIrisTx
}

func (d Tx) PkKvPair() map[string]interface{} {
	return bson.M{"height": d.Height, "tx_index": d.TxIndex}
}

func (d Tx) EnsureIndexes() {
	var indexes []mgo.Index
	indexes = append(indexes,
		mgo.Index{
			Key:        []string{"-height", "-tx_index"},
			Unique:     true,
			Background: true},
		mgo.Index{
			Key:        []string{"-tx_hash"},
			Background: true},
		mgo.Index{
			Key:        []string{"-type"},
			Background: true},
	)

	db.EnsureIndexes(d.Name(), indexes)
}

type Coin struct {
	Denom  string `bson:"denom"`
	Amount string `bson:"amount"`
}

type Coins []*Coin

type Fee struct {
	Amount []Coin `bson:"amount"`
	Gas    int64  `bson:"gas"`
}

type ActualFee struct {
	Denom  string `bson:"denom"`
	Amount string `bson:"amount"`
}

func BuildDocCoins(coins sdk.Coins) []Coin {
	var (
		res []Coin
	)
	if len(coins) > 0 {
		for _, v := range coins {
			c := Coin{
				Denom:  v.Denom,
				Amount: v.Amount.String(),
			}
			res = append(res, c)
		}
	}

	return res
}

func BuildDocCoin(coin sdk.Coin) Coin {
	return Coin{
		Denom:  coin.Denom,
		Amount: coin.Amount.String(),
	}
}
