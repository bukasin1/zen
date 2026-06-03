package framework

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestConcurrentRequests(t *testing.T) {
	app := New()

	app.Route("/test").
		Get(func(c *Context) {
			c.SuccessOK("ok")
		})

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			req := httptest.NewRequest(
				http.MethodGet,
				"/test",
				nil,
			)

			rec := httptest.NewRecorder()

			app.ServeHTTP(rec, req)
		}()
	}

	wg.Wait()
}
