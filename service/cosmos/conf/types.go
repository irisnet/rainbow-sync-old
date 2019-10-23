package conf

import (
	"github.com/irisnet/rainbow-sync/service/cosmos/constant"
	"github.com/irisnet/rainbow-sync/service/cosmos/logger"
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

)

// get value of env var
func init() {
	var err error

	nodeUrl, found := os.LookupEnv(constant.EnvNameSerNetworkFullNode_COSMOS)
	if found {
		BlockChainMonitorUrl = strings.Split(nodeUrl, ",")
	}
	logger.Info("Env Value", logger.Any(constant.EnvNameSerNetworkFullNode_COSMOS, BlockChainMonitorUrl))

	workerNumCreateTask, found := os.LookupEnv(constant.EnvNameWorkerNumCreateTask_COSMOS)
	if found {
		WorkerNumCreateTask, err = strconv.Atoi(workerNumCreateTask)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(constant.EnvNameWorkerNumCreateTask_COSMOS, workerNumCreateTask))
		}
	}
	logger.Info("Env Value", logger.Int(constant.EnvNameWorkerNumCreateTask_COSMOS, WorkerNumCreateTask))

	workerNumExecuteTask, found := os.LookupEnv(constant.EnvNameWorkerNumExecuteTask_COSMOS)
	if found {
		WorkerNumExecuteTask, err = strconv.Atoi(workerNumExecuteTask)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(constant.EnvNameWorkerNumExecuteTask_COSMOS, workerNumCreateTask))
		}
	}
	logger.Info("Env Value", logger.Int(constant.EnvNameWorkerNumExecuteTask_COSMOS, WorkerNumExecuteTask))

	workerMaxSleepTime, found := os.LookupEnv(constant.EnvNameWorkerMaxSleepTime_COSMOS)
	if found {
		WorkerMaxSleepTime, err = strconv.Atoi(workerMaxSleepTime)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(constant.EnvNameWorkerMaxSleepTime_COSMOS, workerMaxSleepTime))
		}
	}
	logger.Info("Env Value", logger.Int(constant.EnvNameWorkerMaxSleepTime_COSMOS, WorkerMaxSleepTime))

	blockNumPerWorkerHandle, found := os.LookupEnv(constant.EnvNameBlockNumPerWorkerHandle_COSMOS)
	if found {
		BlockNumPerWorkerHandle, err = strconv.Atoi(blockNumPerWorkerHandle)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(constant.EnvNameBlockNumPerWorkerHandle_COSMOS, blockNumPerWorkerHandle))
		}
	}
	logger.Info("Env Value", logger.Int(constant.EnvNameBlockNumPerWorkerHandle_COSMOS, BlockNumPerWorkerHandle))
}
