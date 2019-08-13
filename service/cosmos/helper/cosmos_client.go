package cosmoshelper

import (
	"fmt"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
	"github.com/irisnet/rainbow-sync/service/cosmos/logger"
	"time"
)

type CosmosClient struct {
	Id string
	rpcClient.Client
}

func newClient(addr string) *CosmosClient {
	return &CosmosClient{
		Id:     generateId(addr),
		Client: rpcClient.NewHTTP(addr, "/websocket"),
	}
}

// get client from pool
func GetCosmosClient() *CosmosClient {
	c, err := cosmos_pool.BorrowObject(ctx)
	for err != nil {
		logger.Error("GetClient failed,will try again after 3 seconds", logger.String("err", err.Error()))
		time.Sleep(3 * time.Second)
		c, err = cosmos_pool.BorrowObject(ctx)
	}

	return c.(*CosmosClient)
}

// release client
func (c *CosmosClient) Release() {
	err := cosmos_pool.ReturnObject(ctx, c)
	if err != nil {
		logger.Error(err.Error())
	}
}

func (c *CosmosClient) HeartBeat() error {
	http := c.Client.(*rpcClient.HTTP)
	_, err := http.Health()
	return err
}

func generateId(address string) string {
	return fmt.Sprintf("peer[%s]", address)
}
