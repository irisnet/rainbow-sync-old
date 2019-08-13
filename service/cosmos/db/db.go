package db

import (
	"fmt"
	conf "github.com/irisnet/rainbow-sync/service/cosmos/conf/db"
	"github.com/irisnet/rainbow-sync/service/cosmos/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"strings"
	"time"
)

var (
	session *mgo.Session
)

func Start() {
	addrs := strings.Split(conf.Addrs, ",")
	dialInfo := &mgo.DialInfo{
		Addrs:     addrs,
		Database:  conf.Database,
		Username:  conf.User,
		Password:  conf.Passwd,
		Direct:    true,
		Timeout:   time.Second * 10,
		PoolLimit: 4096, // Session.SetPoolLimit
	}

	var err error
	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		logger.Fatal("connect db fail", logger.String("err", err.Error()))
	}
	session.SetMode(mgo.Strong, true)
	logger.Info("init db success")
}

func Stop() {
	logger.Info("release resource :mongoDb")
	session.Close()
}

func getSession() *mgo.Session {
	// max session num is 4096
	return session.Clone()
}

// get collection object
func ExecCollection(collectionName string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(conf.Database).C(collectionName)
	return s(c)
}

func Save(h Docs) error {
	save := func(c *mgo.Collection) error {
		pk := h.PkKvPair()
		n, _ := c.Find(pk).Count()
		if n >= 1 {
			return fmt.Errorf("record exist")
		}
		return c.Insert(h)
	}
	return ExecCollection(h.Name(), save)
}

func Update(h Docs) error {
	update := func(c *mgo.Collection) error {
		key := h.PkKvPair()
		return c.Update(key, h)
	}
	return ExecCollection(h.Name(), update)
}

func Delete(h Docs) error {
	remove := func(c *mgo.Collection) error {
		key := h.PkKvPair()
		return c.Remove(key)
	}
	return ExecCollection(h.Name(), remove)
}

//mgo transaction method
//detail to see: https://godoc.org/gopkg.in/mgo.v2/txn
func Txn(ops []txn.Op) error {
	session := getSession()
	defer session.Close()

	c := session.DB(conf.Database).C(CollectionNameTxn)
	runner := txn.NewRunner(c)

	txObjectId := bson.NewObjectId()
	err := runner.Run(ops, txObjectId, nil)
	if err != nil {
		if err == txn.ErrAborted {
			err = runner.Resume(txObjectId)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
