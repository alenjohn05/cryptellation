package otel

import (
	"context"
	"os"

	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogsgrpc"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"github.com/agoda-com/otelzap"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logs struct {
	clientCleanUp func(ctx context.Context) error
	exporter      *otlplogs.Exporter
	provider      *sdk.LoggerProvider
	logger        *zap.Logger
}

func newLogs(ctx context.Context, serviceName, url string) (logs, error) {
	// Create exporter
	client := otlplogsgrpc.NewClient(otlplogsgrpc.WithEndpoint(url), otlplogsgrpc.WithInsecure())
	exporter, err := otlplogs.NewExporter(ctx, otlplogs.WithClient(client))
	if err != nil {
		return logs{}, err
	}

	// Create resource
	resource := resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(serviceName))

	// Create provider
	provider := sdk.NewLoggerProvider(
		sdk.WithBatcher(exporter),
		sdk.WithResource(resource),
	)

	// Set opentelemetry logger provider globally
	// otellogs.SetLoggerProvider(provider)

	// Create cores
	cores := []zapcore.Core{
		otelzap.NewOtelCore(provider),
	}

	// Set console loggers
	if os.Getenv("DEV_MODE") != "" {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		consoleDebugging := zapcore.Lock(os.Stdout)
		consoleErrors := zapcore.Lock(os.Stderr)

		// Set priorities
		highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel && lvl >= logLevelFromEnv()
		})

		cores = append(cores,
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
		)
	}

	// Create a new logger
	logger := zap.New(zapcore.NewTee(cores...))
	// undo := zap.ReplaceGlobals(logger),

	return logs{
		clientCleanUp: client.Stop,
		exporter:      exporter,
		provider:      provider,
		logger:        logger,
	}, nil
}

func logLevelFromEnv() zapcore.Level {
	envLogLevel := os.Getenv("LOG_LEVEL")
	if envLogLevel == "DEBUG" {
		return zapcore.DebugLevel
	} else if envLogLevel == "INFO" {
		return zapcore.InfoLevel
	} else if envLogLevel == "WARN" {
		return zapcore.WarnLevel
	}
	return zapcore.InfoLevel
}

func (l logs) close(ctx context.Context) {
	_ = l.logger.Sync()
	_ = l.provider.Shutdown(ctx)
	_ = l.exporter.Shutdown(ctx)
	_ = l.clientCleanUp(ctx)
}
