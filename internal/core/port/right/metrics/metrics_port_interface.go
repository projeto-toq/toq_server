package metricsport

import (
	"context"
	"time"
)

// MetricsPortInterface define a interface para coleta de métricas do sistema
// Segue arquitetura hexagonal com separação entre port e adapter
type MetricsPortInterface interface {
	// HTTP Metrics
	IncrementHTTPRequests(method, path, status string)
	ObserveHTTPDuration(method, path string, duration time.Duration)
	SetHTTPRequestsInFlight(count int64)
	IncrementHTTPRequestsInFlight()
	DecrementHTTPRequestsInFlight()
	ObserveHTTPResponseSize(method, path string, size int64)

	// Business Metrics
	SetActiveSessions(count int64)
	IncrementDatabaseQueries(operation, table string)
	IncrementCacheOperations(operation, result string)

	// Email change flow metrics
	IncrementEmailChangeRequest(result string)
	IncrementEmailChangeConfirm(result string)
	IncrementEmailChangeResend(result string)

	// Phone change flow metrics
	IncrementPhoneChangeRequest(result string)
	IncrementPhoneChangeConfirm(result string)
	IncrementPhoneChangeResend(result string)

	// System Metrics
	SetSystemUptime(duration time.Duration)
	IncrementErrors(component, errorType string)

	// Lifecycle
	Initialize(ctx context.Context) error
	Shutdown(ctx context.Context) error
	GetMetricsHandler() interface{} // Retorna handler para endpoint /metrics
}
