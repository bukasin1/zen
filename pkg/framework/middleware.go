package framework

import (
	"fmt"
	"net/http"
	"time"

	frameworkErrors "github.com/Danieljosh-uduma/zen/pkg/framework/errors"
	"github.com/Danieljosh-uduma/zen/pkg/framework/internal/utils"
	"github.com/Danieljosh-uduma/zen/pkg/framework/logger"
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
					c.LogError("Request panicked", logger.Fields{
						"path":      c.Request.URL.Path,
						"method":    c.Request.Method,
						"requestID": c.RequestID(),
						"panic":     rec,
					})

					if c.responseCommitted.Load() {
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

// RequestLogger is a middleware that logs requests to the console.
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
func RequestLogger() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {

			c.AfterResponse(func(c *Context) {
				status := c.StatusCode()
				if status == 0 {
					status = http.StatusOK
				}

				c.LogInfo("request completed", logger.Fields{
					"status":    status,
					"duration":  utils.FormatDuration(c.Duration()),
					"size":      c.ResponseSize(),
					"ip":        utils.GetClientIP(c.Request),
					"method":    c.Request.Method,
					"path":      c.Request.URL.Path,
					"requestID": c.RequestID(),
				})
			})

			next(c)

		}
	}
}

// MaxBodySize is a middleware that limits the size of the request body.
//
// It reads the request body using http.MaxBytesReader.
//
// Parameters:
//   - limit: The maximum size of the request body in bytes. (limit = 0 is same as not calling middleware)
//
// Example:
//
//	app.Use(framework.MaxBodySize(10 * 1024 * 1024)) // 10MB limit
//
// Important: This middleware should be placed before any middleware that reads the request body.
func MaxBodySize(limit int64) Middleware {
	if limit < 0 {
		panic("MaxBodySize: limit cannot be negative")
	}

	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			if limit == 0 {
				// explicitly disable limiting
				next(c)
				return
			}

			if c.Request.ContentLength != -1 && c.Request.ContentLength > limit {
				c.Error(http.StatusRequestEntityTooLarge, "Request body too large", frameworkErrors.ErrRequestBodyTooLarge, nil)
				return
			}
			if c.Request.Body != nil {
				c.Request.Body = http.MaxBytesReader(
					c.Writer,
					c.Request.Body,
					limit,
				)
			}

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
