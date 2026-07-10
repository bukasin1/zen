package zencore

import (
	"fmt"
	"strings"
)

func FormatPrometheusMetrics(
	snapshot MetricsSnapshot,
) string {
	var builder strings.Builder

	builder.WriteString(
		fmt.Sprintf(
			"framework_requests_total %d\n",
			snapshot.TotalRequests,
		),
	)

	builder.WriteString(
		fmt.Sprintf(
			"framework_active_requests %d\n",
			snapshot.ActiveRequests,
		),
	)

	builder.WriteString(
		fmt.Sprintf(
			"framework_panics_total %d\n",
			snapshot.TotalPanics,
		),
	)

	for _, route := range snapshot.Routes {
		builder.WriteString(
			fmt.Sprintf(
				`framework_route_requests_total{method="%s",route="%s"} %d`,
				route.Method,
				route.Route,
				route.RequestCount,
			),
		)

		builder.WriteString("\n")

		builder.WriteString(
			fmt.Sprintf(
				`framework_route_avg_duration_ns{method="%s",route="%s"} %d`,
				route.Method,
				route.Route,
				uint64(route.AverageDuration.Nanoseconds()),
			),
		)

		builder.WriteString("\n")
	}

	return builder.String()
}
