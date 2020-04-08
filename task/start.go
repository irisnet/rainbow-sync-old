package task

func Start() {
	synctask := new(TaskZoneService)
	go synctask.StartCreateTask()
	go synctask.StartExecuteTask()
}
