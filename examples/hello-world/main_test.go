package main

import (
	"net/http"
	"testing"

	"github.com/bukasin1/zen"
)

func TestHelloWorld(t *testing.T) {
	app := zen.New()

	app.Route("/").
		Get(func(c *zen.Context) {
			c.JSON(http.StatusOK, map[string]string{
				"message": "Welcome to Zen!",
			})
		})

	res := zen.PerformTestRequest(
		app,
		http.MethodGet,
		"/",
		nil,
		map[string]string{},
	)

	if !zen.HasStatus(res, http.StatusOK) {
		t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
	}

	var body map[string]string

	if err := zen.DecodeJSONResponse(res, &body); err != nil {
		t.Fatal(err)
	}

	if body["message"] != "Welcome to Zen!" {
		t.Fatalf("unexpected message: %q", body["message"])
	}
}
