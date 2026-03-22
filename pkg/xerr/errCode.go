package xerr

const (
	ServerCommonError = 100001
	RequestParamError = 100002
	TokenExpireError  = 100003
	DbError           = 100004

	// user 业务错误码 (2xxxxx)
	UserPhoneIsRegister  = 200001
	UserPhoneNotRegister = 200002
	UserPasswordError    = 200003
	UserNotExist         = 200004
)

// IsBusinessError 是否为业务错误码
func IsBusinessError(errCode int) bool {
	// 业务错误码通常在特定范围内，或者排除掉系统级错误码
	if errCode > 200000 {
		return true
	}
	// 也可以根据需要补充具体的业务错误码判断
	return false
}
