package utils

const (
	TxStatusSuccess     = "success"
	TxStatusFail        = "fail"
	NoSupportMsgTypeTag = "no support msg parse"

	//cannot find transaction 601bf70ccdee4dde1c8be0d2_f018677a in queue for document {sync_task ObjectIdHex(\"601bdb0ccdee4dd7c214d167\")}
	ErrDbNotFindTransaction = "cannot find transaction"

	IbcRecvPacketEventTypeDenomTrace         = "denomination_trace"
	IbcRecvPacketEventAttrKeyDenomTrace      = "denom"
	IbcTransferEventTypeSendPacket           = "send_packet"
	IbcTransferEventAttriKeyPacketSequence   = "packet_sequence"
	IbcTransferEventAttriKeyPacketScPort     = "packet_src_port"
	IbcTransferEventAttriKeyPacketScChannel  = "packet_src_channel"
	IbcTransferEventAttriKeyPacketDcPort     = "packet_dst_port"
	IbcTransferEventAttriKeyPacketDcChannels = "packet_dst_channel"
)
