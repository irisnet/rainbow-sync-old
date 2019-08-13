package task

func Start() {
	synctask := new(TaskCosmosService)
	go synctask.StartCreateTask()
	go synctask.StartExecuteTask()
}
