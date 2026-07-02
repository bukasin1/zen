package zencore

import (
	"strconv"
	"sync"
	"time"

	"github.com/bukasin1/zen/pkg/zencore/internal/utils"
)

// limiterEntry holds the count and expiration time for a key.
type limiterEntry struct {
	count     int
	expiresAt time.Time
}

// RateLimiter is a simple in-memory rate limiter.
type RateLimiter struct {
	mu            sync.Mutex
	limit         int
	window        time.Duration
	clients       map[string]*limiterEntry
	lastCleanup   time.Time
	cleanupPeriod time.Duration
}

// NewRateLimiter creates a new RateLimiter.
//
// limit is the maximum number of requests allowed in a window.
// window is the time duration for the rate limit.
//
// Safe defaults are used for cleanupPeriod.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:         limit,
		window:        window,
		clients:       make(map[string]*limiterEntry),
		cleanupPeriod: window,
	}
}

func defaultKeyFn(c *Context) string {
	return utils.GetClientIP(c.Request)
}

// Allow checks if a request with the given key is allowed.
//
// It returns true if the request is allowed, and the entry for the key.
// If the request is not allowed, it returns false and the entry for the key.
//
// This method is used internally by the rate limiting middleware.
func (rl *RateLimiter) Allow(key string) (bool, *limiterEntry) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// cleanup happens inline (cheap + safe)
	rl.cleanup(now)

	entry, exists := rl.clients[key]

	if !exists || now.After(entry.expiresAt) {
		entry = &limiterEntry{
			count:     1,
			expiresAt: now.Add(rl.window),
		}
		rl.clients[key] = entry
		return true, entry
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

	if len(rl.clients) == 0 {
		rl.lastCleanup = now
		return
	}

	for key, entry := range rl.clients {
		if now.After(entry.expiresAt) {
			delete(rl.clients, key)
		}
	}

	rl.lastCleanup = now
}

// RateLimit is a middleware that rate limits requests to a handler.
//
// It uses a custom key function to determine the key for each request.
//
// This can allow ratelimiting of a specific route or endpoint request
// Use [RateLimitIP] for IP-based rate limiting.
func RateLimit(rl *RateLimiter, keyFn func(*Context) string) Middleware {
	if keyFn == nil {
		keyFn = defaultKeyFn
	}

	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			key := keyFn(c)
			res := rl.AllowDetailed(key)

			if !res.Allowed {
				c.SetHeader("X-RateLimit-Reset", res.ResetAt.UTC().Format(time.RFC1123))
				c.SetHeader("Retry-After", strconv.Itoa(res.RetryAfter))
				c.Fail(429, "rate limit exceeded")
				return
			}

			c.SetHeader("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
			c.SetHeader("X-RateLimit-Reset", res.ResetAt.UTC().Format(time.RFC1123))
			// c.Writer.Header().Set("X-RateLimit-Reset", strconv.FormatInt(res.ResetAt.Unix(), 10))

			next(c)
		}
	}
}

// RateLimitIP is a middleware that rate limits requests to a handler based on IP address.
//
// It uses a default key function that uses the client IP address as the key for each reques.
// You can use [RateLimitIP] for a simple IP-based rate limiting.
//
// To use a different key function, use [RateLimit].
func RateLimitIP(rl *RateLimiter) Middleware {
	return RateLimit(rl, defaultKeyFn)
}

// RateLimitResult is the result of a rate limit check.
type RateLimitResult struct {
	Allowed    bool
	Remaining  int
	ResetAt    time.Time
	RetryAfter int
}

func (rl *RateLimiter) AllowDetailed(key string) RateLimitResult {
	allowed, entry := rl.Allow(key)

	if allowed {
		remaining := max(rl.limit-entry.count, 0)

		return RateLimitResult{
			Allowed:   true,
			Remaining: remaining,
			ResetAt:   entry.expiresAt,
		}
	}
	return RateLimitResult{
		Allowed:    false,
		Remaining:  0,
		ResetAt:    entry.expiresAt,
		RetryAfter: int(time.Until(entry.expiresAt).Seconds()),
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
