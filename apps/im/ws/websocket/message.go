package websocket

type Message struct {
	Method string      `json:"method"`
	FromId string      `json:"fromId"`
	Data   interface{} `json:"data"`
}

func NewMessage(fromId string, data interface{}) *Message {
	return &Message{FromId: fromId, Data: data}
}
