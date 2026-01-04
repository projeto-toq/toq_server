package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"
	"fmt"

	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// fetchManagedComplex executa query single-row para complexos e preenche Kind informado.
// Retorna sql.ErrNoRows quando não há correspondência.
func (a *PropertyCoverageAdapter) fetchManagedComplex(ctx context.Context, tx *sql.Tx, query string, args []any, kind propertycoveragemodel.CoverageKind) (propertycoverageentities.ManagedComplexEntity, error) {
	logger := utils.LoggerFromContext(ctx)

	row := a.QueryRowContext(ctx, tx, "select", query, args...)
	var entity propertycoverageentities.ManagedComplexEntity
	entity.Kind = kind

	if err := row.Scan(
		&entity.ID,
		&entity.Name,
		&entity.ZipCode,
		&entity.Street,
		&entity.Number,
		&entity.Neighborhood,
		&entity.City,
		&entity.State,
		&entity.ReceptionPhone,
		&entity.Sector,
		&entity.MainRegistration,
		&entity.PropertyTypes,
	); err != nil {
		if err != sql.ErrNoRows {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.property_coverage.admin.fetch.scan_error", "kind", kind, "err", err)
		}
		return entity, err
	}

	return entity, nil
}

// execInsert executa INSERT usando InstrumentedAdapter e retorna LastInsertId; retorna erro com span marcado.
func (a *PropertyCoverageAdapter) execInsert(ctx context.Context, tx *sql.Tx, query string, args ...any) (int64, error) {
	logger := utils.LoggerFromContext(ctx)

	result, err := a.ExecContext(ctx, tx, "insert", query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.admin.insert_error", "err", err)
		return 0, fmt.Errorf("insert managed complex data: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.admin.last_insert_id_error", "err", err)
		return 0, fmt.Errorf("managed complex last insert id: %w", err)
	}

	return id, nil
}

// execUpdate executa UPDATE/DELETE e retorna linhas afetadas; retorna sql.ErrNoRows quando 0 linhas afetadas.
func (a *PropertyCoverageAdapter) execUpdate(ctx context.Context, tx *sql.Tx, op string, query string, args ...any) (int64, error) {
	logger := utils.LoggerFromContext(ctx)

	result, err := a.ExecContext(ctx, tx, op, query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.admin.update_error", "op", op, "err", err)
		return 0, fmt.Errorf("%s managed complex data: %w", op, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.admin.rows_affected_error", "op", op, "err", err)
		return 0, fmt.Errorf("managed complex %s rows affected: %w", op, err)
	}

	if affected == 0 {
		return 0, sql.ErrNoRows
	}

	return affected, nil
}
