package framework

import (
	"fmt"
	"net/http"
	"time"
)

type Middleware func(HandlerFunc) HandlerFunc

func chainMiddlewares(h HandlerFunc, middlewares []Middleware) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}

func Recovery() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			defer func() {
				if err := recover(); err != nil {
					_ = c.JSON(500, map[string]any{
						"error":   "internal server error",
						"details": err,
					})
				}
			}()

			next(c)
		}
	}
}

func Logger() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			// start := time.Now()

			next(c)

			status := c.StatusCode()
			if status == 0 {
				status = http.StatusOK
			}

			duration := formatDuration(c.Duration())
			size := c.ResponseSize()
			method := c.Request.Method
			path := c.Request.URL.Path
			ip := getClientIP(c.Request)
			rid := c.RequestID()

			statusCol := statusColor(status)

			fmt.Printf(
				"[%s] %s%s%s | %s%-3d%s | %s%-6s%s %s | %s | %s | %dB\n",
				rid,
				colorGray,
				c.StartTime().Format("15:04:05"),
				colorReset,

				statusCol, status, colorReset,

				colorBlue, method, colorReset,
				path,

				duration,
				ip,
				size,
			)
		}
	}
}

func Logger1() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			start := time.Now()

			next(c)

			status := c.StatusCode()
			if status == 0 {
				status = http.StatusOK
			}

			fmt.Printf(
				"%v: %-7s %s, %dB .....................................%d ---- %v\n",
				time.Now().Format("2006-01-02 15:04:05"),
				c.Request.Method,
				c.Request.URL.Path,
				c.ResponseSize(),
				status,
				time.Since(start),
			)
		}
	}
}
