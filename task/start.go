package task

func Start() {
	synctask := new(TaskIrisService)
	go synctask.StartCreateTask()
	go synctask.StartExecuteTask()
}
