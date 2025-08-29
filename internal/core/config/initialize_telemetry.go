package config

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func (c *config) InitializeTelemetry() (func(), error) {
	// Check if OTEL should be disabled via env.yaml or environment variable
	if !c.env.TELEMETRY.Enabled || os.Getenv("DISABLE_OTEL") == "true" {
		slog.Info("OpenTelemetry disabled by configuration")
		return func() {}, nil
	}

	// OTLP trace exporter
	otlpEndpoint := c.env.TELEMETRY.OTLP.Endpoint
	if otlpEndpoint == "" {
		// Check if running in Docker by looking for /.dockerenv file
		if _, err := os.Stat("/.dockerenv"); err == nil {
			otlpEndpoint = "jaeger:4318" // HTTP endpoint dentro do Docker (porta interna do Jaeger)
		} else {
			otlpEndpoint = "localhost:14318" // HTTP endpoint fora do Docker (porta externa mapeada)
		}
	}

	// Allow disabling OTLP via environment variable (for debugging)
	if !c.env.TELEMETRY.OTLP.Enabled || os.Getenv("DISABLE_OTLP") == "true" {
		slog.Info("OTLP tracing disabled by configuration")
		return func() {
			slog.Info("OpenTelemetry cleanup (OTLP disabled)")
		}, nil
	}

	// Try to create OTLP exporter with timeout
	ctx, cancel := context.WithTimeout(c.context, 5*time.Second)
	defer cancel()

	// Configure OTLP exporter options based on env.yaml
	var exporterOptions []otlptracehttp.Option
	exporterOptions = append(exporterOptions, otlptracehttp.WithEndpoint(otlpEndpoint))

	// Use insecure connection if configured in env.yaml
	if c.env.TELEMETRY.OTLP.Insecure {
		exporterOptions = append(exporterOptions, otlptracehttp.WithInsecure())
	}

	traceExporter, err := otlptracehttp.New(ctx, exporterOptions...)
	if err != nil {
		slog.Warn("failed to create OTLP trace exporter, continuing without distributed tracing", "endpoint", otlpEndpoint, "err", err)
		// Return empty cleanup function if OTLP fails
		return func() {
			slog.Info("OpenTelemetry cleanup (no OTLP exporter)")
		}, nil
	}

	// Resource
	res, err := resource.New(c.context,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("toq_server"),
			semconv.ServiceVersionKey.String("v2.1-http"),
		),
	)
	if err != nil {
		slog.Error("failed to create resource", "err", err)
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	// Prometheus exporter
	options := prometheus.WithoutScopeInfo()
	prometheusExporter, err := prometheus.New(options)
	if err != nil {
		slog.Error("failed to create Prometheus exporter", "err", err)
		return nil, fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	// Meter provider
	mp := metric.NewMeterProvider(
		metric.WithReader(prometheusExporter),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	// Start Prometheus HTTP server
	go func() {
		metricsPort := c.env.TELEMETRY.METRICS.Port
		if metricsPort == "" {
			metricsPort = ":4318" // Default port
		}
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(metricsPort, nil); err != nil {
			slog.Error("failed to start listen and serve", "err", err)
			// Em uma goroutine, não podemos retornar erro, então apenas logamos
		}
	}()

	return func() {
		// Create a context with timeout for graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		slog.Info("Shutting down OpenTelemetry...")

		// Shutdown TracerProvider with timeout
		if err := tp.Shutdown(shutdownCtx); err != nil {
			slog.Error("failed to shutdown TracerProvider", "err", err)
		} else {
			slog.Info("TracerProvider shutdown completed")
		}

		// Shutdown MeterProvider with timeout
		if err := mp.Shutdown(shutdownCtx); err != nil {
			slog.Error("failed to shutdown MeterProvider", "err", err)
		} else {
			slog.Info("MeterProvider shutdown completed")
		}
	}, nil
}
