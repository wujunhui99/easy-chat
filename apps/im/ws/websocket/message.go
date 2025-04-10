package websocket
type FrameType uint8
const(
	FrameDate FrameType = 0x0
	FramePing FrameType = 0x1
)
type Message struct {
	FrameType FrameType `json:"frameType"`
	Method string `json:"method,omitempty"`
	UserId string `json:"userId,omitempty"`
	FromId string `json:"fromId,omitempty"`
	Data interface{} `json:"data",omitempty"`
}

func NewMessage(fid string, data interface{}) *Message {
	return &Message{
		FrameType: FrameDate,
		FromId: fid,
		Data: data,
	}
}
func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameDate,
		Data: err.Error(),
	}
}