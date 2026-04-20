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

func TestMiddleware2(handler framework.HandlerFunc) framework.HandlerFunc {
	return func(c *framework.Context) {
		fmt.Println("test api 2 middleware")
		handler(c)
	}
}

func main() {
	app := framework.New()

	// System middlewares are auto installed
	// app.Use(framework.Recovery())
	// app.Use(framework.Logger())

	api := app.Group("/api")
	api.Use(TestMiddleware)
	v1 := api.Group("/v1")
	v1.Use(TestMiddleware2)

	api.Get("/health", func(c *framework.Context) {
		_ = c.JSON(200, map[string]string{
			"status": "api running",
		})
	})

	v1.Get("/posts/*", func(c *framework.Context) {
		_ = c.JSON(200, map[string]string{
			"message": "posts endpoint",
			"path":    c.Request.URL.Path,
			"params":  fmt.Sprintf("%#v", c.Params()),
		})
	})

	v1.Get("/posts/:id", func(c *framework.Context) {
		_ = c.JSON(200, map[string]string{
			"postId": c.Param("id"),
			"params": fmt.Sprintf("%#v", c.Params()),
		})
	})

	v1.Get("/users", func(c *framework.Context) {
		_ = c.JSON(200, map[string]string{
			"message": "users endpoint",
		})
	})

	v1.Get("/users/:id", func(c *framework.Context) {
		fmt.Println("user id endpoint")
		c.JSON(200, map[string]string{
			"id":     c.Param("id"),
			"params": fmt.Sprintf("%#v", c.Params()),
		})
	})

	v1.Get("/users/me", func(c *framework.Context) {
		fmt.Println("user me endpoint")
		c.JSON(200, map[string]string{
			"message": "user me endpoint",
		})
	})

	app.Static("/static/*", "./cmd/example/static")
	app.StaticOld("/static2", "./cmd/example/static")
	// app.Static("/", "./cmd/example/public")

	app.Get("/home", func(c *framework.Context) {
		c.JSON(200, map[string]string{
			"status": "Welcome to Zen sample use!",
		})
	})

	app.Get("/health", func(c *framework.Context) {
		c.JSON(200, map[string]string{
			"status": "server running",
		})
	})

	app.Post("/posts/:postId/comments/:commentId", func(c *framework.Context) {
		c.JSON(201, map[string]string{
			"postId":    c.Param("postId"),
			"commentId": c.Param("commentId"),
		})
	})

	app.Get("/users/:id", func(c *framework.Context) {
		c.JSON(200, map[string]string{
			"id": c.Param("id"),
		})
	})

	app.Post("/posts", func(c *framework.Context) {
		c.JSON(200, map[string]string{
			"status": "server running post",
		})
	})

	type CreateUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	app.Post("/users", func(c *framework.Context) {
		var req CreateUserRequest

		if err := c.BindJSON(&req); err != nil {
			c.Fail(400, err.Error())
			return
		}
		fmt.Println("bindjson body unmarshalled:", req)

		// rawReqBody, _ := c.Body()
		// fmt.Println("body:", rawReqBody, string(rawReqBody))

		// if err := json.Unmarshal(rawReqBody, &req); err != nil {
		// 	c.Error(400, err.Error())
		// 	return
		// }
		// fmt.Println("raw body unmarshalled:", req)

		c.JSON(201, map[string]any{
			"message": "user created",
			"user":    req,
		})
	})

	app.Get("/error", func(c *framework.Context) {
		c.Fail(500, "something went wrong")
	})

	app.Get("/panic", func(c *framework.Context) {
		panic("something went wrong get")
	})

	app.Post("/panic", func(c *framework.Context) {
		panic("something went wrong post")
	})

	app.Get("/search", func(c *framework.Context) {
		query := c.Query("q")
		page := c.QueryDefault("page", "1")
		auth := c.Header("Authorization")

		c.JSON(200, map[string]string{
			"query": query,
			"page":  page,
			"auth":  auth,
		})
	})

	log.Println("server starting on :8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
