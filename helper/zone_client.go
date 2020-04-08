package helper

import (
	"fmt"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
	"github.com/irisnet/rainbow-sync/logger"
	"time"
)

type ZoneClient struct {
	Id string
	rpcClient.Client
}

func newClient(addr string) *ZoneClient {
	rpc, err := rpcClient.NewHTTP(addr, "/websocket")
	if err != nil {
		logger.Error("failted to get client", logger.String("err", err.Error()))
		panic(err.Error())
	}
	return &ZoneClient{
		Id:     generateId(addr),
		Client: rpc,
	}
}

// get client from pool
func GetTendermintClient() *ZoneClient {
	c, err := zoneclient_pool.BorrowObject(ctx)
	for err != nil {
		logger.Error("GetClient failed,will try again after 3 seconds", logger.String("err", err.Error()))
		time.Sleep(3 * time.Second)
		c, err = zoneclient_pool.BorrowObject(ctx)
	}

	return c.(*ZoneClient)
}

// release client
func (c *ZoneClient) Release() {
	err := zoneclient_pool.ReturnObject(ctx, c)
	if err != nil {
		logger.Error(err.Error())
	}
}

func (c *ZoneClient) HeartBeat() error {
	http := c.Client.(*rpcClient.HTTP)
	_, err := http.Health()
	return err
}

func generateId(address string) string {
	return fmt.Sprintf("peer[%s]", address)
}
