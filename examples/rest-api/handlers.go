package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bukasin1/zen"
)

// Server holds the application's dependencies.
type Server struct {
	store *BookStore
}

// listBooks returns all books.
func (s *Server) listBooks(c *zen.Context) {
	books := s.store.List()

	c.SuccessOK(books)
}

// getBook returns a single book by ID.
func (s *Server) getBook(c *zen.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.ErrorBadRequest("invalid book id")
		return
	}

	book, err := s.store.Get(id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			c.ErrorNotFound("book not found")
			return
		}

		c.ErrorInternalServer("internal server error")
		return
	}

	c.SuccessOK(book)
}

// createBook creates a new book.
func (s *Server) createBook(c *zen.Context) {
	var req CreateBookRequest

	if err := c.BindJSON(&req); err != nil {
		c.ErrorBadRequest(err.Error())
		return
	}

	book := s.store.Create(req)

	c.SuccessCreated(book)
}

// updateBook updates an existing book.
func (s *Server) updateBook(c *zen.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.ErrorBadRequest("invalid book id")
		return
	}

	var req UpdateBookRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		// c.ErrorBadRequest(err.Error())
		return
	}

	book, err := s.store.Update(id, req)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			c.ErrorNotFound("book not found")
			return
		}

		c.ErrorInternalServer("internal server error")
		return
	}

	c.JSON(http.StatusCreated, book)
}

// deleteBook removes a book.
func (s *Server) deleteBook(c *zen.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.ErrorBadRequest("invalid book id")
		return
	}

	if err := s.store.Delete(id); err != nil {
		if errors.Is(err, ErrBookNotFound) {
			c.ErrorNotFound("book not found")
			return
		}

		c.ErrorInternalServer("internal server error")
		return
	}

	c.NoContent()
}
