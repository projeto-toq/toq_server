package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"
	"fmt"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteManagedComplex deletes a complex by id/kind; returns sql.ErrNoRows when no row is affected.
func (a *PropertyCoverageAdapter) DeleteManagedComplex(ctx context.Context, tx *sql.Tx, id int64, kind propertycoveragemodel.CoverageKind) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

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

	return a.execUpdate(ctx, tx, "delete", query, args...)
}
