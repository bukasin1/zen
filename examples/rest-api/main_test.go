package main

import (
	"net/http"
	"testing"

	"github.com/bukasin1/zen"
)

func newTestApp() *zen.App {
	app := zen.New()

	server := &Server{
		store: NewBookStore(),
	}

	api := app.Group("/api")

	books := api.Route("/books")

	books.Get(server.listBooks)
	books.Post(server.createBook)

	book := api.Route("/books/:id")

	book.Get(server.getBook)
	book.Put(server.updateBook)
	book.Delete(server.deleteBook)

	return app
}

func TestListBooks(t *testing.T) {
	app := newTestApp()

	res := zen.PerformTestRequest(
		app,
		http.MethodGet,
		"/api/books",
		nil,
		map[string]string{},
	)

	if !zen.HasStatus(res, http.StatusOK) {
		t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
	}
}

func TestGetBook(t *testing.T) {
	app := newTestApp()

	res := zen.PerformTestRequest(
		app,
		http.MethodGet,
		"/api/books/1",
		nil,
		map[string]string{},
	)

	if !zen.HasStatus(res, http.StatusOK) {
		t.Fatalf("expected status %d got %d", http.StatusOK, res.Code)
	}
}

func TestCreateBook(t *testing.T) {
	app := newTestApp()

	// body := bytes.NewBufferString(`{
	// 	"title":"Domain-Driven Design",
	// 	"author":"Eric Evans"
	// }`)
	body := []byte(`{
		"title":"Domain-Driven Design",
		"author":"Eric Evans"
	}`)

	res := zen.PerformTestRequest(
		app,
		http.MethodPost,
		"/api/books",
		body,
		map[string]string{},
	)

	if !zen.HasStatus(res, http.StatusCreated) {
		t.Fatalf("expected status %d got %d", http.StatusCreated, res.Code)
	}
}

func TestUpdateBook(t *testing.T) {
	app := newTestApp()

	// body := bytes.NewBufferString(`{
	// 	"title":"Updated Title",
	// 	"author":"Updated Author"
	// }`)
	body := []byte(`{
		"title":"Updated Title",
		"author":"Updated Author"
	}`)

	res := zen.PerformTestRequest(
		app,
		http.MethodPut,
		"/api/books/1",
		body,
		map[string]string{},
	)

	if !zen.HasStatus(res, http.StatusCreated) {
		t.Fatalf("expected status %d got %d", http.StatusCreated, res.Code)
	}
}

func TestDeleteBook(t *testing.T) {
	app := newTestApp()

	res := zen.PerformTestRequest(
		app,
		http.MethodDelete,
		"/api/books/1",
		nil,
		map[string]string{},
	)

	if !zen.HasStatus(res, http.StatusNoContent) {
		t.Fatalf("expected status %d got %d", http.StatusNoContent, res.Code)
	}
}
