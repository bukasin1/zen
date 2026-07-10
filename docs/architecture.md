# Zen Architecture

This document describes the architectural principles, design decisions, and internal structure of Zen.

It is intended for contributors, maintainers, and developers who want to understand how the framework is designed and why particular implementation decisions were made.

This document describes architecture rather than APIs. For usage examples and API references, see the other documentation in the `docs/` directory.

---

# Design Goals

Zen was built around a simple objective:

> Provide the infrastructure required for building production-ready Go services while remaining as close as possible to the Go standard library.

Every architectural decision is evaluated against this objective.

The framework intentionally favors explicit implementations, predictable behavior, and operational maturity over abstraction and convenience.

---

# Core Principles

The following principles guide every subsystem within Zen.

## Standard Library First

Zen extends the Go standard library rather than replacing it.

Existing standard library types such as `http.Server`, `http.Handler`, `context.Context`, `http.Request`, and `http.ResponseWriter` remain central to the framework.

Where the standard library already provides a good solution, Zen does not introduce another abstraction.

---

## Zero External Runtime Dependencies

Zen depends only on Go's standard library.

This provides several advantages:

* Smaller dependency graph
* Easier upgrades
* Faster builds
* Better long-term maintenance
* Reduced security exposure

Applications built with Zen are not coupled to large third-party ecosystems.

---

## Explicit Over Magic

Framework behavior should always be understandable by reading the source code.

Zen intentionally avoids:

* Reflection-based routing
* Annotation-driven behavior
* Hidden dependency injection
* Automatic controller discovery
* Runtime code generation
* Convention-based execution

Every route, middleware, service, and configuration object is registered explicitly.

---

## Correctness Over Convenience

Convenient APIs should never compromise predictable behavior.

Examples include:

* Centralized response lifecycle management
* Protection against double response writes
* Cooperative request timeouts
* Explicit middleware ordering
* Compile-once application initialization

Reliability is preferred over clever implementation techniques.

---

## Production Realism

Features are added because they solve production problems, not because they are commonly found in frameworks.

Examples include:

* Graceful shutdown
* Structured logging
* Metrics
* Operational endpoints
* Panic recovery
* Rate limiting
* Request body limits
* HTTP caching
* Gzip compression

Each feature exists because it improves production deployments.

---

# Architectural Overview

Zen is intentionally organized into two logical layers.

```text
Application
        │
        ▼
package zen
        │
        ▼
package internal
        │
        ▼
Go Standard Library
```

The public `zen` package forms the stable API exposed to applications.

Framework implementation resides within the `internal` package.

This separation allows implementation details to evolve without affecting public APIs.

---

# Public API

Applications interact exclusively with:

```go
import "github.com/bukasin1/zen"
```

The root package is the public contract of the framework.

Only symbols intentionally re-exported by the root package become part of the supported public API.

This approach prevents accidental exposure of implementation details while allowing internal refactoring without breaking users.

---

# Internal Implementation

The framework implementation exists within a single Go package.

Although source files are organized by subsystem, they intentionally remain part of the same package.

Examples include:

```text
router.go
router_trie.go
router_compile.go

context.go
context_query.go
context_binding.go

middleware.go
middleware_auth.go
middleware_timeout.go

response.go
response_json.go
response_file.go
```

Maintaining a single implementation package avoids unnecessary package fragmentation and eliminates import-cycle complexity while keeping subsystem boundaries clear through file organization.

---

# Repository Layout

The repository follows the structure below.

```text
zen/
│
├── internal/
│
├── cmd/
│
├── docs/
│
├── examples/
│
├── zen.go
│
├── README.md
├── CHANGELOG.md
├── LICENSE
├── go.mod
└── go.sum
```

Each top-level directory has a single responsibility.

The repository intentionally avoids unnecessary nesting and scaffolding.

---

# Request Lifecycle

Every incoming request follows the same execution pipeline.

```text
Incoming Request
        │
        ▼
Application Compile Check
        │
        ▼
Context Creation
        │
        ▼
Global Middleware
        │
        ▼
Route Group Middleware
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
Response Complete
```

Every request creates exactly one `Context`.

The same context instance is passed throughout the lifetime of the request.

Context recreation is intentionally forbidden to ensure request-scoped state remains consistent throughout the pipeline.

---

# Application Compilation

Route registration is separated from request execution.

Before serving the first request, Zen compiles the routing and middleware pipeline into an optimized runtime handler.

Compilation occurs exactly once using `sync.Once`.

This guarantees:

* Thread-safe initialization
* Consistent runtime behavior
* No race conditions during concurrent startup
* Minimal per-request overhead

Applications should complete all route registration before calling `Run()`.

---

# Router Responsibilities

The router has a deliberately narrow responsibility.

Its purpose is to:

* Resolve routes
* Match HTTP methods
* Extract route parameters
* Dispatch the appropriate handler

The router intentionally does **not** perform:

* Authentication
* Validation
* Metrics collection
* Logging
* Documentation generation
* Business logic

Those concerns belong elsewhere in the framework.

---

# Middleware Pipeline

Middleware exists exclusively for cross-cutting concerns.

Examples include:

* Authentication
* Authorization
* Logging
* Recovery
* Rate limiting
* Compression
* Timeouts
* Metrics
* CORS

Middleware should never contain business logic.

Execution order is deterministic.

Global middleware executes before group middleware.

Group middleware executes before route middleware.

Finally, the route handler executes.

This predictable ordering makes request flow easy to reason about and debug.

---

# Context

The `Context` type is the central object passed throughout request processing.

It provides a unified interface for:

* HTTP request access
* Response writing
* Route parameters
* Query parameters
* Headers
* Request-scoped storage
* Authentication information
* Services
* Validation
* Logging

Internally, it also tracks request lifecycle state to protect response integrity.

---

# Response Lifecycle

Zen centralizes all response writing.

Rather than allowing arbitrary writes throughout the framework, responses pass through a single managed lifecycle.

This design protects against:

* Double writes
* Partial responses
* Timeout corruption
* Panic corruption

Once a response has been committed, additional writes are ignored or rejected as appropriate.

This behavior provides predictable HTTP responses even under exceptional conditions.

---

# Error Handling

Errors are represented explicitly.

Recoverable application errors should be returned using the framework's error facilities.

Unexpected failures are handled by the panic recovery middleware.

Panic recovery:

* Prevents server crashes
* Produces controlled HTTP responses
* Logs diagnostic information
* Preserves server availability

Framework panics are treated as exceptional conditions rather than normal control flow.

---

# Concurrency Model

Zen is designed to operate safely under concurrent workloads.

Important concurrency guarantees include:

* Route compilation occurs exactly once.
* Service registration is concurrency-safe.
* Metrics collection is concurrency-safe.
* Request contexts are isolated per request.
* No shared mutable request state exists between requests.

Applications remain responsible for protecting their own shared data structures.

---

# Performance Philosophy

Performance is considered throughout the framework but never at the expense of maintainability.

Zen favors:

* Predictable allocations
* Minimal runtime overhead
* Efficient routing
* Explicit execution paths

Micro-optimizations that significantly reduce code clarity are generally avoided unless supported by measurable performance improvements.

---

# Extensibility

Zen is designed to be extensible without encouraging framework bloat.

Applications may extend the framework through:

* Custom middleware
* Services
* Route metadata
* Configuration
* Logging implementations

Core framework behavior should remain intentionally small and focused.

New framework features are expected to satisfy the following criteria:

* Solve a real production problem
* Integrate naturally with existing APIs
* Avoid introducing unnecessary abstractions
* Preserve backward compatibility whenever practical

---

# API Stability

The public API consists exclusively of symbols exposed through the root `zen` package.

Implementation details within the `internal` package are not considered part of the public API and may evolve as the framework develops.

This separation allows internal improvements without affecting application code.

---

# Contributor Guidelines

When contributing to Zen, contributors should follow these architectural principles.

## Prefer Explicit Code

Code should be understandable without hidden behavior or surprising abstractions.

---

## Avoid Reflection

Reflection should only be introduced when no practical alternative exists and the benefits clearly outweigh the added complexity.

---

## Preserve the Standard Library

Prefer extending standard library types over replacing them.

---

## Keep Components Focused

Each subsystem should have a clearly defined responsibility.

Avoid coupling unrelated concerns together.

---

## Design for Production

Every feature should justify its existence through practical production value rather than framework completeness.

---

## Favor Stability

Backward compatibility should be preserved whenever practical.

Breaking changes should only occur when they significantly improve correctness, maintainability, or long-term framework quality.

---

# Summary

Zen is intentionally conservative.

Rather than maximizing feature count, it focuses on providing a carefully designed foundation for building production-ready Go services.

Its architecture emphasizes:

* Simplicity
* Explicit behavior
* Correctness
* Operational maturity
* Long-term maintainability

These principles guide every architectural decision within the framework and serve as the foundation for future development.
