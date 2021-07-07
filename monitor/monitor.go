package monitor

import (
	"context"
	"github.com/irisnet/rainbow-sync/conf"
	"github.com/irisnet/rainbow-sync/lib/logger"
	"github.com/irisnet/rainbow-sync/lib/pool"
	"github.com/irisnet/rainbow-sync/model"
	"github.com/irisnet/rainbow-sync/monitor/metrics"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	NodeStatusNotReachable = 0
	NodeStatusSyncing      = 1
	NodeStatusCatchingUp   = 2

	SyncTaskFollowing  = 1
	SyncTaskCatchingUp = 0
)

type clientNode struct {
	nodeStatus  metrics.Guage
	nodeHeight  metrics.Guage
	dbHeight    metrics.Guage
	nodeTimeGap metrics.Guage
	syncWorkWay metrics.Guage
}

func NewMetricNode(server metrics.Monitor) clientNode {
	nodeHeightMetric := metrics.NewGuage(
		"sync",
		"status",
		"node_height",
		"full node latest block height",
		nil,
	)
	dbHeightMetric := metrics.NewGuage(
		"sync",
		"status",
		"db_height",
		"sync system database max block height",
		nil,
	)
	nodeStatusMetric := metrics.NewGuage(
		"sync",
		"status",
		"node_status",
		"full node status(0:NotReachable,1:Syncing,2:CatchingUp)",
		nil,
	)
	nodeTimeGapMetric := metrics.NewGuage(
		"sync",
		"status",
		"node_seconds_gap",
		"the seconds gap between running env current time with sync db block time",
		nil,
	)
	syncWorkwayMetric := metrics.NewGuage(
		"sync",
		"",
		"task_working_status",
		"sync task working status(0:CatchingUp 1:Following)",
		nil,
	)
	server.RegisterMetrics(nodeHeightMetric, dbHeightMetric, nodeStatusMetric, nodeTimeGapMetric, syncWorkwayMetric)
	nodeHeight, _ := metrics.CovertGuage(nodeHeightMetric)
	dbHeight, _ := metrics.CovertGuage(dbHeightMetric)
	nodeStatus, _ := metrics.CovertGuage(nodeStatusMetric)
	nodeTimeGap, _ := metrics.CovertGuage(nodeTimeGapMetric)
	syncWorkway, _ := metrics.CovertGuage(syncWorkwayMetric)
	return clientNode{
		nodeStatus:  nodeStatus,
		nodeHeight:  nodeHeight,
		dbHeight:    dbHeight,
		nodeTimeGap: nodeTimeGap,
		syncWorkWay: syncWorkway,
	}
}

func (node *clientNode) Report() {
	for {
		t := time.NewTimer(time.Duration(10) * time.Second)
		select {
		case <-t.C:
			node.nodeStatusReport()
		}
	}
}
func (node *clientNode) nodeStatusReport() {
	client, err := pool.GetClientWithTimeout(10 * time.Second)
	if err != nil {
		logger.Error("rpc node connection exception", logger.String("error", err.Error()))
		node.nodeStatus.Set(float64(NodeStatusNotReachable))
		return
	}
	defer func() {
		client.Release()
	}()

	block, err := new(model.Block).GetMaxBlockHeight()
	if err != nil {
		logger.Error("query block exception", logger.String("error", err.Error()))
	}
	node.dbHeight.Set(float64(block.Height))
	status, err := client.Status(context.Background())
	if err != nil {
		logger.Error("rpc node connection exception", logger.String("error", err.Error()))
		node.nodeStatus.Set(float64(NodeStatusNotReachable))
		//return
	} else {
		if status.SyncInfo.CatchingUp {
			node.nodeStatus.Set(float64(NodeStatusCatchingUp))
		} else {
			node.nodeStatus.Set(float64(NodeStatusSyncing))
		}
		node.nodeHeight.Set(float64(status.SyncInfo.LatestBlockHeight))
	}

	follow, err := new(model.SyncTask).QueryValidFollowTasks()
	if err != nil {
		logger.Error("query valid follow task exception", logger.String("error", err.Error()))
		return
	}
	if follow && block.CreateTime > 0 {
		timeGap := time.Now().Unix() - block.CreateTime
		node.nodeTimeGap.Set(float64(timeGap))
	}

	if follow {
		node.syncWorkWay.Set(float64(SyncTaskFollowing))
	} else {
		node.syncWorkWay.Set(float64(SyncTaskCatchingUp))
	}
	return
}

func Start() {
	c := make(chan os.Signal)
	//monitor system signal
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// start monitor
	server := metrics.NewMonitor(conf.SvrConf.PromethousPort)
	node := NewMetricNode(server)

	server.Report(func() {
		go node.Report()
	})
	<-c
}
