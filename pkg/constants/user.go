package constants

// UserStatus 用户状态：1正常 2禁用
type UserStatus int

const (
	EnableStatus UserStatus = iota + 1
	DisableStatus
)
