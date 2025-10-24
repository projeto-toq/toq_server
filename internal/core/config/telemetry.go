package config

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	globallog "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

// TelemetryManager gerencia a inicialização e shutdown do OpenTelemetry
type TelemetryManager struct {
	traceProvider  *sdktrace.TracerProvider
	metricProvider *sdkmetric.MeterProvider
	logProvider    *sdklog.LoggerProvider
	env            globalmodel.Environment
	runtimeEnv     string
}

// NewTelemetryManager cria uma nova instância do gerenciador de telemetria
func NewTelemetryManager(env globalmodel.Environment, runtimeEnv string) *TelemetryManager {
	return &TelemetryManager{
		env:        env,
		runtimeEnv: runtimeEnv,
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

	tracingActive := tm.env.TELEMETRY.TRACES.Enabled && tm.env.TELEMETRY.OTLP.Enabled
	metricsActive := tm.env.TELEMETRY.METRICS.Enabled && tm.env.TELEMETRY.OTLP.Enabled
	logsActive := tm.env.TELEMETRY.LOGS.EXPORT.Enabled && tm.env.TELEMETRY.OTLP.Enabled

	if tracingActive {
		if err := tm.initializeTracing(ctx, res); err != nil {
			return nil, fmt.Errorf("failed to initialize tracing: %w", err)
		}
	} else {
		slog.Info("Tracing pipeline not started (configuration disabled)")
	}

	if metricsActive {
		if err := tm.initializeMetrics(ctx, res); err != nil {
			return nil, fmt.Errorf("failed to initialize metrics: %w", err)
		}
	} else {
		slog.Info("Metrics pipeline not started (configuration disabled)")
	}

	if logsActive {
		if err := tm.initializeLogging(ctx, res); err != nil {
			return nil, fmt.Errorf("failed to initialize logs: %w", err)
		}
	} else {
		slog.Info("Log export pipeline not started (configuration disabled)")
	}

	if !tracingActive && !metricsActive && !logsActive {
		slog.Info("OpenTelemetry exporters disabled; no pipeline initialized")
		return func() {}, nil
	}

	// Configurar propagators apenas quando há pipeline ativo
	tm.configurePropagators()

	slog.Info("OpenTelemetry initialized",
		"tracing_enabled", tracingActive,
		"metrics_enabled", metricsActive,
		"logs_export_enabled", logsActive,
		"otlp_endpoint", tm.env.TELEMETRY.OTLP.Endpoint)

	// Retornar função de shutdown
	return tm.shutdown, nil
}

// initializeLogging configura a exportação de logs via OpenTelemetry
func (tm *TelemetryManager) initializeLogging(ctx context.Context, res *resource.Resource) error {
	if !tm.env.TELEMETRY.LOGS.EXPORT.Enabled || !tm.env.TELEMETRY.OTLP.Enabled {
		return nil
	}

	endpoint, err := normalizeOTLPEndpoint(tm.env.TELEMETRY.OTLP.Endpoint)
	if err != nil {
		return fmt.Errorf("invalid OTLP log endpoint: %w", err)
	}

	options := []otlploghttp.Option{
		otlploghttp.WithEndpoint(endpoint),
	}
	if tm.env.TELEMETRY.OTLP.Insecure {
		options = append(options, otlploghttp.WithInsecure())
	}

	exporter, err := otlploghttp.New(ctx, options...)
	if err != nil {
		return fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}

	processor := sdklog.NewBatchProcessor(exporter)
	logProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(processor),
	)

	globallog.SetLoggerProvider(logProvider)
	tm.logProvider = logProvider

	baseHandler := slog.Default().Handler()
	var otelHandler slog.Handler = otelslog.NewHandler("toq_server",
		otelslog.WithLoggerProvider(logProvider),
		otelslog.WithVersion(globalmodel.AppVersion),
		otelslog.WithSource(true),
	)

	if runtimeCfg, ok := globalmodel.GetLoggingRuntimeConfig(); ok {
		otelHandler = newLevelFilterHandler(otelHandler, runtimeCfg.Level)
	}

	combinedHandler := newTeeHandler(baseHandler, otelHandler)
	if combinedHandler == nil {
		combinedHandler = otelHandler
	}
	slog.SetDefault(slog.New(combinedHandler))
	slog.Info("OpenTelemetry logs initialized", "endpoint", tm.env.TELEMETRY.OTLP.Endpoint)
	return nil
}

// createResource cria o resource comum para tracing e métricas
func (tm *TelemetryManager) createResource() (*resource.Resource, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
	}
	instanceID := fmt.Sprintf("%s-%d", hostname, os.Getpid())

	attributes := []attribute.KeyValue{
		semconv.ServiceName("toq_server"),
		semconv.ServiceVersion(globalmodel.AppVersion),
		semconv.ServiceInstanceID(instanceID),
	}

	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			resource.Default().SchemaURL(),
			attributes...),
	)
}

// initializeTracing configura o tracing OpenTelemetry
func (tm *TelemetryManager) initializeTracing(ctx context.Context, res *resource.Resource) error {
	if !tm.env.TELEMETRY.TRACES.Enabled || !tm.env.TELEMETRY.OTLP.Enabled {
		return nil
	}

	// Normalizar endpoint (aceita tanto host:port quanto http(s)://host:port)
	endpoint, err := normalizeOTLPEndpoint(tm.env.TELEMETRY.OTLP.Endpoint)
	if err != nil {
		return fmt.Errorf("invalid OTLP trace endpoint: %w", err)
	}

	// Configurar opções do exporter
	options := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
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
	if !tm.env.TELEMETRY.METRICS.Enabled || !tm.env.TELEMETRY.OTLP.Enabled {
		return nil
	}

	// Normalizar endpoint (aceita tanto host:port quanto http(s)://host:port)
	endpoint, err := normalizeOTLPEndpoint(tm.env.TELEMETRY.OTLP.Endpoint)
	if err != nil {
		return fmt.Errorf("invalid OTLP metric endpoint: %w", err)
	}

	// Configurar opções do exporter
	options := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(endpoint),
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

	// Garantir que erros do OpenTelemetry sejam registrados como ERROR (ex: falha ao enviar métricas)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		slog.Error("OpenTelemetry error", "component", "otlp.metrics", "endpoint", tm.env.TELEMETRY.OTLP.Endpoint, "error", err)
	}))

	slog.Info("OpenTelemetry metrics initialized", "endpoint", tm.env.TELEMETRY.OTLP.Endpoint)
	return nil
}

// normalizeOTLPEndpoint aceita formatos com e sem esquema e retorna host:port
func normalizeOTLPEndpoint(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("empty endpoint")
	}

	// Se já contém esquema, parsear e extrair Host
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") || strings.HasPrefix(raw, "grpc://") {
		u, err := url.Parse(raw)
		if err != nil {
			return "", fmt.Errorf("failed to parse endpoint: %w", err)
		}
		if u.Host == "" {
			return "", fmt.Errorf("parsed endpoint has empty host")
		}
		return u.Host, nil
	}

	// Caso já esteja no formato host:port, validar minimamente
	if !strings.Contains(raw, ":") {
		return "", fmt.Errorf("endpoint must include port: %s", raw)
	}
	return raw, nil
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

	if tm.logProvider != nil {
		if err := tm.logProvider.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown log provider", "error", err)
		} else {
			slog.Info("OpenTelemetry log provider shutdown completed")
		}
	}

	slog.Info("OpenTelemetry shutdown completed")
}

// newTeeHandler cria um handler que replica registros para todos os handlers fornecidos.
func newTeeHandler(handlers ...slog.Handler) slog.Handler {
	filtered := make([]slog.Handler, 0, len(handlers))
	for _, h := range handlers {
		if h != nil {
			filtered = append(filtered, h)
		}
	}
	switch len(filtered) {
	case 0:
		return nil
	case 1:
		return filtered[0]
	default:
		return &teeHandler{handlers: filtered}
	}
}

type teeHandler struct {
	handlers []slog.Handler
}

func (t *teeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range t.handlers {
		if h != nil && h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (t *teeHandler) Handle(ctx context.Context, record slog.Record) error {
	activeCount := 0
	for _, h := range t.handlers {
		if h != nil && h.Enabled(ctx, record.Level) {
			activeCount++
		}
	}

	if activeCount == 0 {
		return nil
	}

	var firstErr error
	handled := 0
	for _, h := range t.handlers {
		if h == nil || !h.Enabled(ctx, record.Level) {
			continue
		}

		handled++
		rec := record
		if handled < activeCount {
			// Clonamos o registro quando múltiplos handlers ativos precisam recebê-lo.
			rec = record.Clone()
		}
		if err := h.Handle(ctx, rec); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (t *teeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(t.handlers))
	for i, h := range t.handlers {
		if h != nil {
			newHandlers[i] = h.WithAttrs(attrs)
		}
	}
	return &teeHandler{handlers: newHandlers}
}

func (t *teeHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(t.handlers))
	for i, h := range t.handlers {
		if h != nil {
			newHandlers[i] = h.WithGroup(name)
		}
	}
	return &teeHandler{handlers: newHandlers}
}

type levelFilterHandler struct {
	minLevel slog.Level
	next     slog.Handler
}

func newLevelFilterHandler(handler slog.Handler, minLevel slog.Level) slog.Handler {
	return &levelFilterHandler{
		minLevel: minLevel,
		next:     handler,
	}
}

func (l *levelFilterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if level < l.minLevel {
		return false
	}
	return l.next.Enabled(ctx, level)
}

func (l *levelFilterHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level < l.minLevel {
		return nil
	}
	return l.next.Handle(ctx, record)
}

func (l *levelFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &levelFilterHandler{
		minLevel: l.minLevel,
		next:     l.next.WithAttrs(attrs),
	}
}

func (l *levelFilterHandler) WithGroup(name string) slog.Handler {
	return &levelFilterHandler{
		minLevel: l.minLevel,
		next:     l.next.WithGroup(name),
	}
}
