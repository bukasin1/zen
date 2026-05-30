package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Danieljosh-uduma/zen/pkg/framework"
	frameworkErrors "github.com/Danieljosh-uduma/zen/pkg/framework/errors"
	"github.com/Danieljosh-uduma/zen/pkg/framework/logger"
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
	fmt.Println("New context...", c.Request.RemoteAddr)
	c.JSON(200, data)
}

func main() {
	app := framework.New()

	cfg := framework.LoadConfigFromEnv()

	err := cfg.Validate()
	if err != nil {
		// log.Fatalf("Failed to load config: %v", err)
	}

	app.SetAppConfig(framework.Config{
		AppName: "Zen",
		HTTP: framework.HTTPConfig{
			Addr: ":8080",
			// ShutdownTimeout: time.Second * 2,
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
	app.Use(framework.CORS(framework.DefaultCORSConfig()))

	validator := &MyValidator{}

	rateLimiter := framework.NewRateLimiter(5, time.Second*30)

	// attach auth parser globally
	app.Use(framework.AuthMiddleware(validator))

	app.Use(framework.Timeout(time.Second * 2))

	// docs endpoints
	app.MountJSONDocs("/docs.json", framework.RouteDocOptions{IncludeInternal: true})
	app.MountHTMLDocs("/docs.html", framework.RouteDocOptions{IncludeInternal: true})
	// app.MountSwagger("/swagger", "zen", "1.0.0")

	// operational routes
	app.RegisterOperationalRoutes()
	// app.Route("/metrics").Get(func(c *framework.Context) {
	// 	snapshot := app.MetricsSnapshot()

	// 	output := framework.FormatPrometheusMetrics(
	// 		snapshot,
	// 	)

	// 	c.Text(http.StatusOK, output)
	// })

	api := app.Group("/api")
	api.Use(TestMiddleware)
	v1 := api.Group("/v1")
	v1.Use(TestMiddleware2)

	api.Get("/health", func(c *framework.Context) {
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

	api.Get("/timeout", func(c *framework.Context) {
		// Simulates a 3-second operation
		// time.Sleep(3 * time.Second)
		ct := &Context{
			Context: c,
		}
		c.AfterResponse(func(c *framework.Context) {
			log.Printf("Response sent: %d, %s, %s", c.StatusCode(), c.Request.URL.Path, c.Request.Method)
		})

		// _ = c.JSON(200, map[string]string{
		// 	"status": "api running",
		// })
		select {
		case <-c.Done():
			if errors.Is(c.Err(), context.DeadlineExceeded) {
				fmt.Println("context error in sample:", c.Err())
				// or http.StatusGatewayTimeout
				c.Error(http.StatusRequestTimeout, "Request timed out", frameworkErrors.ErrRequestTimeout, nil)
			}
		case <-time.After(3 * time.Second):
			fmt.Println("timeout in sample:", c.Err())
			ct.SuccessOK(map[string]string{
				"status": "api running",
			})
		}
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

	app.Route("/home").Use(func(next framework.HandlerFunc) framework.HandlerFunc {
		return func(c *framework.Context) {
			fmt.Println("home middleware")
			next(c)
		}
	}).Get(func(c *framework.Context) {
		c.JSON(200, map[string]string{
			"status": "Welcome to Zen sample use!",
		})
	})

	app.Route("/health").Use(framework.RateLimit(rateLimiter, nil)).Get(func(c *framework.Context) {
		c.JSON(200, map[string]string{
			"status": "server running",
		})
	})

	// protected route
	app.Route("/me").Use(framework.RequireAuth()).Get(func(c *framework.Context) {
		user := c.MustUser()
		c.SuccessOK(user)
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

	type CreateUserDTO struct {
		Email string `json:"email" validate:"required,email" msg:"Invalid email address"`
		Age   int    `json:"age" validate:"required,min=18,max=100" msg:"Age must be between 18 and 100"`
	}

	app.RegisterService("cache", func() interface{} {
		return new(CreateUserDTO)
	})

	app.Route("/users").
		Use(framework.MaxBodySize(55)).
		Post(func(c *framework.Context) {
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

	app.Route("/upload").
		Post(func(c *framework.Context) {

			err := c.ParseMultipartForm(
				10 << 20, // 10 MB
			)

			if err != nil {
				c.BadRequest("invalid multipart form")
				return
			}

			defer c.RemoveMultipartFiles()

			file, header, err := c.FormFile("file")
			if err != nil {
				c.BadRequest("file is required")
				return
			}

			defer file.Close()

			c.SaveUploadedFile(header, header.Filename)

			c.JSON(http.StatusOK, map[string]any{
				"filename": header.Filename,
				"size":     header.Size,
			})
		})

	app.Route("/profile").
		Get(func(c *framework.Context) {

			response := []byte(`{"name":"john"}`)

			fmt.Println("byte respone:", response)

			etag := framework.GenerateETag(
				response,
			)

			if c.IsETagMatch(etag) {
				fmt.Println("returning not modified")
				c.NotModified()
				return
			}

			c.SetETag(etag)

			c.SetCacheControl("private, max-age=60")

			c.JSON(http.StatusOK, map[string]string{
				"name": "john",
			})
			// c.JSON(http.StatusOK, response)
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
