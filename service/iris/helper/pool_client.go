//init client from clientPool.
//client is httpClient of tendermint

package helper

import (
	"fmt"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
	"github.com/irisnet/rainbow-sync/service/iris/logger"
	"time"
)

type Client struct {
	Id string
	rpcClient.Client
}

func newClient(addr string) *Client {
	return &Client{
		Id:     generateId(addr),
		Client: rpcClient.NewHTTP(addr, "/websocket"),
	}
}

// get client from pool
func GetClient() *Client {
	c, err := pool.BorrowObject(ctx)
	for err != nil {
		logger.Error("GetClient failed,will try again after 3 seconds", logger.String("err", err.Error()))
		time.Sleep(3 * time.Second)
		c, err = pool.BorrowObject(ctx)
	}

	return c.(*Client)
}

// release client
func (c *Client) Release() {
	err := pool.ReturnObject(ctx, c)
	if err != nil {
		logger.Error(err.Error())
	}
}

func (c *Client) HeartBeat() error {
	http := c.Client.(*rpcClient.HTTP)
	_, err := http.Health()
	return err
}

func generateId(address string) string {
	return fmt.Sprintf("peer[%s]", address)
}
