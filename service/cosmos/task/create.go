package task

import (
	"gopkg.in/mgo.v2/txn"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
	model "github.com/irisnet/rainbow-sync/service/cosmos/db"
	"github.com/irisnet/rainbow-sync/service/cosmos/logger"
	cmodel "github.com/irisnet/rainbow-sync/service/cosmos/model"
	"github.com/irisnet/rainbow-sync/service/cosmos/conf"
	"github.com/irisnet/rainbow-sync/service/cosmos/block"
)

type TaskCosmosService struct {
	blockType       block.Cosmos_Block
	syncCosmosModel cmodel.SyncCosmosTask
}

func (s *TaskCosmosService) StartCreateTask() {
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

func (s *TaskCosmosService) createTask(blockNumPerWorkerHandle int64, chanLimit chan bool) {
	var (
		syncCosmosTasks   []*cmodel.SyncCosmosTask
		ops               []txn.Op
		invalidFollowTask cmodel.SyncCosmosTask
		logMsg            string
	)

	defer func() {
		if err := recover(); err != nil {
			logger.Error("Create  cosmos task failed", logger.Any("err", err),
				logger.String("Chain Block", s.blockType.Name()))
		}
		<-chanLimit
	}()

	// check valid follow task if exist
	// status of valid follow task is unhandled or underway
	validFollowTasks, err := s.syncCosmosModel.QueryAll(
		[]string{
			model.SyncTaskStatusUnHandled,
			model.SyncTaskStatusUnderway,
		}, model.SyncTaskTypeFollow)
	if err != nil {
		logger.Error("Query sync cosmos task failed", logger.String("err", err.Error()),
			logger.String("Chain Block", s.blockType.Name()))
		return
	}
	if len(validFollowTasks) == 0 {
		// get max end_height from sync_task
		maxEndHeight, err := s.syncCosmosModel.GetMaxBlockHeight()
		if err != nil {
			logger.Error("Get Cosmos max endBlock failed", logger.String("err", err.Error()),
				logger.String("Chain Block", s.blockType.Name()))
			return
		}

		blockChainLatestHeight, err := getCosmosBlockChainLatestHeight()
		if err != nil {
			logger.Error("Get Cosmos current block height failed", logger.String("err", err.Error()))
			return
		}

		if maxEndHeight+blockNumPerWorkerHandle <= blockChainLatestHeight {
			syncCosmosTasks = createCosmosCatchUpTask(maxEndHeight, blockNumPerWorkerHandle, blockChainLatestHeight)
			logMsg = fmt.Sprintf("Create cosmos catch up task during follow task not exist,from-to:%v-%v,Chain Block:%v",
				maxEndHeight+1, blockChainLatestHeight, s.blockType.Name())
		} else {
			finished, err := s.assertAllCatchUpCosmosTaskFinished()
			if err != nil {
				logger.Error("AssertAllCatchUpTaskFinished failed", logger.String("err", err.Error()))
				return
			}
			if finished {
				syncCosmosTasks = createFollowCosmosTask(maxEndHeight, blockNumPerWorkerHandle, blockChainLatestHeight)
				logMsg = fmt.Sprintf("Create follow cosmos task during follow task not exist,from-to:%v-%v,Chain Block:%v",
					maxEndHeight+1, blockChainLatestHeight, s.blockType.Name())
			}
		}
	} else {
		followTask := validFollowTasks[0]
		followedHeight := followTask.CurrentHeight
		if followedHeight == 0 {
			followedHeight = followTask.StartHeight - 1
		}

		blockChainLatestHeight, err := getCosmosBlockChainLatestHeight()
		if err != nil {
			logger.Error("Get Cosmos blockChain latest height failed", logger.String("err", err.Error()),
				logger.String("Chain Block", s.blockType.Name()))
			return
		}

		if followedHeight+blockNumPerWorkerHandle <= blockChainLatestHeight {
			syncCosmosTasks = createCosmosCatchUpTask(followedHeight, blockNumPerWorkerHandle, blockChainLatestHeight)

			invalidFollowTask = followTask
			logMsg = fmt.Sprintf("Create cosmos catch up task during follow task exist,from-to:%v-%v,invalidFollowTaskId:%v,invalidFollowTaskCurHeight:%v,Chain Block:%v",
				followedHeight+1, blockChainLatestHeight, invalidFollowTask.ID.Hex(), invalidFollowTask.CurrentHeight, s.blockType.Name())
		}
	}

	// bulk insert or remove use transaction
	ops = make([]txn.Op, 0, len(syncCosmosTasks)+1)
	if len(syncCosmosTasks) > 0 {
		for _, v := range syncCosmosTasks {
			objectId := bson.NewObjectId()
			v.ID = objectId
			op := txn.Op{
				C:      cmodel.CollectionNameSyncCosmosTask,
				Id:     objectId,
				Assert: nil,
				Insert: v,
			}

			ops = append(ops, op)
		}
	}

	if invalidFollowTask.ID.Valid() {
		op := txn.Op{
			C:  cmodel.CollectionNameSyncCosmosTask,
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
			logger.Warn("Create Cosmos sync task fail", logger.String("err", err.Error()),
				logger.String("Chain Block", s.blockType.Name()))
		} else {
			logger.Info(fmt.Sprintf("Create sync Cosmos task success,%v", logMsg), logger.String("Chain Block", s.blockType.Name()))
		}
	}

	time.Sleep(1 * time.Second)
}

func createCosmosCatchUpTask(maxEndHeight, blockNumPerWorker, currentBlockHeight int64) []*cmodel.SyncCosmosTask {
	var (
		syncTasks []*cmodel.SyncCosmosTask
	)
	logger.Info("createCosmosCatchUpTask", logger.Int64("maxEndHeight", maxEndHeight),
		logger.Int64("blockNumPerWorker", blockNumPerWorker), logger.Int64("currentBlockHeight", currentBlockHeight))

	if length := currentBlockHeight - (maxEndHeight + blockNumPerWorker); length > 0 {
		syncTasks = make([]*cmodel.SyncCosmosTask, 0, length+1)
	}

	for maxEndHeight+blockNumPerWorker <= currentBlockHeight {
		syncTask := cmodel.SyncCosmosTask{
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

func (s *TaskCosmosService) assertAllCatchUpCosmosTaskFinished() (bool, error) {
	var (
		allCatchUpTaskFinished = false
	)

	// assert all catch up task whether finished
	tasks, err := s.syncCosmosModel.QueryAll(
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

func createFollowCosmosTask(maxEndHeight, blockNumPerWorker, currentBlockHeight int64) []*cmodel.SyncCosmosTask {
	var (
		syncCosmosTasks []*cmodel.SyncCosmosTask
	)
	syncCosmosTasks = make([]*cmodel.SyncCosmosTask, 0, 1)

	if maxEndHeight+blockNumPerWorker > currentBlockHeight {
		syncTask := cmodel.SyncCosmosTask{
			StartHeight:    maxEndHeight + 1,
			EndHeight:      0,
			Status:         model.SyncTaskStatusUnHandled,
			LastUpdateTime: time.Now().Unix(),
		}

		syncCosmosTasks = append(syncCosmosTasks, &syncTask)
	}

	return syncCosmosTasks
}
