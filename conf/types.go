package conf

import (
	"github.com/irisnet/rainbow-sync/logger"
	"github.com/irisnet/rainbow-sync/utils"
	"os"
	"strconv"
	"strings"
)

var (
	SvrConf              *ServerConf
	blockChainMonitorUrl = []string{"tcp://192.168.150.31:16657"}

	workerNumCreateTask     = 1
	workerNumExecuteTask    = 30
	workerMaxSleepTime      = 2 * 60
	blockNumPerWorkerHandle = 50

	initConnectionNum = 50  // fast init num of tendermint client pool
	maxConnectionNum  = 100 // max size of tendermint client pool
)

type ServerConf struct {
	NodeUrls                []string
	WorkerNumCreateTask     int
	WorkerNumExecuteTask    int
	WorkerMaxSleepTime      int
	BlockNumPerWorkerHandle int

	MaxConnectionNum  int
	InitConnectionNum int
}

const (
	EnvNameDbAddr     = "DB_ADDR"
	EnvNameDbUser     = "DB_USER"
	EnvNameDbPassWd   = "DB_PASSWD"
	EnvNameDbDataBase = "DB_DATABASE"

	EnvNameSerNetworkFullNodes     = "SER_BC_FULL_NODES"
	EnvNameWorkerNumExecuteTask    = "WORKER_NUM_EXECUTE_TASK"
	EnvNameWorkerMaxSleepTime      = "WORKER_MAX_SLEEP_TIME"
	EnvNameBlockNumPerWorkerHandle = "BLOCK_NUM_PER_WORKER_HANDLE"
)

// get value of env var
func init() {
	var err error

	nodeUrl, found := os.LookupEnv(EnvNameSerNetworkFullNodes)
	if found {
		blockChainMonitorUrl = strings.Split(nodeUrl, ",")
	}

	if v, found := os.LookupEnv(EnvNameWorkerNumExecuteTask); found {
		workerNumExecuteTask, err = strconv.Atoi(v)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(EnvNameWorkerNumExecuteTask, v))
		}
	}

	if v, found := os.LookupEnv(EnvNameWorkerMaxSleepTime); found {
		workerMaxSleepTime, err = strconv.Atoi(v)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(EnvNameWorkerMaxSleepTime, v))
		}
	}

	if v, found := os.LookupEnv(EnvNameBlockNumPerWorkerHandle); found {
		blockNumPerWorkerHandle, err = strconv.Atoi(v)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(EnvNameBlockNumPerWorkerHandle, v))
		}
	}
	SvrConf = &ServerConf{
		NodeUrls:                blockChainMonitorUrl,
		WorkerNumCreateTask:     workerNumCreateTask,
		WorkerNumExecuteTask:    workerNumExecuteTask,
		WorkerMaxSleepTime:      workerMaxSleepTime,
		BlockNumPerWorkerHandle: blockNumPerWorkerHandle,

		MaxConnectionNum:  maxConnectionNum,
		InitConnectionNum: initConnectionNum,
	}
	logger.Debug("print server config", logger.String("serverConf", utils.MarshalJsonIgnoreErr(SvrConf)))
}
