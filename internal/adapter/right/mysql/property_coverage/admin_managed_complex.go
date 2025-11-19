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

	query := fmt.Sprintf("SELECT * FROM (%s) managed ORDER BY managed.coverage_kind, managed.id DESC", strings.Join(clauses, " UNION ALL "))

	if params.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, params.Limit)
	}

	if params.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, params.Offset)
	}

	rows, err := a.QueryContext(ctx, tx, "select", query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.admin.list.query_error", "err", err)
		return nil, fmt.Errorf("list managed complexes: %w", err)
	}
	defer rows.Close()

	result := make([]propertycoveragemodel.ManagedComplexInterface, 0)

	for rows.Next() {
		var entity propertycoverageentities.ManagedComplexEntity
		var coverageKind string
		if err := rows.Scan(
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
		); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.property_coverage.admin.list.scan_error", "err", err)
			return nil, fmt.Errorf("scan managed complex row: %w", err)
		}

		entity.Kind = propertycoveragemodel.CoverageKind(coverageKind)
		domain := propertycoverageconverters.ManagedComplexEntityToDomain(entity)
		result = append(result, domain)
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.admin.list.rows_error", "err", err)
		return nil, fmt.Errorf("iterate managed complex rows: %w", err)
	}

	return result, nil
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

func (a *PropertyCoverageAdapter) GetManagedComplex(ctx context.Context, tx *sql.Tx, id int64, kind propertycoveragemodel.CoverageKind) (propertycoveragemodel.ManagedComplexInterface, error) {
	switch kind {
	case propertycoveragemodel.CoverageKindVertical:
		return a.getVerticalComplexByID(ctx, tx, id)
	case propertycoveragemodel.CoverageKindHorizontal:
		return a.getHorizontalComplexByID(ctx, tx, id)
	case propertycoveragemodel.CoverageKindStandalone:
		return a.getStandaloneComplexByID(ctx, tx, id)
	default:
		return nil, sql.ErrNoRows
	}
}

func (a *PropertyCoverageAdapter) getVerticalComplexByID(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.ManagedComplexInterface, error) {
	const query = `
		SELECT vc.id, vc.name, vc.zip_code, vc.street, vc.number, vc.neighborhood, vc.city, vc.state,
		       vc.reception_phone, vc.sector, vc.main_registration, vc.type
		FROM vertical_complexes vc
		WHERE vc.id = ?
		LIMIT 1;
	`

	entity, err := a.fetchManagedComplex(ctx, tx, query, []any{id}, propertycoveragemodel.CoverageKindVertical)
	if err != nil {
		return nil, err
	}

	domain := propertycoverageconverters.ManagedComplexEntityToDomain(entity)

	towers, err := a.listTowers(ctx, tx, entity.ID)
	if err != nil {
		return nil, err
	}
	domain.SetTowers(towers)

	sizes, err := a.listSizes(ctx, tx, entity.ID)
	if err != nil {
		return nil, err
	}
	domain.SetSizes(sizes)

	return domain, nil
}

func (a *PropertyCoverageAdapter) getHorizontalComplexByID(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.ManagedComplexInterface, error) {
	const query = `
		SELECT hc.id, hc.name, hc.zip_code, hc.street, hc.number, hc.neighborhood, hc.city, hc.state,
		       hc.reception_phone, hc.sector, hc.main_registration, hc.type
		FROM horizontal_complexes hc
		WHERE hc.id = ?
		LIMIT 1;
	`

	entity, err := a.fetchManagedComplex(ctx, tx, query, []any{id}, propertycoveragemodel.CoverageKindHorizontal)
	if err != nil {
		return nil, err
	}

	domain := propertycoverageconverters.ManagedComplexEntityToDomain(entity)

	zipCodes, err := a.listZipCodes(ctx, tx, entity.ID)
	if err != nil {
		return nil, err
	}
	domain.SetZipCodes(zipCodes)

	return domain, nil
}

func (a *PropertyCoverageAdapter) getStandaloneComplexByID(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.ManagedComplexInterface, error) {
	const query = `
		SELECT nc.id, NULL AS name, nc.zip_code, NULL AS street, NULL AS number,
		       nc.neighborhood, nc.city, nc.state, NULL AS reception_phone, nc.sector,
		       NULL AS main_registration, nc.type
		FROM no_complex_zip_codes nc
		WHERE nc.id = ?
		LIMIT 1;
	`

	entity, err := a.fetchManagedComplex(ctx, tx, query, []any{id}, propertycoveragemodel.CoverageKindStandalone)
	if err != nil {
		return nil, err
	}

	return propertycoverageconverters.ManagedComplexEntityToDomain(entity), nil
}

func (a *PropertyCoverageAdapter) fetchManagedComplex(ctx context.Context, tx *sql.Tx, query string, args []any, kind propertycoveragemodel.CoverageKind) (propertycoverageentities.ManagedComplexEntity, error) {
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
		return entity, err
	}

	return entity, nil
}

func (a *PropertyCoverageAdapter) CreateManagedComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error) {
	switch entity.Kind() {
	case propertycoveragemodel.CoverageKindVertical:
		return a.insertVerticalComplex(ctx, tx, entity)
	case propertycoveragemodel.CoverageKindHorizontal:
		return a.insertHorizontalComplex(ctx, tx, entity)
	case propertycoveragemodel.CoverageKindStandalone:
		return a.insertStandaloneComplex(ctx, tx, entity)
	default:
		return 0, fmt.Errorf("unsupported coverage kind %s", entity.Kind())
	}
}

func (a *PropertyCoverageAdapter) insertVerticalComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error) {
	const query = `
		INSERT INTO vertical_complexes (
			name, zip_code, street, number, neighborhood, city, state,
			reception_phone, sector, main_registration, type
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	return a.execInsert(ctx, tx, query,
		entity.Name(),
		entity.ZipCode(),
		entity.Street(),
		entity.Number(),
		nullableString(entity.Neighborhood()),
		entity.City(),
		entity.State(),
		nullableString(entity.ReceptionPhone()),
		entity.Sector(),
		nullableString(entity.MainRegistration()),
		entity.PropertyTypes(),
	)
}

func (a *PropertyCoverageAdapter) insertHorizontalComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error) {
	const query = `
		INSERT INTO horizontal_complexes (
			name, zip_code, street, number, neighborhood, city, state,
			reception_phone, sector, main_registration, type
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	return a.execInsert(ctx, tx, query,
		entity.Name(),
		entity.ZipCode(),
		entity.Street(),
		nullableString(entity.Number()),
		nullableString(entity.Neighborhood()),
		entity.City(),
		entity.State(),
		nullableString(entity.ReceptionPhone()),
		entity.Sector(),
		nullableString(entity.MainRegistration()),
		entity.PropertyTypes(),
	)
}

func (a *PropertyCoverageAdapter) insertStandaloneComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error) {
	const query = `
		INSERT INTO no_complex_zip_codes (
			zip_code, neighborhood, city, state, sector, type
		) VALUES (?, ?, ?, ?, ?, ?);
	`

	return a.execInsert(ctx, tx, query,
		entity.ZipCode(),
		nullableString(entity.Neighborhood()),
		entity.City(),
		entity.State(),
		entity.Sector(),
		entity.PropertyTypes(),
	)
}

func (a *PropertyCoverageAdapter) execInsert(ctx context.Context, tx *sql.Tx, query string, args ...any) (int64, error) {
	logger := utils.LoggerFromContext(ctx)
	result, err := a.ExecContext(ctx, tx, "insert", query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.admin.insert_error", "err", err)
		return 0, fmt.Errorf("insert managed complex: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.admin.last_insert_id_error", "err", err)
		return 0, fmt.Errorf("managed complex last insert id: %w", err)
	}

	return id, nil
}

func (a *PropertyCoverageAdapter) UpdateManagedComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error) {
	switch entity.Kind() {
	case propertycoveragemodel.CoverageKindVertical:
		return a.updateVerticalComplex(ctx, tx, entity)
	case propertycoveragemodel.CoverageKindHorizontal:
		return a.updateHorizontalComplex(ctx, tx, entity)
	case propertycoveragemodel.CoverageKindStandalone:
		return a.updateStandaloneComplex(ctx, tx, entity)
	default:
		return 0, fmt.Errorf("unsupported coverage kind %s", entity.Kind())
	}
}

func (a *PropertyCoverageAdapter) updateVerticalComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error) {
	const query = `
		UPDATE vertical_complexes SET
			name = ?,
			zip_code = ?,
			street = ?,
			number = ?,
			neighborhood = ?,
			city = ?,
			state = ?,
			reception_phone = ?,
			sector = ?,
			main_registration = ?,
			type = ?
		WHERE id = ?
		LIMIT 1;
	`

	return a.execUpdate(ctx, tx, query,
		entity.Name(),
		entity.ZipCode(),
		entity.Street(),
		entity.Number(),
		nullableString(entity.Neighborhood()),
		entity.City(),
		entity.State(),
		nullableString(entity.ReceptionPhone()),
		entity.Sector(),
		nullableString(entity.MainRegistration()),
		entity.PropertyTypes(),
		entity.ID(),
	)
}

func (a *PropertyCoverageAdapter) updateHorizontalComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error) {
	const query = `
		UPDATE horizontal_complexes SET
			name = ?,
			zip_code = ?,
			street = ?,
			number = ?,
			neighborhood = ?,
			city = ?,
			state = ?,
			reception_phone = ?,
			sector = ?,
			main_registration = ?,
			type = ?
		WHERE id = ?
		LIMIT 1;
	`

	return a.execUpdate(ctx, tx, query,
		entity.Name(),
		entity.ZipCode(),
		entity.Street(),
		nullableString(entity.Number()),
		nullableString(entity.Neighborhood()),
		entity.City(),
		entity.State(),
		nullableString(entity.ReceptionPhone()),
		entity.Sector(),
		nullableString(entity.MainRegistration()),
		entity.PropertyTypes(),
		entity.ID(),
	)
}

func (a *PropertyCoverageAdapter) updateStandaloneComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error) {
	const query = `
		UPDATE no_complex_zip_codes SET
			zip_code = ?,
			neighborhood = ?,
			city = ?,
			state = ?,
			sector = ?,
			type = ?
		WHERE id = ?
		LIMIT 1;
	`

	return a.execUpdate(ctx, tx, query,
		entity.ZipCode(),
		nullableString(entity.Neighborhood()),
		entity.City(),
		entity.State(),
		entity.Sector(),
		entity.PropertyTypes(),
		entity.ID(),
	)
}

func (a *PropertyCoverageAdapter) execUpdate(ctx context.Context, tx *sql.Tx, query string, args ...any) (int64, error) {
	logger := utils.LoggerFromContext(ctx)
	result, err := a.ExecContext(ctx, tx, "update", query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.admin.update_error", "err", err)
		return 0, fmt.Errorf("update managed complex: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.admin.rows_affected_error", "err", err)
		return 0, fmt.Errorf("managed complex rows affected: %w", err)
	}

	return affected, nil
}

func (a *PropertyCoverageAdapter) DeleteManagedComplex(ctx context.Context, tx *sql.Tx, id int64, kind propertycoveragemodel.CoverageKind) (int64, error) {
	var (
		query string
		args  = []any{id}
	)

	switch kind {
	case propertycoveragemodel.CoverageKindVertical:
		query = "DELETE FROM vertical_complexes WHERE id = ? LIMIT 1;"
	case propertycoveragemodel.CoverageKindHorizontal:
		query = "DELETE FROM horizontal_complexes WHERE id = ? LIMIT 1;"
	case propertycoveragemodel.CoverageKindStandalone:
		query = "DELETE FROM no_complex_zip_codes WHERE id = ? LIMIT 1;"
	default:
		return 0, fmt.Errorf("unsupported coverage kind %s", kind)
	}

	logger := utils.LoggerFromContext(ctx)
	result, err := a.ExecContext(ctx, tx, "delete", query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.admin.delete_error", "err", err)
		return 0, fmt.Errorf("delete managed complex: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.admin.delete_rows_error", "err", err)
		return 0, fmt.Errorf("managed complex delete rows affected: %w", err)
	}

	return affected, nil
}

func (a *PropertyCoverageAdapter) GetVerticalComplexByZipNumber(ctx context.Context, tx *sql.Tx, zipCode, number string) (propertycoveragemodel.ManagedComplexInterface, error) {
	const query = `
		SELECT vc.id, vc.name, vc.zip_code, vc.street, vc.number, vc.neighborhood, vc.city, vc.state,
		       vc.reception_phone, vc.sector, vc.main_registration, vc.type
		FROM vertical_complexes vc
		WHERE vc.zip_code = ?
		  AND (
		        UPPER(REPLACE(TRIM(vc.number), ' ', '')) = ?
		     OR FIND_IN_SET(
		          ?,
		          REPLACE(REPLACE(UPPER(TRIM(vc.number)), ' ', ''), ';', ',')
		        ) > 0
		      )
		LIMIT 1;
	`

	entity, err := a.fetchManagedComplex(ctx, tx, query, []any{zipCode, number, number}, propertycoveragemodel.CoverageKindVertical)
	if err != nil {
		return nil, err
	}

	domain := propertycoverageconverters.ManagedComplexEntityToDomain(entity)

	towers, err := a.listTowers(ctx, tx, entity.ID)
	if err != nil {
		return nil, err
	}
	domain.SetTowers(towers)

	sizes, err := a.listSizes(ctx, tx, entity.ID)
	if err != nil {
		return nil, err
	}
	domain.SetSizes(sizes)

	return domain, nil
}

func (a *PropertyCoverageAdapter) GetHorizontalComplexByZip(ctx context.Context, tx *sql.Tx, zipCode string) (propertycoveragemodel.ManagedComplexInterface, error) {
	const query = `
		SELECT hc.id, hc.name, hc.zip_code, hc.street, hc.number, hc.neighborhood, hc.city, hc.state,
		       hc.reception_phone, hc.sector, hc.main_registration, hc.type
		FROM horizontal_complexes hc
		WHERE hc.zip_code = ?
		LIMIT 1;
	`

	entity, err := a.fetchManagedComplex(ctx, tx, query, []any{zipCode}, propertycoveragemodel.CoverageKindHorizontal)
	if err != nil {
		return nil, err
	}

	domain := propertycoverageconverters.ManagedComplexEntityToDomain(entity)

	zipCodes, err := a.listZipCodes(ctx, tx, entity.ID)
	if err != nil {
		return nil, err
	}
	domain.SetZipCodes(zipCodes)

	return domain, nil
}
