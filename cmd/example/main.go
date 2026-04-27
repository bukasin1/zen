package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Danieljosh-uduma/zen/pkg/framework"
	"github.com/Danieljosh-uduma/zen/pkg/framework/share/logger"
)

type MyValidator struct{}

type AuthUser struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

func (v *MyValidator) Validate(ctx context.Context, token string) (any, error) {
	// decode JWT, validate signature, etc.
	// return user struct

	return AuthUser{
		ID:   "123",
		Role: "admin",
	}, nil
}

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

// custom context
type Context struct {
	*framework.Context
}

func (c *Context) SuccessOK(data any) {
	fmt.Println("New context...")
	c.JSON(200, data)
}

func main() {
	app := framework.New()

	app.SetAppConfig(framework.Config{
		AppName: "Zen",
		HTTP: framework.HTTPConfig{
			Addr: ":8080",
		},
		Log: framework.LogConfig{
			Level:      "debug",
			Pretty:     true,
			EnableJSON: true,
		},
	})

	// System middlewares are auto installed
	// app.Use(framework.Recovery())
	// app.Use(framework.Logger())

	validator := &MyValidator{}

	// attach auth parser globally
	app.Use(framework.AuthMiddleware(validator))

	api := app.Group("/api")
	api.Use(TestMiddleware)
	v1 := api.Group("/v1")
	v1.Use(TestMiddleware2)

	api.Get("/health", func(c *framework.Context) {
		// Simulates a 3-second operation
		time.Sleep(3 * time.Second)
		ct := &Context{
			Context: c,
		}
		c.AfterResponse(func(c *framework.Context) {
			log.Printf("Response sent: %d, %s, %s", c.StatusCode(), c.Request.URL.Path, c.Request.Method)
		})

		ct.SuccessOK(map[string]string{
			"status": "api running",
		})
		// _ = c.JSON(200, map[string]string{
		// 	"status": "api running",
		// })
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

	// protected route
	app.Get("/me",
		framework.RequireAuth()(
			func(c *framework.Context) {
				user := c.MustUser()
				c.SuccessOK(user)
			},
		),
	)

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

	type CreateUserDTO struct {
		Email string `json:"email" validate:"required,email" msg:"Invalid email address"`
		Age   int    `json:"age" validate:"required,min=18,max=100" msg:"Age must be between 18 and 100"`
	}

	app.RegisterService("cache", func() interface{} {
		return new(CreateUserDTO)
	})

	app.Post("/users", func(c *framework.Context) {
		var req CreateUserDTO

		c.MustBindAndValidate(&req)

		// if err := c.BindAndValidate(&req); err != nil {
		// 	c.Fail(400, err.Error())
		// 	return
		// }
		fmt.Println("bindjson body unmarshalled:", req)

		// rawReqBody, _ := c.Body()
		// fmt.Println("body:", rawReqBody, string(rawReqBody))

		// if err := json.Unmarshal(rawReqBody, &req); err != nil {
		// 	c.Error(400, err.Error())
		// 	return
		// }
		// fmt.Println("raw body unmarshalled:", req)

		cache := framework.GetService[*CreateUserDTO](app, "cache")

		c.JSON(201, map[string]any{
			"message": "user created",
			"user":    req,
			"cache":   cache,
		})
	})

	app.Get("/error", func(c *framework.Context) {
		c.AfterResponse(func(c *framework.Context) {
			log.Printf("Response sent in get error: %d, %s, %s", c.StatusCode(), c.Request.URL.Path, c.Request.Method)
		})
		c.Fail(500, "something went wrong")
	})

	app.Get("/panic", func(c *framework.Context) {
		panic("something went wrong get")
	})

	app.Post("/panic", func(c *framework.Context) {
		c.AfterResponse(func(c *framework.Context) {
			log.Printf("Response sent in post panic: %d, %s, %s", c.StatusCode(), c.Request.URL.Path, c.Request.Method)
		})
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

	app.OnStart(func(ctx context.Context) error {
		// modify based on development or production
		app.SetLogger(logger.NewDevConsoleLogger())
		// app.SetLogger(logger.NewConsoleLogger(false))
		return nil
	})

	app.OnShutdown(func(ctx context.Context) error {
		fmt.Println("application shutting down...")
		// close DB, queues, etc.
		return nil
	})

	if err := app.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
