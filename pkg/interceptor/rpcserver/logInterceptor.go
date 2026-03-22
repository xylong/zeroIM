package rpcserver

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zerr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()

	// Panic Recovery
	defer func() {
		if r := recover(); r != nil {
			err = status.Errorf(codes.Internal, "panic: %v", r)
			logx.WithContext(ctx).Errorf("【RPC SRV PANIC】 method: %s, err: %v", info.FullMethod, r)
		}
	}()

	resp, err = handler(ctx, req)

	duration := time.Since(start)

	if err == nil {
		logx.WithContext(ctx).Infof("【RPC SRV OK】 method: %s, duration: %v", info.FullMethod, duration)
		return resp, nil
	}

	// 错误转换逻辑
	causeErr := errors.Cause(err)
	if e, ok := causeErr.(*zerr.CodeMsg); ok {
		// 转换业务错误码为 gRPC 状态错误，保持业务 code 传递
		err = status.Error(codes.Code(e.Code), e.Msg)
	} else if _, ok := status.FromError(causeErr); !ok {
		// 如果不是 gRPC status 错误，也不是业务错误，包装为 Internal 错误
		err = status.Error(codes.Internal, err.Error())
	}

	logx.WithContext(ctx).Errorf("【RPC SRV ERR】 method: %s, duration: %v, err: %v", info.FullMethod, duration, err)

	return resp, err
}
