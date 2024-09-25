package config

import (
	"os"
	"strconv"
)

const (
	loggerLogLevelEnvName          = "LOGGER_LOG_LEVEL"
	loggerRollingFilenameEnvName   = "LOGGER_ROLLING_FILENAME"
	loggerRollingMaxSizeEnvName    = "LOGGER_ROLLING_FILE_MAX_SIZE"
	loggerRollingMaxBackupsEnvName = "LOGGER_ROLLING_FILE_MAX_BACKUPS"
	loggerRollingMaxAgeEnvName     = "LOGGER_ROLLING_FILE_MAX_AGE"
)

// LoggerConfig logger config
type LoggerConfig interface {
	GetLogLevel() string
	GetRollingConfig() RollingConfig
}

type loggerConfig struct {
	LogLevel string
	Rolling  RollingConfig
}

// RollingConfig rolling config
type RollingConfig struct {
	Filename           string
	MaxSizeInMegabytes int
	MaxBackups         int
	MaxAgeInDays       int
}

func NewLoggerConfig() *loggerConfig {
	return &loggerConfig{
		LogLevel: getEnv(loggerLogLevelEnvName, "info"),
		Rolling: RollingConfig{
			Filename:           getEnv(loggerRollingFilenameEnvName, "logs/app.log"),
			MaxSizeInMegabytes: getEnvInt(loggerRollingMaxSizeEnvName, 10),
			MaxBackups:         getEnvInt(loggerRollingMaxBackupsEnvName, 3),
			MaxAgeInDays:       getEnvInt(loggerRollingMaxAgeEnvName, 7),
		},
	}
}

func (c *loggerConfig) GetLogLevel() string {
	return c.LogLevel
}

func (c *loggerConfig) GetRollingConfig() RollingConfig {
	return c.Rolling
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvInt(key string, defaultValue int) int {
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
