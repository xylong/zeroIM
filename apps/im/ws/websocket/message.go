package websocket

type FrameType uint8

const (
	FrameData FrameType = 0x0 // 数据帧
	FramePing FrameType = 0x1 // 心跳包
)

// Message ws消息
type Message struct {
	FrameType `json:"frameType"`
	Method    string      `json:"method"`
	FromId    string      `json:"fromId"`
	UserId    string      `json:"userId"`
	Data      interface{} `json:"data"`
}

func NewMessage(fromId, userId string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FromId:    fromId,
		UserId:    userId,
		Data:      data}
}
