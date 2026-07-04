package main

import (
	"log"

	"github.com/bukasin1/zen"
)

func main() {
	app := zen.New()

	// Global middleware.
	app.Use(
		zen.RequestLogger(),
		zen.Recovery(),
	)

	server := &Server{
		store: NewBookStore(),
	}

	api := app.Group("/api")

	booksRoute := api.Route("/books")
	bookRoute := api.Route("/books/:id")

	booksRoute.
		Summary("List books").
		Description("Returns all books.").
		Tags("Books").
		Get(server.listBooks)

	booksRoute.
		Summary("Create a book").
		Description("Creates a new book.").
		Tags("Books").
		Post(server.createBook)

	bookRoute.
		Summary("Get a book").
		Description("Returns a single book by its ID.").
		Tags("Books").
		Get(server.getBook)

	bookRoute.
		Summary("Update a book").
		Description("Updates an existing book.").
		Tags("Books").
		Put(server.updateBook)

	bookRoute.
		Summary("Delete a book").
		Description("Deletes a book.").
		Tags("Books").
		Delete(server.deleteBook)

	log.Println("REST API listening on :8080")

	if err := app.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
