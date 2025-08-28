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
	// Check if OTEL should be disabled
	if os.Getenv("DISABLE_OTEL") == "true" {
		slog.Info("OpenTelemetry disabled by configuration")
		return func() {}, nil
	}

	// OTLP trace exporter
	otlpEndpoint := os.Getenv("OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		// Check if running in Docker by looking for /.dockerenv file
		if _, err := os.Stat("/.dockerenv"); err == nil {
			otlpEndpoint = "otel-collector:4317" // Dentro do Docker
		} else {
			otlpEndpoint = "localhost:4317" // Fora do Docker
		}
	}

	// Try to create OTLP exporter with timeout
	ctx, cancel := context.WithTimeout(c.context, 5*time.Second)
	defer cancel()

	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(fmt.Sprintf("http://%s", otlpEndpoint)),
	)
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
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":4318", nil); err != nil {
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
