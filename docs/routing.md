# Routing

Zen's routing system is designed to be explicit, predictable, and efficient.

Routes are registered directly in code using a fluent API that encourages readability while remaining close to Go's standard library.

This document explains route registration, route groups, dynamic parameters, metadata, static files, and recommended routing practices.

---

# Registering Routes

Routes are registered through the application's `Route()` method.

Example:

```go
app.Route("/users").
	Get(func(c *zen.Context) {
		c.JSON(users)
	})
```

Unlike many frameworks, Zen intentionally avoids individual registration methods such as:

```go
app.Get(...)
app.Post(...)
app.Put(...)
```

Instead, every route begins with:

```go
app.Route(...)
```

This provides a single location for attaching metadata, middleware, and additional route configuration while keeping the API consistent.

---

# Supported HTTP Methods

Zen supports all standard HTTP methods.

Example:

```go
app.Route("/users").Get(handler)
app.Route("/users").Post(handler)
app.Route("/users/:id").Put(handler)
app.Route("/users/:id").Patch(handler)
app.Route("/users/:id").Delete(handler)

app.Route("/users").Head(handler)
app.Route("/users").Options(handler)
```

Each method registers a distinct route for the specified path.

---

# Route Paths

Route paths should begin with a forward slash.

Examples:

```text
/users
/users/:id
/api/v1/users
/files/*
```

Zen automatically normalizes redundant slashes during registration.

Example:

```text
//users///profile
```

becomes:

```text
/users/profile
```

This helps eliminate accidental duplicate routes caused by inconsistent path formatting.

---

# Route Parameters

Dynamic route parameters are defined using a leading colon.

Example:

```go
app.Route("/users/:id").
	Get(func(c *zen.Context) {
		id := c.Param("id")

		c.JSON(map[string]string{
			"id": id,
		})
	})
```

Request:

```text
GET /users/42
```

Result:

```go
id == "42"
```

Parameter names are case-sensitive.

---

# Multiple Parameters

Routes may contain multiple parameters.

Example:

```go
app.Route("/users/:userID/orders/:orderID").
	Get(func(c *zen.Context) {

		userID := c.Param("userID")
		orderID := c.Param("orderID")

		// ...
	})
```

Request:

```text
/users/15/orders/200
```

Produces:

```text
userID = 15
orderID = 200
```

---

# Wildcard Routes

Wildcard routes match everything after the specified path.

Example:

```go
app.Route("/assets/*").
	Get(handler)
```

Requests such as:

```text
/assets/css/site.css
/assets/images/logo.png
/assets/js/app.js
```

are all matched by the same route.

Wildcard routes should generally be reserved for file serving or similar use cases.

---

# Route Groups

Related routes can be organized into groups.

Example:

```go
api := app.Group("/api")

api.Route("/users").
	Get(usersHandler)

api.Route("/products").
	Get(productsHandler)
```

This produces:

```text
/api/users
/api/products
```

Groups improve readability while reducing repeated path prefixes.

---

# Nested Groups

Groups may be nested.

Example:

```go
api := app.Group("/api")

v1 := api.Group("/v1")

v1.Route("/users").
	Get(usersHandler)
```

Result:

```text
/api/v1/users
```

Nested groups inherit middleware from their parent groups.

---

# API Versioning

Versioning is implemented naturally using groups.

Example:

```go
v1 := app.Group("/api/v1")

v2 := app.Group("/api/v2")
```

Each version can evolve independently without affecting existing APIs.

Zen intentionally avoids introducing specialized versioning abstractions.

---

# Route Middleware

Middleware may be attached directly to individual routes.

Example:

```go
app.Route("/profile").
	Use(AuthMiddleware()).
	Get(profileHandler)
```

Route middleware executes after global and group middleware.

This allows authentication or other behavior to be applied only where required.

---

# Route Metadata

Zen allows routes to include descriptive metadata.

Example:

```go
app.Route("/users").
	Summary("List users").
	Description("Returns all registered users.").
	Tags("Users").
	Name("users.list").
	Version("v1").
	OperationID("listUsers").
	Get(handler)
```

Route metadata is used by Zen's documentation and introspection facilities.

Metadata does not affect request routing.

---

# Route Naming

Routes may be assigned unique names.

Example:

```go
app.Route("/users/:id").
	Name("users.show").
	Get(showUserHandler)
```

Route names provide a stable identifier that is independent of the URL.

This is useful when:

* Looking up registered routes
* Documentation generation
* Metrics and instrumentation
* Administrative tooling

Route names should be unique within an application.

---

# Route Tags

Tags group related endpoints together.

Example:

```go
app.Route("/users").
	Tags("Users").
	Get(listUsersHandler)

app.Route("/users/:id").
	Tags("Users").
	Get(showUserHandler)
```

Tags improve organization when generating documentation.

A route may belong to multiple tags.

Example:

```go
.Tags("Users", "Administration")
```

---

# Route Descriptions

Routes may include both a short summary and a detailed description.

Example:

```go
app.Route("/users/:id").
	Summary("Retrieve a user").
	Description("Returns a single user identified by its unique ID.").
	Get(handler)
```

The summary should be brief.

The description should provide additional context when necessary.

---

# Internal Routes

Routes intended only for internal use may be marked as internal.

Example:

```go
app.Route("/debug/cache").
	Internal().
	Get(handler)
```

Documentation generators may choose to exclude internal routes from public documentation.

---

# Deprecated Routes

Routes that remain available for backward compatibility but should no longer be used can be marked as deprecated.

Example:

```go
app.Route("/v1/users").
	Deprecated().
	Get(handler)
```

Deprecation metadata does not change runtime behavior.

Applications remain responsible for deciding when deprecated endpoints should eventually be removed.

---

# Custom Metadata

Applications may attach additional metadata to routes.

Example:

```go
app.Route("/orders").
	Meta("owner", "Payments Team").
	Meta("service", "Billing").
	Get(handler)
```

Custom metadata allows applications to associate additional information with routes without affecting request handling.

---

# Static File Serving

Zen supports serving static files.

Example:

```go
app.Static("/assets", "./public")
```

Requests such as:

```text
/assets/logo.png
/assets/css/site.css
/assets/app.js
```

are served directly from the configured directory.

Static file serving should generally be used only for development or small deployments.

In larger production environments, static assets are commonly served by dedicated web servers or CDNs.

---

# Route Matching

When an incoming request is received, Zen performs route resolution using its compiled routing structure.

Matching considers:

1. HTTP method
2. Path
3. Dynamic parameters
4. Wildcards

Only a single matching route is selected.

If no route matches, the configured 404 handler is executed.

If the path exists but the HTTP method does not, the configured 405 handler is returned where applicable.

---

# Route Compilation

Route registration occurs during application startup.

Before the application begins serving requests, Zen compiles all registered routes into an optimized runtime structure.

Compilation is performed only once.

Benefits include:

* Faster request dispatch
* Thread-safe initialization
* Consistent runtime behavior
* Reduced per-request work

Applications should complete all route registration before calling:

```go
app.Run(":8080")
```

---

# Route Lookup

Zen provides route introspection APIs that allow applications and tooling to inspect registered routes.

Typical use cases include:

* Documentation generation
* Administrative dashboards
* Debugging
* Route validation
* Metrics

These APIs operate on route definitions rather than incoming requests.

---

# Best Practices

## Register Routes During Startup

Routes should be registered during application initialization.

Avoid registering routes dynamically while the server is running.

---

## Group Related Endpoints

Instead of repeatedly writing:

```text
/api/v1/users
/api/v1/orders
/api/v1/products
```

create a group:

```go
api := app.Group("/api/v1")
```

This improves readability and maintainability.

---

## Keep Handlers Focused

Route handlers should primarily coordinate request processing.

Business logic should live in dedicated services or application packages rather than inside handlers.

---

## Use Metadata

Attach summaries, descriptions, tags, and names to routes.

Metadata improves generated documentation and administrative tooling while making route definitions easier to understand.

---

## Prefer Explicit Routes

Register routes intentionally.

Avoid patterns that automatically discover or register handlers.

Explicit route registration makes applications easier to read, debug, and maintain.

---

## Keep URL Structures Consistent

Use consistent naming conventions throughout an application.

For example:

```text
/users
/users/:id
/users/:id/orders
/users/:id/orders/:orderID
```

A predictable URL structure improves API usability.

---

# Summary

Zen's routing system is designed around explicit registration, predictable behavior, and efficient request dispatch.

Rather than relying on reflection or automatic discovery, routes are declared directly in code, making applications easier to understand and maintain.

By combining route groups, metadata, middleware, dynamic parameters, and a compiled routing structure, Zen provides a routing system that remains lightweight while supporting the needs of production applications.
