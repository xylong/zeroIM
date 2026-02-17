package websocket

import (
	"fmt"
	"net/http"
	"time"
)

// Authentication 鉴权
type Authentication interface {
	// Auth 鉴权
	Auth(w http.ResponseWriter, r *http.Request) bool
	// UserId 获取用户id
	UserId(r *http.Request) string
}

type authentication struct{}

func (a *authentication) Auth(w http.ResponseWriter, r *http.Request) bool {
	// todo 握手阶段鉴权
	// header带token或者直接传token
	return true
}

func (a *authentication) UserId(r *http.Request) string {
	query := r.URL.Query()
	if query != nil && query["userId"] != nil {
		return fmt.Sprintf("%v", query["userId"])
	}

	return fmt.Sprintf("%v", time.Now().UnixMilli())
}
