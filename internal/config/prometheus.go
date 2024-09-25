package config

import (
	"net"
	"os"
	"time"

	"github.com/nqxcode/auth_microservice/internal/utils"
	"github.com/pkg/errors"
)

const (
	prometheusHostEnvName        = "PROMETHEUS_HOST"
	prometheusPortEnvName        = "PROMETHEUS_PORT"
	prometheusMetricsPathEnvName = "PROMETHEUS_METRICS_PATH"
	prometheusReadHeaderTimeout  = "PROMETHEUS_READ_HEADER_TIMEOUT"
)

const (
	defaultMetricsPath = "/metrics"
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
	host := os.Getenv(prometheusHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("prometheus host not found")
	}

	port := os.Getenv(prometheusPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("prometheus port not found")
	}

	return &prometheusConfig{
		host:              host,
		port:              port,
		metricsPath:       utils.GetEnv(prometheusMetricsPathEnvName, defaultMetricsPath),
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
