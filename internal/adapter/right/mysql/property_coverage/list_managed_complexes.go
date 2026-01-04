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

// ListManagedComplexes returns managed complexes (vertical/horizontal/standalone) with admin filters.
// Returns empty slice when no rows; never returns sql.ErrNoRows.
func (a *PropertyCoverageAdapter) ListManagedComplexes(ctx context.Context, tx *sql.Tx, params propertycoveragerepository.ListManagedComplexesParams) ([]propertycoveragemodel.ManagedComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	clauses, args := buildListClauses(params)
	if len(clauses) == 0 {
		return []propertycoveragemodel.ManagedComplexInterface{}, nil
	}

	query := fmt.Sprintf(
		"SELECT managed.coverage_kind, managed.id, managed.name, managed.zip_code, managed.street, managed.number, managed.neighborhood, managed.city, managed.state, managed.reception_phone, managed.sector, managed.main_registration, managed.property_types FROM (%s) managed ORDER BY managed.coverage_kind, managed.id DESC",
		strings.Join(clauses, " UNION ALL "),
	)
	query, args = applyPagination(query, args, params.Limit, params.Offset)

	rows, qErr := a.QueryContext(ctx, tx, "select", query, args...)
	if qErr != nil {
		utils.SetSpanError(ctx, qErr)
		logger.Error("mysql.property_coverage.admin.list.query_error", "err", qErr)
		return nil, fmt.Errorf("list managed complexes: %w", qErr)
	}
	defer rows.Close()

	result := make([]propertycoveragemodel.ManagedComplexInterface, 0)
	for rows.Next() {
		var entity propertycoverageentities.ManagedComplexEntity
		var coverageKind string
		if scanErr := rows.Scan(
			&coverageKind,
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
		); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.property_coverage.admin.list.scan_error", "err", scanErr)
			return nil, fmt.Errorf("scan managed complex row: %w", scanErr)
		}

		entity.Kind = propertycoveragemodel.CoverageKind(coverageKind)
		result = append(result, propertycoverageconverters.ManagedComplexEntityToDomain(entity))
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.admin.list.rows_error", "err", err)
		return nil, fmt.Errorf("iterate managed complex rows: %w", err)
	}

	return result, nil
}

// applyPagination appends LIMIT/OFFSET when provided keeping args aligned.
func applyPagination(query string, args []any, limit int, offset int) (string, []any) {
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	if offset > 0 {
		query += " OFFSET ?"
		args = append(args, offset)
	}

	return query, args
}

func buildListClauses(params propertycoveragerepository.ListManagedComplexesParams) ([]string, []any) {
	includeVertical := params.Kind == nil || *params.Kind == propertycoveragemodel.CoverageKindVertical
	includeHorizontal := params.Kind == nil || *params.Kind == propertycoveragemodel.CoverageKindHorizontal
	includeStandalone := params.Kind == nil || *params.Kind == propertycoveragemodel.CoverageKindStandalone

	clauses := make([]string, 0, 3)
	args := make([]any, 0)

	if includeVertical {
		clause, clauseArgs := buildVerticalClause(params)
		clauses = append(clauses, clause)
		args = append(args, clauseArgs...)
	}

	if includeHorizontal {
		clause, clauseArgs := buildHorizontalClause(params)
		clauses = append(clauses, clause)
		args = append(args, clauseArgs...)
	}

	if includeStandalone {
		clause, clauseArgs := buildStandaloneClause(params)
		clauses = append(clauses, clause)
		args = append(args, clauseArgs...)
	}

	return clauses, args
}

func buildVerticalClause(params propertycoveragerepository.ListManagedComplexesParams) (string, []any) {
	var builder strings.Builder
	args := make([]any, 0)

	builder.WriteString("SELECT 'VERTICAL' AS coverage_kind, vc.id, vc.name, vc.zip_code, vc.street, vc.number, vc.neighborhood, vc.city, vc.state, vc.reception_phone, vc.sector, vc.main_registration, vc.type FROM vertical_complexes vc WHERE 1=1")
	appendCommonFilters(&builder, &args, "vc", params)

	if params.Name != "" {
		builder.WriteString(" AND vc.name LIKE ?")
		args = append(args, likePattern(params.Name))
	}

	if params.Number != "" {
		builder.WriteString(" AND vc.number LIKE ?")
		args = append(args, likePattern(params.Number))
	}

	return builder.String(), args
}

func buildHorizontalClause(params propertycoveragerepository.ListManagedComplexesParams) (string, []any) {
	var builder strings.Builder
	args := make([]any, 0)

	builder.WriteString("SELECT 'HORIZONTAL' AS coverage_kind, hc.id, hc.name, hc.zip_code, hc.street, hc.number, hc.neighborhood, hc.city, hc.state, hc.reception_phone, hc.sector, hc.main_registration, hc.type FROM horizontal_complexes hc WHERE 1=1")
	appendCommonFilters(&builder, &args, "hc", params)

	if params.Name != "" {
		builder.WriteString(" AND hc.name LIKE ?")
		args = append(args, likePattern(params.Name))
	}

	if params.Number != "" {
		builder.WriteString(" AND hc.number LIKE ?")
		args = append(args, likePattern(params.Number))
	}

	return builder.String(), args
}

func buildStandaloneClause(params propertycoveragerepository.ListManagedComplexesParams) (string, []any) {
	var builder strings.Builder
	args := make([]any, 0)

	builder.WriteString("SELECT 'STANDALONE' AS coverage_kind, nc.id, NULL AS name, nc.zip_code, NULL AS street, NULL AS number, nc.neighborhood, nc.city, nc.state, NULL AS reception_phone, nc.sector, NULL AS main_registration, nc.type FROM no_complex_zip_codes nc WHERE 1=1")
	appendCommonFilters(&builder, &args, "nc", params)

	return builder.String(), args
}

func appendCommonFilters(builder *strings.Builder, args *[]any, alias string, params propertycoveragerepository.ListManagedComplexesParams) {
	if params.ZipCode != "" {
		builder.WriteString(" AND ")
		builder.WriteString(alias)
		builder.WriteString(".zip_code = ?")
		*args = append(*args, params.ZipCode)
	}

	if params.City != "" {
		builder.WriteString(" AND ")
		builder.WriteString(alias)
		builder.WriteString(".city LIKE ?")
		*args = append(*args, likePattern(params.City))
	}

	if params.State != "" {
		builder.WriteString(" AND ")
		builder.WriteString(alias)
		builder.WriteString(".state = ?")
		*args = append(*args, params.State)
	}

	if params.Sector != nil {
		builder.WriteString(" AND ")
		builder.WriteString(alias)
		builder.WriteString(".sector = ?")
		*args = append(*args, *params.Sector)
	}

	if params.PropertyType != nil {
		builder.WriteString(" AND ")
		builder.WriteString(alias)
		builder.WriteString(".type = ?")
		*args = append(*args, *params.PropertyType)
	}
}
