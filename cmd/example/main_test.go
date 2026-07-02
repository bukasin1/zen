package main

import (
	"net/http"
	"testing"

	"github.com/bukasin1/zen"
)

type Response struct {
	Status string `json:"status"`
}

func TestHealthRoute(t *testing.T) {

	app := zen.New()

	app.Get("/health", func(c *zen.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	rec := zen.PerformTestRequest(
		app,
		http.MethodGet,
		"/health",
		nil,
		nil,
	)

	resp, err := zen.DecodeJSONResponseAs[Response](rec)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "ok" {
		t.Fatalf("expected ok got %s", resp.Status)
	}
}

func TestAuthMiddleware(t *testing.T) {

	ctx, rec := zen.NewTestContext(
		http.MethodGet,
		"/protected",
		nil,
	)

	middleware := zen.RequireAuth()

	called := false

	handler := middleware(func(c *zen.Context) {
		called = true
	})

	handler(ctx)

	if called {
		t.Fatal("handler should not execute")
	}

	if rec.Code != http.StatusUnauthorized {
		t.Fatal("expected 401")
	}
}
