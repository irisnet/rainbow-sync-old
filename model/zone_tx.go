package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
	"github.com/irisnet/rainbow-sync/conf"
	"gopkg.in/mgo.v2"
	"github.com/irisnet/rainbow-sync/db"
)

type (
	ZoneTx struct {
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
		Events        []Event    `bson:"events"`
		Msgs          []DocTxMsg `bson:"msgs"`
		IBCPacketHash string     `bson:"ibc_packet_hash"`
	}

	DocTxMsg struct {
		Type string `bson:"type"`
		Msg  Msg    `bson:"msg"`
	}

	Msg interface {
		Type() string
		BuildMsg(msg interface{})
	}
)

const (
	CollectionNameZoneTx = "sync_%v_tx"
	ZoneTxHashTag        = "tx_hash"
	ZoneIbcPacketHashTag = "ibc_packet_hash"
	ZoneTypeTag          = "type"
	ZoneFromTag          = "from"
	ZoneToTag            = "to"
	ZoneHeightTag        = "height"
	ZoneStatusTag        = "status"
	ZoneTnitiatorTag     = "initiator"
)

func (d ZoneTx) Name() string {
	return fmt.Sprintf(CollectionNameZoneTx, conf.ZoneName)
}

func (d ZoneTx) PkKvPair() map[string]interface{} {
	return bson.M{ZoneTxHashTag: d.TxHash}
}

func (d ZoneTx) EnsureIndexes() {
	var indexes []mgo.Index
	indexes = append(indexes,
		mgo.Index{
			Key:        []string{ZoneTxHashTag},
			Unique:     true,
			Background: true,
		}, mgo.Index{
			Key:        []string{ZoneTypeTag},
			Background: true,
		}, mgo.Index{
			Key:        []string{ZoneTnitiatorTag},
			Background: true,
		}, mgo.Index{
			Key:        []string{ZoneIbcPacketHashTag},
			Background: true,
		}, mgo.Index{
			Key:        []string{ZoneStatusTag},
			Background: true,
		}, mgo.Index{
			Key:        []string{ZoneFromTag},
			Background: true,
		}, mgo.Index{
			Key:        []string{ZoneToTag, ZoneHeightTag},
			Background: true,
		})
	db.EnsureIndexes(d.Name(), indexes)
}

type Coin struct {
	Denom  string  `bson:"denom" `
	Amount float64 `bson:"amount"`
}

type Coins []*Coin

type Fee struct {
	Amount []*Coin `bson:"amount"`
	Gas    int64   `bson:"gas"`
}

type Tag map[string]string

type RawLog struct {
	MsgIndex int    `bson:"msg_index"`
	Success  bool   `bson:"success"`
	Log      string `bson:"log"`
}

type Event struct {
	Type       string            `bson:"type"`
	Attributes map[string]string `bson:"attributes" `
}
