package framework

import (
	"sync"
	"time"
)

type limiterEntry struct {
	count     int
	expiresAt time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	clients map[string]*limiterEntry
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:   limit,
		window:  window,
		clients: make(map[string]*limiterEntry),
	}
}

func clientKey(c *Context) string {
	return c.Request.RemoteAddr
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	entry, exists := rl.clients[key]

	if !exists || now.After(entry.expiresAt) {
		rl.clients[key] = &limiterEntry{
			count:     1,
			expiresAt: now.Add(rl.window),
		}
		return true
	}

	if entry.count >= rl.limit {
		return false
	}

	entry.count++
	return true
}

func RateLimit(rl *RateLimiter) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {

			key := c.Request.RemoteAddr

			if !rl.Allow(key) {
				c.Fail(429, "rate limit exceeded")
				return
			}

			next(c)
		}
	}
}
