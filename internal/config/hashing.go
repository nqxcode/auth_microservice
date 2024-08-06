package config

import (
	"os"

	"github.com/pkg/errors"
)

const hashingSaltEnvName = "HASHING_SALT"

// HashingConfig hashing config
type HashingConfig interface {
	Salt() string
}

type hashingConfig struct {
	salt string
}

// NewHashingConfig create new hashing config
func NewHashingConfig() (HashingConfig, error) {
	salt := os.Getenv(hashingSaltEnvName)
	if len(salt) == 0 {
		return nil, errors.New("salt is empty")
	}

	return &hashingConfig{
		salt: salt,
	}, nil
}

// Salt get salt for password hashing
func (hc *hashingConfig) Salt() string {
	return hc.salt
}
