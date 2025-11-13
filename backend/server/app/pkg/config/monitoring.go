package config

import "os"

// MonitoringConfig holds configuration for Prometheus/Grafana integration.
type MonitoringConfig struct {
	Enabled          bool
	DBURL            string
	MetricsNamespace string
}

// LoadMonitoringConfig reads monitoring-related environment variables.
func LoadMonitoringConfig() MonitoringConfig {
	enabled := os.Getenv("MONITORING_ENABLED")
	cfg := MonitoringConfig{
		Enabled:          enabled == "" || enabled == "true" || enabled == "1",
		DBURL:            os.Getenv("MONITORING_DB_URL"),
		MetricsNamespace: os.Getenv("METRICS_NAMESPACE"),
	}

	if cfg.MetricsNamespace == "" {
		cfg.MetricsNamespace = "woragis"
	}

	return cfg
}
