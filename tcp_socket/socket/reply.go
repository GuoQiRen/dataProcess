package socket

import "dataProcess/tools/utils/jsonx"

func ReplyMessage(reply jsonx.Marshaler, id string, seq uint32, doc []byte, cancelTask chan string) []byte {
	cancelTask <- id
	return MarshalRespond(seq, reply, doc) // 发回的消息格式
}
