package websocket

import "strings"

type ServerOption func(*serverOption)

type serverOption struct {
	auth    Authentication
	pattern string
}

func newServerOption(options ...ServerOption) *serverOption {
	so := &serverOption{
		auth:    &authentication{},
		pattern: "/ws",
	}

	for _, option := range options {
		option(so)
	}
	return so
}

func WithAuthentication(auth Authentication) ServerOption {
	return func(o *serverOption) {
		if auth != nil {
			o.auth = auth
		}
	}
}

func WithPattern(pattern string) ServerOption {
	return func(o *serverOption) {
		if pattern == "" {
			return
		}
		if !strings.HasPrefix(pattern, "/") {
			pattern = "/" + strings.Trim(pattern, "/")
		}
		o.pattern = pattern
	}
}

// todo 鉴权、心跳间隔、最大连接数、是否支持多端、是否启用压缩
