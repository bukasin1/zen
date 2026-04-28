package framework

import (
	"fmt"
	"net/http"
	"time"

	frameworkErrors "github.com/Danieljosh-uduma/zen/pkg/framework/internal/errors"
	"github.com/Danieljosh-uduma/zen/pkg/framework/internal/utils"
	"github.com/Danieljosh-uduma/zen/pkg/framework/share/logger"
)

type Middleware func(HandlerFunc) HandlerFunc

func chainMiddlewares(h HandlerFunc, middlewares []Middleware) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}

// Recovery is a middleware that recovers from panics and returns a 500 error.
func Recovery() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			defer func() {
				if rec := recover(); rec != nil {
					if c.responseCommitted {
						return
					}

					switch err := rec.(type) {

					case *frameworkErrors.AppError:
						c.Error(err.Status, err.Message, err.Code, err.Details)

					case error:
						c.Error(500, err.Error(), frameworkErrors.ErrInternal, nil)

					default:
						c.Error(500, "internal server error", frameworkErrors.ErrInternal, rec)
					}
				}
			}()

			next(c)
		}
	}
}

// Logger is a middleware that logs requests to the console.
//
// It logs the following information:
// - Request ID
// - Timestamp
// - HTTP method
// - HTTP status code
// - Request path
// - Request duration
// - Request size
// - Client IP address
func Logger() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {

			c.AfterResponse(func(c *Context) {
				status := c.StatusCode()
				if status == 0 {
					status = http.StatusOK
				}

				statusCol := utils.StatusColor(status)

				c.LogInfo("request completed", logger.Fields{
					"status":    fmt.Sprintf("%s%d%s", statusCol, status, utils.ColorReset),
					"duration":  utils.FormatDuration(c.Duration()),
					"size":      c.ResponseSize(),
					"ip":        utils.GetClientIP(c.Request),
					"method":    fmt.Sprintf("%s%s%s", utils.ColorBlue, c.Request.Method, utils.ColorReset),
					"path":      c.Request.URL.Path,
					"requestID": c.RequestID(),
				})
			})

			next(c)

		}
	}
}

func Logger1() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			c.AfterResponse(func(c *Context) {
				status := c.StatusCode()
				if status == 0 {
					status = http.StatusOK
				}

				duration := utils.FormatDuration(c.Duration())
				size := c.ResponseSize()
				method := c.Request.Method
				path := c.Request.URL.Path
				ip := utils.GetClientIP(c.Request)
				rid := c.RequestID()

				statusCol := utils.StatusColor(status)

				fmt.Printf(
					"[%s] %s%s%s | %s%-3d%s | %s%-6s%s %s | %s | %s | %dB\n",
					rid,
					utils.ColorGray,
					c.StartTime().Format("15:04:05"),
					utils.ColorReset,

					statusCol, status, utils.ColorReset,

					utils.ColorBlue, method, utils.ColorReset,
					path,

					duration,
					ip,
					size,
				)

			})

			next(c)
		}
	}
}

func LoggerOld() Middleware {
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
