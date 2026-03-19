package cachex

import (
	"math/rand"
	"time"
)

type Options struct {
	TTL       time.Duration
	RandomTTL time.Duration // 随机抖动，防雪崩
}

func (o *Options) getTTL() time.Duration {
	if o.RandomTTL <= 0 {
		return o.TTL
	}
	return o.TTL + time.Duration(rand.Int63n(int64(o.RandomTTL)))
}
