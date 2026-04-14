package main

import (
	"log"

	"github.com/Danieljosh-uduma/zen/pkg/framework"
)

func main() {
	app := framework.New()

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

	log.Println("server starting on :8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
