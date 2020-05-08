package constant

const (
	TxTypeTransfer          = "Transfer"
	TxTypeIBCBankTransfer   = "IBCBankTransfer"
	TxTypeIBCBankMsgPacket  = "IBCMsgPacket"
	TxTypeIBCBankMsgTimeout = "IBCMsgTimeout"

	TxMsgTypeTransfer          = "Transfer"
	TxMsgTypeIBCBankTransfer   = "IBCBankTransfer"
	TxMsgTypeIBCBankMsgPacket  = "IBCMsgPacket"
	TxMsgTypeIBCBankMsgTimeout = "IBCMsgTimeout"

	TxTypeAddLiquidity    = "AddLiquidity"
	TxTypeRemoveLiquidity = "RemoveLiquidity"
	TxTypeSwapOrder       = "SwapOrder"

	TxStatusSuccess                    = "success"
	TxStatusFail                       = "fail"
	EventTypeSendPacket                = "send_packet"
	EventAttributesKeyPacket           = "packet_data"
	EventAttributesKeySequence         = "packet_sequence"
	EventAttributesKeyDstPort          = "packet_dst_port"
	EventAttributesKeyDstChannel       = "packet_dst_channel"
	EventAttributesKeySrcPort          = "packet_src_port"
	EventAttributesKeySrcChannel       = "packet_src_channel"
	EventAttributesKeyTimeoutHeight    = "packet_timeout_height"
	EventAttributesKeyTimeoutTimestamp = "packet_timeout_timestamp"
	EventTypeCoinSwapTransfer          = "transfer"

	EnvNameZoneChainId                  = "ZONE_CHAIN_ID"
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
