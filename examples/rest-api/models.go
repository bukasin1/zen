package main

// Book represents a single book in the API.
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
}

// CreateBookRequest represents the payload used to create a book.
type CreateBookRequest struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
}

// UpdateBookRequest represents the payload used to update a book.
type UpdateBookRequest struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
}
