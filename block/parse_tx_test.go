package block

import (
	"encoding/json"
	"github.com/irisnet/rainbow-sync/conf"
	"github.com/irisnet/rainbow-sync/helper"
	"github.com/irisnet/rainbow-sync/logger"
	"testing"
	"github.com/irisnet/rainbow-sync/model"
	"github.com/irisnet/rainbow-sync/db"
	"fmt"
)

func TestParseCosmosTxModel(t *testing.T) {
	helper.Init(conf.BlockChainMonitorUrl, conf.MaxConnectionNum, conf.InitConnectionNum)
	client := helper.GetTendermintClient()
	defer func() {
		client.Release()
		logger.Info("Release tm client")
	}()
	type args struct {
		b      int64
		client *helper.RpcClient
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test parse cosmos tx",
			args: args{
				b:      502803,
				client: client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cosmos := ZoneBlock{}
			res, err := cosmos.ParseZoneTxs(tt.args.b, tt.args.client)
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

func TestZoneBlock_ParseBlock(t *testing.T) {
	helper.Init(conf.BlockChainMonitorUrl, conf.MaxConnectionNum, conf.InitConnectionNum)
	client := helper.GetTendermintClient()
	db.Start()
	defer func() {
		client.Release()
		logger.Info("Release tm client")
	}()
	type args struct {
		b      int64
		client *helper.RpcClient
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test parse cosmos tx",
			args: args{
				b:      325227,
				client: client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cosmos := ZoneBlock{}
			block, res, err := cosmos.ParseBlock(tt.args.b, tt.args.client)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println("length txs:", len(res))
			if err := cosmos.SaveDocsWithTxn(block, res, model.SyncZoneTask{}); err != nil {
				t.Fatal(err.Error())
			}
		})
	}
}