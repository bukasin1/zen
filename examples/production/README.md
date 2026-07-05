# Production

This example demonstrates the recommended structure for a production-ready Zen application.

It intentionally focuses on operational readiness rather than business logic, showing how to configure, bootstrap, and run a Zen service using only the Go standard library and Zen itself.

The application includes configuration loading, middleware registration, operational endpoints, graceful startup, static asset serving, and automated tests.

---

# Features Demonstrated

This example demonstrates:

* Typed application configuration
* Environment variable configuration
* Global middleware registration
* Request logging
* Panic recovery
* Request timeouts
* Gzip compression
* CORS
* Static file serving
* Operational routes
* JSON API responses
* Graceful server startup
* Integration testing

---

# Project Structure

```text
production/
├── config.go
├── handlers.go
├── main.go
├── main_test.go
├── .env.example
├── public/
│   ├── index.html
│   └── style.css
└── README.md
```

### File Responsibilities

| File           | Responsibility                               |
| -------------- | -------------------------------------------- |
| `main.go`      | Application bootstrap and route registration |
| `config.go`    | Application configuration loading            |
| `handlers.go`  | HTTP request handlers                        |
| `main_test.go` | Integration tests                            |
| `.env.example` | Example environment configuration            |
| `public/`      | Static assets                                |

This separation keeps each file focused on a single responsibility while remaining simple enough for small and medium-sized services.

---

# Running the Example

From the repository root:

```bash
cd examples/production
go run .
```

The application starts on:

```text
http://localhost:8080
```

---

# Available Endpoints

| Method | Endpoint        | Description              |
| ------ | --------------- | ------------------------ |
| GET    | `/`             | Application landing page |
| GET    | `/api/hello`    | Sample JSON endpoint     |
| GET    | `/health`       | Health check             |
| GET    | `/metrics`      | Runtime metrics          |
| GET    | `/runtime/info` | Runtime information      |

---

# Configuration

Configuration is loaded from environment variables.

Example configuration is provided in:

```text
.env.example
```

Copy it before starting the application:

```bash
cp .env.example .env
```

Then adjust values as required for your environment.

---

# Middleware

The example demonstrates registering global middleware during application startup.

Typical production middleware includes:

* Request logging
* Panic recovery
* Request timeout
* Compression
* CORS

Middleware registration occurs before any routes are registered, ensuring consistent behavior across the application.

---

# Operational Routes

Zen provides built-in operational endpoints for production deployments.

These endpoints are intended for monitoring, orchestration platforms, and operational tooling.

| Route           | Purpose             |
| --------------- | ------------------- |
| `/health`       | Health checks       |
| `/metrics`      | Application metrics |
| `/runtime/info` | Runtime diagnostics |

---

# Testing

Run the integration tests:

```bash
go test
```

The tests verify:

* Application startup
* Static page serving
* JSON API responses
* Operational endpoints

---

# Recommended Project Growth

As your application grows, keep business logic separate from HTTP handlers.

A typical production project may evolve into a structure similar to:

```text
production/
├── config.go
├── handlers.go
├── services.go
├── repositories.go
├── models.go
├── middleware.go
├── routes.go
└── main.go
```

Zen does not require any particular project structure. Organize your application in the way that best suits your team while keeping responsibilities clearly separated.

---

# Next Steps

After understanding this example, explore the other examples included with Zen:

* `examples/hello-world`
* `examples/rest-api`
* `examples/file-server`

Together, they cover the core capabilities of the framework, from a minimal application to a production-oriented service.
