package config

import (
	"net"
	"os"

	"github.com/pkg/errors"
)

const (
	prometheusHostEnvName        = "PROMETHEUS_HOST"
	prometheusPortEnvName        = "PROMETHEUS_PORT"
	prometheusMetricsPathEnvName = "PROMETHEUS_METRICS_PATH"
)

// PrometheusConfig swagger config
type PrometheusConfig interface {
	Address() string
	MetricsPath() string
}

type prometheusConfig struct {
	host        string
	port        string
	metricsPath string
}

func NewPrometheusConfig() (PrometheusConfig, error) {
	host := os.Getenv(prometheusHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("swagger host not found")
	}

	port := os.Getenv(prometheusPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("swagger port not found")
	}

	metricsPath := os.Getenv(prometheusMetricsPathEnvName)
	if len(metricsPath) == 0 {
		return nil, errors.New("prometheus metrics path not found")
	}

	return &prometheusConfig{
		host:        host,
		port:        port,
		metricsPath: metricsPath,
	}, nil
}

func (cfg *prometheusConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *prometheusConfig) MetricsPath() string {
	return cfg.metricsPath
}
