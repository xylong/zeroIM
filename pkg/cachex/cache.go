package cachex

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"time"
)

type Cache struct {
	rdb *redis.Client
	sf  *singleflight.Group
}

func NewCache(rdb *redis.Client) *Cache {
	return &Cache{
		rdb: rdb,
		sf:  &singleflight.Group{},
	}
}

func GetWithCache[T any](
	c *Cache,
	ctx context.Context,
	key string,
	opt Options,
	query func() (*T, error),
) (*T, error) {

	// 1. 查缓存
	if c.rdb != nil {
		val, err := c.rdb.Get(ctx, key).Result()
		if err == nil {
			var wrapper Wrapper
			if json.Unmarshal([]byte(val), &wrapper) == nil {

				if !wrapper.Exist {
					return nil, gorm.ErrRecordNotFound
				}

				var t T
				if json.Unmarshal(wrapper.Data, &t) == nil {
					return &t, nil
				}
			}
		}
	}

	// 2. singleflight 防击穿
	res, err, _ := c.sf.Do(key, func() (interface{}, error) {
		data, err := query()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.setAsync(ctx, key, Wrapper{Exist: false}, opt.NilTTL)
				return nil, gorm.ErrRecordNotFound
			}
			return nil, err
		}

		b, _ := json.Marshal(data)
		c.setAsync(ctx, key, Wrapper{
			Exist: true,
			Data:  b,
		}, opt.getTTL())

		return data, nil
	})

	if err != nil {
		return nil, err
	}

	v, ok := res.(*T)
	if !ok {
		return nil, errors.New("cachex: type assert fail")
	}

	return v, nil
}

func (c *Cache) setAsync(ctx context.Context, key string, val Wrapper, ttl time.Duration) {
	if c.rdb == nil {
		return
	}

	go func() {
		b, _ := json.Marshal(val)

		cctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_ = c.rdb.Set(cctx, key, b, ttl).Err()
	}()
}

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	if c.rdb == nil {
		return nil
	}
	return c.rdb.Del(ctx, keys...).Err()
}

func (c *Cache) DelPattern(ctx context.Context, pattern string) error {
	if c.rdb == nil {
		return nil
	}

	iter := c.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		_ = c.rdb.Del(ctx, iter.Val()).Err()
	}

	return iter.Err()
}
