package model

import (
	"github.com/irisnet/rainbow-sync/db"
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/mgo.v2"
	"fmt"
	"github.com/irisnet/rainbow-sync/conf"
)

const (
	CollectionNameSyncZoneTask = "sync_%v_task"

	SyncZoneTaskStartHeightTag = "start_height"
	SyncZoneTaskEndHeightTag   = "end_height"
	SyncZoneTaskStatusTag      = "status"
)

type (
	SyncZoneTask struct {
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

func (d *SyncZoneTask) Name() string {
	return fmt.Sprintf(CollectionNameSyncZoneTask, conf.ZoneName)
}

func (d *SyncZoneTask) PkKvPair() map[string]interface{} {
	return bson.M{SyncZoneTaskStartHeightTag: d.StartHeight, SyncZoneTaskEndHeightTag: d.EndHeight}
}

func (d *SyncZoneTask) EnsureIndexes() {
	var indexes []mgo.Index
	indexes = append(indexes,
		mgo.Index{
			Key:        []string{SyncZoneTaskStartHeightTag, SyncZoneTaskEndHeightTag},
			Unique:     true,
			Background: true,
		}, mgo.Index{
			Key:        []string{SyncZoneTaskStatusTag},
			Background: true,
		})
	db.EnsureIndexes(d.Name(), indexes)
}

// get max block height in sync task
func (d *SyncZoneTask) GetMaxBlockHeight() (int64, error) {
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
func (d *SyncZoneTask) QueryAll(status []string, taskType string) ([]SyncZoneTask, error) {
	var syncTasks []SyncZoneTask
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

func (d *SyncZoneTask) GetExecutableTask(maxWorkerSleepTime int64) ([]SyncZoneTask, error) {
	var tasks []SyncZoneTask

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

	ret := make([]SyncZoneTask, 0, len(tasks))
	//filter the task which last_update_time >= now
	for _, task := range tasks {
		if task.LastUpdateTime >= t && task.Status == db.SyncTaskStatusUnderway {
			continue
		}
		ret = append(ret, task)
	}
	//fmt.Println("SyncZoneTask GetExecutableTask ret:",ret)
	//fmt.Println("SyncZoneTask GetExecutableTask:",tasks)

	return ret, nil
}

func (d *SyncZoneTask) GetTaskById(id bson.ObjectId) (SyncZoneTask, error) {
	var task SyncZoneTask

	fn := func(c *mgo.Collection) error {
		return c.FindId(id).One(&task)
	}

	err := db.ExecCollection(d.Name(), fn)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (d *SyncZoneTask) GetTaskByIdAndWorker(id bson.ObjectId, worker string) (SyncZoneTask, error) {
	var task SyncZoneTask

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
func (d *SyncZoneTask) TakeOverTask(task SyncZoneTask, workerId string) error {
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
func (d *SyncZoneTask) UpdateLastUpdateTime(task SyncZoneTask) error {
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
