package resultx

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zrpcErr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
	"net/http"
	"zeroIM/pkg/xerr"
)

// Response api响应
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Success 成功
func Success(data interface{}) *Response {
	return &Response{
		Code: http.StatusOK,
		Msg:  "",
		Data: data,
	}
}

// Fail 失败
func Fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Msg:  err,
		Data: nil,
	}
}

func OkHandler(_ context.Context, v interface{}) any {
	return Success(v)
}

func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		errCode := xerr.ServerCommonError
		errMsg := xerr.ErrMsg(errCode)

		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*zrpcErr.CodeMsg); ok {
			errCode = e.Code
			errMsg = e.Msg
		} else {
			if s, ok := status.FromError(causeErr); ok {
				errCode = int(s.Code())
				errMsg = s.Message()
			}
		}

		// 日志记录
		logx.WithContext(ctx).Errorf("【%s】 err %v", name, err)

		return http.StatusBadRequest, Fail(errCode, errMsg)
	}
}
