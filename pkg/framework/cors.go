package framework

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		MaxAge: 300,
	}
}

func applyCORSHeaders(c *Context, config CORSConfig, origin string) {
	// check if origin is allowed
	if len(config.AllowOrigins) > 0 {
		if config.AllowOrigins[0] == "*" {
			c.SetHeader("Access-Control-Allow-Origin", "*")
		} else if slices.Contains(config.AllowOrigins, origin) {
			c.SetHeader("Access-Control-Allow-Origin", origin)
		}
	}

	if config.AllowMethods != nil {
		c.SetHeader("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
	}

	if config.AllowHeaders != nil {
		c.SetHeader("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
	}

	if config.ExposeHeaders != nil {
		c.SetHeader(
			"Access-Control-Expose-Headers",
			strings.Join(config.ExposeHeaders, ", "),
		)
	}

	if config.AllowCredentials {
		c.SetHeader("Access-Control-Allow-Credentials", "true")
	}

	if config.MaxAge > 0 {
		c.SetHeader("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
	}

}

func CORS(config CORSConfig) Middleware {
	// if credentials are allowed, we cannot use wildcard origin
	if config.AllowCredentials && slices.Contains(config.AllowOrigins, "*") {
		panic("CORS: wildcard origin cannot be used with credentials")
	}

	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {

			origin := c.Header("Origin")

			// only apply cors headers if origin is present
			if origin != "" {
				// prevent preflight request on origin requests
				if c.Request.Method == http.MethodOptions {
					c.NoContent()
					return
				}

				applyCORSHeaders(c, config, origin)
			}

			next(c)
		}
	}
}
