package cosmos

import (
	"github.com/irisnet/rainbow-sync/service/cosmos/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNameSyncCosmosTask = "sync_cosmos_task"
)

type (
	SyncCosmosTask struct {
		ID             bson.ObjectId `bson:"_id"`
		StartHeight    int64         `bson:"start_height"`     // task start height
		EndHeight      int64         `bson:"end_height"`       // task end height
		CurrentHeight  int64         `bson:"current_height"`   // task current height
		Status         string        `bson:"status"`           // task status
		WorkerId       string        `bson:"worker_id"`        // worker id
		WorkerLogs     []WorkerLog   `bson:"worker_logs"`      // worker logs
		LastUpdateTime int64         `bson:"last_update_time"` // unix timestamp
	}

	WorkerLog struct {
		WorkerId  string    `bson:"worker_id"`  // worker id
		BeginTime time.Time `bson:"begin_time"` // time which worker begin handle this task
	}
)

func (d *SyncCosmosTask) Name() string {
	return CollectionNameSyncCosmosTask
}

func (d *SyncCosmosTask) PkKvPair() map[string]interface{} {
	return bson.M{"start_height": d.CurrentHeight, "end_height": d.EndHeight}
}

// get max block height in sync task
func (d *SyncCosmosTask) GetMaxBlockHeight() (int64, error) {
	type maxHeightRes struct {
		MaxHeight int64 `bson:"max"`
	}
	var res []maxHeightRes

	q := []bson.M{
		{
			"$group": bson.M{
				"_id": nil,
				"max": bson.M{"$max": "$end_height"},
			},
		},
	}

	getMaxBlockHeightFn := func(c *mgo.Collection) error {
		return c.Pipe(q).All(&res)
	}
	err := db.ExecCollection(d.Name(), getMaxBlockHeightFn)

	if err != nil {
		return 0, err
	}
	if len(res) > 0 {
		return res[0].MaxHeight, nil
	}

	return 0, nil
}

// query record by status
func (d *SyncCosmosTask) QueryAll(status []string, taskType string) ([]SyncCosmosTask, error) {
	var syncTasks []SyncCosmosTask
	q := bson.M{}

	if len(status) > 0 {
		q["status"] = bson.M{
			"$in": status,
		}
	}

	switch taskType {
	case db.SyncTaskTypeCatchUp:
		q["end_height"] = bson.M{
			"$ne": 0,
		}
		break
	case db.SyncTaskTypeFollow:
		q["end_height"] = bson.M{
			"$eq": 0,
		}
		break
	}

	fn := func(c *mgo.Collection) error {
		return c.Find(q).All(&syncTasks)
	}

	err := db.ExecCollection(d.Name(), fn)

	if err != nil {
		return syncTasks, err
	}

	return syncTasks, nil
}

func (d *SyncCosmosTask) GetExecutableTask(maxWorkerSleepTime int64) ([]SyncCosmosTask, error) {
	var tasks []SyncCosmosTask

	t := time.Now().Add(time.Duration(-maxWorkerSleepTime) * time.Second).Unix()
	q := bson.M{
		"status": bson.M{
			"$in": []string{db.SyncTaskStatusUnHandled, db.SyncTaskStatusUnderway},
		},
	}

	fn := func(c *mgo.Collection) error {
		return c.Find(q).Sort("-status").Limit(1000).All(&tasks)
	}

	err := db.ExecCollection(d.Name(), fn)

	if err != nil {
		return tasks, err
	}

	ret := make([]SyncCosmosTask, 0, len(tasks))
	//filter the task which last_update_time >= now
	for _, task := range tasks {
		if task.LastUpdateTime >= t && task.Status == db.SyncTaskStatusUnderway {
			continue
		}
		ret = append(ret, task)
	}
	//fmt.Println("SyncCosmosTask GetExecutableTask ret:",ret)
	//fmt.Println("SyncCosmosTask GetExecutableTask:",tasks)

	return ret, nil
}

func (d *SyncCosmosTask) GetTaskById(id bson.ObjectId) (SyncCosmosTask, error) {
	var task SyncCosmosTask

	fn := func(c *mgo.Collection) error {
		return c.FindId(id).One(&task)
	}

	err := db.ExecCollection(d.Name(), fn)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (d *SyncCosmosTask) GetTaskByIdAndWorker(id bson.ObjectId, worker string) (SyncCosmosTask, error) {
	var task SyncCosmosTask

	fn := func(c *mgo.Collection) error {
		q := bson.M{
			"_id":       id,
			"worker_id": worker,
		}

		return c.Find(q).One(&task)
	}

	err := db.ExecCollection(d.Name(), fn)
	if err != nil {
		return task, err
	}
	return task, nil
}

// take over a task
// update status, worker_id, worker_logs and last_update_time
func (d *SyncCosmosTask) TakeOverTask(task SyncCosmosTask, workerId string) error {
	// multiple goroutine attempt to update same record,
	// use this selector to ensure only one goroutine can update success at same time
	fn := func(c *mgo.Collection) error {
		selector := bson.M{
			"_id":              task.ID,
			"last_update_time": task.LastUpdateTime,
		}

		task.Status = db.SyncTaskStatusUnderway
		task.WorkerId = workerId
		task.LastUpdateTime = time.Now().Unix()
		task.WorkerLogs = append(task.WorkerLogs, WorkerLog{
			WorkerId:  workerId,
			BeginTime: time.Now(),
		})

		return c.Update(selector, task)
	}

	return db.ExecCollection(d.Name(), fn)
}

// update task last update time
func (d *SyncCosmosTask) UpdateLastUpdateTime(task SyncCosmosTask) error {
	fn := func(c *mgo.Collection) error {
		selector := bson.M{
			"_id":       task.ID,
			"worker_id": task.WorkerId,
		}

		task.LastUpdateTime = time.Now().Unix()

		return c.Update(selector, task)
	}

	return db.ExecCollection(d.Name(), fn)
}
