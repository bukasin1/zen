# REST API

This example demonstrates how to build a small RESTful API with Zen.

Unlike the `hello-world` example, this project is intentionally organized into multiple files to illustrate a structure that scales naturally as an application grows.

The example implements a simple in-memory Book API and showcases many of the features developers use every day.

---

# Project Structure

```text
rest-api/
├── handlers.go
├── main.go
├── main_test.go
├── models.go
├── README.md
└── store.go
```

| File           | Purpose                                                                |
| -------------- | ---------------------------------------------------------------------- |
| `main.go`      | Application bootstrap, middleware registration, and route registration |
| `handlers.go`  | HTTP handlers                                                          |
| `models.go`    | Request and response models                                            |
| `store.go`     | Thread-safe in-memory data store                                       |
| `main_test.go` | Example tests using Zen's testing utilities                            |

---

# Features Demonstrated

This example demonstrates:

* Route groups
* RESTful routing
* Path parameters
* Request binding
* Request validation
* JSON responses
* Proper HTTP status codes
* Global middleware
* Route metadata
* Thread-safe application state
* Unit testing

---

# Running the Example

From the repository root:

```bash
cd examples/rest-api
go run .
```

The server starts on:

```text
http://localhost:8080
```

---

# Available Endpoints

| Method | Endpoint         | Description            |
| ------ | ---------------- | ---------------------- |
| GET    | `/api/books`     | List all books         |
| GET    | `/api/books/:id` | Retrieve a single book |
| POST   | `/api/books`     | Create a book          |
| PUT    | `/api/books/:id` | Update a book          |
| DELETE | `/api/books/:id` | Delete a book          |

---

# Example Requests

## List Books

```http
GET /api/books
```

---

## Retrieve a Book

```http
GET /api/books/1
```

---

## Create a Book

```http
POST /api/books
Content-Type: application/json

{
    "title": "Learning Go",
    "author": "Jon Bodner"
}
```

---

## Update a Book

```http
PUT /api/books/1
Content-Type: application/json

{
    "title": "The Go Programming Language",
    "author": "Alan A. A. Donovan"
}
```

---

## Delete a Book

```http
DELETE /api/books/1
```

---

# Testing

Run the example tests:

```bash
go test
```

The tests demonstrate Zen's testing helpers for exercising HTTP handlers without starting a real server.

---

# Notes

The application intentionally stores data in memory.

This keeps the example focused on demonstrating Zen rather than introducing database configuration or external dependencies.

The structure shown here provides a good starting point for small and medium-sized applications. As an application grows, additional packages such as services, repositories, and configuration can be introduced without changing the overall organization demonstrated by this example.
