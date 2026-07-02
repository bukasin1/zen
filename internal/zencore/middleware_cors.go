package zencore

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type CORSConfig struct {
	AllowOrigins     []string // Can be empty, but if not, it will be used. If empty, it will use [DefaultCORSConfig]
	AllowMethods     []string // Can be empty, but if not, it will be used. If empty, it will use [DefaultCORSConfig]
	AllowHeaders     []string // Can be empty, but if not, it will be used. If empty, it will use [DefaultCORSConfig]
	ExposeHeaders    []string // Can be empty, but if not, it will be used. If empty, it will use [DefaultCORSConfig]
	AllowCredentials bool     // Can be empty, but if not, it will be used. If empty, it will use [DefaultCORSConfig]
	MaxAge           int      // Can be empty, but if not, it will be used. If empty, it will use [DefaultCORSConfig]
}

// DefaultCORSConfig returns the default CORS configuration.
// It allows all origins, methods, and headers.
// It also sets the MaxAge to 300 seconds (5 minutes).
// Use [DefaultCORSConfig] when you want to allow all origins, methods, and headers.
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

func normalizeOrigin(origin string) string {
	// trim whitespace and trailing slash
	return strings.TrimSuffix(strings.TrimSpace(origin), "/")
}

// normalizeOrigins normalizes the origins by trimming whitespace and trailing slashes
func normalizeOrigins(origins []string) []string {
	for i, origin := range origins {
		origins[i] = normalizeOrigin(origin)
	}
	return origins
}

// normalizes the CORS configuration by setting default values if not provided.
// It also normalizes origins by trimming whitespace and trailing slashes.
func normalizeCORSConfig(config CORSConfig) CORSConfig {
	// set default values
	cfg := DefaultCORSConfig()

	// override with provided values
	if len(config.AllowOrigins) > 0 {
		cfg.AllowOrigins = config.AllowOrigins
	}
	if len(config.AllowMethods) > 0 {
		cfg.AllowMethods = config.AllowMethods
	}
	if len(config.AllowHeaders) > 0 {
		cfg.AllowHeaders = config.AllowHeaders
	}
	if len(config.ExposeHeaders) > 0 {
		cfg.ExposeHeaders = config.ExposeHeaders
	}
	if config.MaxAge > 0 {
		cfg.MaxAge = config.MaxAge
	}
	if config.AllowCredentials {
		cfg.AllowCredentials = config.AllowCredentials
	}

	cfg.AllowOrigins = normalizeOrigins(cfg.AllowOrigins)

	return cfg
}

func applyCORSHeaders(c *Context, config CORSConfig, origin string) {
	normalizedOrigin := normalizeOrigin(origin)
	// check if origin is allowed
	if len(config.AllowOrigins) > 0 {
		if config.AllowOrigins[0] == "*" {
			c.SetHeader("Access-Control-Allow-Origin", "*")
		} else if slices.Contains(config.AllowOrigins, normalizedOrigin) {
			c.SetHeader("Access-Control-Allow-Origin", normalizedOrigin)
			c.AddHeader("Vary", "Origin")
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

// CORS returns a middleware that applies CORS headers to the response.
//
// Note:
//   - [DefaultCORSConfig] should be used when you want to allow all origins, methods, and headers.
//   - [DefaultCORSConfig] is used if no config is provided.
//   - When custom config is provided, it merges with default values if the field is empty/zero value
//     otherwise it overrides the default config fields.
//   - If credentials are allowed, wildcard origin cannot be used.
//   - Wildcard origin ("*") must be the first and only entry in AllowOrigins if it's included.
//     Mixing wildcard and specific origins is not allowed.
//
// Example:
//
//	// use default cors config
//	app.Use(zencore.CORS(zencore.DefaultCORSConfig()))
//	// This works as well - with default config
//	app.Use(zencore.CORS(zencore.CORSConfig{}))
//
//	// use custom cors config
//	app.Use(zencore.CORS(zencore.CORSConfig{
//		AllowOrigins:     []string{"http://localhost:3000"},
//		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
//		AllowHeaders:     []string{"Content-Type", "Authorization"},
//		ExposeHeaders:    []string{"Content-Length"},
//		AllowCredentials: true,
//		MaxAge:           86400,
//	}))
func CORS(config CORSConfig) Middleware {
	// normalize config
	cfg := normalizeCORSConfig(config)

	// validate wildcard origin - must be first and only wildcard entry
	if slices.Contains(cfg.AllowOrigins, "*") && cfg.AllowOrigins[0] != "*" {
		panic(`CORS: "*" origin must be the first and only wildcard entry`)
	}

	// validate wildcard origin - cannot be combined with specific origins
	if cfg.AllowOrigins[0] == "*" && len(cfg.AllowOrigins) > 1 {
		panic(`CORS: wildcard origin "*" cannot be combined with specific origins`)
	}

	// if credentials are allowed, we cannot use wildcard origin
	if cfg.AllowCredentials && slices.Contains(cfg.AllowOrigins, "*") {
		panic("CORS: wildcard origin cannot be used with credentials")
	}

	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {

			origin := c.Header("Origin")

			// only apply cors headers if origin is present
			if origin != "" {
				applyCORSHeaders(c, cfg, origin)

				// prevent preflight request on origin requests
				if c.Request.Method == http.MethodOptions {
					c.NoContent()
					return
				}
			}

			next(c)
		}
	}
}
