package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// InitLogger initializes the OpenTelemetry Logger Provider, registers it globally,
// and sets up the default standard slog logger to use the OTel bridge.
func InitLogger(ctx context.Context, serviceName, env string) (*sdklog.LoggerProvider, error) {
	var processors []sdklog.Processor

	// 1a. Only initialize the OTLP exporter if we are in production or if an OTel endpoint is explicitly configured.
	// This prevents infinite connection-refused loops/warnings during standard local development.
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if env == "production" || otlpEndpoint != "" {
		// Create OTLP/HTTP log exporter.
		// We use WithInsecure() to default to HTTP instead of HTTPS (ideal for Datadog Agent / local OTel Collector).
		otlpExporter, err := otlploghttp.New(ctx, otlploghttp.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
		}
		processors = append(processors, sdklog.NewBatchProcessor(otlpExporter))
		slog.Info("OpenTelemetry OTLP log exporter enabled", "endpoint", otlpEndpoint)
	}

	// 1b. Create stdout log exporter for beautiful terminal console output (always enabled)
	stdoutExporter, err := stdoutlog.New(stdoutlog.WithPrettyPrint())
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout log exporter: %w", err)
	}
	processors = append(processors, sdklog.NewSimpleProcessor(stdoutExporter))

	// 2. Define OTel resource attributes (Service Name, Env)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String(env),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTel resource: %w", err)
	}

	// 3. Create LoggerProvider registering our active processors.
	opts := []sdklog.LoggerProviderOption{
		sdklog.WithResource(res),
	}
	for _, p := range processors {
		opts = append(opts, sdklog.WithProcessor(p))
	}
	provider := sdklog.NewLoggerProvider(opts...)

	// 4. Set global logger provider
	global.SetLoggerProvider(provider)

	// 5. Configure standard library's slog to use the OTel bridge logger
	slog.SetDefault(otelslog.NewLogger(serviceName))

	slog.Info("OpenTelemetry logging initialized successfully", "service", serviceName, "env", env)

	return provider, nil
}
