package constants

type HandlerResult int

const (
	// NoHandlerResult 未处理
	NoHandlerResult HandlerResult = iota + 1
	// PassHandlerResult 处理通过
	PassHandlerResult
	// RejectHandlerResult 处理拒绝
	RejectHandlerResult
	// CancelHandlerResult 处理取消
	CancelHandlerResult
)
