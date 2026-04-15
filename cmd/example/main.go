package main

import (
	"fmt"
	"log"

	"github.com/Danieljosh-uduma/zen/pkg/framework"
)

func TestMiddleware(handler framework.HandlerFunc) framework.HandlerFunc {
	return func(c *framework.Context) {
		fmt.Println("test api middleware")
		handler(c)
	}
}

func main() {
	app := framework.New()

	app.Use(framework.Recovery())
	app.Use(framework.Logger())

	api := app.Group("/api")
	api.Use(TestMiddleware)
	v1 := api.Group("/v1")

	api.Get("/health", func(c *framework.Context) {
		_ = c.JSON(200, map[string]string{
			"status": "api running",
		})
	})

	v1.Get("/users", func(c *framework.Context) {
		_ = c.JSON(200, map[string]string{
			"message": "users endpoint",
		})
	})

	v1.Get("/users/:id", func(c *framework.Context) {
		c.JSON(200, map[string]string{
			"id": c.Param("id"),
		})
	})

	app.Get("/health", func(c *framework.Context) {
		c.JSON(200, map[string]string{
			"status": "server running",
		})
	})

	app.Get("/users/:id", func(c *framework.Context) {
		c.JSON(200, map[string]string{
			"id": c.Param("id"),
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
