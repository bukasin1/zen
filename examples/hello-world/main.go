package main

import (
	"net/http"

	"github.com/bukasin1/zen"
)

func main() {
	app := zen.New()

	app.Route("/").
		Get(func(c *zen.Context) {
			c.JSON(http.StatusOK, map[string]string{
				"message": "Welcome to Zen!",
			})
		})

	app.Run(":8080")
}
