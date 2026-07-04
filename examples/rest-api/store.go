package main

import (
	"errors"
	"sync"
)

// ErrBookNotFound is returned when a requested book does not exist.
var ErrBookNotFound = errors.New("book not found")

// BookStore provides a concurrency-safe in-memory store for books.
type BookStore struct {
	mu     sync.RWMutex
	nextID int
	books  map[int]Book
}

// NewBookStore creates a new book store.
// It initializes the store with some sample books.
func NewBookStore() *BookStore {
	store := &BookStore{
		nextID: 3,
		books: map[int]Book{
			1: {
				ID:     1,
				Title:  "The Go Programming Language",
				Author: "Alan A. A. Donovan",
			},
			2: {
				ID:     2,
				Title:  "Introducing Go",
				Author: "Caleb Doxsey",
			},
		},
	}

	return store
}

// List returns all books.
func (s *BookStore) List() []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	books := make([]Book, 0, len(s.books))

	for _, book := range s.books {
		books = append(books, book)
	}

	return books
}

// Get returns a book by its ID.
func (s *BookStore) Get(id int) (Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, ok := s.books[id]
	if !ok {
		return Book{}, ErrBookNotFound
	}

	return book, nil
}

// Create stores a new book and returns it.
func (s *BookStore) Create(req CreateBookRequest) Book {
	s.mu.Lock()
	defer s.mu.Unlock()

	book := Book{
		ID:     s.nextID,
		Title:  req.Title,
		Author: req.Author,
	}

	s.books[book.ID] = book
	s.nextID++

	return book
}

// Update replaces an existing book.
func (s *BookStore) Update(id int, req UpdateBookRequest) (Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book, ok := s.books[id]
	if !ok {
		return Book{}, ErrBookNotFound
	}

	book.Title = req.Title
	book.Author = req.Author

	s.books[id] = book

	return book, nil
}

// Delete removes a book from the store.
func (s *BookStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[id]; !ok {
		return ErrBookNotFound
	}

	delete(s.books, id)

	return nil
}
