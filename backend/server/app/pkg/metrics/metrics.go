package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestTotal counts the total number of HTTP requests
	HTTPRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTPRequestDuration tracks the duration of HTTP requests in seconds
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets, // Default buckets: .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10
		},
		[]string{"method", "endpoint"},
	)

	// HTTPRequestsInFlight tracks the number of HTTP requests currently being processed
	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)

	// DatabaseQueryDuration tracks the duration of database queries in seconds
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		},
		[]string{"operation", "table"},
	)

	// DatabaseConnectionsActive tracks the number of active database connections
	DatabaseConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
	)

	// ExternalAPIRequestsTotal counts the total number of external API requests
	ExternalAPIRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "external_api_requests_total",
			Help: "Total number of external API requests",
		},
		[]string{"service", "endpoint", "status"},
	)

	// ExternalAPIDuration tracks the duration of external API calls in seconds
	ExternalAPIDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "external_api_duration_seconds",
			Help:    "External API call duration in seconds",
			Buckets: []float64{.01, .05, .1, .25, .5, 1, 2.5, 5, 10, 30},
		},
		[]string{"service", "endpoint"},
	)
)

// RecordHTTPRequest records an HTTP request metric
func RecordHTTPRequest(method, endpoint, status string, duration float64) {
	HTTPRequestTotal.WithLabelValues(method, endpoint, status).Inc()
	HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// IncHTTPRequestsInFlight increments the in-flight requests counter
func IncHTTPRequestsInFlight() {
	HTTPRequestsInFlight.Inc()
}

// DecHTTPRequestsInFlight decrements the in-flight requests counter
func DecHTTPRequestsInFlight() {
	HTTPRequestsInFlight.Dec()
}

// RecordDatabaseQuery records a database query metric
func RecordDatabaseQuery(operation, table string, duration float64) {
	DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration)
}

// SetDatabaseConnectionsActive sets the number of active database connections
func SetDatabaseConnectionsActive(count float64) {
	DatabaseConnectionsActive.Set(count)
}

// RecordExternalAPIRequest records an external API request metric
func RecordExternalAPIRequest(service, endpoint, status string, duration float64) {
	ExternalAPIRequestsTotal.WithLabelValues(service, endpoint, status).Inc()
	ExternalAPIDuration.WithLabelValues(service, endpoint).Observe(duration)
}
