package websocket

type Message struct {
	Method string      `json:"method"`
	FromId string      `json:"fromId"`
	Data   interface{} `json:"data"`
}
