package mysqladapter

import (
	"context"
	"database/sql"
	"strings"
	"time"

	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// InstrumentedAdapter centraliza o envio de métricas para queries SQL executadas pelos adapters.
type InstrumentedAdapter struct {
	db       *Database
	metrics  metricsport.MetricsPortInterface
	executor SQLExecutor
}

// NewInstrumentedAdapter cria uma estrutura auxiliar reutilizável pelos adapters MySQL.
func NewInstrumentedAdapter(db *Database, metrics metricsport.MetricsPortInterface) InstrumentedAdapter {
	return InstrumentedAdapter{
		db:       db,
		metrics:  metrics,
		executor: NewSQLExecutor(db, metrics),
	}
}

// DB retorna o wrapper Database configurado para o adapter.
func (a InstrumentedAdapter) DB() *Database {
	return a.db
}

// Executor retorna o executor compartilhado configurado para o adapter.
func (a *InstrumentedAdapter) Executor() SQLExecutor {
	return a.executor
}

// ExecContext encaminha comandos de escrita para o executor compartilhado.
func (a *InstrumentedAdapter) ExecContext(ctx context.Context, tx *sql.Tx, operation, query string, args ...any) (sql.Result, error) {
	return a.executor.ExecContext(ctx, tx, operation, query, args...)
}

// QueryContext encaminha comandos de leitura que retornam múltiplas linhas.
func (a *InstrumentedAdapter) QueryContext(ctx context.Context, tx *sql.Tx, operation, query string, args ...any) (*sql.Rows, error) {
	return a.executor.QueryContext(ctx, tx, operation, query, args...)
}

// QueryRowContext executa comandos de leitura que retornam uma única linha.
func (a *InstrumentedAdapter) QueryRowContext(ctx context.Context, tx *sql.Tx, operation, query string, args ...any) *sql.Row {
	return a.executor.QueryRowContext(ctx, tx, operation, query, args...)
}

// PrepareContext cria statements preparados com instrumentação padronizada.
func (a *InstrumentedAdapter) PrepareContext(ctx context.Context, tx *sql.Tx, operation, query string) (*sql.Stmt, func(), error) {
	return a.executor.PrepareContext(ctx, tx, operation, query)
}

// Observe registra contagem e duração de uma operação SQL.
func (a InstrumentedAdapter) Observe(operation, query string, duration time.Duration) {
	if a.metrics == nil {
		return
	}

	table := extractTableName(query)
	if table == "" {
		table = "unknown"
	}

	a.metrics.IncrementDatabaseQueries(operation, table)
	a.metrics.ObserveDatabaseQueryDuration(operation, table, duration)
}

// ObserveOnComplete retorna uma função de defer para registrar métricas ao final da operação.
func (a InstrumentedAdapter) ObserveOnComplete(operation, query string) func() {
	if a.metrics == nil {
		return func() {}
	}

	startedAt := time.Now()
	return func() {
		a.Observe(operation, query, time.Since(startedAt))
	}
}

// extractTableName tenta inferir o nome da tabela principal a partir de um statement SQL simples.
func extractTableName(query string) string {
	if query == "" {
		return ""
	}

	lower := strings.ToLower(query)
	lower = strings.TrimSpace(lower)

	tokens := strings.Fields(lower)
	if len(tokens) == 0 {
		return ""
	}

	keywords := map[string]struct{}{
		"from":   {},
		"join":   {},
		"into":   {},
		"update": {},
		"table":  {},
	}

	for idx, token := range tokens {
		if _, ok := keywords[token]; !ok {
			continue
		}

		if token == "table" {
			if idx+1 < len(tokens) {
				return sanitizeTableToken(tokens[idx+1])
			}
			continue
		}

		if idx+1 >= len(tokens) {
			continue
		}

		candidate := tokens[idx+1]

		if candidate == "set" && token == "update" {
			continue
		}

		if strings.HasPrefix(candidate, "(") {
			continue
		}

		return sanitizeTableToken(candidate)
	}

	return ""
}

func sanitizeTableToken(token string) string {
	if token == "" {
		return ""
	}

	sanitized := strings.Trim(token, "`,()")

	if strings.Contains(sanitized, "as") {
		parts := strings.SplitN(sanitized, "as", 2)
		sanitized = strings.TrimSpace(parts[0])
	}

	sanitized = strings.Fields(sanitized)[0]

	if idx := strings.LastIndex(sanitized, "."); idx != -1 && idx+1 < len(sanitized) {
		sanitized = sanitized[idx+1:]
	}

	sanitized = strings.Trim(sanitized, "`,()")

	return sanitized
}
