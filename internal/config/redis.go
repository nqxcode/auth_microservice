package config

import (
	"net"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	redisHostEnvName              = "REDIS_HOST"
	redisPortEnvName              = "REDIS_PORT"
	redisPasswordEnvName          = "REDIS_PASSWORD"
	redisConnectionTimeoutEnvName = "REDIS_CONNECTION_TIMEOUT_SEC"
	redisMaxIdleEnvName           = "REDIS_MAX_IDLE"
	redisIdleTimeoutEnvName       = "REDIS_IDLE_TIMEOUT_SEC"
	redisDbEnvName                = "REDIS_DB"
)

type redisConfig struct {
	host string
	port string

	password string

	connectionTimeout time.Duration

	maxIdle     int
	idleTimeout time.Duration

	db int
}

func NewRedisConfig() (*redisConfig, error) {
	host := os.Getenv(redisHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("redis host not found")
	}

	port := os.Getenv(redisPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("redis port not found")
	}

	password := os.Getenv(redisPasswordEnvName)
	if len(password) == 0 {
		return nil, errors.New("redis password not found")
	}

	connectionTimeoutStr := os.Getenv(redisConnectionTimeoutEnvName)
	if len(connectionTimeoutStr) == 0 {
		return nil, errors.New("redis connection timeout not found")
	}

	connectionTimeout, err := strconv.ParseInt(connectionTimeoutStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse connection timeout")
	}

	maxIdleStr := os.Getenv(redisMaxIdleEnvName)
	if len(maxIdleStr) == 0 {
		return nil, errors.New("redis max idle not found")
	}

	maxIdle, err := strconv.Atoi(maxIdleStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse max idle")
	}

	idleTimeoutStr := os.Getenv(redisIdleTimeoutEnvName)
	if len(idleTimeoutStr) == 0 {
		return nil, errors.New("redis idle timeout not found")
	}

	idleTimeout, err := strconv.ParseInt(idleTimeoutStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse idle timeout")
	}

	dbStr := os.Getenv(redisDbEnvName)
	if len(dbStr) == 0 {
		return nil, errors.New("redis db not found")
	}

	db, err := strconv.ParseInt(dbStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse db")
	}

	return &redisConfig{
		host:              host,
		port:              port,
		password:          password,
		connectionTimeout: time.Duration(connectionTimeout) * time.Second,
		maxIdle:           maxIdle,
		idleTimeout:       time.Duration(idleTimeout) * time.Second,
		db:                int(db),
	}, nil
}

func (cfg *redisConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *redisConfig) Password() string {
	return cfg.password
}

func (cfg *redisConfig) ConnectionTimeout() time.Duration {
	return cfg.connectionTimeout
}

func (cfg *redisConfig) MaxIdle() int {
	return cfg.maxIdle
}

func (cfg *redisConfig) IdleTimeout() time.Duration {
	return cfg.idleTimeout
}

func (cfg *redisConfig) DB() int {
	return cfg.db
}
