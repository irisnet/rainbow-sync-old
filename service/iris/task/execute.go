package task

import (
	"fmt"
	"github.com/irisnet/rainbow-sync/service/iris/logger"
	imodel "github.com/irisnet/rainbow-sync/service/iris/model"
	"github.com/irisnet/rainbow-sync/service/iris/conf"
	model "github.com/irisnet/rainbow-sync/service/iris/db"
	"github.com/irisnet/rainbow-sync/service/iris/helper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"time"
)

func (s *TaskIrisService) StartExecuteTask() {
	var (
		blockNumPerWorkerHandle = int64(conf.BlockNumPerWorkerHandle)
		workerMaxSleepTime      = int64(conf.WorkerMaxSleepTime)
	)
	if workerMaxSleepTime <= 1*60 {
		logger.Fatal("workerMaxSleepTime shouldn't less than 1 minute")
	}

	logger.Info("Start execute task", logger.String("Chain Block", s.blockType.Name()))

	// buffer channel to limit goroutine num
	chanLimit := make(chan bool, conf.WorkerNumExecuteTask)
	helper.Init(conf.BlockChainMonitorUrl, conf.MaxConnectionNum, conf.InitConnectionNum)
	defer func() {
		helper.ClosePool()
	}()

	for {
		chanLimit <- true
		go s.executeTask(blockNumPerWorkerHandle, workerMaxSleepTime, chanLimit)
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func (s *TaskIrisService) executeTask(blockNumPerWorkerHandle, maxWorkerSleepTime int64, chanLimit chan bool) {
	var (
		workerId, taskType     string
		blockChainLatestHeight int64
	)
	genWorkerId := func() string {
		// generate worker id use hostname@xxx
		hostname, _ := os.Hostname()
		return fmt.Sprintf("%v@%v", hostname, bson.NewObjectId().Hex())
	}

	healthCheckQuit := make(chan bool)
	workerId = genWorkerId()
	client := helper.GetClient()

	defer func() {
		if r := recover(); r != nil {
			logger.Error("execute task fail", logger.Any("err", r))
		}
		close(healthCheckQuit)
		<-chanLimit
		client.Release()
	}()

	// check whether exist executable task
	// status = unhandled or
	// status = underway and now - lastUpdateTime > confTime
	tasks, err := s.syncIrisModel.GetExecutableTask(maxWorkerSleepTime)
	if err != nil {
		logger.Error("Get executable task fail", logger.String("err", err.Error()))
	}
	if len(tasks) == 0 {
		// there is no executable tasks
		return
	}

	// take over sync task
	// attempt to update status, worker_id and worker_logs
	task := tasks[0]
	err = s.syncIrisModel.TakeOverTask(task, workerId)
	if err != nil {
		if err == mgo.ErrNotFound {
			// this task has been take over by other goroutine
			logger.Info("Task has been take over by other goroutine", logger.String("Chain Block", s.blockType.Name()))
		} else {
			logger.Error("Take over task fail", logger.String("Chain Block", s.blockType.Name()), logger.String("err", err.Error()))
		}
		return
	} else {
		// task over task success, update task worker to current worker
		task.WorkerId = workerId
	}

	if task.EndHeight != 0 {
		taskType = model.SyncTaskTypeCatchUp
	} else {
		taskType = model.SyncTaskTypeFollow
	}
	logger.Info("worker begin execute task", logger.String("Chain Block", s.blockType.Name()),
		logger.String("curWorker", workerId), logger.Any("taskId", task.ID),
		logger.String("from-to", fmt.Sprintf("%v-%v", task.StartHeight, task.EndHeight)))

	// worker health check, if worker is alive, then update last update time every minute.
	// health check will exit in follow conditions:
	// 1. task is not owned by current worker
	// 2. task is invalid
	workerHealthCheck := func(taskId bson.ObjectId, currentWorker string) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("worker health check err", logger.String("Chain Block", s.blockType.Name()), logger.Any("err", r))
			}
		}()

		func() {
			for {
				select {
				case <-healthCheckQuit:
					logger.Info("get health check quit signal, now exit health check", logger.String("Chain Block", s.blockType.Name()))
					return
				default:
					task, err := s.syncIrisModel.GetTaskByIdAndWorker(taskId, workerId)
					if err == nil {
						if _, valid := assertTaskValid(task, blockNumPerWorkerHandle); valid {
							// update task last update time
							if err := s.syncIrisModel.UpdateLastUpdateTime(task); err != nil {
								logger.Error("update last update time fail", logger.String("Chain Block", s.blockType.Name()), logger.String("err", err.Error()))
							}
						} else {
							logger.Info("task is invalid, exit health check", logger.String("Chain Block", s.blockType.Name()), logger.String("taskId", taskId.Hex()))
							return
						}
					} else {
						if err == mgo.ErrNotFound {
							logger.Info("task may be task over by other goroutine, exit health check", logger.String("Chain Block", s.blockType.Name()),
								logger.String("taskId", taskId.Hex()), logger.String("curWorker", workerId))
							return
						} else {
							logger.Error("get task by id and worker fail", logger.String("Chain Block", s.blockType.Name()), logger.String("taskId", taskId.Hex()),
								logger.String("curWorker", workerId))
						}
					}
				}
				time.Sleep(1 * time.Minute)
			}
		}()
	}
	go workerHealthCheck(task.ID, workerId)

	// check task is valid
	// valid catch up task: current_height < end_height
	// valid follow task: current_height + blockNumPerWorkerHandle > blockChainLatestHeight
	blockChainLatestHeight, isValid := assertTaskValid(task, blockNumPerWorkerHandle)
	for isValid {
		var inProcessBlock int64
		if task.CurrentHeight == 0 {
			inProcessBlock = task.StartHeight
		} else {
			inProcessBlock = task.CurrentHeight + 1
		}

		// if inProcessBlock > blockChainLatestHeight, should wait blockChainLatestHeight update
		if taskType == model.SyncTaskTypeFollow && inProcessBlock > blockChainLatestHeight {
			logger.Info("wait blockChain latest height update", logger.String("Chain Block", s.blockType.Name()),
				logger.Int64("curSyncedHeight", inProcessBlock-1),
				logger.Int64("blockChainLatestHeight", blockChainLatestHeight))
			time.Sleep(2 * time.Second)
			// continue to assert task is valid
			blockChainLatestHeight, isValid = assertTaskValid(task, blockNumPerWorkerHandle)
			continue
		}

		// parse data from block
		blockDoc, assetDetailDocs, txDocs, err := s.blockType.ParseBlock(inProcessBlock, client)
		if err != nil {
			logger.Error("Parse block fail", logger.String("Chain Block", s.blockType.Name()), logger.Int64("block", inProcessBlock),
				logger.String("err", err.Error()))
		}

		// check task owner
		workerUnchanged, err := assertTaskWorkerUnchanged(task.ID, task.WorkerId)
		if err != nil {
			logger.Error("assert task worker is unchanged fail", logger.String("Chain Block", s.blockType.Name()), logger.String("err", err.Error()))
		}
		if workerUnchanged {
			// save data and update sync task
			taskDoc := task
			taskDoc.CurrentHeight = inProcessBlock
			taskDoc.LastUpdateTime = time.Now().Unix()
			taskDoc.Status = model.SyncTaskStatusUnderway
			if inProcessBlock == task.EndHeight {
				taskDoc.Status = model.SyncTaskStatusCompleted
			}

			err := s.blockType.SaveDocsWithTxn(blockDoc, assetDetailDocs, txDocs, taskDoc)
			if err != nil {
				logger.Error("save docs fail", logger.String("Chain Block", s.blockType.Name()), logger.String("err", err.Error()))
			} else {
				task.CurrentHeight = inProcessBlock
			}

			// continue to assert task is valid
			blockChainLatestHeight, isValid = assertTaskValid(task, blockNumPerWorkerHandle)
		} else {
			logger.Info("task worker changed", logger.String("Chain Block", s.blockType.Name()), logger.Any("task_id", task.ID),
				logger.String("origin worker", workerId), logger.String("current worker", task.WorkerId))
			return
		}
	}

	logger.Info("worker finish execute task", logger.String("Chain Block", s.blockType.Name()),
		logger.String("task_worker", task.WorkerId), logger.Any("task_id", task.ID),
		logger.String("from-to-current", fmt.Sprintf("%v-%v-%v", task.StartHeight, task.EndHeight, task.CurrentHeight)))
}

// assert task is valid
// valid catch up task: current_height < end_height
// valid follow task: current_height + blockNumPerWorkerHandle > blockChainLatestHeight
func assertTaskValid(task imodel.SyncTask, blockNumPerWorkerHandle int64) (int64, bool) {
	var (
		taskType               string
		flag                   = false
		blockChainLatestHeight int64
		err                    error
	)
	if task.EndHeight != 0 {
		taskType = model.SyncTaskTypeCatchUp
	} else {
		taskType = model.SyncTaskTypeFollow
	}
	currentHeight := task.CurrentHeight
	if currentHeight == 0 {
		currentHeight = task.StartHeight - 1
	}

	switch taskType {
	case model.SyncTaskTypeCatchUp:
		if currentHeight < task.EndHeight {
			flag = true
		}
		break
	case model.SyncTaskTypeFollow:
		blockChainLatestHeight, err = getBlockChainLatestHeight()
		if err != nil {
			logger.Error("get blockChain latest height err", logger.String("err", err.Error()))
			return blockChainLatestHeight, flag
		}
		if currentHeight+blockNumPerWorkerHandle > blockChainLatestHeight {
			flag = true
		}
		break
	}
	return blockChainLatestHeight, flag
}

// assert task worker unchanged
func assertTaskWorkerUnchanged(taskId bson.ObjectId, workerId string) (bool, error) {
	var (
		syncTaskModel imodel.SyncTask
	)
	// check task owner
	task, err := syncTaskModel.GetTaskById(taskId)
	if err != nil {
		return false, err
	}

	if task.WorkerId == workerId {
		return true, nil
	} else {
		return false, nil
	}
}

// get current block height
func getBlockChainLatestHeight() (int64, error) {
	client := helper.GetClient()
	defer func() {
		client.Release()
	}()
	status, err := client.Status()
	if err != nil {
		return 0, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}
