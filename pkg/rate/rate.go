package rate

import (
	"sync"
	"time"
)

// RateLimiter struct with token bucket mechanism
type RateLimiter struct {
	rate   int          // rate of requests per second
	burst  int          // max burst size
	tokens int          // current number of tokens
	mu     sync.Mutex   // to make the rate limiter thread-safe
	ticker *time.Ticker // ticker to replenish tokens
}

// NewRateLimiter initializes a new rate limiter
func NewRateLimiter(rate int, burst int) *RateLimiter {
	rl := &RateLimiter{
		rate:   rate,
		burst:  burst,
		tokens: burst,
		ticker: time.NewTicker(time.Second / time.Duration(rate)),
	}
	go rl.start()
	return rl
}

// start replenishes tokens at the specified rate
func (rl *RateLimiter) start() {
	for range rl.ticker.C {
		rl.mu.Lock()
		if rl.tokens < rl.burst {
			rl.tokens++
		}
		rl.mu.Unlock()
	}
}

// Allow checks if a request is allowed
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}
