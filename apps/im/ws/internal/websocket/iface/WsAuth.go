package iface

import "net/http"

// Authentication 鉴权
type Authentication interface {
	// Auth 鉴权
	Auth(w http.ResponseWriter, r *http.Request) bool
	// UserId 获取用户id
	UserId(r *http.Request) string
}
