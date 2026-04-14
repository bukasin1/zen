package main

import (
	"log"

	"github.com/Danieljosh-uduma/zen/pkg/framework"
)

func main() {
	app := framework.New()

	app.Use(framework.Recovery())
	app.Use(framework.Logger())

	app.Get("/health", func(c *framework.Context) {
		c.JSON(200, map[string]string{
			"status": "server running",
		})
	})

	app.Post("/healths", func(c *framework.Context) {
		c.JSON(200, map[string]string{
			"status": "server running post",
		})
	})

	app.Get("/error", func(c *framework.Context) {
		c.Error(500, "something went wrong")
	})

	app.Get("/panic", func(c *framework.Context) {
		panic("something went wrong get")
	})

	app.Post("/panic", func(c *framework.Context) {
		panic("something went wrong post")
	})

	log.Println("server starting on :8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
