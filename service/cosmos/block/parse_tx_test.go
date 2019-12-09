package block

import (
	"testing"
	"github.com/irisnet/rainbow-sync/service/cosmos/logger"
	cosmosConf "github.com/irisnet/rainbow-sync/service/cosmos/conf"
	"github.com/irisnet/rainbow-sync/service/cosmos/helper"
	"encoding/json"
)

func TestParseCosmosTxModel(t *testing.T) {
	cosmoshelper.Init(cosmosConf.BlockChainMonitorUrl, cosmosConf.MaxConnectionNum, cosmosConf.InitConnectionNum)
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
				b:      174871,
				client: client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cosmos := Cosmos_Block{}
			res, err := cosmos.ParseCosmosTxs(tt.args.b, tt.args.client)
			if err != nil {
				t.Fatal(err)
			}
			resBytes, _ := json.MarshalIndent(res, "", "\t")
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
