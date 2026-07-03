package zencore

type TokenValidator interface {
	Validate(ctx *Context, token string) (any, error)
}

func AuthMiddleware(validator TokenValidator) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {

			authHeader := c.Header("Authorization")
			if authHeader == "" {
				next(c)
				return
			}

			// Expect: "Bearer <token>"
			const prefix = "Bearer "
			if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
				c.Unauthorized("invalid authorization header format")
				return
			}

			token := authHeader[len(prefix):]

			user, err := validator.Validate(c, token)
			if err != nil {
				c.Unauthorized("invalid or expired token")
				return
			}

			c.SetUser(user)

			next(c)
		}
	}
}

func RequireAuth() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			if c.User() == nil {
				c.Unauthorized("authentication required")
				return
			}

			next(c)
		}
	}
}

func RequireRole(check func(user any) bool) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			user := c.User()
			if user == nil {
				c.Unauthorized("authentication required")
				return
			}

			if !check(user) {
				c.Forbidden("insufficient permissions")
				return
			}

			next(c)
		}
	}
}
