# Zen Release Checklist

Use this checklist before publishing every new release.

---

# Repository

* [ ] Repository is public.
* [ ] Repository description is up to date.
* [ ] Repository topics are configured.
* [ ] Default branch is `main`.
* [ ] Branch protection is configured.
* [ ] GitHub issue templates are present.
* [ ] Pull request template is present.
* [ ] LICENSE is present.
* [ ] CHANGELOG.md has been updated.
* [ ] README.md reflects the current release.

---

# Documentation

* [ ] Core documentation has been reviewed.
* [ ] Example applications compile.
* [ ] Public API examples are up to date.

---

# Code Quality

* [ ] All code has been formatted.

```bash
go fmt ./...
```

* [ ] Static analysis passes.

```bash
go vet ./...
```

* [ ] All tests pass.

```bash
go test ./...
```

* [ ] Race detector passes.

```bash
go test -race ./...
```

---

# Framework Validation

* [ ] Public API has been reviewed.
* [ ] No unintended exported identifiers exist.
* [ ] Public API documentation has been reviewed.
* [ ] Backward compatibility has been considered.

---

# Examples

* [ ] Hello World example works.
* [ ] REST API example works.
* [ ] File Server example works.
* [ ] Production example works.

---

# Operational Endpoints

Verify the following endpoints:

* [ ] `/health/live`
* [ ] `/health/ready`
* [ ] `/runtime/info`
* [ ] `/metrics`

---

# Release

* [ ] Version number updated.
* [ ] CHANGELOG finalized.
* [ ] Git tag created.
* [ ] GitHub Release drafted.
* [ ] GitHub Release published.

---

# Post Release

* [ ] Verify `go install` / `go get` works.
* [ ] Verify example applications still build against the released version.
* [ ] Announce the release.
