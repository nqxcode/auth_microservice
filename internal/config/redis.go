package config

import (
	"os"
	"strconv"

	"github.com/nqxcode/platform_common/client/cache"
	"github.com/pkg/errors"
)

const (
	addressEnvName  = "REDIS_ADDRESS"
	passwordEnvName = "REDIS_PASSWORD"
	dbEnvName       = "REDIS_DB"
)

type redisConfig struct {
	address  string
	password string
	db       int
}

// NewRedisConfig create new redis config
func NewRedisConfig() (cache.RedisConfig, error) {
	address := os.Getenv(addressEnvName)
	if len(address) == 0 {
		return nil, errors.New("redis address not found")
	}

	password := os.Getenv(passwordEnvName)
	if len(address) == 0 {
		return nil, errors.New("redis password not found")
	}

	dbStr := os.Getenv(dbEnvName)
	if len(dbStr) == 0 {
		return nil, errors.New("redis database not found")
	}
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return nil, errors.New("redis database is not a number")
	}

	return &redisConfig{
		address:  address,
		password: password,
		db:       db,
	}, nil
}

// Address get address
func (cfg *redisConfig) Address() string {
	return cfg.address
}

// Password get password
func (cfg *redisConfig) Password() string {
	return cfg.password
}

// DB get db
func (cfg *redisConfig) DB() int {
	return cfg.db
}
