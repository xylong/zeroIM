package xerr

var codeText = map[int]string{
	ServerCommonError: "服务异常，请稍后处理",
	RequestParamError: "参数不正确",
	TokenExpireError:  "token失效，请重新登陆",
	DbError:           "数据库繁忙,请稍后再试",
}

func ErrMsg(errCode int) string {
	if msg, ok := codeText[errCode]; ok {
		return msg
	}
	return codeText[ServerCommonError]
}
