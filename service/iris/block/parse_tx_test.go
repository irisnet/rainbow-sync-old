package block

import (
	"testing"
	"github.com/irisnet/rainbow-sync/service/iris/helper"
	irisConf "github.com/irisnet/rainbow-sync/service/iris/conf"
	"encoding/json"
)

func TestIris_Block_ParseIrisTx(t *testing.T) {
	helper.Init(irisConf.BlockChainMonitorUrl, irisConf.MaxConnectionNum, irisConf.InitConnectionNum)
	client := helper.GetClient()
	defer func() {
		client.Release()
	}()
	type args struct {
		b      int64
		client *helper.Client
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test parse iris tx",
			args: args{
				b:      108,
				client: client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iris := Iris_Block{}
			res, err := iris.ParseIrisTxs(tt.args.b, tt.args.client)
			if err != nil {
				t.Fatal(err)
			}
			resBytes, _ := json.MarshalIndent(res, "", "\t")
			t.Log(string(resBytes))
		})
	}
}
