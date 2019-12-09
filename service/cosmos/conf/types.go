package conf

import (
	"os"
	"strings"
	"strconv"
	"github.com/irisnet/rainbow-sync/service/cosmos/logger"
)

var (
	BlockChainMonitorUrl = []string{"tcp://10.2.10.140:36657"}

	WorkerNumCreateTask     = 2
	WorkerNumExecuteTask    = 30
	WorkerMaxSleepTime      = 2 * 60
	BlockNumPerWorkerHandle = 50

	InitConnectionNum = 50  // fast init num of tendermint client pool
	MaxConnectionNum  = 100 // max size of tendermint client pool

)

const (
	EnvNameSerNetworkFullNode_COSMOS      = "SER_BC_FULL_NODE_COSMOS"
	EnvNameWorkerNumCreateTask_COSMOS     = "WORKER_NUM_CREATE_TASK_COSMOS"
	EnvNameWorkerNumExecuteTask_COSMOS    = "WORKER_NUM_EXECUTE_TASK_COSMOS"
	EnvNameWorkerMaxSleepTime_COSMOS      = "WORKER_MAX_SLEEP_TIME_COSMOS"
	EnvNameBlockNumPerWorkerHandle_COSMOS = "BLOCK_NUM_PER_WORKER_HANDLE_COSMOS"

	EnvNameDbAddr     = "DB_ADDR"
	EnvNameDbUser     = "DB_USER"
	EnvNameDbPassWd   = "DB_PASSWD"
	EnvNameDbDataBase = "DB_DATABASE"
)

// get value of env var
func init() {
	var err error

	nodeUrl, found := os.LookupEnv(EnvNameSerNetworkFullNode_COSMOS)
	if found {
		BlockChainMonitorUrl = strings.Split(nodeUrl, ",")
	}
	logger.Info("Env Value", logger.Any(EnvNameSerNetworkFullNode_COSMOS, BlockChainMonitorUrl))

	workerNumCreateTask, found := os.LookupEnv(EnvNameWorkerNumCreateTask_COSMOS)
	if found {
		WorkerNumCreateTask, err = strconv.Atoi(workerNumCreateTask)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(EnvNameWorkerNumCreateTask_COSMOS, workerNumCreateTask))
		}
	}
	logger.Info("Env Value", logger.Int(EnvNameWorkerNumCreateTask_COSMOS, WorkerNumCreateTask))

	workerNumExecuteTask, found := os.LookupEnv(EnvNameWorkerNumExecuteTask_COSMOS)
	if found {
		WorkerNumExecuteTask, err = strconv.Atoi(workerNumExecuteTask)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(EnvNameWorkerNumExecuteTask_COSMOS, workerNumCreateTask))
		}
	}
	logger.Info("Env Value", logger.Int(EnvNameWorkerNumExecuteTask_COSMOS, WorkerNumExecuteTask))

	workerMaxSleepTime, found := os.LookupEnv(EnvNameWorkerMaxSleepTime_COSMOS)
	if found {
		WorkerMaxSleepTime, err = strconv.Atoi(workerMaxSleepTime)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(EnvNameWorkerMaxSleepTime_COSMOS, workerMaxSleepTime))
		}
	}
	logger.Info("Env Value", logger.Int(EnvNameWorkerMaxSleepTime_COSMOS, WorkerMaxSleepTime))

	blockNumPerWorkerHandle, found := os.LookupEnv(EnvNameBlockNumPerWorkerHandle_COSMOS)
	if found {
		BlockNumPerWorkerHandle, err = strconv.Atoi(blockNumPerWorkerHandle)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(EnvNameBlockNumPerWorkerHandle_COSMOS, blockNumPerWorkerHandle))
		}
	}
	logger.Info("Env Value", logger.Int(EnvNameBlockNumPerWorkerHandle_COSMOS, BlockNumPerWorkerHandle))
}
