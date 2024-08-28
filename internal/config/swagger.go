package config

import (
	"net"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	swaggerHostEnvName         = "SWAGGER_HOST"
	swaggerPortEnvName         = "SWAGGER_PORT"
	swaggerReadTimeoutEnvName  = "SWAGGER_READ_TIMEOUT"
	swaggerWriteTimeoutEnvName = "SWAGGER_WRITE_TIMEOUT"
	swaggerIdleTimeoutEnvName  = "SWAGGER_IDLE_TIMEOUT"
)

// SwaggerConfig swagger config
type SwaggerConfig interface {
	Address() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	IdleTimeout() time.Duration
}

type swaggerConfig struct {
	host         string
	port         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

// NewSwaggerConfig new swagger config
func NewSwaggerConfig() (SwaggerConfig, error) {
	host := os.Getenv(swaggerHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("swagger host not found")
	}

	port := os.Getenv(swaggerPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("swagger port not found")
	}

	swaggerReadTimeoutStr := os.Getenv(swaggerReadTimeoutEnvName)
	if len(swaggerReadTimeoutStr) == 0 {
		swaggerReadTimeoutStr = "0"
	}
	swaggerReadTimeout, err := strconv.Atoi(swaggerReadTimeoutStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse http read timeout")
	}

	swaggerWriteTimeoutStr := os.Getenv(swaggerWriteTimeoutEnvName)
	if len(swaggerWriteTimeoutStr) == 0 {
		swaggerWriteTimeoutStr = "0"
	}
	swaggerWriteTimeout, err := strconv.Atoi(swaggerWriteTimeoutStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse http write timeout")
	}

	swaggerIdleTimeoutStr := os.Getenv(swaggerIdleTimeoutEnvName)
	if len(swaggerIdleTimeoutStr) == 0 {
		swaggerIdleTimeoutStr = "0"
	}
	swaggerIdleTimeout, err := strconv.Atoi(swaggerIdleTimeoutStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse http idle timeout")
	}

	return &swaggerConfig{
		host:         host,
		port:         port,
		readTimeout:  time.Duration(swaggerReadTimeout) * time.Second,
		writeTimeout: time.Duration(swaggerWriteTimeout) * time.Second,
		idleTimeout:  time.Duration(swaggerIdleTimeout) * time.Second,
	}, nil
}

func (cfg *swaggerConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
func (cfg *swaggerConfig) ReadTimeout() time.Duration {
	return cfg.readTimeout
}
func (cfg *swaggerConfig) WriteTimeout() time.Duration {
	return cfg.writeTimeout
}
func (cfg *swaggerConfig) IdleTimeout() time.Duration {
	return cfg.idleTimeout
}
