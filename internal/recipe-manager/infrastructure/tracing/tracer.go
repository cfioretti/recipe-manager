package tracing

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type TracerConfig struct {
	ServiceName        string
	ServiceVersion     string
	JaegerEndpoint     string
	Environment        string
	SamplingRatio      float64
	BatchTimeout       time.Duration
	MaxExportBatchSize int
}

func DefaultTracerConfig() *TracerConfig {
	return &TracerConfig{
		ServiceName:        getEnvOrDefault("OTEL_SERVICE_NAME", "recipe-manager"),
		ServiceVersion:     getEnvOrDefault("OTEL_SERVICE_VERSION", "1.0.0"),
		JaegerEndpoint:     getEnvOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318"),
		Environment:        getEnvOrDefault("ENVIRONMENT", "development"),
		SamplingRatio:      1.0, // Sample all traces in development
		BatchTimeout:       5 * time.Second,
		MaxExportBatchSize: 512,
	}
}

type TracerProvider struct {
	provider *trace.TracerProvider
	cleanup  func(context.Context) error
}

func NewTracerProvider(config *TracerConfig) (*TracerProvider, error) {
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(config.JaegerEndpoint),
		otlptracehttp.WithInsecure(), // Use HTTP instead of HTTPS for local development
		otlptracehttp.WithURLPath("/v1/traces"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(config.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithBatchTimeout(config.BatchTimeout),
			trace.WithMaxExportBatchSize(config.MaxExportBatchSize),
		),
		trace.WithResource(res),
		trace.WithSampler(trace.TraceIDRatioBased(config.SamplingRatio)),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &TracerProvider{
		provider: tp,
		cleanup: func(ctx context.Context) error {
			return tp.Shutdown(ctx)
		},
	}, nil
}

func (tp *TracerProvider) GetTracer(name string) oteltrace.Tracer {
	return tp.provider.Tracer(name)
}

func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	return tp.cleanup(ctx)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var globalTracerProvider *TracerProvider

func InitTracing(config *TracerConfig) error {
	if config == nil {
		config = DefaultTracerConfig()
	}

	tp, err := NewTracerProvider(config)
	if err != nil {
		return fmt.Errorf("failed to initialize tracing: %w", err)
	}

	globalTracerProvider = tp
	return nil
}

func GetGlobalTracer(name string) oteltrace.Tracer {
	if globalTracerProvider == nil {
		return otel.Tracer(name)
	}
	return globalTracerProvider.GetTracer(name)
}

func ShutdownTracing(ctx context.Context) error {
	if globalTracerProvider == nil {
		return nil
	}
	return globalTracerProvider.Shutdown(ctx)
}
