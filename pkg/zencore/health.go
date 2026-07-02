package zencore

import "net/http"

type HealthResponse struct {
	Status string `json:"status"`
}

// RegisterHealthRoutes registers health check routes.
//
// It adds two endpoints:
//   - GET /health/live - Always returns 200 OK (liveness)
//   - GET /health/ready - Returns 200 OK only when the app is ready and not shutting down (readiness)
//
// Note: health check routes are registered as internal routes and are not
// visible in the generated OpenAPI documentation.
func (a *App) RegisterHealthRoutes() {
	a.Route("/health/live").
		Internal().
		Summary("Liveness probe").
		Get(func(c *Context) {
			c.JSON(http.StatusOK, HealthResponse{
				Status: "alive",
			})
		})

	a.Route("/health/ready").
		Internal().
		Summary("Readiness probe").
		Get(func(c *Context) {

			if a.runtimeState == nil {
				c.JSON(http.StatusServiceUnavailable, HealthResponse{
					Status: "not_ready",
				})
				return
			}

			if a.runtimeState.IsShuttingDown() {
				c.JSON(http.StatusServiceUnavailable, HealthResponse{
					Status: "shutting_down",
				})
				return
			}

			if !a.runtimeState.IsReady() {
				c.JSON(http.StatusServiceUnavailable, HealthResponse{
					Status: "not_ready",
				})
				return
			}

			c.JSON(http.StatusOK, HealthResponse{
				Status: "ready",
			})
		})
}
