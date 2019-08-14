package iris

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	model "github.com/irisnet/rainbow-sync/service/iris/db"
	"time"
)

const (
	CollectionNameSyncTask = "sync_iris_task"
)

type (
	SyncTask struct {
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

func (d SyncTask) Name() string {
	return CollectionNameSyncTask
}

func (d SyncTask) PkKvPair() map[string]interface{} {
	return bson.M{"start_height": d.CurrentHeight, "end_height": d.EndHeight}
}

// get max block height in sync task
func (d SyncTask) GetMaxBlockHeight() (int64, error) {
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
	err := model.ExecCollection(d.Name(), getMaxBlockHeightFn)

	if err != nil {
		return 0, err
	}
	if len(res) > 0 {
		return res[0].MaxHeight, nil
	}

	return 0, nil
}

// query record by status
func (d SyncTask) QueryAll(status []string, taskType string) ([]SyncTask, error) {
	var syncTasks []SyncTask
	q := bson.M{}

	if len(status) > 0 {
		q["status"] = bson.M{
			"$in": status,
		}
	}

	switch taskType {
	case model.SyncTaskTypeCatchUp:
		q["end_height"] = bson.M{
			"$ne": 0,
		}
		break
	case model.SyncTaskTypeFollow:
		q["end_height"] = bson.M{
			"$eq": 0,
		}
		break
	}

	fn := func(c *mgo.Collection) error {
		return c.Find(q).All(&syncTasks)
	}

	err := model.ExecCollection(d.Name(), fn)

	if err != nil {
		return syncTasks, err
	}

	return syncTasks, nil
}

func (d SyncTask) GetExecutableTask(maxWorkerSleepTime int64) ([]SyncTask, error) {
	var tasks []SyncTask

	t := time.Now().Add(time.Duration(-maxWorkerSleepTime) * time.Second).Unix()
	q := bson.M{
		"status": bson.M{
			"$in": []string{model.SyncTaskStatusUnHandled, model.SyncTaskStatusUnderway},
		},
	}

	fn := func(c *mgo.Collection) error {
		return c.Find(q).Sort("-status").Limit(1000).All(&tasks)
	}

	err := model.ExecCollection(d.Name(), fn)

	if err != nil {
		return tasks, err
	}

	ret := make([]SyncTask, 0, len(tasks))
	//filter the task which last_update_time >= now
	for _, task := range tasks {
		if task.LastUpdateTime >= t && task.Status == model.SyncTaskStatusUnderway {
			continue
		}
		ret = append(ret, task)
	}

	return ret, nil
}

func (d SyncTask) GetTaskById(id bson.ObjectId) (SyncTask, error) {
	var task SyncTask

	fn := func(c *mgo.Collection) error {
		return c.FindId(id).One(&task)
	}

	err := model.ExecCollection(d.Name(), fn)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (d SyncTask) GetTaskByIdAndWorker(id bson.ObjectId, worker string) (SyncTask, error) {
	var task SyncTask

	fn := func(c *mgo.Collection) error {
		q := bson.M{
			"_id":       id,
			"worker_id": worker,
		}

		return c.Find(q).One(&task)
	}

	err := model.ExecCollection(d.Name(), fn)
	if err != nil {
		return task, err
	}
	return task, nil
}

// take over a task
// update status, worker_id, worker_logs and last_update_time
func (d SyncTask) TakeOverTask(task SyncTask, workerId string) error {
	// multiple goroutine attempt to update same record,
	// use this selector to ensure only one goroutine can update success at same time
	fn := func(c *mgo.Collection) error {
		selector := bson.M{
			"_id":              task.ID,
			"last_update_time": task.LastUpdateTime,
		}

		task.Status = model.SyncTaskStatusUnderway
		task.WorkerId = workerId
		task.LastUpdateTime = time.Now().Unix()
		task.WorkerLogs = append(task.WorkerLogs, WorkerLog{
			WorkerId:  workerId,
			BeginTime: time.Now(),
		})

		return c.Update(selector, task)
	}

	return model.ExecCollection(d.Name(), fn)
}

// update task last update time
func (d SyncTask) UpdateLastUpdateTime(task SyncTask) error {
	fn := func(c *mgo.Collection) error {
		selector := bson.M{
			"_id":       task.ID,
			"worker_id": task.WorkerId,
		}

		task.LastUpdateTime = time.Now().Unix()

		return c.Update(selector, task)
	}

	return model.ExecCollection(d.Name(), fn)
}
