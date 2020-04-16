package constant

const (
	TxTypeTransfer        = "Transfer"
	TxTypeIBCBankTransfer = "IBCBankTransfer"

	TxMsgTypeTransfer          = "Transfer"
	TxMsgTypeIBCBankTransfer   = "IBCBankTransfer"
	TxMsgTypeIBCBankMsgPacket  = "IBCMsgPacket"
	TxMsgTypeIBCBankMsgTimeout = "IBCMsgTimeout"

	TxTypeAddLiquidity    = "AddLiquidity"
	TxTypeRemoveLiquidity = "RemoveLiquidity"
	TxTypeSwapOrder       = "SwapOrder"

	TxStatusSuccess          = "success"
	TxStatusFail             = "fail"
	EventTypeSendPacket      = "send_packet"
	EventAttributesKeyPacket = "packet_data"

	EnvNameZoneName                     = "ZONE_NAME"
	EnvNameSerNetworkFullNode_ZONE      = "SER_BC_FULL_NODE_ZONE"
	EnvNameWorkerNumExecuteTask_ZONE    = "WORKER_NUM_EXECUTE_TASK_ZONE"
	EnvNameWorkerMaxSleepTime_ZONE      = "WORKER_MAX_SLEEP_TIME_ZONE"
	EnvNameBlockNumPerWorkerHandle_ZONE = "BLOCK_NUM_PER_WORKER_HANDLE_ZONE"

	EnvNameDbAddr     = "DB_ADDR"
	EnvNameDbUser     = "DB_USER"
	EnvNameDbPassWd   = "DB_PASSWD"
	EnvNameDbDataBase = "DB_DATABASE"

	BatchLimit = 1000
)
