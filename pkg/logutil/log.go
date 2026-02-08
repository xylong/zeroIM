package logutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

func BizError(ctx context.Context, op string, err error, fields map[string]interface{}) {
	if err == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	op = strings.TrimSpace(op)
	if op == "" {
		op = "unknown_op"
	}
	logger := logx.WithContext(ctx)
	msg := fmt.Sprintf("%s: %s (%T)", op, err.Error(), err)
	if len(fields) > 0 {
		msg = fmt.Sprintf("%s | %s", msg, joinFields(fields))
	}
	logger.Errorf("%s", msg)
}

func BizErrorJSON(ctx context.Context, op string, err error, fields map[string]interface{}) {
	if err == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	op = strings.TrimSpace(op)
	if op == "" {
		op = "unknown_op"
	}
	logger := logx.WithContext(ctx)
	data := map[string]interface{}{
		"op":       op,
		"err":      err.Error(),
		"err_type": fmt.Sprintf("%T", err),
	}
	if fields != nil {
		for k, v := range fields {
			kl := strings.ToLower(k)
			if kl == "op" || kl == "err" || kl == "err_type" {
				continue
			}
			if shouldRedact(kl) {
				data[k] = "***"
				continue
			}
			if s, ok := v.(string); ok {
				data[k] = truncateString(s, 256)
			} else {
				data[k] = v
			}
		}
	}
	logger.Error(data)
}

func shouldRedact(k string) bool {
	if strings.Contains(k, "pass") {
		return true
	}
	if strings.Contains(k, "password") {
		return true
	}
	if strings.Contains(k, "token") {
		return true
	}
	if strings.Contains(k, "secret") {
		return true
	}
	if strings.Contains(k, "authorization") {
		return true
	}
	return false
}

func truncateString(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func joinFields(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return ""
	}
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	// sort for deterministic output
	for i := 0; i < len(keys)-1; i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[j] < keys[i] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	var b strings.Builder
	for idx, k := range keys {
		if idx > 0 {
			b.WriteString(", ")
		}
		kl := strings.ToLower(k)
		if shouldRedact(kl) {
			b.WriteString(k)
			b.WriteString("=***")
			continue
		}
		if s, ok := fields[k].(string); ok {
			b.WriteString(fmt.Sprintf("%s=%s", k, truncateString(s, 256)))
		} else {
			b.WriteString(fmt.Sprintf("%s=%v", k, fields[k]))
		}
	}
	return b.String()
}
