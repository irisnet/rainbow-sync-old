package helper

import (
	"context"
	"github.com/irisnet/rainbow-sync/logger"
	commonPool "github.com/jolestar/go-commons-pool"
	"math/rand"
	"sync"
)

type (
	PoolFactory struct {
		peersMap sync.Map
	}
	EndPoint struct {
		Address   string
		Available bool
	}
)

var (
	zoneclient_poolFactory PoolFactory
	zoneclient_pool        *commonPool.ObjectPool
	ctx                    = context.Background()
)

func Init(BlockChainMonitorUrl []string, MaxConnectionNum, InitConnectionNum int) {
	var syncMap sync.Map
	for _, url := range BlockChainMonitorUrl {
		key := generateId(url)
		endPoint := EndPoint{
			Address:   url,
			Available: true,
		}

		syncMap.Store(key, endPoint)
	}
	zoneclient_poolFactory = PoolFactory{
		peersMap: syncMap,
	}

	config := commonPool.NewDefaultPoolConfig()
	config.MaxTotal = MaxConnectionNum
	config.MaxIdle = InitConnectionNum
	config.MinIdle = InitConnectionNum
	config.TestOnBorrow = true
	config.TestOnCreate = true
	config.TestWhileIdle = true

	logger.Info("PoolConfig", logger.Int("config.MaxTotal", config.MaxTotal),
		logger.Int("config.MaxIdle", config.MaxIdle))
	zoneclient_pool = commonPool.NewObjectPool(ctx, &zoneclient_poolFactory, config)
	zoneclient_pool.PreparePool(ctx)
}

func ClosePool() {
	zoneclient_pool.Close(ctx)
}

func (f *PoolFactory) MakeObject(ctx context.Context) (*commonPool.PooledObject, error) {
	endpoint := f.GetEndPoint()
	return commonPool.NewPooledObject(newClient(endpoint.Address)), nil
}

func (f *PoolFactory) DestroyObject(ctx context.Context, object *commonPool.PooledObject) error {
	c := object.Object.(*ZoneClient)
	if c.IsRunning() {
		c.Stop()
	}
	return nil
}

func (f *PoolFactory) ValidateObject(ctx context.Context, object *commonPool.PooledObject) bool {
	// do validate
	c := object.Object.(*ZoneClient)
	if c.HeartBeat() != nil {
		value, ok := f.peersMap.Load(c.Id)
		if ok {
			endPoint := value.(EndPoint)
			endPoint.Available = true
			f.peersMap.Store(c.Id, endPoint)
		}
		return false
	}
	return true
}

func (f *PoolFactory) ActivateObject(ctx context.Context, object *commonPool.PooledObject) error {
	return nil
}

func (f *PoolFactory) PassivateObject(ctx context.Context, object *commonPool.PooledObject) error {
	return nil
}

func (f *PoolFactory) GetEndPoint() EndPoint {
	var (
		keys        []string
		selectedKey string
	)

	f.peersMap.Range(func(k, value interface{}) bool {
		key := k.(string)
		endPoint := value.(EndPoint)
		if endPoint.Available {
			keys = append(keys, key)
		}
		selectedKey = key

		return true
	})

	if len(keys) > 0 {
		index := rand.Intn(len(keys))
		selectedKey = keys[index]
	}
	value, ok := f.peersMap.Load(selectedKey)
	if ok {
		return value.(EndPoint)
	} else {
		logger.Error("Can't get selected end point", logger.String("selectedKey", selectedKey))
	}
	return EndPoint{}
}
