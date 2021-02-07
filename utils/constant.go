package utils

const (
	TxStatusSuccess     = "success"
	TxStatusFail        = "fail"
	NoSupportMsgTypeTag = "no support msg parse"
	//unable to resolve type URL /cosmos.bank.v1beta1.MsgSend
	ErrNoSupportTxPrefix = "unable to resolve type URL"

	//cannot find transaction 601bf70ccdee4dde1c8be0d2_f018677a in queue for document {sync_task ObjectIdHex(\"601bdb0ccdee4dd7c214d167\")}
	ErrDbNotFindTransaction = "cannot find transaction"
)
