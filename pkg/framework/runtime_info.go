package framework

import (
	"net/http"
	"runtime"
	"time"
)

type RuntimeInfo struct {
	GoVersion         string    `json:"go_version"`
	NumCPU            int       `json:"num_cpu"`
	NumGoroutine      int       `json:"num_goroutine"`
	MemoryAllocBytes  uint64    `json:"memory_alloc_bytes"`
	MemorySystemBytes uint64    `json:"memory_system_bytes"`
	Uptime            string    `json:"uptime"`
	StartedAt         time.Time `json:"started_at"`
}

// RegisterRuntimeRoutes registers the runtime diagnostics route.
//
// It adds one endpoint:
//   - GET /runtime/info - Returns runtime diagnostics (Go version, CPU count, memory usage, etc.)
//
// Note: The runtime diagnostics route is registered as an internal route and is not
// visible in the generated OpenAPI documentation.
func (a *App) RegisterRuntimeRoutes() {
	a.Route("/runtime/info").
		Internal().
		Summary("Runtime diagnostics").
		Get(func(c *Context) {

			var memory runtime.MemStats

			runtime.ReadMemStats(&memory)

			startedAt := time.Time{}
			uptime := time.Duration(0)

			if a.runtimeState != nil {
				startedAt = a.runtimeState.StartedAt()
				uptime = a.runtimeState.Uptime()
			}

			c.JSON(http.StatusOK, RuntimeInfo{
				GoVersion:         runtime.Version(),
				NumCPU:            runtime.NumCPU(),
				NumGoroutine:      runtime.NumGoroutine(),
				MemoryAllocBytes:  memory.Alloc,
				MemorySystemBytes: memory.Sys,
				Uptime:            uptime.String(),
				StartedAt:         startedAt,
			})
		})
}

// RegisterOperationalRoutes registers all operational routes (health, runtime, metrics).
//
// It adds four endpoints:
//   - GET /health/live - Always returns 200 OK (liveness)
//   - GET /health/ready - Returns 200 OK only when the app is ready and not shutting down (readiness)
//   - GET /runtime/info - Returns runtime diagnostics (Go version, CPU count, memory usage, etc.)
//   - GET /metrics - Returns metrics in Prometheus format
//
// Note: All operational routes are registered as internal routes and are not
// visible in the generated OpenAPI documentation.
func (a *App) RegisterOperationalRoutes() {
	a.RegisterHealthRoutes()
	a.RegisterRuntimeRoutes()
	a.RegisterMetricsRoute()
}
