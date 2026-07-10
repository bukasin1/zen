# Changelog

All notable changes to this project will be documented in this file.

The format is based on **Keep a Changelog**, and this project follows **Semantic Versioning**.

---

## [0.1.0] - 2026-07-05

### Added

#### Core Framework

* HTTP application framework
* Custom router with route groups and versioning
* Dynamic route parameters
* Request context abstraction
* Middleware pipeline
* Centralized response handling
* Panic recovery and classification
* Graceful shutdown
* Configuration system
* Environment variable helpers
* Service registry
* Authentication and authorization middleware
* Rate limiting
* Request body size limiting
* CORS middleware
* Cooperative request timeout middleware

#### Responses

* Standardized JSON responses
* Redirect helpers
* File serving helpers
* Download helpers

#### Validation

* Request binding
* Struct validation
* Validation tags
* Application error integration

#### Logging

* Console logger
* Development console logger
* Request logging

#### Observability

* Metrics collection
* Runtime information
* Health endpoints
* Route instrumentation

#### HTTP Features

* Multipart form support
* Gzip compression
* HTTP caching support
* ETag handling
* Cache-Control support

#### Testing

* HTTP testing helpers
* JSON testing helpers
* Benchmark helpers

#### Documentation

* Route introspection
* HTML documentation export
* JSON documentation export

#### Examples

* Hello World
* REST API
* File Server
* Production Application

#### Documentation

* Framework architecture
* Routing guide
* Middleware guide
* Deployment guide

### Performance

* Router optimization using trie-based routing
* Route template metrics to prevent cardinality explosion
* Unified production and testing request pipeline

### Stability

* Concurrency-safe application compilation using `sync.Once`
* Race detector validation
* Sustained load testing
* Memory leak validation
* Goroutine leak validation

### Release

Initial public release of **Zen**.
