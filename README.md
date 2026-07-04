# Zen

> **A production-oriented Go backend framework focused on simplicity, explicitness, and operational maturity.**

Zen is a lightweight backend framework for Go that embraces the Go standard library instead of replacing it. It provides the essential building blocks for building production-ready APIs and backend services while remaining explicit, predictable, and easy to understand.

Unlike many frameworks that rely on reflection, hidden runtime behavior, or large dependency trees, Zen is intentionally conservative. Every subsystem is designed around the principles of correctness, simplicity, and long-term maintainability.

Zen is built for developers who want the productivity of a framework without giving up the transparency and control of the Go standard library.

---

## Why Zen?

Modern backend frameworks often become increasingly complex over time. Features are layered on top of one another, abstractions grow, and eventually understanding the framework becomes almost as difficult as understanding the application itself.

Zen takes a different approach.

Every feature included in the framework must justify its existence. If a problem can already be solved well using the Go standard library, Zen will not replace it simply for the sake of abstraction.

The result is a framework that remains lightweight, explicit, and production-oriented without sacrificing developer experience.

---

## Why Another Go Framework?

The Go standard library already provides an excellent foundation for building HTTP services. Zen does not attempt to replace it.

Instead, Zen builds upon it by providing a carefully selected set of production-oriented capabilities that developers repeatedly implement themselves, including:

- HTTP routing
- Middleware composition
- Request context helpers
- Validation
- Authentication
- Structured logging
- Configuration management
- Metrics
- Operational endpoints
- Testing utilities

Zen intentionally avoids becoming an application platform.

There are no hidden runtime behaviors, reflection-based APIs, code generation requirements, or dependency injection containers.

The goal is simple:

Provide the infrastructure that almost every backend service eventually needs while keeping applications as close as possible to idiomatic Go.

---

## Philosophy

Zen is built around a small set of principles that guide every architectural decision.

* **Standard Library First**
  Whenever possible, Zen builds on top of the Go standard library instead of replacing it.

* **Zero External Runtime Dependencies**
  The framework itself depends only on Go's standard library.

* **Explicit Over Magic**
  Framework behavior should always be obvious from reading the code.

* **Correctness Over Convenience**
  Reliable behavior is preferred over clever shortcuts.

* **Production Realism**
  Every feature should solve real production problems rather than demonstrate framework capabilities.

* **Operational Maturity**
  Logging, metrics, graceful shutdown, testing, and observability are treated as first-class concerns.

* **Performance Conscious**
  Features should introduce minimal overhead while remaining understandable.

* **Maintainable Architecture**
  Simplicity today should not become technical debt tomorrow.

---

## Features

Zen currently includes:

* High-performance trie-based HTTP router
* Route groups and API versioning
* Dynamic route parameters
* Middleware pipeline
* Request context abstraction
* Structured request and application logging
* Request validation
* Authentication and authorization middleware
* Request body binding
* Multipart form and file upload support
* HTTP response helpers
* Centralized error handling
* Panic recovery
* Rate limiting
* Request body size limits
* CORS middleware
* Cooperative request timeout middleware
* Health and operational endpoints
* Metrics and instrumentation foundation
* HTTP caching (ETag and Cache-Control)
* Gzip compression
* Configuration management
* Dependency injection through a service registry
* Comprehensive testing helpers
* Built-in benchmarking utilities
* Route documentation and introspection

---

## Installation

**Requirements**

Zen currently requires:

- Go 1.24 or later
- A supported operating system (Linux, macOS, or Windows)

Zen has no external runtime dependencies.

**Install the latest version:**

```bash
go get github.com/bukasin1/zen
```

**Import the framework:**

```go
import "github.com/bukasin1/zen"
```

---

## Quick Start

Create a simple HTTP server:

```go
package main

import (
	"github.com/bukasin1/zen"
)

func main() {
	app := zen.New()

	app.Route("/").
		Get(func(c *zen.Context) {
			c.JSON(map[string]string{
				"message": "Welcome to Zen!",
			})
		})

	app.Run(":8080")
}
```

Run your application:

```bash
go run .
```

Visit:

```
http://localhost:8080
```

Response:

```json
{
    "message": "Welcome to Zen!"
}
```

---

## Example Applications

The `examples/` directory contains complete, runnable applications demonstrating Zen's features and recommended project structure.

| Example | Description |
|---------|-------------|
| [`examples/basic-web-app`](examples/basic-web-app/) | Basic web application with HTML template rendering |
| [`examples/hello-world`](examples/hello-world/) | Minimal HTTP server |
| [`examples/rest-api`](examples/rest-api/) | REST API with route groups, middleware, validation, and JSON responses |
| `examples/file-upload` | Multipart file uploads |
| `examples/production-ready` | Production-ready web application |
<!-- | `examples/auth` | Authentication and authorization |
| `examples/config` | Configuration loading from environment |
| `examples/graceful-shutdown` | Graceful server shutdown |
| `examples/observability` | Health, runtime information, and metrics endpoints | -->

Each example is self-contained, can be executed independently and can be started with:

```bash
cd examples/<example-name>
go run .
```

The examples are intended to demonstrate production-oriented usage patterns and are kept in sync with the framework as it evolves.
---

## Project Structure

Zen follows a deliberately minimal repository layout.

```
zen/
│
├── internal/          # Framework implementation
├── cmd/               # Executable entry points
├── docs/              # Project documentation
├── examples/          # Runnable examples
│
├── zen.go             # Public API surface
│
├── README.md
├── CHANGELOG.md
├── LICENSE
├── go.mod
└── go.sum
```

Only the root `zen` package forms the public API. Framework implementation details remain inside the `internal` package, allowing the public API to remain stable while internal components evolve independently.

---

## Documentation

Detailed documentation is available in the `docs/` directory.

| Document | Description |
|----------|-------------|
| [`docs/architecture.md`](docs/architecture.md) | Framework architecture and design principles |
| [`docs/routing.md`](docs/routing.md) | Routing system and route groups |
| [`docs/middleware.md`](docs/middleware.md) | Middleware pipeline and custom middleware |
| [`docs/deployment.md`](docs/deployment.md) | Production deployment recommendations |
<!-- | `docs/context.md` | Request context APIs |
| `docs/responses.md` | Response helpers |
| `docs/validation.md` | Request binding and validation |
| `docs/authentication.md` | Authentication and authorization |
| `docs/configuration.md` | Configuration management |
| `docs/services.md` | Service registry |
| `docs/logging.md` | Structured logging |
| `docs/metrics.md` | Metrics and instrumentation |
| `docs/testing.md` | Testing utilities |
| `docs/releases.md` | Versioning and release process |
| `docs/contributing.md` | Contribution guidelines | -->


The README provides a quick introduction. The `docs/` directory contains comprehensive reference documentation for every major subsystem.

> Additional documentation will be added as Zen evolves, with each new subsystem accompanied by dedicated documentation.

---

## Operational Features

Zen includes several built-in operational capabilities designed to make applications production-ready without requiring additional libraries.

### Health Endpoint

Applications can expose a health endpoint suitable for load balancers, orchestrators, and uptime monitoring.

Example:

```text
GET /health
```

---

### Runtime Information

Expose runtime information useful during development and operations.

Example:

```text
GET /runtime/info
```

---

### Metrics

Zen provides a lightweight metrics foundation using only the Go standard library.

Metrics include route-level request counts and latency information while avoiding high-cardinality metric labels.

Example:

```text
GET /metrics
```

No external metrics libraries or exporters are required by the framework.

---

## Design Principles

Zen intentionally avoids several common framework patterns.

### No Reflection-Based Routing

Routes are registered explicitly.

There are no controller scanners or runtime route discovery.

---

### No Hidden Runtime Behavior

Application startup, middleware execution, routing, validation, and request handling all happen explicitly.

There is no hidden lifecycle.

---

### No Code Generation

Zen does not require generated code.

Projects remain ordinary Go applications.

---

### No Runtime Dependency Injection

Dependencies are registered explicitly through the service registry.

There are no reflection containers or runtime object graphs.

---

### No Framework Magic

Reading the source code should always explain what the framework is doing.

Unexpected behavior is considered a design failure.

---

## Performance

Zen is designed to provide predictable performance while remaining understandable and maintainable.

The framework has been validated using benchmarks, race detection, sustained load testing, and memory monitoring.

Validation performed during development includes:

* Unit testing
* Race detection
* Benchmarking
* Sustained HTTP load testing
* Memory usage monitoring
* Goroutine leak verification

The framework is intended for production workloads including:

* REST APIs
* Internal services
* SaaS backends
* Business applications
* Microservices

---

## Compatibility

Zen currently supports:

* Go 1.24+
* Linux
* macOS
* Windows

Because Zen relies exclusively on the Go standard library, applications remain highly portable across supported platforms.

---

## Versioning

Zen follows Semantic Versioning.

Current release:

```text
v0.1.0
```

Version `v0.x` indicates that the framework is considered production-capable but is still evolving toward a long-term stable API.

Breaking API changes may occur before `v1.0.0`, although they will be made conservatively and documented in the changelog.

Once `v1.0.0` is released, public API stability becomes a primary compatibility guarantee.

---

## Roadmap

Future development will continue to focus on production maturity rather than feature quantity.

Areas of ongoing improvement include:

* Performance optimization
* Additional middleware
* Expanded documentation
* Improved observability
* Continued API refinement
* Long-term API stabilization

Features are added only when they solve practical production problems and align with Zen's design philosophy.

---

## Contributing

Contributions are welcome.

If you discover a bug, have an improvement to suggest, or would like to contribute code, please open an issue before beginning significant work so the proposed change can be discussed.

When contributing:

* Keep changes focused and well-scoped.
* Follow existing coding conventions.
* Prefer explicit implementations over abstraction.
* Avoid introducing external runtime dependencies.
* Preserve backward compatibility whenever practical.
* Include tests for new functionality.
* Update documentation when behavior changes.

The project's architecture intentionally favors simplicity over cleverness. Contributions should follow the same philosophy.

---

## License

Zen is released under the MIT License.

See the `LICENSE` file for the complete license text.

---

## Acknowledgements

Zen was built around a simple belief:

> A framework should make building software easier without making the framework itself difficult to understand.

Instead of competing on feature count or abstraction layers, Zen focuses on providing a carefully designed foundation for building reliable, maintainable, and production-ready Go services.

If Zen helps you build something useful, we're glad it could be part of your journey.
