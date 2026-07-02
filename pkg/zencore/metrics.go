package zencore

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Metrics_Snapshot.go

type RouteMetricsSnapshot struct {
	Route  string
	Method string

	RequestCount uint64

	ActiveRequests int64

	Status1xx uint64
	Status2xx uint64
	Status3xx uint64
	Status4xx uint64
	Status5xx uint64

	PanicCount uint64

	TotalDuration   time.Duration
	AverageDuration time.Duration
	MinDuration     time.Duration
	MaxDuration     time.Duration
}

type MetricsSnapshot struct {
	TotalRequests  uint64
	ActiveRequests int64

	TotalPanics uint64

	Routes map[string]RouteMetricsSnapshot
}

// Metrics_Collector.go

type MetricsCollector interface {
	OnRequestStart(c *Context)

	OnRequestFinish(
		c *Context,
		duration time.Duration,
	)

	OnPanic(
		c *Context,
		panicInfo PanicInfo,
	)

	Snapshot() MetricsSnapshot
}

type routeMetrics struct {
	route  string
	method string

	requestCount atomic.Uint64

	activeRequests atomic.Int64

	status1xx atomic.Uint64
	status2xx atomic.Uint64
	status3xx atomic.Uint64
	status4xx atomic.Uint64
	status5xx atomic.Uint64

	panicCount atomic.Uint64

	totalDurationNs atomic.Uint64

	minDurationNs atomic.Uint64
	maxDurationNs atomic.Uint64
}

type InMemoryMetricsCollector struct {
	totalRequests  atomic.Uint64
	activeRequests atomic.Int64
	totalPanics    atomic.Uint64

	mu sync.RWMutex

	routes map[string]*routeMetrics
}

func NewInMemoryMetricsCollector() *InMemoryMetricsCollector {
	return &InMemoryMetricsCollector{
		routes: make(map[string]*routeMetrics),
	}
}

func (m *InMemoryMetricsCollector) routeKey(
	method string,
	route string,
) string {
	return fmt.Sprintf("%s:%s", method, route)
}

func (m *InMemoryMetricsCollector) getOrCreateRouteMetrics(
	method string,
	route string,
) *routeMetrics {
	key := m.routeKey(method, route)

	m.mu.RLock()
	existing, ok := m.routes[key]
	m.mu.RUnlock()

	if ok {
		return existing
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	existing, ok = m.routes[key]
	if ok {
		return existing
	}

	rm := &routeMetrics{
		route:  route,
		method: method,
	}

	m.routes[key] = rm

	return rm
}

func (m *InMemoryMetricsCollector) OnRequestStart(c *Context) {
	m.totalRequests.Add(1)
	m.activeRequests.Add(1)

	// route := c.FullPath()
	route := c.Request.URL.Path

	if route == "" {
		route = "UNKNOWN"
	}

	rm := m.getOrCreateRouteMetrics(
		c.Request.Method,
		route,
	)

	rm.requestCount.Add(1)
	rm.activeRequests.Add(1)
}

func (m *InMemoryMetricsCollector) OnRequestFinish(
	c *Context,
	duration time.Duration,
) {
	m.activeRequests.Add(-1)

	// route := c.FullPath()
	route := c.Request.URL.Path

	if route == "" {
		route = "UNKNOWN"
	}

	rm := m.getOrCreateRouteMetrics(
		c.Request.Method,
		route,
	)

	rm.activeRequests.Add(-1)

	statusCode := c.StatusCode()

	switch {
	case statusCode >= 100 && statusCode < 200:
		rm.status1xx.Add(1)

	case statusCode >= 200 && statusCode < 300:
		rm.status2xx.Add(1)

	case statusCode >= 300 && statusCode < 400:
		rm.status3xx.Add(1)

	case statusCode >= 400 && statusCode < 500:
		rm.status4xx.Add(1)

	case statusCode >= 500:
		rm.status5xx.Add(1)
	}

	durationNs := uint64(duration.Nanoseconds())

	rm.totalDurationNs.Add(durationNs)

	m.updateMinDuration(rm, durationNs)
	m.updateMaxDuration(rm, durationNs)
}

func (m *InMemoryMetricsCollector) updateMinDuration(
	rm *routeMetrics,
	durationNs uint64,
) {
	for {
		current := rm.minDurationNs.Load()

		if current != 0 && current <= durationNs {
			return
		}

		if rm.minDurationNs.CompareAndSwap(
			current,
			durationNs,
		) {
			return
		}
	}
}

func (m *InMemoryMetricsCollector) updateMaxDuration(
	rm *routeMetrics,
	durationNs uint64,
) {
	for {
		current := rm.maxDurationNs.Load()

		if current >= durationNs {
			return
		}

		if rm.maxDurationNs.CompareAndSwap(
			current,
			durationNs,
		) {
			return
		}
	}
}

func (m *InMemoryMetricsCollector) OnPanic(
	c *Context,
	panicInfo PanicInfo,
) {
	m.totalPanics.Add(1)

	// route := c.FullPath()
	route := c.Request.URL.Path

	if route == "" {
		route = "UNKNOWN"
	}

	rm := m.getOrCreateRouteMetrics(
		c.Request.Method,
		route,
	)

	rm.panicCount.Add(1)
}

func (m *InMemoryMetricsCollector) Snapshot() MetricsSnapshot {
	snapshot := MetricsSnapshot{
		TotalRequests:  m.totalRequests.Load(),
		ActiveRequests: m.activeRequests.Load(),
		TotalPanics:    m.totalPanics.Load(),
		Routes:         make(map[string]RouteMetricsSnapshot),
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for key, rm := range m.routes {
		requestCount := rm.requestCount.Load()

		totalDurationNs := rm.totalDurationNs.Load()

		var averageDuration time.Duration

		if requestCount > 0 {
			averageDuration = time.Duration(
				totalDurationNs / requestCount,
			)
		}

		snapshot.Routes[key] = RouteMetricsSnapshot{
			Route:  rm.route,
			Method: rm.method,

			RequestCount: requestCount,

			ActiveRequests: rm.activeRequests.Load(),

			Status1xx: rm.status1xx.Load(),
			Status2xx: rm.status2xx.Load(),
			Status3xx: rm.status3xx.Load(),
			Status4xx: rm.status4xx.Load(),
			Status5xx: rm.status5xx.Load(),

			PanicCount: rm.panicCount.Load(),

			TotalDuration:   time.Duration(totalDurationNs),
			AverageDuration: averageDuration,
			MinDuration:     time.Duration(rm.minDurationNs.Load()),
			MaxDuration:     time.Duration(rm.maxDurationNs.Load()),
		}
	}

	return snapshot
}

// --------- APP Methods ---------

func (a *App) SetMetricsCollector(
	collector MetricsCollector,
) {
	if collector == nil {
		return
	}

	a.metricsCollector = collector
}

func (a *App) MetricsSnapshot() MetricsSnapshot {
	if a.metricsCollector == nil {
		return MetricsSnapshot{}
	}

	return a.metricsCollector.Snapshot()
}

// RegisterMetricsRoute registers the metrics endpoint.
//
// It adds one endpoint:
//   - GET /metrics - Returns metrics in Prometheus format
//
// Note: The metrics endpoint is registered as an internal route and is not
// visible in the generated OpenAPI documentation.
func (a *App) RegisterMetricsRoute() {
	a.Route("/metrics").
		Internal().
		Summary("Metrics endpoint").
		Get(func(c *Context) {

			if a.metricsCollector == nil {
				c.Text(http.StatusServiceUnavailable, "metrics disabled")
				return
			}

			snapshot := a.metricsCollector.Snapshot()

			output := FormatPrometheusMetrics(snapshot)

			c.SetHeader("Content-Type", "text/plain; version=0.0.4")

			c.Text(http.StatusOK, output)
		})
}
