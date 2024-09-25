package tracing

import (
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

// Init initializes tracing
func Init(logger *zap.Logger, serviceName, localAgentHostPort string) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: localAgentHostPort,
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("failed to init tracing", zap.Error(err))
	}
}
