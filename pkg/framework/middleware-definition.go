package framework

type MiddlewareDefinition struct {
	Name string
	Func Middleware
}

// NamedMiddleware returns a MiddlewareDefinition for the given name and middleware.
//
// Example:
//
//	loggerMiddleware := framework.NamedMiddleware("logger", framework.Logger())
//	app.UseNamed(loggerMiddleware)
func NamedMiddleware(name string, middleware Middleware) MiddlewareDefinition {
	return MiddlewareDefinition{
		Name: name,
		Func: middleware,
	}
}
