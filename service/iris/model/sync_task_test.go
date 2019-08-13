package iris

import (
	"testing"
	"encoding/json"
	model "github.com/irisnet/rainbow-sync/service/iris/db"
)

func TestMain(m *testing.M) {
	model.Start()
	m.Run()
}

func TestSyncTask_GetExecutableTask(t *testing.T) {
	d := SyncTask{}

	if res, err := d.GetExecutableTask(120); err != nil {
		t.Fatal(err)
	} else {
		resBytes, _ := json.Marshal(res)
		t.Log(string(resBytes))
	}
}
