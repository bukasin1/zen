package framework

import (
	"context"
	"time"
)

// Timeout returns a middleware that injects a timeout-aware context
// into the request lifecycle.
//
// IMPORTANT:
//
// This middleware DOES NOT forcibly terminate handler execution.
// Go does not support safely killing goroutines.
//
// Instead, this middleware:
//
//   - creates a request-scoped timeout context
//   - propagates cancellation via c.Request.Context()
//   - allows downstream services (DBs, HTTP clients, etc.) to stop work
//   - preserves framework response lifecycle integrity
//   - preserves panic recovery behavior
//
// Handlers and services are expected to respect:
//
//	select {
//	case <-ctx.Done():
//	    return
//	default:
//	}
//
// Example:
//
//	app.Use(framework.Timeout(5 * time.Second))
//
//	db.QueryContext(c.StdContext(), ...)
//
//	http.NewRequestWithContext(c.StdContext(), ...)
func Timeout(duration time.Duration) Middleware {
	if duration <= 0 {
		panic("Timeout: duration must be greater than zero")
	}
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			// create timeout-aware request context
			ctx, cancel := context.WithTimeout(
				c.Request.Context(),
				duration,
			)
			defer cancel()

			// attach updated context to request
			c.Request = c.Request.WithContext(ctx)

			// continue request lifecycle normally
			next(c)

			// below are things we considered but did not implement (Forced Timeout)
			// done := make(chan struct{})

			// go func() {
			// 	defer close(done)
			// 	next(c)
			// }()

			// select {
			// case <-done:
			// 	return

			// case <-ctx.Done():
			// 	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			// 		// or http.StatusGatewayTimeout
			// 		c.Error(http.StatusRequestTimeout, "Request timed out", frameworkErrors.ErrRequestTimeout, nil)
			// 		// c.Abort()
			// 	}
			// }

			// // Create a channel to receive the result
			// resultChan := make(chan bool)

			// go func() {
			// 	defer func() {
			// 		if r := recover(); r != nil {
			// 			resultChan <- false // indicate panic to timeout goroutine
			// 		}
			// 	}()
			// 	next(c)
			// 	resultChan <- true // indicate success
			// }()

			// select {
			// case <-ctx.Done():
			// 	if ctx.Err() == context.DeadlineExceeded {
			// 		c.Error(http.StatusGatewayTimeout, "Request timed out", frameworkErrors.ErrRequestTimeout, nil)
			// 	}
			// 	// c.Abort()
			// case <-resultChan:
			// 	// continue normal flow only if next(c) returned
			// 	return
			// }

		}
	}
}
