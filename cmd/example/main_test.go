package main

import (
	"net/http"
	"testing"

	"github.com/bukasin1/zen/pkg/framework"
)

type Response struct {
	Status string `json:"status"`
}

func TestHealthRoute(t *testing.T) {

	app := framework.New()

	app.Get("/health", func(c *framework.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	rec := framework.PerformTestRequest(
		app,
		http.MethodGet,
		"/health",
		nil,
		nil,
	)

	resp, err := framework.DecodeJSONResponseAs[Response](rec)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "ok" {
		t.Fatalf("expected ok got %s", resp.Status)
	}
}

func TestAuthMiddleware(t *testing.T) {

	ctx, rec := framework.NewTestContext(
		http.MethodGet,
		"/protected",
		nil,
	)

	middleware := framework.RequireAuth()

	called := false

	handler := middleware(func(c *framework.Context) {
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
