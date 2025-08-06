package config

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func (c *config) InitializeTelemetry() func() {
	// OTLP trace exporter
	otlpEndpoint := os.Getenv("OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "otel-collector:4317" // Use o nome do serviço Docker e a porta OTLP gRPC
	}
	traceExporter, err := otlptracegrpc.New(c.context,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otlpEndpoint),
	)
	if err != nil {
		slog.Error("failed to create OTLP trace exporter", "err", err)
		os.Exit(1)
	}

	// Resource
	res, err := resource.New(c.context,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("toq_server"),
			semconv.ServiceVersionKey.String("v2.1-grpc"),
		),
	)
	if err != nil {
		slog.Error("failed to create resource", "err", err)
		os.Exit(1)
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
		os.Exit(1)
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
		slog.Error("failed to start listen and serve", "err", http.ListenAndServe(":4318", nil))
		os.Exit(1)
	}()

	return func() {
		if err := tp.Shutdown(c.context); err != nil {
			slog.Error("failed to shutdown TracerProvider", "err", err)
			os.Exit(1)
		}
		if err := mp.Shutdown(c.context); err != nil {
			slog.Error("failed to shutdown MeterProvider", "err", err)
			os.Exit(1)
		}
	}
}
