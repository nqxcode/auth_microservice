package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

var globalEnabled bool

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

	globalEnabled = true
}

// InitNoop initializes noop tracer
func InitNoop() {
	globalEnabled = false
}

// StartSpanFromContext starts a new span with opentracing.StartSpanFromContext
func StartSpanFromContext(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	if !globalEnabled {
		return nil, ctx
	}

	return opentracing.StartSpanFromContext(ctx, operationName)
}
