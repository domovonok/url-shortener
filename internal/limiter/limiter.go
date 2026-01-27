package limiter

import (
	"sync"
	"time"

	"github.com/domovonok/url-shortener/internal/config"
)

type TokenBucket struct {
	capacity   int
	tokens     int
	refillRate int
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(cfg config.RateLimitConfig) *TokenBucket {
	return &TokenBucket{
		capacity:   cfg.Capacity,
		tokens:     cfg.Capacity,
		refillRate: cfg.RefillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

func (tb *TokenBucket) Remaining() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	return tb.tokens
}

func (tb *TokenBucket) Capacity() int {
	return tb.capacity
}

func (tb *TokenBucket) RefillRate() int {
	return tb.refillRate
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	tokensToAdd := int(elapsed.Milliseconds()) * tb.refillRate / 1000
	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}
}
