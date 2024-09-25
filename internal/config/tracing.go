package config

import (
	"github.com/pkg/errors"
	"net"
	"os"
)

const (
	tracingHostEnvName = "TRACING_HOST"
	tracingPortEnvName = "TRACING_PORT"
)

// TracingConfig tracing config
type TracingConfig interface {
	Address() string
}

type tracingConfig struct {
	host string
	port string
}

// NewTracingConfig new tracing config
func NewTracingConfig() (TracingConfig, error) {
	host := os.Getenv(tracingHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("tracing host not found")
	}

	port := os.Getenv(tracingPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("tracing port not found")
	}

	return &tracingConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *tracingConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
