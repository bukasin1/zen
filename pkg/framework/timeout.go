package framework

import (
	"context"
	"errors"
	"net/http"
	"time"

	frameworkErrors "github.com/Danieljosh-uduma/zen/pkg/framework/errors"
)

func Timeout(duration time.Duration) Middleware {
	if duration <= 0 {
		panic("Timeout: duration must be greater than zero")
	}
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			// timeout := time.After(duration)
			ctx, cancel := context.WithTimeout(
				c.Request.Context(),
				duration,
			)
			defer cancel()

			c.Request = c.Request.WithContext(ctx)

			done := make(chan struct{})

			go func() {
				defer close(done)
				next(c)
			}()

			select {
			case <-done:
				return

			case <-ctx.Done():
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					c.Error(http.StatusGatewayTimeout, "Request timed out", frameworkErrors.ErrRequestTimeout, nil)
					// c.Abort()
				}
			}

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
