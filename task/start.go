package task

import "github.com/irisnet/rainbow-sync/cron"

func Start() {
	synctask := new(TaskIrisService)
	go synctask.StartCreateTask()
	go synctask.StartExecuteTask()
	go new(cron.CronService).StartCronService()
}
