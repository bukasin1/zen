package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/bukasin1/zen"
)

func newTestApp() *zen.App {
	app := zen.New()

	server := &Server{}

	app.Static("/static", staticDir)
	app.Static("/uploads", uploadDir)

	app.Route("/").
		Get(server.index)

	app.Route("/upload").
		Post(server.upload)

	return app
}

func TestIndexPage(t *testing.T) {
	app := newTestApp()

	res := zen.PerformTestRequest(
		app,
		http.MethodGet,
		"/",
		nil,
		map[string]string{},
	)

	if !zen.HasStatus(res, http.StatusOK) {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, res.Code)
	}
}

func TestFileUpload(t *testing.T) {
	app := newTestApp()

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		t.Fatal(err)
	}

	filename := "test-upload.txt"

	defer os.Remove(filepath.Join(uploadDir, filename))

	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.WriteString(part, "Hello from Zen!"); err != nil {
		t.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"/upload",
		&body,
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	res := zen.PerformTestRequestFromRequest(
		app,
		req,
	)

	if !zen.HasStatus(res, http.StatusCreated) {
		t.Fatalf("expected status code %d, got %d", http.StatusCreated, res.Code)
	}

	if _, err := os.Stat(filepath.Join(uploadDir, filename)); err != nil {
		t.Fatalf("expected uploaded file to exist: %v", err)
	}
}
