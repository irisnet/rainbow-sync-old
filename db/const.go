package db

const (
	// value of status
	SyncTaskStatusUnHandled = "unhandled"
	SyncTaskStatusUnderway  = "underway"
	SyncTaskStatusCompleted = "completed"

	// only for follow task
	// when current_height of follow task add blockNumPerWorkerHandle
	// less than blockchain current_height, this follow task's status should be set invalid
	FollowTaskStatusInvalid = "invalid"

	// taskType
	SyncTaskTypeCatchUp = "catch_up"
	SyncTaskTypeFollow  = "follow"
)
