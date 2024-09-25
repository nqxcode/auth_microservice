package config

import (
	"github.com/nqxcode/auth_microservice/internal/utils"
	"net"
	"time"
)

const (
	prometheusHostEnvName        = "PROMETHEUS_HOST"
	prometheusPortEnvName        = "PROMETHEUS_PORT"
	prometheusMetricsPathEnvName = "PROMETHEUS_METRICS_PATH"
	prometheusReadHeaderTimeout  = "PROMETHEUS_READ_HEADER_TIMEOUT"
)

// PrometheusConfig swagger config
type PrometheusConfig interface {
	Address() string
	MetricsPath() string
	ReadHeaderTimeout() time.Duration
}

type prometheusConfig struct {
	host              string
	port              string
	metricsPath       string
	readHeaderTimeout time.Duration
}

// NewPrometheusConfig new prometheus config
func NewPrometheusConfig() (PrometheusConfig, error) {
	return &prometheusConfig{
		host:              utils.GetEnv(prometheusHostEnvName, "localhost"),
		port:              utils.GetEnv(prometheusPortEnvName, "2112"),
		metricsPath:       utils.GetEnv(prometheusMetricsPathEnvName, "/metrics"),
		readHeaderTimeout: utils.GetEnvDuration(prometheusReadHeaderTimeout, 10*time.Second),
	}, nil
}

func (cfg *prometheusConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *prometheusConfig) MetricsPath() string {
	return cfg.metricsPath
}

func (cfg *prometheusConfig) ReadHeaderTimeout() time.Duration {
	return cfg.readHeaderTimeout
}
