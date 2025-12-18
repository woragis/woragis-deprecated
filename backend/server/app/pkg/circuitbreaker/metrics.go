package circuitbreaker

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sony/gobreaker"
)

var (
	// CircuitBreakerState tracks the current state of circuit breakers
	CircuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "circuit_breaker_state",
			Help: "Circuit breaker state (0=closed, 1=half-open, 2=open)",
		},
		[]string{"name"},
	)

	// CircuitBreakerTransitions counts state transitions
	CircuitBreakerTransitions = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker_transitions_total",
			Help: "Total number of circuit breaker state transitions",
		},
		[]string{"name", "from", "to"},
	)

	// CircuitBreakerRequestsRejected counts rejected requests (circuit open)
	CircuitBreakerRequestsRejected = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker_requests_rejected_total",
			Help: "Total number of requests rejected due to open circuit",
		},
		[]string{"name"},
	)

	// CircuitBreakerRequestsAllowed counts allowed requests
	CircuitBreakerRequestsAllowed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker_requests_allowed_total",
			Help: "Total number of requests allowed through circuit breaker",
		},
		[]string{"name"},
	)
)

// stateToFloat converts circuit breaker state to float64 for metrics
func stateToFloat(state gobreaker.State) float64 {
	switch state {
	case gobreaker.StateClosed:
		return 0
	case gobreaker.StateHalfOpen:
		return 1
	case gobreaker.StateOpen:
		return 2
	default:
		return -1
	}
}

// RecordStateChange records circuit breaker state change in metrics
func RecordStateChange(name string, from, to gobreaker.State) {
	CircuitBreakerState.WithLabelValues(name).Set(stateToFloat(to))
	CircuitBreakerTransitions.WithLabelValues(name, from.String(), to.String()).Inc()
}

// RecordRequestRejected records a rejected request
func RecordRequestRejected(name string) {
	CircuitBreakerRequestsRejected.WithLabelValues(name).Inc()
}

// RecordRequestAllowed records an allowed request
func RecordRequestAllowed(name string) {
	CircuitBreakerRequestsAllowed.WithLabelValues(name).Inc()
}
