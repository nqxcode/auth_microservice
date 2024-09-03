package config

import (
	"net"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	httpHostEnvName         = "HTTP_HOST"
	httpPortEnvName         = "HTTP_PORT"
	httpReadTimeoutEnvName  = "HTTP_READ_TIMEOUT"
	httpWriteTimeoutEnvName = "HTTP_WRITE_TIMEOUT" // nolint: gosec
	httpIdleTimeoutEnvName  = "HTTP_IDLE_TIMEOUT"
)

// HTTPConfig http config
type HTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	IdleTimeout() time.Duration
}

type httpConfig struct {
	host         string
	port         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

// NewHTTPConfig new http config
func NewHTTPConfig() (HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("http host not found")
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}

	httpReadTimeoutStr := os.Getenv(httpReadTimeoutEnvName)
	if len(httpReadTimeoutStr) == 0 {
		httpReadTimeoutStr = "0"
	}
	httpReadTimeout, err := strconv.Atoi(httpReadTimeoutStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse http read timeout")
	}

	httpWriteTimeoutStr := os.Getenv(httpWriteTimeoutEnvName)
	if len(httpWriteTimeoutStr) == 0 {
		httpWriteTimeoutStr = "0"
	}
	httpWriteTimeout, err := strconv.Atoi(httpWriteTimeoutStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse http write timeout")
	}

	httpIdleTimeoutStr := os.Getenv(httpIdleTimeoutEnvName)
	if len(httpIdleTimeoutStr) == 0 {
		httpIdleTimeoutStr = "0"
	}
	httpIdleTimeout, err := strconv.Atoi(httpIdleTimeoutStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse http idle timeout")
	}

	return &httpConfig{
		host:         host,
		port:         port,
		readTimeout:  time.Duration(httpReadTimeout) * time.Second,
		writeTimeout: time.Duration(httpWriteTimeout) * time.Second,
		idleTimeout:  time.Duration(httpIdleTimeout) * time.Second,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
func (cfg *httpConfig) ReadTimeout() time.Duration {
	return cfg.readTimeout
}
func (cfg *httpConfig) WriteTimeout() time.Duration {
	return cfg.writeTimeout
}
func (cfg *httpConfig) IdleTimeout() time.Duration {
	return cfg.idleTimeout
}
