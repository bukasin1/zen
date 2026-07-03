# Deployment

This document provides recommendations for deploying Zen applications in production.

Zen is designed to work with standard Go deployment practices. It does not require specialized runtime environments, external processes, or framework-specific infrastructure.

The recommendations in this document are intended to help applications remain reliable, observable, and maintainable in production.

---

# Deployment Philosophy

Zen follows the same philosophy in deployment as it does in application development:

* Prefer simplicity.
* Build on the Go standard library.
* Avoid unnecessary infrastructure dependencies.
* Make operational behavior explicit.
* Favor predictable systems over complex automation.

Applications should remain ordinary Go binaries that can be deployed anywhere Go applications are supported.

---

# Build Applications for Production

Build optimized binaries before deployment.

Example:

```bash
go build -o app .
```

For Linux deployments:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .
```

For ARM64 environments:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o app .
```

Producing a single executable simplifies deployment and reduces operational complexity.

---

# Configuration

Configuration should be supplied through environment variables.

Avoid embedding environment-specific values directly in source code.

Typical configuration includes:

* Server port
* Environment
* Database connection strings
* API keys
* Authentication secrets
* Timeout values
* Logging configuration

This allows the same binary to be deployed across multiple environments without modification.

---

# Graceful Shutdown

Always enable graceful shutdown.

Zen provides built-in graceful shutdown support to allow requests already in progress to complete before the application exits.

Graceful shutdown helps prevent:

* Interrupted client requests
* Partial responses
* Data inconsistency
* Abrupt connection termination

Applications should terminate only after cleanup has completed or the configured shutdown timeout has elapsed.

---

# Health Checks

Applications should expose a health endpoint.

Example:

```text
GET /health
```

Health endpoints allow orchestration systems, load balancers, and monitoring tools to verify that an application is available.

Health checks should remain lightweight and fast.

Avoid performing expensive operations during every health request.

---

# Runtime Information

Zen can expose runtime information for operational visibility.

Example:

```text
GET /runtime/info
```

Typical information includes:

* Go version
* Operating system
* CPU count
* Memory statistics
* Goroutine count
* Uptime

This endpoint is useful during development and production troubleshooting.

Sensitive information should not be exposed publicly.

---

# Metrics

Applications should expose metrics.

Example:

```text
GET /metrics
```

Zen's metrics system provides insight into application behavior while remaining lightweight and dependency-free.

Metrics currently include:

* Request counts
* Route-level metrics
* Request latency

Metrics should be monitored continuously in production environments.

---

# Logging

Structured logging should be enabled in production.

Logs should include enough information to reconstruct request execution and diagnose failures.

Recommended fields include:

* Timestamp
* Request ID
* HTTP method
* Route
* Status code
* Response time
* Client IP
* Error details (when applicable)

Avoid logging sensitive information such as:

* Passwords
* Authentication tokens
* API secrets
* Personal information

Logs should assist operators without compromising security.

---

# Timeouts

Production servers should always configure sensible timeout values.

Recommended areas include:

* Read timeout
* Write timeout
* Idle timeout
* Shutdown timeout

Timeouts help protect applications from slow clients and resource exhaustion.

Applications should also ensure long-running handlers respect request context cancellation.

---

# Request Limits

Limit request body sizes where appropriate.

Large or unbounded request bodies can increase memory usage and expose applications to denial-of-service attacks.

Zen provides middleware for request body size enforcement.

Choose limits appropriate for the application's expected workload.

---

# Rate Limiting

Internet-facing applications should enable rate limiting.

Rate limiting helps protect services against:

* Abuse
* Automated attacks
* Resource exhaustion
* Accidental client overload

Limits should be selected based on application requirements rather than arbitrary values.

---

# Compression

Enable gzip compression for responses where appropriate.

Compression reduces bandwidth usage and improves response times for text-based content.

Binary formats such as images, videos, and compressed archives generally should not be compressed again.

Zen's compression middleware automatically avoids unnecessary compression based on content type and configured thresholds.

---

# HTTP Caching

Applications serving cacheable resources should enable HTTP caching.

Zen provides middleware supporting:

* ETag generation
* Cache-Control headers
* Conditional requests
* `304 Not Modified` responses

Correct caching reduces:

* Network traffic
* Response latency
* Server load

Cache policies should be chosen carefully based on the type of resource being served.

---

# Authentication

Authentication secrets should never be hardcoded into the application.

Instead, supply secrets through environment variables or a secure secret management solution.

Examples include:

* JWT signing keys
* API secrets
* OAuth credentials
* Encryption keys

Secrets should not be committed to version control or written to application logs.

---

# HTTPS

Production applications should always be served over HTTPS.

TLS may be terminated by:

* A reverse proxy
* A load balancer
* The application itself

Whichever approach is chosen, all client traffic should be encrypted.

Redirect plain HTTP traffic to HTTPS where appropriate.

---

# Reverse Proxies

Zen applications commonly run behind reverse proxies such as:

* NGINX
* Caddy
* HAProxy
* Cloud load balancers

A reverse proxy can provide:

* TLS termination
* Compression
* Static asset delivery
* Request buffering
* Rate limiting
* Access logging

Zen works well both with and without a reverse proxy.

---

# Static Assets

For small applications, Zen can serve static files directly.

For larger production deployments, consider serving static assets through:

* A reverse proxy
* A CDN
* An object storage service

Separating static content from application traffic often improves scalability and cache efficiency.

---

# Monitoring

Production services should be monitored continuously.

At a minimum, monitor:

* Request rate
* Error rate
* Latency
* Memory usage
* CPU usage
* Goroutine count
* Application uptime

Monitoring helps detect problems before they affect users.

---

# Backups

If the application manages persistent data, establish a regular backup strategy.

Backup frequency and retention policies depend on business requirements.

Recovery procedures should be tested periodically rather than assumed to work.

---

# Resource Management

Provision resources according to expected workload.

Monitor:

* CPU utilization
* Memory consumption
* Disk usage
* Network throughput

Avoid relying solely on synthetic benchmarks when sizing production infrastructure.

Measure real application behavior under representative workloads.

---

# Deployment Checklist

Before deploying a Zen application to production, verify the following:

* Production build created successfully
* Environment variables configured
* Graceful shutdown enabled
* Health endpoint available
* Metrics endpoint available
* Structured logging enabled
* Appropriate timeout values configured
* Request body limits configured
* Rate limiting enabled where required
* Compression enabled
* HTTP caching configured where appropriate
* HTTPS enabled
* Authentication secrets stored securely
* Monitoring configured
* Backup strategy defined

Completing this checklist helps ensure applications are ready for production workloads.

---

# Troubleshooting

## High Memory Usage

Possible causes include:

* Large request bodies
* Excessive in-memory caching
* Application-level memory leaks
* High request concurrency

Use Go's runtime profiling tools and Zen's runtime information endpoint to investigate memory usage.

---

## Increasing Goroutine Count

A continuously increasing goroutine count may indicate:

* Goroutine leaks
* Blocked operations
* Long-running background tasks
* Ignored context cancellation

Monitor goroutine counts during sustained load and ensure they return to a stable baseline after traffic subsides.

---

## Slow Responses

Investigate:

* Database performance
* External service dependencies
* Long-running handlers
* Timeout configuration
* Middleware overhead

Measure before optimizing.

Changes should be guided by profiling and benchmarking rather than assumptions.

---

## Unexpected Errors

Review:

* Structured application logs
* Request IDs
* Runtime information
* Metrics
* Panic recovery logs

These sources typically provide sufficient information to diagnose production issues.

---

# Summary

Zen applications are deployed like ordinary Go applications.

The framework deliberately avoids specialized deployment requirements and integrates naturally with standard operational practices.

By combining graceful shutdown, structured logging, metrics, health checks, request limits, compression, HTTP caching, and explicit configuration, Zen provides a strong foundation for building reliable production services.

The framework's goal is not to dictate deployment strategy, but to provide the tools necessary for applications to operate predictably and safely in production environments.
