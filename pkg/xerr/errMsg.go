package xerr

var codeText = map[int]string{
	ServerCommonError: "服务异常，请稍后处理",
	RequestParamError: "参数不正确",
	TokenExpireError:  "token失效，请重新登陆",
	DbError:           "数据库繁忙,请稍后再试",

	UserPhoneIsRegister:  "手机号已注册",
	UserPhoneNotRegister: "手机号未注册",
	UserPasswordError:    "密码错误",
	UserNotExist:         "用户不存在",
}

func ErrMsg(errCode int) string {
	if msg, ok := codeText[errCode]; ok {
		return msg
	}
	return codeText[ServerCommonError]
}
