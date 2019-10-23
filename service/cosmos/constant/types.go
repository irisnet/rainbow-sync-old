package constant

const (
	TxTypeTransfer                  = "Transfer"
	TxTypeIBCBankTransfer           = "IBCBankTransfer"
	TxTypeIBCBankRecvTransferPacket = "IBCRecvTransferPacket"

	TxMsgTypeTransfer                  = "Transfer"
	TxMsgTypeIBCBankTransfer           = "IBCBankTransfer"
	TxMsgTypeIBCBankRecvTransferPacket = "IBCRecvTransferPacket"

	TxStatusSuccess          = "success"
	TxStatusFail             = "fail"
	EventTypeSendPacket      = "send_packet"
	EventAttributesKeyPacket = "Packet"

	EnvNameSerNetworkFullNode_COSMOS      = "SER_BC_FULL_NODE_COSMOS"
	EnvNameWorkerNumCreateTask_COSMOS     = "WORKER_NUM_CREATE_TASK_COSMOS"
	EnvNameWorkerNumExecuteTask_COSMOS    = "WORKER_NUM_EXECUTE_TASK_COSMOS"
	EnvNameWorkerMaxSleepTime_COSMOS      = "WORKER_MAX_SLEEP_TIME_COSMOS"
	EnvNameBlockNumPerWorkerHandle_COSMOS = "BLOCK_NUM_PER_WORKER_HANDLE_COSMOS"

	EnvNameDbAddr     = "DB_ADDR"
	EnvNameDbUser     = "DB_USER"
	EnvNameDbPassWd   = "DB_PASSWD"
	EnvNameDbDataBase = "DB_DATABASE"
)
