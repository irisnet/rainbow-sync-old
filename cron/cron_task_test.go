package cron

import (
	"testing"
	"time"
	"github.com/irisnet/rainbow-sync/db"
	"github.com/irisnet/rainbow-sync/helper"
	"github.com/irisnet/rainbow-sync/conf"
)

func TestGetUnknownTxsByPage(t *testing.T) {
	helper.Init(conf.BlockChainMonitorUrl, conf.MaxConnectionNum, conf.InitConnectionNum)
	db.Start()
	defer func() {
		db.Stop()
	}()
	GetUnknownTxsByPage(0, 2)
}

func TestCronService_StartCronService(t *testing.T) {
	helper.Init(conf.BlockChainMonitorUrl, conf.MaxConnectionNum, conf.InitConnectionNum)
	db.Start()
	defer func() {
		db.Stop()
	}()
	new(CronService).StartCronService()
	time.Sleep(1 * time.Minute)
}
