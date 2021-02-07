package block

import (
	"encoding/json"
	"github.com/irisnet/rainbow-sync/lib/pool"
	"testing"
)

func TestIris_Block_ParseIrisTx(t *testing.T) {
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
				b:      4435,
				client: client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, res, msg, err := ParseBlock(tt.args.b, tt.args.client)
			if err != nil {
				t.Fatal(err)
			}
			resBytes, _ := json.Marshal(block)
			t.Log(string(resBytes))
			resBytes, _ = json.Marshal(res)
			t.Log(string(resBytes))
			resBytes, _ = json.Marshal(msg)
			t.Log(string(resBytes))
		})
	}
}
