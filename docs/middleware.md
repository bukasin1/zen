# Middleware

Middleware provides a mechanism for executing logic before and after a request reaches its final route handler.

Zen uses middleware to implement cross-cutting concerns such as logging, authentication, rate limiting, compression, metrics, and panic recovery.

Middleware should remain focused on infrastructure concerns and should not contain business logic.

---

# What Is Middleware?

Middleware is a function that wraps another request handler.

It allows applications to inspect, modify, or terminate request processing before the final handler executes.

Common middleware responsibilities include:

* Authentication
* Authorization
* Logging
* Panic recovery
* Request metrics
* Rate limiting
* CORS
* Gzip compression
* Request body limits
* Timeouts
* HTTP caching

Business logic should remain inside route handlers or application services.

---

# Middleware Execution Order

Middleware executes in a predictable order.

```text
Incoming Request
        │
        ▼
Global Middleware
        │
        ▼
Group Middleware
        │
        ▼
Route Middleware
        │
        ▼
Route Handler
        │
        ▼
AfterResponse Hooks
        │
        ▼
Outgoing Response
```

Execution always follows this sequence.

This deterministic ordering makes applications easier to understand and debug.

---

# Global Middleware

Global middleware applies to every request handled by the application.

Example:

```go
app.Use(
	LoggingMiddleware(),
	RecoveryMiddleware(),
	MetricsMiddleware(),
)
```

Global middleware should be used for concerns shared by all routes.

Typical examples include:

* Request logging
* Panic recovery
* Metrics
* Request IDs

---

# Group Middleware

Middleware may also be attached to route groups.

Example:

```go
api := app.Group("/api")

api.Use(
	AuthMiddleware(),
)
```

All routes within the group automatically inherit the group's middleware.

Nested groups also inherit middleware from their parent groups.

This makes group middleware ideal for API versioning or protected sections of an application.

---

# Route Middleware

Middleware may be attached to individual routes.

Example:

```go
app.Route("/profile").
	Use(RequireAuth()).
	Get(profileHandler)
```

Route middleware executes after global and group middleware.

Use route middleware for behavior that applies only to specific endpoints.

---

# Middleware Composition

Middleware layers naturally.

Example:

```text
Global
    ↓
API Group
    ↓
Admin Group
    ↓
Route
    ↓
Handler
```

Each layer contributes additional behavior while preserving predictable execution order.

---

# Middleware Short-Circuiting

Middleware may stop request processing.

For example:

```text
Request
    │
Authentication Middleware
    │
Authentication Failed
    │
401 Unauthorized
```

When middleware writes a response and terminates processing, subsequent middleware and the route handler are not executed.

This behavior is commonly used for:

* Authentication failures
* Authorization failures
* Rate limiting
* Invalid requests
* Request size violations

---

# Writing Custom Middleware

Custom middleware follows the standard middleware signature used throughout Zen.

Example:

```go
func RequestTimer() zen.Middleware {
	return func(next zen.HandlerFunc) zen.HandlerFunc {
		return func(c *zen.Context) {

			start := time.Now()

			next(c)

			duration := time.Since(start)

			log.Printf("request completed in %s", duration)
		}
	}
}
```

Middleware should generally perform work before and/or after calling `next`.

Calling `next` passes control to the next middleware or the final route handler.

---

# Built-in Middleware

Zen includes a collection of production-oriented middleware that can be composed as needed.

Built-in middleware includes:

* Panic recovery
* Structured request logging
* Authentication
* Authorization
* Rate limiting
* Request body size limiting
* CORS
* Cooperative request timeouts
* Gzip compression
* HTTP caching (ETag and Cache-Control)
* Metrics instrumentation

Middleware is intentionally modular. Applications enable only the middleware they require.

---

# Named Middleware

Middleware may be registered with a descriptive name.

Example:

```go
app.UseNamed(
	"RequestLogger",
	LoggingMiddleware(),
)

app.UseNamed(
	"Recovery",
	RecoveryMiddleware(),
)
```

Naming middleware improves:

* Route introspection
* Generated documentation
* Debugging
* Administrative tooling

Middleware names should be unique and descriptive.

---

# Authentication and Authorization

Authentication and authorization middleware have an intentional execution order.

Authentication must execute before authorization.

Example:

```go
app.Use(AuthMiddleware())

app.Route("/admin").
	Use(RequireRole("admin")).
	Get(adminHandler)
```

Execution order:

```text
AuthMiddleware
        │
        ▼
RequireRole
        │
        ▼
Route Handler
```

`RequireRole()` expects an authenticated user to already exist in the request context.

Running authorization before authentication results in an invalid middleware configuration.

---

# Cooperative Timeouts

Zen implements cooperative request timeouts.

Timeout middleware signals that a request has exceeded its configured duration.

Long-running handlers are expected to cooperate by respecting cancellation through the request context.

Zen intentionally does not forcibly terminate goroutines.

Forced interruption risks:

* Partial responses
* Response corruption
* Resource leaks
* Unpredictable application state

Cooperative cancellation preserves response integrity and aligns with Go's concurrency model.

---

# Response Lifecycle

Middleware should never bypass Zen's response management.

Responses should always be written through the `Context`.

Example:

```go
c.JSON(data)
```

instead of writing directly to the underlying `http.ResponseWriter`.

This ensures:

* Response lifecycle protection
* Correct status handling
* Consistent headers
* Double-write prevention
* Compatibility with middleware such as compression and caching

---

# Panic Recovery

Applications should include panic recovery middleware near the beginning of the global middleware chain.

Example:

```go
app.Use(
	RecoveryMiddleware(),
	LoggingMiddleware(),
)
```

Recovery middleware prevents unexpected panics from terminating the HTTP server.

Recovered panics are converted into controlled HTTP responses while preserving diagnostic information for logging.

---

# Error Handling

Middleware should return appropriate HTTP responses when terminating request processing.

Examples include:

* `400 Bad Request`
* `401 Unauthorized`
* `403 Forbidden`
* `404 Not Found`
* `413 Payload Too Large`
* `429 Too Many Requests`
* `500 Internal Server Error`

Errors should be explicit and meaningful.

Avoid exposing internal implementation details in client-facing responses.

---

# Middleware Performance

Middleware executes for every matching request.

Keep middleware lightweight.

Recommendations:

* Avoid unnecessary allocations.
* Minimize blocking operations.
* Reuse immutable configuration where practical.
* Avoid performing expensive work that belongs in application services.

Middleware should primarily coordinate request flow rather than perform heavy computation.

---

# Best Practices

## Keep Middleware Focused

Each middleware should have a single responsibility.

Good examples include:

* Logging
* Authentication
* Compression

Avoid combining unrelated responsibilities into a single middleware.

---

## Prefer Composition

Multiple small middleware components are generally easier to understand and maintain than one large middleware that performs many unrelated tasks.

---

## Keep Business Logic Out of Middleware

Middleware should implement infrastructure concerns.

Business rules belong in application services or route handlers.

---

## Respect Request Context Cancellation

Long-running middleware should observe request context cancellation and terminate work when appropriate.

This becomes especially important for:

* External service calls
* Database operations
* File processing

---

## Do Not Modify Shared State

Middleware executes concurrently across many requests.

Avoid mutating shared data without appropriate synchronization.

---

## Register Middleware During Startup

Middleware should be registered before the application begins serving requests.

Avoid dynamically modifying middleware pipelines while requests are being processed.

---

# Common Mistakes

## Registering Authorization Before Authentication

Incorrect:

```go
app.Route("/admin").
	Use(RequireRole("admin")).
	Get(adminHandler)
```

Correct:

```go
app.Use(AuthMiddleware())

app.Route("/admin").
	Use(RequireRole("admin")).
	Get(adminHandler)
```

---

## Writing Directly to the ResponseWriter

Incorrect:

```go
w.Write(...)
```

Preferred:

```go
c.JSON(...)
```

Using the `Context` preserves Zen's managed response lifecycle.

---

## Ignoring Context Cancellation

Long-running operations should periodically check whether the request has been cancelled.

Ignoring cancellation may waste resources after the client has disconnected or a timeout has occurred.

---

## Performing Heavy Business Logic

Middleware should coordinate request processing.

Complex business workflows should remain within application services.

---

# Summary

Zen's middleware system is intentionally simple, explicit, and composable.

By separating cross-cutting concerns from business logic and enforcing a predictable execution order, middleware remains easy to understand, test, and maintain.

The framework encourages small, focused middleware components that cooperate naturally to build reliable, production-ready request pipelines.
