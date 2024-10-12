package memorycache

import (
	"errors"
	"time"
)

type CacheValue struct {
	value     int
	expiresAt time.Time
}

type MemoryCache struct {
	store map[string]CacheValue
}

func New() (*MemoryCache, error) {
	client := &MemoryCache{
		store: make(map[string]CacheValue),
	}

	err := client.Ping()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *MemoryCache) Ping() error {
	return nil
}

func (c *MemoryCache) Close() error {
	return nil
}

func (c *MemoryCache) Increment(key string) (int64, error) {
	var curr CacheValue

	prev, ok := c.store[key]
	if !ok {
		curr = CacheValue{value: 1, expiresAt: time.Now().Add(24 * time.Hour)}
	} else {
		curr = CacheValue{value: prev.value + 1, expiresAt: prev.expiresAt}
	}

	if curr.expiresAt.Before(time.Now()) {
		curr.value = 1
		curr.expiresAt = time.Now().Add(24 * time.Hour)
	}

	c.store[key] = curr
	return int64(curr.value), nil
}

func (c *MemoryCache) Expire(key string, ttl time.Duration) error {
	record, ok := c.store[key]
	if !ok {
		return errors.New("key not found")
	}

	record.expiresAt = time.Now().Add(ttl)
	c.store[key] = record
	return nil
}
