package main

import (
	"log"
	"net/http"

	"github.com/Danieljosh-uduma/zen/pkg/framework"
)

func main() {
	app := framework.New()

	app.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server running"))
	})

	app.Post("/healths", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server running post"))
	})

	log.Println("server starting on :8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
