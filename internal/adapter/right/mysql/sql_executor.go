package mysqladapter

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

type sqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// SQLExecutor centraliza a execução instrumentada de comandos SQL para todos os adapters MySQL.
type SQLExecutor struct {
	db      *Database
	metrics metricsport.MetricsPortInterface
}

// NewSQLExecutor cria um executor instrumentado compartilhado.
func NewSQLExecutor(db *Database, metrics metricsport.MetricsPortInterface) SQLExecutor {
	return SQLExecutor{db: db, metrics: metrics}
}

// ExecContext executa comandos que não retornam linhas (INSERT/UPDATE/DELETE).
func (e SQLExecutor) ExecContext(ctx context.Context, tx *sql.Tx, operation, query string, args ...any) (sql.Result, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	observer := e.observe(operation, query)
	defer observer()

	executor := e.pickExecutor(tx)
	result, err := executor.ExecContext(ctx, query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.executor.exec_error", "query", query, "err", err)
		return nil, err
	}

	return result, nil
}

// QueryContext executa comandos SELECT que retornam múltiplas linhas.
func (e SQLExecutor) QueryContext(ctx context.Context, tx *sql.Tx, operation, query string, args ...any) (*sql.Rows, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	observer := e.observe(operation, query)
	defer observer()

	executor := e.pickExecutor(tx)
	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.executor.query_error", "query", query, "err", err)
		return nil, err
	}

	return rows, nil
}

// QueryRowContext executa um SELECT que retorna uma única linha.
func (e SQLExecutor) QueryRowContext(ctx context.Context, tx *sql.Tx, operation, query string, args ...any) *sql.Row {
	ctxWithTracer, spanEnd, err := utils.GenerateTracer(ctx)
	if err == nil {
		ctx = ctxWithTracer
	} else {
		slog.Warn("mysql.executor.query_row.tracer_error", "err", err)
		spanEnd = func() {}
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	observer := e.observe(operation, query)
	defer observer()
	defer spanEnd()

	executor := e.pickExecutor(tx)
	row := executor.QueryRowContext(ctx, query, args...)
	if row == nil {
		err := sql.ErrNoRows
		utils.SetSpanError(ctx, err)
		logger.Warn("mysql.executor.query_row.nil_row", "query", query)
	}

	return row
}

// PrepareContext cria um statement preparado com instrumentação compartilhada.
func (e SQLExecutor) PrepareContext(ctx context.Context, tx *sql.Tx, operation, query string) (*sql.Stmt, func(), error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, nil, err
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	observer := e.observe(operation, query)
	executor := e.pickExecutor(tx)

	stmt, err := executor.PrepareContext(ctx, query)
	if err != nil {
		spanEnd()
		observer()
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.executor.prepare_error", "query", query, "err", err)
		return nil, nil, err
	}

	cleanup := func() {
		spanEnd()
		observer()
		if closeErr := stmt.Close(); closeErr != nil {
			logger.Warn("mysql.executor.prepare_close_error", "err", closeErr)
		}
	}

	return stmt, cleanup, nil
}

func (e SQLExecutor) QueryRowContextWithScan(ctx context.Context, tx *sql.Tx, operation, query string, scan func(*sql.Row) error, args ...any) error {
	row := e.QueryRowContext(ctx, tx, operation, query, args...)
	if row == nil {
		return sql.ErrNoRows
	}
	return scan(row)
}

func (e SQLExecutor) pickExecutor(tx *sql.Tx) sqlExecutor {
	if tx != nil {
		return tx
	}
	return e.db.GetDB()
}

func (e SQLExecutor) observe(operation, query string) func() {
	if e.metrics == nil {
		return func() {}
	}

	table := extractTableName(query)
	if table == "" {
		table = "unknown"
	}

	start := time.Now()
	return func() {
		e.metrics.IncrementDatabaseQueries(operation, table)
		e.metrics.ObserveDatabaseQueryDuration(operation, table, time.Since(start))
	}
}
