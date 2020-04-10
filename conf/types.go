package conf

import (
	"github.com/irisnet/rainbow-sync/constant"
	"github.com/irisnet/rainbow-sync/logger"
	"os"
	"strconv"
	"strings"
)

var (
	BlockChainMonitorUrl = []string{"tcp://192.168.150.31:56657"}

	WorkerNumCreateTask     = 1
	WorkerNumExecuteTask    = 30
	WorkerMaxSleepTime      = 2 * 60
	BlockNumPerWorkerHandle = 50

	InitConnectionNum = 50  // fast init num of tendermint client pool
	MaxConnectionNum  = 100 // max size of tendermint client pool

	ZoneName = "cosmos"
)

// get value of env var
func init() {
	var err error

	nodeUrl, found := os.LookupEnv(constant.EnvNameSerNetworkFullNode_ZONE)
	if found {
		BlockChainMonitorUrl = strings.Split(nodeUrl, ",")
	}
	logger.Info("Env Value", logger.Any(constant.EnvNameSerNetworkFullNode_ZONE, BlockChainMonitorUrl))

	//workerNumCreateTask, found := os.LookupEnv(constant.EnvNameWorkerNumCreateTask_ZONE)
	//if found {
	//	WorkerNumCreateTask, err = strconv.Atoi(workerNumCreateTask)
	//	if err != nil {
	//		logger.Fatal("Can't convert str to int", logger.String(constant.EnvNameWorkerNumCreateTask_ZONE, workerNumCreateTask))
	//	}
	//}
	//logger.Info("Env Value", logger.Int(constant.EnvNameWorkerNumCreateTask_ZONE, WorkerNumCreateTask))

	workerNumExecuteTask, found := os.LookupEnv(constant.EnvNameWorkerNumExecuteTask_ZONE)
	if found {
		WorkerNumExecuteTask, err = strconv.Atoi(workerNumExecuteTask)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(constant.EnvNameWorkerNumExecuteTask_ZONE, workerNumExecuteTask))
		}
	}
	logger.Info("Env Value", logger.Int(constant.EnvNameWorkerNumExecuteTask_ZONE, WorkerNumExecuteTask))

	workerMaxSleepTime, found := os.LookupEnv(constant.EnvNameWorkerMaxSleepTime_ZONE)
	if found {
		WorkerMaxSleepTime, err = strconv.Atoi(workerMaxSleepTime)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(constant.EnvNameWorkerMaxSleepTime_ZONE, workerMaxSleepTime))
		}
	}
	logger.Info("Env Value", logger.Int(constant.EnvNameWorkerMaxSleepTime_ZONE, WorkerMaxSleepTime))

	blockNumPerWorkerHandle, found := os.LookupEnv(constant.EnvNameBlockNumPerWorkerHandle_ZONE)
	if found {
		BlockNumPerWorkerHandle, err = strconv.Atoi(blockNumPerWorkerHandle)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(constant.EnvNameBlockNumPerWorkerHandle_ZONE, blockNumPerWorkerHandle))
		}
	}
	logger.Info("Env Value", logger.Int(constant.EnvNameBlockNumPerWorkerHandle_ZONE, BlockNumPerWorkerHandle))

	zoneName, found := os.LookupEnv(constant.EnvNameZoneName)
	if found {
		ZoneName = zoneName
	}
	logger.Info("Env Value", logger.String(constant.EnvNameZoneName, ZoneName))

}
