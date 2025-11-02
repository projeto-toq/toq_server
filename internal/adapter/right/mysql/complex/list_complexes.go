package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	complexrepoconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	repository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/complex_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) ListComplexes(ctx context.Context, tx *sql.Tx, params repository.ListComplexesParams) ([]complexmodel.ComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	var builder strings.Builder
	builder.WriteString("SELECT id, name, zip_code, street, number, neighborhood, city, state, reception_phone, sector, main_registration, type FROM complex WHERE 1=1")
	args := make([]any, 0)

	if params.Name != "" {
		builder.WriteString(" AND name LIKE ?")
		args = append(args, fmt.Sprintf("%%%s%%", params.Name))
	}

	if params.ZipCode != "" {
		builder.WriteString(" AND zip_code = ?")
		args = append(args, params.ZipCode)
	}

	if params.City != "" {
		builder.WriteString(" AND city LIKE ?")
		args = append(args, fmt.Sprintf("%%%s%%", params.City))
	}

	if params.State != "" {
		builder.WriteString(" AND state = ?")
		args = append(args, params.State)
	}

	if params.Sector != nil {
		builder.WriteString(" AND sector = ?")
		args = append(args, *params.Sector)
	}

	if params.PropertyType != nil {
		builder.WriteString(" AND type = ?")
		args = append(args, *params.PropertyType)
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

	query := builder.String()

	rows, err := ca.QueryContext(ctx, tx, "select", query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.list.read_error", "error", err, "params", params)
		return nil, fmt.Errorf("list complexes query: %w", err)
	}
	defer rows.Close()

	entities, err := rowsToEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.list.scan_error", "error", err, "params", params)
		return nil, fmt.Errorf("scan complexes rows: %w", err)
	}

	complexes := make([]complexmodel.ComplexInterface, 0, len(entities))

	for _, entity := range entities {
		complexDomain, errConv := complexrepoconverters.ComplexEntityToDomain(entity)
		if errConv != nil {
			utils.SetSpanError(ctx, errConv)
			logger.Error("mysql.complex.list.convert_error", "error", errConv)
			return nil, fmt.Errorf("convert complex entity: %w", errConv)
		}

		complexes = append(complexes, complexDomain)
	}

	return complexes, nil
}
