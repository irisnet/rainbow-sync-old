package model

import (
	"encoding/json"
	model "github.com/irisnet/rainbow-sync/db"
	"testing"
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

func TestErrTx_Save(t *testing.T) {
	errtx := ErrTx{
		Log:    "",
		Height: 7056,
		Repair: 1,
		TxHash: "",
	}
	errtx.Save()
}

func TestErrTx_Find(t *testing.T) {
	if res, err := new(ErrTx).Find(0, 5); err != nil {
		t.Fatal(err)
	} else {
		resBytes, _ := json.Marshal(res)
		t.Log(string(resBytes))
	}
}
