package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	propertycoveragerepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/property_coverage_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PropertyCoverageAdapter) CreateHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, zip propertycoveragemodel.HorizontalComplexZipCodeInterface) (int64, error) {
	const query = `
		INSERT INTO horizontal_complex_zip_codes (
			horizontal_complex_id, zip_code
		) VALUES (?, ?);
	`

	return a.execInsert(ctx, tx, query, zip.HorizontalComplexID(), zip.ZipCode())
}

func (a *PropertyCoverageAdapter) UpdateHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, zip propertycoveragemodel.HorizontalComplexZipCodeInterface) (int64, error) {
	const query = `
		UPDATE horizontal_complex_zip_codes SET
			horizontal_complex_id = ?,
			zip_code = ?
		WHERE id = ?
		LIMIT 1;
	`

	return a.execUpdate(ctx, tx, query, zip.HorizontalComplexID(), zip.ZipCode(), zip.ID())
}

func (a *PropertyCoverageAdapter) DeleteHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, id int64) (int64, error) {
	const query = "DELETE FROM horizontal_complex_zip_codes WHERE id = ? LIMIT 1;"
	return a.execUpdate(ctx, tx, query, id)
}

func (a *PropertyCoverageAdapter) GetHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.HorizontalComplexZipCodeInterface, error) {
	const query = `
		SELECT id, horizontal_complex_id, zip_code
		FROM horizontal_complex_zip_codes
		WHERE id = ?
		LIMIT 1;
	`

	row := a.QueryRowContext(ctx, tx, "select", query, id)
	var entity propertycoverageentities.HorizontalComplexZipCodeEntity
	if err := row.Scan(&entity.ID, &entity.HorizontalComplexID, &entity.ZipCode); err != nil {
		return nil, err
	}

	return propertycoverageconverters.HorizontalComplexZipCodeEntityToDomain(entity), nil
}

func (a *PropertyCoverageAdapter) ListHorizontalComplexZipCodes(ctx context.Context, tx *sql.Tx, params propertycoveragerepository.ListHorizontalComplexZipCodesParams) ([]propertycoveragemodel.HorizontalComplexZipCodeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	var builder strings.Builder
	args := []any{params.HorizontalComplexID}

	builder.WriteString("SELECT id, horizontal_complex_id, zip_code FROM horizontal_complex_zip_codes WHERE horizontal_complex_id = ?")

	if strings.TrimSpace(params.ZipCode) != "" {
		builder.WriteString(" AND zip_code = ?")
		args = append(args, params.ZipCode)
	}

	builder.WriteString(" ORDER BY id DESC")

	if params.Limit > 0 {
		builder.WriteString(" LIMIT ?")
		args = append(args, params.Limit)
	}

	if params.Offset > 0 {
		builder.WriteString(" OFFSET ?")
		args = append(args, params.Offset)
	}

	rows, err := a.QueryContext(ctx, tx, "select", builder.String(), args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.zip.list.query_error", "err", err)
		return nil, fmt.Errorf("list horizontal zip codes: %w", err)
	}
	defer rows.Close()

	result := make([]propertycoveragemodel.HorizontalComplexZipCodeInterface, 0)

	for rows.Next() {
		var entity propertycoverageentities.HorizontalComplexZipCodeEntity
		if err := rows.Scan(&entity.ID, &entity.HorizontalComplexID, &entity.ZipCode); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.property_coverage.zip.list.scan_error", "err", err)
			return nil, fmt.Errorf("scan zip code row: %w", err)
		}
		result = append(result, propertycoverageconverters.HorizontalComplexZipCodeEntityToDomain(entity))
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.zip.list.rows_error", "err", err)
		return nil, fmt.Errorf("iterate zip code rows: %w", err)
	}

	return result, nil
}

func (a *PropertyCoverageAdapter) listZipCodes(ctx context.Context, tx *sql.Tx, complexID int64) ([]propertycoveragemodel.HorizontalComplexZipCodeInterface, error) {
	params := propertycoveragerepository.ListHorizontalComplexZipCodesParams{
		HorizontalComplexID: complexID,
	}
	return a.ListHorizontalComplexZipCodes(ctx, tx, params)
}
