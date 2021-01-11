package block

import (
	"encoding/json"
	irisConf "github.com/irisnet/rainbow-sync/conf"
	"github.com/irisnet/rainbow-sync/lib/pool"
	"testing"
)

func TestIris_Block_ParseIrisTx(t *testing.T) {
	pool.Init(irisConf.SvrConf.NodeUrls, irisConf.SvrConf.MaxConnectionNum, irisConf.SvrConf.InitConnectionNum)
	client := pool.GetClient()
	defer func() {
		client.Release()
	}()
	type args struct {
		b      int64
		client *pool.Client
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test parse iris tx",
			args: args{
				b:      857,
				client: client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _, err := ParseTxs(tt.args.b, tt.args.client)
			if err != nil {
				t.Fatal(err)
			}
			resBytes, _ := json.Marshal(res)
			t.Log(string(resBytes))
		})
	}
}
