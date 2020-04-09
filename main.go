package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/irisnet/rainbow-sync/logger"
	"github.com/irisnet/rainbow-sync/db"
	"github.com/irisnet/rainbow-sync/task"
	"runtime"
	"github.com/irisnet/rainbow-sync/model"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	c := make(chan os.Signal)

	defer func() {
		logger.Info("System Exit")

		db.Stop()

		if err := recover(); err != nil {
			logger.Error("", logger.Any("err", err))
			os.Exit(1)
		}
	}()

	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	logger.Info("Start sync Program")

	db.Start()
	model.CheckIndex()
	task.Start()

	<-c
}
