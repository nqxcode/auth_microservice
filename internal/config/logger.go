package config

import (
	"github.com/nqxcode/auth_microservice/internal/utils"
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

// NewLoggerConfig new logger config
func NewLoggerConfig() *loggerConfig {
	return &loggerConfig{
		LogLevel: utils.GetEnv(loggerLogLevelEnvName, "info"),
		Rolling: RollingConfig{
			Filename:           utils.GetEnv(loggerRollingFilenameEnvName, "logs/app.log"),
			MaxSizeInMegabytes: utils.GetEnvInt(loggerRollingMaxSizeEnvName, 10),
			MaxBackups:         utils.GetEnvInt(loggerRollingMaxBackupsEnvName, 3),
			MaxAgeInDays:       utils.GetEnvInt(loggerRollingMaxAgeEnvName, 7),
		},
	}
}

func (c *loggerConfig) GetLogLevel() string {
	return c.LogLevel
}

func (c *loggerConfig) GetRollingConfig() RollingConfig {
	return c.Rolling
}
