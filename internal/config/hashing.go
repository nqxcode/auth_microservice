package config

import (
	"github.com/pkg/errors"
	"os"
)

const hashingSaltEnvName = "HASHING_SALT"

type HashingConfig interface {
	Salt() string
}

type hashingConfig struct {
	salt string
}

func NewHashingConfig() (HashingConfig, error) {
	salt := os.Getenv(hashingSaltEnvName)
	if len(salt) == 0 {
		return nil, errors.New("salt is empty")
	}

	return &hashingConfig{
		salt: salt,
	}, nil
}

func (hc *hashingConfig) Salt() string {
	return hc.salt
}
