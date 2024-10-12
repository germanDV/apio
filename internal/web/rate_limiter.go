package web

import (
	"fmt"
	"strings"
	"time"

	"github.com/germandv/apio/internal/cache"
	"github.com/germandv/apio/internal/errs"
)

type rateLimiter struct {
	cacheClient cache.Client
	limit       int64
	window      time.Duration
}

// newRateLimiter returns a rate limiter that limits requests by IP address.
func newRateLimiter(cacheClient cache.Client, reqsPerMin int64, window time.Duration) rateLimiter {
	return rateLimiter{
		cacheClient: cacheClient,
		limit:       reqsPerMin,
		window:      window,
	}
}

func (rl rateLimiter) check(ip string) error {
	if strings.HasPrefix(ip, "[::1]") {
		ip = "localhost"
	}
	key := fmt.Sprintf("request_count_%s", ip)

	current, err := rl.cacheClient.Increment(key)
	if err != nil {
		return fmt.Errorf("%w: %s", errs.ErrIncCacheCounter, err)
	}

	if current > rl.limit {
		return fmt.Errorf("%w: %d", errs.ErrTooManyRequests, current)
	} else if current == 1 {
		// First request, set expiration to the given window.
		err = rl.cacheClient.Expire(key, rl.window)
		if err != nil {
			return fmt.Errorf("%w: %s", errs.ErrSettingCacheTTL, err)
		}
	}

	return nil
}
