package model

import (
	"github.com/irisnet/rainbow-sync/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNameBlock = "sync_iris_block"
)

type (
	Block struct {
		Height     int64 `bson:"height"`
		CreateTime int64 `bson:"create_time"`
	}
)

func (d Block) Name() string {
	return CollectionNameBlock
}

func (d Block) EnsureIndexes() {
	var indexes []mgo.Index
	indexes = append(indexes, mgo.Index{
		Key:        []string{"-height"},
		Unique:     true,
		Background: true,
	})
	db.EnsureIndexes(d.Name(), indexes)
}

func (d Block) PkKvPair() map[string]interface{} {
	return bson.M{"height": d.Height}
}

func (d Block) GetMaxBlockHeight() (Block, error) {
	var result Block

	getMaxBlockHeightFn := func(c *mgo.Collection) error {
		return c.Find(nil).Select(bson.M{"height": 1, "create_time": 1}).Sort("-height").Limit(1).One(&result)
	}

	if err := db.ExecCollection(d.Name(), getMaxBlockHeightFn); err != nil {
		return result, err
	}

	return result, nil
}
