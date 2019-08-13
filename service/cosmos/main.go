package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/irisnet/rainbow-sync/service/cosmos/logger"
	model "github.com/irisnet/rainbow-sync/service/cosmos/db"
	"github.com/irisnet/rainbow-sync/service/cosmos/task"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	c := make(chan os.Signal)

	defer func() {
		logger.Info("System Exit")

		model.Stop()

		if err := recover(); err != nil {
			logger.Error("", logger.Any("err", err))
			os.Exit(1)
		}
	}()

	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	logger.Info("Start sync Program")

	model.Start()
	task.Start()

	<-c
}
