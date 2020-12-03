package cron

import (
	"github.com/irisnet/rainbow-sync/conf"
	"github.com/irisnet/rainbow-sync/db"
	"github.com/irisnet/rainbow-sync/lib/pool"
	"testing"
	"time"
)

func TestGetUnknownTxsByPage(t *testing.T) {
	pool.Init(conf.SvrConf.NodeUrls, conf.SvrConf.MaxConnectionNum, conf.SvrConf.InitConnectionNum)
	db.Start()
	defer func() {
		db.Stop()
	}()
	GetUnknownTxsByPage(0, 2)
}

func TestCronService_StartCronService(t *testing.T) {
	pool.Init(conf.SvrConf.NodeUrls, conf.SvrConf.MaxConnectionNum, conf.SvrConf.InitConnectionNum)
	db.Start()
	defer func() {
		db.Stop()
	}()
	new(CronService).StartCronService()
	time.Sleep(1 * time.Minute)
}
