package main

import (
	"net/http"
	"testing"

	"github.com/bukasin1/zen"
)

func TestHomePage(t *testing.T) {
	app := newApp()

	res := zen.PerformTestRequest(
		app,
		http.MethodGet,
		"/",
		nil,
		map[string]string{},
	)

	if !zen.HasStatus(res, http.StatusOK) {
		t.Errorf("Expected status OK, got %d", res.Code)
	}
}

func TestHelloEndpoint(t *testing.T) {
	app := newApp()

	res := zen.PerformTestRequest(
		app,
		http.MethodGet,
		"/api/hello",
		nil,
		map[string]string{},
	)

	if !zen.HasStatus(res, http.StatusOK) {
		t.Errorf("Expected status OK, got %d", res.Code)
	}
}

func TestHealthEndpoint(t *testing.T) {
	app := newApp()

	res := zen.PerformTestRequest(
		app,
		http.MethodGet,
		"/health/live",
		nil,
		map[string]string{},
	)

	if !zen.HasStatus(res, http.StatusOK) {
		t.Errorf("Expected status OK, got %d", res.Code)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	app := newApp()

	res := zen.PerformTestRequest(
		app,
		http.MethodGet,
		"/metrics",
		nil,
		map[string]string{},
	)

	if !zen.HasStatus(res, http.StatusOK) {
		t.Errorf("Expected status OK, got %d", res.Code)
	}
}

func TestRuntimeInfoEndpoint(t *testing.T) {
	app := newApp()

	res := zen.PerformTestRequest(
		app,
		http.MethodGet,
		"/runtime/info",
		nil,
		map[string]string{},
	)

	if !zen.HasStatus(res, http.StatusOK) {
		t.Errorf("Expected status OK, got %d", res.Code)
	}
}
