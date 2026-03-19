package cachex

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"
)

var ErrNotFound = errors.New("not found")

type Wrapper[T any] struct {
	Data *T `json:"data,omitempty"`
}

type Cache struct {
	rc *rockscache.Client
}

func NewCache(rdb redis.UniversalClient) *Cache {
	return &Cache{
		rc: rockscache.NewClient(rdb, rockscache.NewDefaultOptions()),
	}
}

// GetWithCache 获取缓存数据
func GetWithCache[T any](
	ctx context.Context,
	c *Cache,
	key string,
	opt Options,
	query func(ctx context.Context) (*T, error),
) (*T, error) {

	val, err := c.rc.Fetch(key, opt.getTTL(), func() (string, error) {
		data, err := query(ctx)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				res, e := json.Marshal(Wrapper[T]{Data: nil})
				if e != nil {
					return "", e
				}
				return string(res), nil
			}
			return "", err
		}

		res, err := json.Marshal(Wrapper[T]{Data: data})
		if err != nil {
			return "", err
		}
		return string(res), nil
	})

	if err != nil {
		return nil, err
	}

	var wrapper Wrapper[T]
	if err := json.Unmarshal([]byte(val), &wrapper); err != nil {
		return nil, err
	}

	if wrapper.Data == nil {
		return nil, ErrNotFound
	}

	return wrapper.Data, nil
}

func (c *Cache) Del(keys ...string) error {
	for _, key := range keys {
		if err := c.rc.TagAsDeleted(key); err != nil {
			return err
		}
	}
	return nil
}
