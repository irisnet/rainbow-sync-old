package conf

import (
	"github.com/irisnet/rainbow-sync/logger"
	"os"
	"strconv"
	"strings"
)

var (
	BlockChainMonitorUrl = []string{"tcp://192.168.150.31:26657"}

	IrisNetwork             = "testnet"
	WorkerNumCreateTask     = 1
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

	EnvNameSerNetworkFullNodes     = "SER_BC_FULL_NODES"
	EnvNameWorkerNumExecuteTask    = "WORKER_NUM_EXECUTE_TASK"
	EnvNameWorkerMaxSleepTime      = "WORKER_MAX_SLEEP_TIME"
	EnvNameBlockNumPerWorkerHandle = "BLOCK_NUM_PER_WORKER_HANDLE"
	EnvNameIrisNetwork             = "IRIS_NETWORK"
)

// get value of env var
func init() {
	var err error

	nodeUrl, found := os.LookupEnv(EnvNameSerNetworkFullNodes)
	if found {
		BlockChainMonitorUrl = strings.Split(nodeUrl, ",")
	}
	logger.Info("Env Value", logger.Any(EnvNameSerNetworkFullNodes, BlockChainMonitorUrl))


	workerNumExecuteTask, found := os.LookupEnv(EnvNameWorkerNumExecuteTask)
	if found {
		WorkerNumExecuteTask, err = strconv.Atoi(workerNumExecuteTask)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(EnvNameWorkerNumExecuteTask, workerNumExecuteTask))
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
