package helper

import (
	"fmt"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
	"github.com/irisnet/rainbow-sync/logger"
	"time"
)

type RpcClient struct {
	Id string
	rpcClient.Client
}

func newClient(addr string) *RpcClient {
	rpc, err := rpcClient.NewHTTP(addr, "/websocket")
	if err != nil {
		logger.Error("failted to get client", logger.String("err", err.Error()))
		panic(err.Error())
	}
	return &RpcClient{
		Id:     generateId(addr),
		Client: rpc,
	}
}

// get client from pool
func GetTendermintClient() *RpcClient {
	c, err := rpcClientPool.BorrowObject(ctx)
	for err != nil {
		logger.Error("GetClient failed,will try again after 3 seconds", logger.String("err", err.Error()))
		time.Sleep(3 * time.Second)
		c, err = rpcClientPool.BorrowObject(ctx)
	}

	return c.(*RpcClient)
}

// release client
func (c *RpcClient) Release() {
	err := rpcClientPool.ReturnObject(ctx, c)
	if err != nil {
		logger.Error(err.Error())
	}
}

func (c *RpcClient) HeartBeat() error {
	http := c.Client.(*rpcClient.HTTP)
	_, err := http.Health()
	return err
}

func generateId(address string) string {
	return fmt.Sprintf("peer[%s]", address)
}
