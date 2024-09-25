package config

import (
	"os"

	"github.com/pkg/errors"
)

const (
	serviceNameEnvName = "SERVICE_NAME"
)

// AppConfig app config
type AppConfig interface {
	GetName() string
}

type appConfig struct {
	name string
}

// NewAppConfig new app config
func NewAppConfig() (AppConfig, error) {
	serviceName := os.Getenv(serviceNameEnvName)
	if len(serviceName) == 0 {
		return nil, errors.New("app name not found")
	}
	return &appConfig{
		name: serviceName,
	}, nil
}

// GetName get app name
func (c *appConfig) GetName() string {
	return c.name
}
