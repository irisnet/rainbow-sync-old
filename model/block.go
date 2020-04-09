package model

import (
	"github.com/irisnet/rainbow-sync/conf"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"github.com/irisnet/rainbow-sync/db"
)

const (
	CollectionNameBlock = "sync_%v_block"
)

type (
	Block struct {
		Id         bson.ObjectId `bson:"_id"`
		Height     int64         `bson:"height"`
		CreateTime int64         `bson:"create_time"`
	}
)

func (b *Block) zoneName() string {
	return conf.ZoneName
}

func (b *Block) Name() string {
	return fmt.Sprintf(CollectionNameBlock, b.zoneName())
}

func (b *Block) PkKvPair() map[string]interface{} {
	return bson.M{"height": b.Height}
}

func (b *Block) EnsureIndexes() {
	var indexes []mgo.Index
	indexes = append(indexes,
		mgo.Index{
			Key:        []string{"height"},
			Unique:     true,
			Background: true,
		})
	db.EnsureIndexes(b.Name(), indexes)
}
