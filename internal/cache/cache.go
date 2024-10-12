package cache

import (
	"time"
)

type Client interface {
	Ping() error
	Close() error
	Increment(key string) (int64, error)
	Expire(key string, ttl time.Duration) error
}
