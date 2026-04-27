package framework

import (
	"strconv"
	"sync"
	"time"

	"github.com/Danieljosh-uduma/zen/pkg/framework/internal/utils"
)

type limiterEntry struct {
	count     int
	expiresAt time.Time
}

type RateLimiter struct {
	mu            sync.Mutex
	limit         int
	window        time.Duration
	clients       map[string]*limiterEntry
	lastCleanup   time.Time
	cleanupPeriod time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:         limit,
		window:        window,
		clients:       make(map[string]*limiterEntry),
		cleanupPeriod: window, // safe default
	}
}

func defaultKeyFn(c *Context) string {
	return utils.GetClientIP(c.Request)
}

func (rl *RateLimiter) Allow(key string) (bool, *limiterEntry) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// cleanup happens inline (cheap + safe)
	rl.cleanup(now)

	entry, exists := rl.clients[key]

	if !exists || now.After(entry.expiresAt) {
		rl.clients[key] = &limiterEntry{
			count:     1,
			expiresAt: now.Add(rl.window),
		}
		return true, rl.clients[key]
	}

	if entry.count >= rl.limit {
		return false, entry
	}

	entry.count++
	return true, entry
}

func (rl *RateLimiter) cleanup(now time.Time) {
	if now.Sub(rl.lastCleanup) < rl.cleanupPeriod {
		return
	}

	for key, entry := range rl.clients {
		if now.After(entry.expiresAt) {
			delete(rl.clients, key)
		}
	}

	rl.lastCleanup = now
}

func RateLimitWithKey(rl *RateLimiter, keyFn func(*Context) string) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			key := keyFn(c)
			res := rl.AllowDetailed(key)

			if !res.Allowed {
				c.Fail(429, "rate limit exceeded")
				return
			}

			c.Writer.Header().Set("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
			c.Writer.Header().Set("X-RateLimit-Reset", strconv.FormatInt(res.ResetAt.Unix(), 10))

			next(c)
		}
	}
}

func RateLimit(rl *RateLimiter) Middleware {
	return RateLimitWithKey(rl, defaultKeyFn)
}

type RateLimitResult struct {
	Allowed   bool
	Remaining int
	ResetAt   time.Time
}

func (rl *RateLimiter) AllowDetailed(key string) RateLimitResult {
	allowed, entry := rl.Allow(key)

	if allowed {
		return RateLimitResult{
			Allowed:   true,
			Remaining: rl.limit - entry.count,
			ResetAt:   entry.expiresAt,
		}
	}
	return RateLimitResult{
		Allowed:   false,
		Remaining: 0,
		ResetAt:   entry.expiresAt,
	}
}

func (rl *RateLimiter) AllowDetailed2(key string) RateLimitResult {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// cleanup happens inline (cheap + safe)
	rl.cleanup(now)

	var rlResult RateLimitResult
	entry, exists := rl.clients[key]

	if !exists || now.After(entry.expiresAt) {
		rl.clients[key] = &limiterEntry{
			count:     1,
			expiresAt: now.Add(rl.window),
		}
		rlResult.Allowed = true
		rlResult.Remaining = rl.limit - entry.count
		rlResult.ResetAt = entry.expiresAt
		return rlResult
	}

	if entry.count >= rl.limit {
		rlResult.Allowed = false
		rlResult.Remaining = 0
		rlResult.ResetAt = entry.expiresAt
		return rlResult
	}

	entry.count++
	rlResult.Allowed = true
	rlResult.Remaining = rl.limit - entry.count
	rlResult.ResetAt = entry.expiresAt
	return rlResult
}
