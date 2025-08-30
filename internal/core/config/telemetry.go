package config

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

// TelemetryManager gerencia a inicialização e shutdown do OpenTelemetry
type TelemetryManager struct {
	traceProvider  *sdktrace.TracerProvider
	metricProvider *sdkmetric.MeterProvider
	env            globalmodel.Environment
}

// NewTelemetryManager cria uma nova instância do gerenciador de telemetria
func NewTelemetryManager(env globalmodel.Environment) *TelemetryManager {
	return &TelemetryManager{
		env: env,
	}
}

// Initialize configura e inicializa o OpenTelemetry
func (tm *TelemetryManager) Initialize(ctx context.Context) (func(), error) {
	if !tm.env.TELEMETRY.Enabled {
		slog.Info("OpenTelemetry disabled by configuration")
		return func() {}, nil
	}

	// Criar resource comum
	res, err := tm.createResource()
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Inicializar tracing
	if err := tm.initializeTracing(ctx, res); err != nil {
		return nil, fmt.Errorf("failed to initialize tracing: %w", err)
	}

	// Inicializar métricas
	if err := tm.initializeMetrics(ctx, res); err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	// Configurar propagators
	tm.configurePropagators()

	slog.Info("OpenTelemetry initialized successfully",
		"tracing_enabled", true,
		"metrics_enabled", true,
		"endpoint", tm.env.TELEMETRY.OTLP.Endpoint)

	// Retornar função de shutdown
	return tm.shutdown, nil
}

// createResource cria o resource comum para tracing e métricas
func (tm *TelemetryManager) createResource() (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			resource.Default().SchemaURL(),
			semconv.ServiceName("toq_server"),
			semconv.ServiceVersion("1.0.0"),
			semconv.ServiceInstanceID("toq_server_instance_1"),
		),
	)
}

// initializeTracing configura o tracing OpenTelemetry
func (tm *TelemetryManager) initializeTracing(ctx context.Context, res *resource.Resource) error {
	if !tm.env.TELEMETRY.OTLP.Enabled {
		slog.Info("OTLP tracing disabled")
		return nil
	}

	// Configurar opções do exporter
	options := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(tm.env.TELEMETRY.OTLP.Endpoint),
	}

	if tm.env.TELEMETRY.OTLP.Insecure {
		options = append(options, otlptracehttp.WithInsecure())
	}

	// Criar exporter OTLP HTTP
	exporter, err := otlptracehttp.New(ctx, options...)
	if err != nil {
		return fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// Criar trace provider
	tm.traceProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Definir como provider global
	otel.SetTracerProvider(tm.traceProvider)

	slog.Info("OpenTelemetry tracing initialized", "endpoint", tm.env.TELEMETRY.OTLP.Endpoint)
	return nil
}

// initializeMetrics configura as métricas OpenTelemetry
func (tm *TelemetryManager) initializeMetrics(ctx context.Context, res *resource.Resource) error {
	if !tm.env.TELEMETRY.OTLP.Enabled {
		slog.Info("OTLP metrics disabled")
		return nil
	}

	// Configurar opções do exporter
	options := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(tm.env.TELEMETRY.OTLP.Endpoint),
	}

	if tm.env.TELEMETRY.OTLP.Insecure {
		options = append(options, otlpmetrichttp.WithInsecure())
	}

	// Criar exporter OTLP HTTP para métricas
	exporter, err := otlpmetrichttp.New(ctx, options...)
	if err != nil {
		return fmt.Errorf("failed to create OTLP metric exporter: %w", err)
	}

	// Criar metric provider
	tm.metricProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter,
			sdkmetric.WithInterval(30*time.Second))),
		sdkmetric.WithResource(res),
	)

	// Definir como provider global
	otel.SetMeterProvider(tm.metricProvider)

	slog.Info("OpenTelemetry metrics initialized", "endpoint", tm.env.TELEMETRY.OTLP.Endpoint)
	return nil
}

// configurePropagators configura os propagators para tracing distribuído
func (tm *TelemetryManager) configurePropagators() {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

// shutdown executa o shutdown graceful do OpenTelemetry
func (tm *TelemetryManager) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if tm.traceProvider != nil {
		if err := tm.traceProvider.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown trace provider", "error", err)
		} else {
			slog.Info("OpenTelemetry trace provider shutdown completed")
		}
	}

	if tm.metricProvider != nil {
		if err := tm.metricProvider.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown metric provider", "error", err)
		} else {
			slog.Info("OpenTelemetry metric provider shutdown completed")
		}
	}

	slog.Info("OpenTelemetry shutdown completed")
}
