package utils

import (
	"os"
	"strconv"
	"time"
)

// GetEnv get environment variable
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

// GetEnvInt get environment variable as int
func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// GetEnvDuration get environment variable as duration
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return time.Duration(defaultValue)
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return time.Duration(defaultValue)
	}

	return duration
}
