package websocket

import "time"

type FrameType uint8

const (
	FrameData  FrameType = 0x0 // 数据帧
	FramePing  FrameType = 0x1 // Ping 帧
	FrameAck   FrameType = 0x2 // Ack 帧
	FrameNoAck FrameType = 0x3 // 无 Ack 帧
	FrameErr   FrameType = 0x9 // 错误帧
)

type Message struct {
	Id        string `json:"id,omitempty"`
	FrameType `json:"frameType,omitempty"`
	AckSeq    int         `json:"ackSeq,omitempty"`
	ackTime   time.Time   `json:"ackTime,omitempty"`
	errCount  int         `json:"errCount,omitempty"`
	Method    string      `json:"method,omitempty"`
	FromId    string      `json:"formId,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func NewMessage(fid string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FromId:    fid,
		Data:      data,
	}
}
func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameData,
		Data:      err.Error(),
	}
}
