package prometheusadapter

import (
	"context"
	"sync/atomic"
	"time"

	metricsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusAdapter implementa MetricsPortInterface usando Prometheus
type PrometheusAdapter struct {
	registry *prometheus.Registry

	// Thread-safe counter for requests in flight
	requestsInFlightCounter int64

	// HTTP Metrics
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration  *prometheus.HistogramVec
	httpRequestsInFlight prometheus.Gauge
	httpResponseSize     *prometheus.HistogramVec

	// Business Metrics (kept)
	activeSessions       prometheus.Gauge
	databaseQueriesTotal *prometheus.CounterVec
	cacheOperationsTotal *prometheus.CounterVec

	// System Metrics
	systemUptime prometheus.Gauge
	errorsTotal  *prometheus.CounterVec
}

// NewPrometheusAdapter cria uma nova instância do adapter Prometheus
func NewPrometheusAdapter() metricsport.MetricsPortInterface {
	registry := prometheus.NewRegistry()

	adapter := &PrometheusAdapter{
		registry:                registry,
		requestsInFlightCounter: 0, // Explicitly initialize counter to 0
	}

	adapter.initializeMetrics()
	adapter.registerMetrics()

	// Ensure the gauge starts at 0
	adapter.httpRequestsInFlight.Set(0)

	return adapter
}

// initializeMetrics inicializa todas as métricas Prometheus
func (p *PrometheusAdapter) initializeMetrics() {
	// HTTP Metrics
	p.httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	p.httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	p.httpRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)

	p.httpResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "path"},
	)

	// Business Metrics
	p.activeSessions = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_sessions_total",
			Help: "Current number of active user sessions",
		},
	)

	p.databaseQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_queries_total",
			Help: "Total number of database queries executed",
		},
		[]string{"operation", "table"},
	)

	p.cacheOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "Total number of cache operations",
		},
		[]string{"operation", "result"},
	)

	// Removed business flow counters (email/phone/password) to avoid duplication with HTTP metrics

	// System Metrics
	p.systemUptime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_uptime_seconds",
			Help: "System uptime in seconds",
		},
	)

	p.errorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors by component and type",
		},
		[]string{"component", "error_type"},
	)
}

// registerMetrics registra todas as métricas no registry
func (p *PrometheusAdapter) registerMetrics() {
	p.registry.MustRegister(
		p.httpRequestsTotal,
		p.httpRequestDuration,
		p.httpRequestsInFlight,
		p.httpResponseSize,
		p.activeSessions,
		p.databaseQueriesTotal,
		p.cacheOperationsTotal,
		p.systemUptime,
		p.errorsTotal,
	)
}

// HTTP Metrics Implementation
func (p *PrometheusAdapter) IncrementHTTPRequests(method, path, status string) {
	p.httpRequestsTotal.WithLabelValues(method, path, status).Inc()
}

func (p *PrometheusAdapter) ObserveHTTPDuration(method, path string, duration time.Duration) {
	p.httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

func (p *PrometheusAdapter) SetHTTPRequestsInFlight(count int64) {
	atomic.StoreInt64(&p.requestsInFlightCounter, count)
	p.httpRequestsInFlight.Set(float64(count))
}

func (p *PrometheusAdapter) IncrementHTTPRequestsInFlight() {
	newCount := atomic.AddInt64(&p.requestsInFlightCounter, 1)
	p.httpRequestsInFlight.Set(float64(newCount))
}

func (p *PrometheusAdapter) DecrementHTTPRequestsInFlight() {
	newCount := atomic.AddInt64(&p.requestsInFlightCounter, -1)
	// Garantir que nunca seja negativo
	if newCount < 0 {
		atomic.StoreInt64(&p.requestsInFlightCounter, 0)
		newCount = 0
	}
	p.httpRequestsInFlight.Set(float64(newCount))
}

func (p *PrometheusAdapter) ObserveHTTPResponseSize(method, path string, size int64) {
	p.httpResponseSize.WithLabelValues(method, path).Observe(float64(size))
}

// Business Metrics Implementation
func (p *PrometheusAdapter) SetActiveSessions(count int64) {
	p.activeSessions.Set(float64(count))
}

func (p *PrometheusAdapter) IncrementDatabaseQueries(operation, table string) {
	p.databaseQueriesTotal.WithLabelValues(operation, table).Inc()
}

func (p *PrometheusAdapter) IncrementCacheOperations(operation, result string) {
	p.cacheOperationsTotal.WithLabelValues(operation, result).Inc()
}

// Removed flow metrics methods (email/phone/password)

// System Metrics Implementation
func (p *PrometheusAdapter) SetSystemUptime(duration time.Duration) {
	p.systemUptime.Set(duration.Seconds())
}

func (p *PrometheusAdapter) IncrementErrors(component, errorType string) {
	p.errorsTotal.WithLabelValues(component, errorType).Inc()
}

// Lifecycle Implementation
func (p *PrometheusAdapter) Initialize(ctx context.Context) error {
	// Métricas já inicializadas no construtor
	return nil
}

func (p *PrometheusAdapter) Shutdown(ctx context.Context) error {
	// Prometheus não requer shutdown específico
	return nil
}

func (p *PrometheusAdapter) GetMetricsHandler() interface{} {
	return promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})
}
