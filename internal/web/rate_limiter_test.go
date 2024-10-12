package web

import (
	"errors"
	"testing"
	"time"

	"github.com/germandv/apio/internal/cache/memorycache"
	"github.com/germandv/apio/internal/errs"
)

func TestRateLimiter(t *testing.T) {
	cc, _ := memorycache.New()
	reqsPerMin := int64(3)
	window := 300 * time.Millisecond
	limiter := newRateLimiter(cc, reqsPerMin, window)

	// Requests up to the limit should be accepted.
	for i := 0; i < int(reqsPerMin); i++ {
		err := limiter.check("127.0.0.1")
		if err != nil {
			t.Fatal(err)
		}
	}

	// Request above the limit should be rejected.
	err := limiter.check("127.0.0.1")
	if !errors.Is(err, errs.ErrTooManyRequests) {
		t.Fatal(err)
	}

	// Request from a different IP should be accepted.
	err = limiter.check("168.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	// Request after the window should be accepted.
	time.Sleep(window)
	err = limiter.check("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
}
