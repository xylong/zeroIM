package cachex

import (
	"context"
	"errors"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/singleflight"
)

// todo 带本地缓存 + batch + 自动续期

// 1 字节，ASCII NUL，无歧义,无反序列化开销
const (
	nilValue          = "\x00"
	nilValueByte byte = 0
)

func isNilCache(val []byte) bool {
	return len(val) == 1 && val[0] == nilValueByte
}

type Cache struct {
	rdb *redis.Client
	sf  singleflight.Group
}

func NewCache(rdb *redis.Client) *Cache {
	if rdb == nil {
		panic("redis client cannot be nil")
	}
	return &Cache{rdb: rdb}
}

// GetWithCache 读取缓存
func GetWithCache[T any](
	ctx context.Context,
	c *Cache, key string,
	opt Options,
	fn func(ctx context.Context) (*T, error),
) (*T, error) {
	// 1. Fast Path: 优先查询缓存，不进 SingleFlight(对热点数据避免内部互斥锁竞争)
	val, err := c.rdb.Get(ctx, key).Bytes()
	if err == nil {
		if isNilCache(val) {
			return nil, nil
		}
		var res T
		if err := sonic.Unmarshal(val, &res); err == nil {
			return &res, nil
		}
		logx.WithContext(ctx).Errorf("cache decode failed, key=%s, err=%v", key, err)
		_ = c.Del(ctx, key)
	} else if !errors.Is(err, redis.Nil) {
		logx.WithContext(ctx).Errorf("redis get error: %v", err)
	}

	// 2. Slow Path: 进 SingleFlight 防止击穿
	v, err, _ := c.sf.Do(key, func() (any, error) {
		// Double Check: 进锁后再查一次缓存，防止并发请求在等待 SingleFlight 时缓存已被回写
		val, err := c.rdb.Get(ctx, key).Bytes()
		if err == nil {
			if isNilCache(val) {
				return nil, nil
			}
			var res T
			if err := sonic.Unmarshal(val, &res); err == nil {
				return &res, nil
			}
			_ = c.Del(ctx, key)
		}

		// ===== 3. 回源 =====
		res, err := fn(ctx)
		if err != nil {
			return nil, err
		}

		// ===== 4. 处理 Nil 缓存 (防穿透) =====
		if res == nil {
			// 是否允许缓存空值
			if !opt.CacheNil {
				return nil, nil
			}
			nilTTL := opt.NilTTL
			if nilTTL <= 0 {
				nilTTL = time.Second * 10 // ✅ 防穿透兜底
			}

			ttl := randTTL(nilTTL, opt.RandomTTL)
			if ttl > 0 {
				if err := c.rdb.Set(ctx, key, nilValue, ttl).Err(); err != nil {
					logx.WithContext(ctx).Errorf("set nil cache failed, key=%s, err=%v", key, err)
				}
			}

			return nil, nil
		}

		// ===== 5. 处理正常缓存 =====
		data, err := sonic.Marshal(res)
		if err != nil {
			return nil, err
		}

		ttl := randTTL(opt.TTL, opt.RandomTTL)
		if ttl > 0 {
			if err := c.rdb.Set(ctx, key, data, ttl).Err(); err != nil {
				logx.WithContext(ctx).Errorf("cache set failed, key=%s, err=%v", key, err)
			}
		}

		return res, nil
	})

	if err != nil {
		return nil, err
	}

	if v == nil {
		return nil, nil
	}

	res, ok := v.(*T)
	if !ok {
		return nil, errors.New("type assertion failed")
	}

	return res, nil
}

func Set[T any](ctx context.Context, c *Cache, key string, val T, ttl time.Duration) error {
	data, err := sonic.Marshal(val)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

// ================= Hash 示例 =================

func HSet[T any](ctx context.Context, c *Cache, key, field string, val T) error {
	data, err := sonic.Marshal(val)
	if err != nil {
		return err
	}
	return c.rdb.HSet(ctx, key, field, data).Err()
}

func HGet[T any](ctx context.Context, c *Cache, key, field string) (*T, error) {
	val, err := c.rdb.HGet(ctx, key, field).Bytes()
	if err != nil {
		return nil, err
	}
	var res T
	if err := sonic.Unmarshal(val, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}
