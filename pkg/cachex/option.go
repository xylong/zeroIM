package cachex

import (
	"math/rand"
	"time"
)

type Options struct {
	TTL       time.Duration
	NilTTL    time.Duration // 防穿透
	RandomTTL time.Duration // 随机抖动，防雪崩
	CacheNil  bool          // 是否缓存nil
}

func randTTL(base, jitter time.Duration) time.Duration {
	if base <= 0 {
		return 0
	}
	if jitter <= 0 {
		return base
	}
	return base + time.Duration(rand.Int63n(int64(jitter)))
}
