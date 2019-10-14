package conf

import (
	"github.com/irisnet/rainbow-sync/service/iris/logger"
	"os"
	"strconv"
	"strings"
)

var (
	BlockChainMonitorUrl = []string{"tcp://192.168.150.31:46657"}

	IrisNetwork             = "testnet"
	WorkerNumCreateTask     = 2
	WorkerNumExecuteTask    = 30
	WorkerMaxSleepTime      = 2 * 60
	BlockNumPerWorkerHandle = 50

	InitConnectionNum = 50  // fast init num of tendermint client pool
	MaxConnectionNum  = 100 // max size of tendermint client pool
)

const (
	EnvNameDbAddr     = "DB_ADDR"
	EnvNameDbUser     = "DB_USER"
	EnvNameDbPassWd   = "DB_PASSWD"
	EnvNameDbDataBase = "DB_DATABASE"

	EnvNameSerNetworkFullNode      = "SER_BC_FULL_NODE"
	EnvNameWorkerNumCreateTask     = "WORKER_NUM_CREATE_TASK"
	EnvNameWorkerNumExecuteTask    = "WORKER_NUM_EXECUTE_TASK"
	EnvNameWorkerMaxSleepTime      = "WORKER_MAX_SLEEP_TIME"
	EnvNameBlockNumPerWorkerHandle = "BLOCK_NUM_PER_WORKER_HANDLE"
	EnvNameIrisNetwork             = "IRIS_NETWORK"
)

// get value of env var
func init() {
	var err error

	nodeUrl, found := os.LookupEnv(EnvNameSerNetworkFullNode)
	if found {
		BlockChainMonitorUrl = strings.Split(nodeUrl, ",")
	}
	logger.Info("Env Value", logger.Any(EnvNameSerNetworkFullNode, BlockChainMonitorUrl))

	workerNumCreateTask, found := os.LookupEnv(EnvNameWorkerNumCreateTask)
	if found {
		WorkerNumCreateTask, err = strconv.Atoi(workerNumCreateTask)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(EnvNameWorkerNumCreateTask, workerNumCreateTask))
		}
	}
	logger.Info("Env Value", logger.Int(EnvNameWorkerNumCreateTask, WorkerNumCreateTask))

	workerNumExecuteTask, found := os.LookupEnv(EnvNameWorkerNumExecuteTask)
	if found {
		WorkerNumExecuteTask, err = strconv.Atoi(workerNumExecuteTask)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(EnvNameWorkerNumExecuteTask, workerNumCreateTask))
		}
	}
	logger.Info("Env Value", logger.Int(EnvNameWorkerNumExecuteTask, WorkerNumExecuteTask))

	workerMaxSleepTime, found := os.LookupEnv(EnvNameWorkerMaxSleepTime)
	if found {
		WorkerMaxSleepTime, err = strconv.Atoi(workerMaxSleepTime)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(EnvNameWorkerMaxSleepTime, workerMaxSleepTime))
		}
	}
	logger.Info("Env Value", logger.Int(EnvNameWorkerMaxSleepTime, WorkerMaxSleepTime))

	blockNumPerWorkerHandle, found := os.LookupEnv(EnvNameBlockNumPerWorkerHandle)
	if found {
		BlockNumPerWorkerHandle, err = strconv.Atoi(blockNumPerWorkerHandle)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(EnvNameBlockNumPerWorkerHandle, blockNumPerWorkerHandle))
		}
	}
	logger.Info("Env Value", logger.Int(EnvNameBlockNumPerWorkerHandle, BlockNumPerWorkerHandle))
	network, found := os.LookupEnv(EnvNameIrisNetwork)
	if found {
		IrisNetwork = network
	} else {
		panic("not found " + EnvNameIrisNetwork)
	}
	logger.Info("Env Value", logger.String(EnvNameIrisNetwork, IrisNetwork))
}
