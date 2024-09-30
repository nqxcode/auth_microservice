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

const (
	defaultLoggerLogLevel          = "info"
	defaultLoggerRollingFilename   = "logs/app.log"
	defaultLoggerRollingMaxSize    = 10
	defaultLoggerRollingMaxBackups = 3
	defaultLoggerRollingMaxAge     = 7
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
		LogLevel: utils.GetEnv(loggerLogLevelEnvName, defaultLoggerLogLevel),
		Rolling: RollingConfig{
			Filename:           utils.GetEnv(loggerRollingFilenameEnvName, defaultLoggerRollingFilename),
			MaxSizeInMegabytes: utils.GetEnvInt(loggerRollingMaxSizeEnvName, defaultLoggerRollingMaxSize),
			MaxBackups:         utils.GetEnvInt(loggerRollingMaxBackupsEnvName, defaultLoggerRollingMaxBackups),
			MaxAgeInDays:       utils.GetEnvInt(loggerRollingMaxAgeEnvName, defaultLoggerRollingMaxAge),
		},
	}
}

func (c *loggerConfig) GetLogLevel() string {
	return c.LogLevel
}

func (c *loggerConfig) GetRollingConfig() RollingConfig {
	return c.Rolling
}
