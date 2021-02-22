//init client from clientPool.
//client is httpClient of tendermint

package pool

import (
	"context"
	"fmt"
	"github.com/irisnet/rainbow-sync/lib/logger"
	rpcClient "github.com/tendermint/tendermint/rpc/client/http"
	"time"
)

type Client struct {
	Id string
	*rpcClient.HTTP
}

func newClient(addr string) (*Client, error) {
	client, err := rpcClient.New(addr, "/websocket")
	return &Client{
		Id:   generateId(addr),
		HTTP: client,
	}, err
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
	http := c.HTTP
	_, err := http.Health(context.Background())
	return err
}

func generateId(address string) string {
	return fmt.Sprintf("peer[%s]", address)
}
