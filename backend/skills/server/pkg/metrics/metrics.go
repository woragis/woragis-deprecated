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

	// HealthCheckTotal counts the total number of health check requests
	HealthCheckTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "health_check_total",
			Help: "Total number of health check requests",
		},
		[]string{"check_type", "status"},
	)

	// HealthCheckDuration tracks the duration of health checks in seconds
	HealthCheckDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "health_check_duration_seconds",
			Help:    "Health check duration in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"check_type"},
	)

	// HealthCheckStatus tracks the current health status (1 = healthy, 0 = unhealthy)
	HealthCheckStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "health_check_status",
			Help: "Current health check status (1 = healthy, 0 = unhealthy)",
		},
		[]string{"check_name", "check_type"},
	)

	// Business Metrics for Auth Service

	// UserRegistrationsTotal counts the total number of user registrations
	UserRegistrationsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_user_registrations_total",
			Help: "Total number of user registrations",
		},
	)

	// UserLoginsTotal counts the total number of user logins
	UserLoginsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_user_logins_total",
			Help: "Total number of user logins",
		},
		[]string{"status"}, // status: success, failure
	)

	// TokenRefreshesTotal counts the total number of token refreshes
	TokenRefreshesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_token_refreshes_total",
			Help: "Total number of token refreshes",
		},
		[]string{"status"}, // status: success, failure
	)

	// TokenRevocationsTotal counts the total number of token revocations
	TokenRevocationsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_token_revocations_total",
			Help: "Total number of token revocations",
		},
	)

	// EmailVerificationsTotal counts the total number of email verifications
	EmailVerificationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_email_verifications_total",
			Help: "Total number of email verifications",
		},
		[]string{"status"}, // status: success, failure
	)

	// PasswordChangesTotal counts the total number of password changes
	PasswordChangesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_password_changes_total",
			Help: "Total number of password changes",
		},
		[]string{"status"}, // status: success, failure
	)

	// RequestTimeoutsTotal counts the total number of request timeouts
	RequestTimeoutsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_request_timeouts_total",
			Help: "Total number of request timeouts",
		},
		[]string{"endpoint"},
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

// RecordHealthCheck records a health check metric
func RecordHealthCheck(checkType, status string, duration float64) {
	HealthCheckTotal.WithLabelValues(checkType, status).Inc()
	HealthCheckDuration.WithLabelValues(checkType).Observe(duration)
}

// SetHealthCheckStatus sets the health check status metric
func SetHealthCheckStatus(checkName, checkType string, isHealthy bool) {
	status := 0.0
	if isHealthy {
		status = 1.0
	}
	HealthCheckStatus.WithLabelValues(checkName, checkType).Set(status)
}

// Business Metrics Functions

// RecordUserRegistration records a user registration
func RecordUserRegistration() {
	UserRegistrationsTotal.Inc()
}

// RecordUserLogin records a user login attempt
func RecordUserLogin(success bool) {
	status := "failure"
	if success {
		status = "success"
	}
	UserLoginsTotal.WithLabelValues(status).Inc()
}

// RecordTokenRefresh records a token refresh attempt
func RecordTokenRefresh(success bool) {
	status := "failure"
	if success {
		status = "success"
	}
	TokenRefreshesTotal.WithLabelValues(status).Inc()
}

// RecordTokenRevocation records a token revocation
func RecordTokenRevocation() {
	TokenRevocationsTotal.Inc()
}

// RecordEmailVerification records an email verification attempt
func RecordEmailVerification(success bool) {
	status := "failure"
	if success {
		status = "success"
	}
	EmailVerificationsTotal.WithLabelValues(status).Inc()
}

// RecordPasswordChange records a password change attempt
func RecordPasswordChange(success bool) {
	status := "failure"
	if success {
		status = "success"
	}
	PasswordChangesTotal.WithLabelValues(status).Inc()
}

// RecordRequestTimeout records a request timeout
func RecordRequestTimeout(endpoint string) {
	RequestTimeoutsTotal.WithLabelValues(endpoint).Inc()
}

