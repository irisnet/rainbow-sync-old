package block

import (
	"encoding/json"
	"github.com/irisnet/rainbow-sync/service/cosmos/conf"
	"github.com/irisnet/rainbow-sync/service/cosmos/helper"
	"github.com/irisnet/rainbow-sync/service/cosmos/logger"
	"testing"
)

func TestParseCosmosTxModel(t *testing.T) {
	cosmoshelper.Init(conf.BlockChainMonitorUrl, conf.MaxConnectionNum, conf.InitConnectionNum)
	client := cosmoshelper.GetCosmosClient()
	defer func() {
		client.Release()
		logger.Info("Release tm client")
	}()
	type args struct {
		b      int64
		client *cosmoshelper.CosmosClient
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test parse cosmos tx",
			args: args{
				b:      17174,
				client: client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cosmos := CosmosBlock{}
			res, err := cosmos.ParseCosmosTxs(tt.args.b, tt.args.client)
			if err != nil {
				t.Fatal(err)
			}
			resBytes, _ := json.Marshal(res)
			t.Log(string(resBytes))
		})
	}
}

func Test_parseRawlog(t *testing.T) {
	rawlog := "[{\"msg_index\":\"0\",\"success\":false,\"log\":\"\"}," +
		"{\"msg_index\":\"1\",\"success\":true,\"log\":\"\"}," +
		"{\"msg_index\":\"2\",\"success\":true,\"log\":\"\"}]"
	ret, err := parseRawlog(rawlog)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}
