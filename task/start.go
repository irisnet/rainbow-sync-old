package task

import "github.com/irisnet/rainbow-sync/monitor"

func Start() {
	synctask := new(TaskIrisService)
	go synctask.StartCreateTask()
	go synctask.StartExecuteTask()
	go monitor.Start()
}
