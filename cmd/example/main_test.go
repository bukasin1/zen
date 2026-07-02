package main

import (
	"net/http"
	"testing"

	"github.com/bukasin1/zen/pkg/zencore"
)

type Response struct {
	Status string `json:"status"`
}

func TestHealthRoute(t *testing.T) {

	app := zencore.New()

	app.Get("/health", func(c *zencore.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	rec := zencore.PerformTestRequest(
		app,
		http.MethodGet,
		"/health",
		nil,
		nil,
	)

	resp, err := zencore.DecodeJSONResponseAs[Response](rec)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "ok" {
		t.Fatalf("expected ok got %s", resp.Status)
	}
}

func TestAuthMiddleware(t *testing.T) {

	ctx, rec := zencore.NewTestContext(
		http.MethodGet,
		"/protected",
		nil,
	)

	middleware := zencore.RequireAuth()

	called := false

	handler := middleware(func(c *zencore.Context) {
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
