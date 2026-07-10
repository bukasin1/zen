# Hello World

This example demonstrates the smallest possible Zen application.

It registers a single route that returns a JSON response.

## Running

From the repository root:

```bash
cd examples/hello-world
go run .
```

The server starts on:

```text
http://localhost:8080
```

Request:

```http
GET /
```

Response:

```json
{
    "message": "Welcome to Zen!"
}
```

## What this example demonstrates

* Creating a new Zen application
* Registering a route
* Returning a JSON response
* Running the HTTP server

This example intentionally keeps the application as small as possible and serves as the recommended starting point for new Zen users.
