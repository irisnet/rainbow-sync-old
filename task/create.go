package task

import (
	"fmt"
	"github.com/irisnet/rainbow-sync/conf"
	model "github.com/irisnet/rainbow-sync/db"
	"github.com/irisnet/rainbow-sync/lib/logger"
	imodel "github.com/irisnet/rainbow-sync/model"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"time"
)

type TaskIrisService struct {
	syncIrisModel imodel.SyncTask
}

const maxRecordNumForBatchInsert = 1000

func (s *TaskIrisService) StartCreateTask() {
	blockNumPerWorkerHandle := int64(conf.SvrConf.BlockNumPerWorkerHandle)

	logger.Info("Start create task")

	// buffer channel to limit goroutine num
	chanLimit := make(chan bool, conf.SvrConf.WorkerNumCreateTask)

	for {
		chanLimit <- true
		go s.createTask(blockNumPerWorkerHandle, chanLimit)
		time.Sleep(time.Duration(1) * time.Minute)
	}
}

func (s *TaskIrisService) createTask(blockNumPerWorkerHandle int64, chanLimit chan bool) {
	var (
		syncIrisTasks     []*imodel.SyncTask
		ops               []txn.Op
		invalidFollowTask imodel.SyncTask
		logMsg            string
	)

	defer func() {
		if err := recover(); err != nil {
			logger.Error("Create task failed", logger.Any("err", err))
		}
		<-chanLimit
	}()
	// check valid follow task if exist
	// status of valid follow task is unhandled or underway
	validFollowTasks, err := s.syncIrisModel.QueryAll(
		[]string{
			model.SyncTaskStatusUnHandled,
			model.SyncTaskStatusUnderway,
		}, model.SyncTaskTypeFollow)
	if err != nil {
		logger.Error("Query sync task failed", logger.String("err", err.Error()))
		return
	}
	if len(validFollowTasks) == 0 {
		// get max end_height from sync_task
		maxEndHeight, err := s.syncIrisModel.GetMaxBlockHeight()
		if err != nil {
			logger.Error("Get max endBlock failed", logger.String("err", err.Error()))
			return
		}

		blockChainLatestHeight, err := getBlockChainLatestHeight()
		if err != nil {
			logger.Error("Get current block height failed", logger.String("err", err.Error()))
			return
		}

		if maxEndHeight+blockNumPerWorkerHandle <= blockChainLatestHeight {
			syncIrisTasks = createCatchUpTask(maxEndHeight, blockNumPerWorkerHandle, blockChainLatestHeight)
			logMsg = fmt.Sprintf("Create  catch up task during follow task not exist,from-to:%v-%v",
				maxEndHeight+1, blockChainLatestHeight)
		} else {
			finished, err := s.assertAllCatchUpTaskFinished()
			if err != nil {
				logger.Error("AssertAllCatchUpTaskFinished failed", logger.String("err", err.Error()))
				return
			}
			if finished {
				syncIrisTasks = createFollowTask(maxEndHeight, blockNumPerWorkerHandle, blockChainLatestHeight)
				logMsg = fmt.Sprintf("Create follow task during follow task not exist,from-to:%v-%v",
					maxEndHeight+1, blockChainLatestHeight)
			}
		}
	} else {
		followTask := validFollowTasks[0]
		followedHeight := followTask.CurrentHeight
		if followedHeight == 0 {
			followedHeight = followTask.StartHeight - 1
		}

		blockChainLatestHeight, err := getBlockChainLatestHeight()
		if err != nil {
			logger.Error("Get blockChain latest height failed", logger.String("err", err.Error()))
			return
		}

		if followedHeight+blockNumPerWorkerHandle <= blockChainLatestHeight {
			syncIrisTasks = createCatchUpTask(followedHeight, blockNumPerWorkerHandle, blockChainLatestHeight)

			invalidFollowTask = followTask
			logMsg = fmt.Sprintf("Create catch up task during follow task exist,from-to:%v-%v,invalidFollowTaskId:%v,invalidFollowTaskCurHeight:%v",
				followedHeight+1, blockChainLatestHeight, invalidFollowTask.ID.Hex(), invalidFollowTask.CurrentHeight)

		}
	}

	// bulk insert or remove use transaction
	ops = make([]txn.Op, 0, len(syncIrisTasks)+1)
	if len(syncIrisTasks) > 0 {
		for _, v := range syncIrisTasks {
			objectId := bson.NewObjectId()
			v.ID = objectId
			op := txn.Op{
				C:      imodel.CollectionNameSyncTask,
				Id:     objectId,
				Assert: nil,
				Insert: v,
			}

			ops = append(ops, op)
		}
	}

	if invalidFollowTask.ID.Valid() {
		op := txn.Op{
			C:  imodel.CollectionNameSyncTask,
			Id: invalidFollowTask.ID,
			Assert: bson.M{
				"current_height":   invalidFollowTask.CurrentHeight,
				"last_update_time": invalidFollowTask.LastUpdateTime,
			},
			Update: bson.M{
				"$set": bson.M{
					"status":           model.FollowTaskStatusInvalid,
					"last_update_time": time.Now().Unix(),
				},
			},
		}
		ops = append(ops, op)
	}

	if len(ops) > 0 {
		err := model.Txn(ops)
		if err != nil {
			logger.Warn("Create sync task fail", logger.String("err", err.Error()))
		} else {
			logger.Info(fmt.Sprintf("Create sync task success,%v", logMsg))
		}
	}

	time.Sleep(1 * time.Second)
}

func createCatchUpTask(maxEndHeight, blockNumPerWorker, currentBlockHeight int64) []*imodel.SyncTask {
	var (
		syncTasks []*imodel.SyncTask
	)
	if length := currentBlockHeight - (maxEndHeight + blockNumPerWorker); length > 0 {
		syncTasks = make([]*imodel.SyncTask, 0, length+1)
	}

	for maxEndHeight+blockNumPerWorker <= currentBlockHeight {
		if len(syncTasks) >= maxRecordNumForBatchInsert {
			break
		}
		syncTask := imodel.SyncTask{
			StartHeight:    maxEndHeight + 1,
			EndHeight:      maxEndHeight + blockNumPerWorker,
			Status:         model.SyncTaskStatusUnHandled,
			LastUpdateTime: time.Now().Unix(),
		}
		syncTasks = append(syncTasks, &syncTask)

		maxEndHeight += blockNumPerWorker
	}

	return syncTasks
}

func (s *TaskIrisService) assertAllCatchUpTaskFinished() (bool, error) {
	var (
		allCatchUpTaskFinished = false
	)

	// assert all catch up task whether finished
	tasks, err := s.syncIrisModel.QueryAll(
		[]string{
			model.SyncTaskStatusUnHandled,
			model.SyncTaskStatusUnderway,
		},
		model.SyncTaskTypeCatchUp)
	if err != nil {
		return false, err
	}

	if len(tasks) == 0 {
		allCatchUpTaskFinished = true
	}

	return allCatchUpTaskFinished, nil
}

func createFollowTask(maxEndHeight, blockNumPerWorker, currentBlockHeight int64) []*imodel.SyncTask {
	var (
		syncIrisTasks []*imodel.SyncTask
	)
	syncIrisTasks = make([]*imodel.SyncTask, 0, 1)

	if maxEndHeight+blockNumPerWorker > currentBlockHeight {
		syncTask := imodel.SyncTask{
			StartHeight:    maxEndHeight + 1,
			EndHeight:      0,
			Status:         model.SyncTaskStatusUnHandled,
			LastUpdateTime: time.Now().Unix(),
		}

		syncIrisTasks = append(syncIrisTasks, &syncTask)
	}

	return syncIrisTasks
}
