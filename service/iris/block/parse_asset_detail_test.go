package block

import (
	"testing"
	irisConf "github.com/irisnet/rainbow-sync/service/iris/conf"
	"encoding/json"
	"github.com/irisnet/rainbow-sync/service/iris/logger"
	"github.com/irisnet/rainbow-sync/service/iris/helper"
)

func TestIris_Block_ParseIrisAssetDetail(t *testing.T) {
	helper.Init(irisConf.BlockChainMonitorUrl, irisConf.MaxConnectionNum, irisConf.InitConnectionNum)
	client := helper.GetClient()
	defer func() {
		client.Release()
		logger.Info("Release tm client")
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
			name: "test parse asset detail",
			args: args{
				b:      19301,
				client: client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			biris := Iris_Block{}
			res, err := biris.ParseIrisAssetDetail(tt.args.b, tt.args.client)
			if err != nil {
				t.Fatal(err)
			}
			resBytes, _ := json.MarshalIndent(res, "", "\t")
			t.Log(string(resBytes))
		})
	}
}
