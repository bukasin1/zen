# Contributing to Zen

Thank you for your interest in contributing to Zen!

Zen is a lightweight, production-oriented Go backend framework built with the Go standard library. Contributions are welcome, whether they improve the framework itself, documentation, examples, testing, or performance.

---

# Philosophy

Before contributing, please understand the core principles of the project.

Zen aims to be:

* Standard library first
* Production oriented
* Explicit over magic
* Lightweight and understandable
* Correct before convenient
* Performance focused
* Dependency free
* Easy to maintain

The framework intentionally avoids:

* Runtime reflection where practical
* Hidden framework behavior
* Code generation
* Heavy abstractions
* Feature bloat

When proposing changes, please consider whether they align with these principles.

---

# Development Setup

Clone the repository:

```bash
git clone https://github.com/bukasin1/zen.git

cd zen
```

Run all tests:

```bash
go test ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

Run the race detector:

```bash
go test -race ./...
```

---

# Project Structure

```text
docs/           Documentation
examples/       Example applications
internal/       Framework implementation
zen.go          Public API
```

The implementation lives under `internal/`.

The root `zen` package exposes the stable public API used by applications.

---

# Coding Guidelines

Please follow these guidelines when contributing:

* Use only the Go standard library.
* Keep APIs explicit and predictable.
* Prefer simplicity over cleverness.
* Avoid unnecessary abstractions.
* Keep functions focused on a single responsibility.
* Write clear, idiomatic Go.
* Preserve backward compatibility whenever possible.

---

# Testing

New functionality should include appropriate tests.

Before opening a pull request, verify:

```bash
go test ./...
go test -race ./...
```

If your changes affect performance, consider adding or updating benchmarks.

---

# Documentation

If your contribution changes the public API or user-facing behavior, update the relevant documentation and examples where appropriate.

---

# Pull Requests

Please keep pull requests focused on a single logical change.

Smaller pull requests are easier to review and discuss than large, unrelated changes.

---

# Reporting Bugs

If you discover a bug, please use the Bug Report issue template and include a minimal reproducible example whenever possible.

---

# Feature Requests

Before proposing a new feature, consider whether it aligns with Zen's design philosophy.

The project intentionally favors a smaller, well-designed API over a large feature set.

---

# Questions

If you have questions about using Zen, please use the Question issue template.

---

# License

By contributing to Zen, you agree that your contributions will be licensed under the MIT License.
