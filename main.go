package main

import (
	"github.com/irisnet/rainbow-sync/db"
	"github.com/irisnet/rainbow-sync/lib/logger"
	"github.com/irisnet/rainbow-sync/lib/pool"
	"github.com/irisnet/rainbow-sync/model"
	"github.com/irisnet/rainbow-sync/task"
	"os"
	"os/signal"
	"syscall"

	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	c := make(chan os.Signal)

	defer func() {
		logger.Info("System Exit")

		db.Stop()
		pool.ClosePool()

		if err := recover(); err != nil {
			logger.Error("", logger.Any("err", err))
			os.Exit(1)
		}
	}()

	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	logger.Info("Start sync Program")

	db.Start()
	model.EnsureDocsIndexes()
	task.Start()

	<-c
}
