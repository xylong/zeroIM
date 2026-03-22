package ctxdata

import (
	"context"
	"encoding/json"
	"strconv"
)

// GetUId 从jwt token获取uid
func GetUId(ctx context.Context) int64 {
	if v := ctx.Value(Identify); v != nil {
		switch val := v.(type) {
		case string:
			id, _ := strconv.ParseInt(val, 10, 64)
			return id
		case float64:
			return int64(val)
		case int64:
			return val
		case int:
			return int64(val)
		case json.Number:
			id, _ := val.Int64()
			return id
		}
	}
	return 0
}

// GetUIdStr 从jwt token获取uid字符串
func GetUIdStr(ctx context.Context) string {
	if v := ctx.Value(Identify); v != nil {
		switch val := v.(type) {
		case string:
			return val
		case float64:
			return strconv.FormatInt(int64(val), 10)
		case int64:
			return strconv.FormatInt(val, 10)
		case int:
			return strconv.Itoa(val)
		case json.Number:
			return val.String()
		}
	}
	return ""
}
