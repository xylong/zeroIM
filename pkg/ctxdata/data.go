package ctxdata

import "context"

// GetUId 从jwt token获取uid
func GetUId(ctx context.Context) string {
	if u, ok := ctx.Value(Identify).(string); ok {
		return u
	}
	return ""
}
