# File Server

This example demonstrates file uploads, static file serving, and server-rendered HTML using Zen and the Go standard library.

It shows how to accept multipart form uploads, save uploaded files to disk, serve them back over HTTP, and render a simple HTML page listing all uploaded files.

---

# Project Structure

```text
file-server/
├── handlers.go
├── main.go
├── main_test.go
├── public/
│   └── style.css
├── templates/
│   └── index.html
├── uploads/
└── README.md
```

---

# Features Demonstrated

This example demonstrates:

* Multipart form uploads
* Saving uploaded files
* Static file serving
* Server-side HTML rendering using `html/template`
* Serving uploaded files
* Zen response helpers
* Unit testing

---

# Running the Example

From the repository root:

```bash
cd examples/file-server
go run .
```

Open your browser:

```text
http://localhost:8080
```

Choose a file and click **Upload File**.

Uploaded files are immediately available at:

```text
http://localhost:8080/uploads/<filename>
```

and are listed on the application's home page.

---

# Running the Tests

```bash
go test
```

The tests verify:

* The upload page is served successfully.
* Multipart uploads succeed.
* Uploaded files are written to disk.

---

# Notes

This example intentionally stores uploaded files on the local filesystem.

It is designed to demonstrate Zen's multipart helpers and static file serving while relying only on Go's standard library for HTML templating.

In a production application, you might replace the local filesystem with cloud object storage or another persistence layer, while keeping the HTTP handling code largely unchanged.
