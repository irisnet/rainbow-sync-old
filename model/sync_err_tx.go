package model

import (
	"github.com/irisnet/rainbow-sync/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ErrTx struct {
	Height int64  `bson:"height"`
	TxHash string `bson:"tx_hash"`
	Repair int    `bson:"repair"`
	Log    string `bson:"log"`
}

const (
	CollectionNameErrTx = "sync_err_tx"
)

func (d ErrTx) Name() string {
	return CollectionNameErrTx
}

func (d ErrTx) PkKvPair() map[string]interface{} {
	return bson.M{"height": d.Height, "tx_hash": d.TxHash}
}

func (d ErrTx) EnsureIndexes() {
	var indexes []mgo.Index
	indexes = append(indexes,
		mgo.Index{
			Key:        []string{"-height", "-tx_hash"},
			Unique:     true,
			Background: true},
	)

	db.EnsureIndexes(d.Name(), indexes)
}

func (d ErrTx) Save() error {
	return db.Save(d)
}

func (d ErrTx) Find(skip, limit int) ([]ErrTx, error) {
	var res []ErrTx
	q := bson.M{
		"repair": 0,
	}
	sorts := []string{"-height"}
	fn := func(c *mgo.Collection) error {
		return c.Find(q).Sort(sorts...).Skip(skip).Limit(limit).All(&res)
	}

	err := db.ExecCollection(d.Name(), fn)
	return res, err
}

func (d ErrTx) Update(one ErrTx) error {
	return db.Update(&one)
}
