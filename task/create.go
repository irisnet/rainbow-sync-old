package task

import (
	"fmt"
	"github.com/irisnet/rainbow-sync/block"
	"github.com/irisnet/rainbow-sync/conf"
	model "github.com/irisnet/rainbow-sync/db"
	"github.com/irisnet/rainbow-sync/logger"
	cmodel "github.com/irisnet/rainbow-sync/model"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"time"
)

type TaskZoneService struct {
	blockType block.ZoneBlock
	syncModel cmodel.SyncZoneTask
}

func (s *TaskZoneService) StartCreateTask() {
	blockNumPerWorkerHandle := int64(conf.BlockNumPerWorkerHandle)

	logger.Info("Start create task", logger.String("Chain Block", s.blockType.Name()))

	// buffer channel to limit goroutine num
	chanLimit := make(chan bool, conf.WorkerNumCreateTask)

	for {
		chanLimit <- true
		go s.createTask(blockNumPerWorkerHandle, chanLimit)
		time.Sleep(time.Duration(1) * time.Minute)
	}
}

func (s *TaskZoneService) createTask(blockNumPerWorkerHandle int64, chanLimit chan bool) {
	var (
		syncZoneTasks     []*cmodel.SyncZoneTask
		ops               []txn.Op
		invalidFollowTask cmodel.SyncZoneTask
		logMsg            string
	)

	defer func() {
		if err := recover(); err != nil {
			logger.Error("Create  task failed", logger.Any("err", err),
				logger.String("Chain Block", s.blockType.Name()))
		}
		<-chanLimit
	}()

	// check valid follow task if exist
	// status of valid follow task is unhandled or underway
	validFollowTasks, err := s.syncModel.QueryAll(
		[]string{
			model.SyncTaskStatusUnHandled,
			model.SyncTaskStatusUnderway,
		}, model.SyncTaskTypeFollow)
	if err != nil {
		logger.Error("Query sync  task failed", logger.String("err", err.Error()),
			logger.String("Chain Block", s.blockType.Name()))
		return
	}
	if len(validFollowTasks) == 0 {
		// get max end_height from sync_task
		maxEndHeight, err := s.syncModel.GetMaxBlockHeight()
		if err != nil {
			logger.Error("Get  max endBlock failed", logger.String("err", err.Error()),
				logger.String("Chain Block", s.blockType.Name()))
			return
		}

		blockChainLatestHeight, err := getBlockChainLatestHeight()
		if err != nil {
			logger.Error("Get  current block height failed", logger.String("err", err.Error()))
			return
		}

		if maxEndHeight+blockNumPerWorkerHandle <= blockChainLatestHeight {
			syncZoneTasks = createCatchUpTask(maxEndHeight, blockNumPerWorkerHandle, blockChainLatestHeight)
			logMsg = fmt.Sprintf("Create  catch up task during follow task not exist,from-to:%v-%v,Chain Block:%v",
				maxEndHeight+1, blockChainLatestHeight, s.blockType.Name())
		} else {
			finished, err := s.assertAllCatchUpTaskFinished()
			if err != nil {
				logger.Error("AssertAllCatchUpTaskFinished failed", logger.String("err", err.Error()))
				return
			}
			if finished {
				syncZoneTasks = createFollowTask(maxEndHeight, blockNumPerWorkerHandle, blockChainLatestHeight)
				logMsg = fmt.Sprintf("Create follow  task during follow task not exist,from-to:%v-%v,Chain Block:%v",
					maxEndHeight+1, blockChainLatestHeight, s.blockType.Name())
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
			logger.Error("Get  blockChain latest height failed", logger.String("err", err.Error()),
				logger.String("Chain Block", s.blockType.Name()))
			return
		}

		if followedHeight+blockNumPerWorkerHandle <= blockChainLatestHeight {
			syncZoneTasks = createCatchUpTask(followedHeight, blockNumPerWorkerHandle, blockChainLatestHeight)

			invalidFollowTask = followTask
			logMsg = fmt.Sprintf("Create  catch up task during follow task exist,from-to:%v-%v,invalidFollowTaskId:%v,invalidFollowTaskCurHeight:%v,Chain Block:%v",
				followedHeight+1, blockChainLatestHeight, invalidFollowTask.ID.Hex(), invalidFollowTask.CurrentHeight, s.blockType.Name())
		}
	}

	// bulk insert or remove use transaction
	ops = make([]txn.Op, 0, len(syncZoneTasks)+1)
	if len(syncZoneTasks) > 0 {
		for _, v := range syncZoneTasks {
			objectId := bson.NewObjectId()
			v.ID = objectId
			op := txn.Op{
				C:      s.syncModel.Name(),
				Id:     objectId,
				Assert: nil,
				Insert: v,
			}

			ops = append(ops, op)
		}
	}

	if invalidFollowTask.ID.Valid() {
		op := txn.Op{
			C:  s.syncModel.Name(),
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
			logger.Warn("Create  sync task fail", logger.String("err", err.Error()),
				logger.String("Chain Block", s.blockType.Name()))
		} else {
			logger.Info(fmt.Sprintf("Create sync  task success,%v", logMsg), logger.String("Chain Block", s.blockType.Name()))
		}
	}

	time.Sleep(1 * time.Second)
}

func createCatchUpTask(maxEndHeight, blockNumPerWorker, currentBlockHeight int64) []*cmodel.SyncZoneTask {
	var (
		syncTasks []*cmodel.SyncZoneTask
	)
	logger.Info("create CatchUpTask", logger.Int64("maxEndHeight", maxEndHeight),
		logger.Int64("blockNumPerWorker", blockNumPerWorker), logger.Int64("currentBlockHeight", currentBlockHeight))

	if length := currentBlockHeight - (maxEndHeight + blockNumPerWorker); length > 0 {
		syncTasks = make([]*cmodel.SyncZoneTask, 0, length+1)
	}

	for maxEndHeight+blockNumPerWorker <= currentBlockHeight {
		syncTask := cmodel.SyncZoneTask{
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

func (s *TaskZoneService) assertAllCatchUpTaskFinished() (bool, error) {
	var (
		allCatchUpTaskFinished = false
	)

	// assert all catch up task whether finished
	tasks, err := s.syncModel.QueryAll(
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

func createFollowTask(maxEndHeight, blockNumPerWorker, currentBlockHeight int64) []*cmodel.SyncZoneTask {
	var (
		syncZoneTasks []*cmodel.SyncZoneTask
	)
	syncZoneTasks = make([]*cmodel.SyncZoneTask, 0, 1)

	if maxEndHeight+blockNumPerWorker > currentBlockHeight {
		syncTask := cmodel.SyncZoneTask{
			StartHeight:    maxEndHeight + 1,
			EndHeight:      0,
			Status:         model.SyncTaskStatusUnHandled,
			LastUpdateTime: time.Now().Unix(),
		}

		syncZoneTasks = append(syncZoneTasks, &syncTask)
	}

	return syncZoneTasks
}
