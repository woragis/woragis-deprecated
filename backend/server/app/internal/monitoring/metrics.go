package monitoring

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
)

type metricsRegistry struct {
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	userRegistrations   prometheus.Counter
	handler             http.Handler
}

func newMetricsRegistry(namespace string) *metricsRegistry {
	httpRequests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "http_requests_total",
		Help:      "Total count of HTTP requests",
	}, []string{"method", "route", "status"})

	durations := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "http_request_duration_seconds",
		Help:      "Histogram of HTTP request durations in seconds",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "route"})

	registrations := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "user_registrations_total",
		Help:      "Total number of user registrations",
	})

	reg := metricsRegistry{
		httpRequestsTotal:   httpRequests,
		httpRequestDuration: durations,
		userRegistrations:   registrations,
	}

	prometheus.MustRegister(httpRequests, durations, registrations)
	reg.handler = promhttp.Handler()
	return &reg
}

func (m *metricsRegistry) metricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		status := strconv.Itoa(c.Response().StatusCode())
		route := c.Route().Path
		if route == "" {
			route = c.Path()
		}

		m.httpRequestsTotal.WithLabelValues(c.Method(), route, status).Inc()
		m.httpRequestDuration.WithLabelValues(c.Method(), route).Observe(time.Since(start).Seconds())

		return err
	}
}

func (m *metricsRegistry) metricsHandler() http.Handler {
	return m.handler
}

func (m *metricsRegistry) incrementRegistrations() {
	m.userRegistrations.Inc()
}

func (m *metricsRegistry) snapshot() (string, error) {
	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	encoder := expfmt.NewEncoder(&buf, expfmt.FmtText)
	for _, mf := range families {
		if err := encoder.Encode(mf); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}
