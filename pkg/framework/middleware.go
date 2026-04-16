package framework

import (
	"fmt"
	"net"
	"net/http"
	"strings"
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

// ANSI colors
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

func statusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return colorGreen
	case status >= 300 && status < 400:
		return colorCyan
	case status >= 400 && status < 500:
		return colorYellow
	default:
		return colorRed
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func getClientIP(r *http.Request) string {
	// common reverse proxy headers
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func Logger() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			start := time.Now()

			next(c)

			status := c.StatusCode()
			if status == 0 {
				status = http.StatusOK
			}

			duration := formatDuration(time.Since(start))
			size := c.ResponseSize()
			method := c.Request.Method
			path := c.Request.URL.Path
			ip := getClientIP(c.Request)

			statusCol := statusColor(status)

			fmt.Printf(
				"%s%s%s | %s%-3d%s | %s%-6s%s %s | %s | %s | %dB\n",
				colorGray,
				time.Now().Format("15:04:05"),
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
