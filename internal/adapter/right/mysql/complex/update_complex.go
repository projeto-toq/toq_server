package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) UpdateComplex(ctx context.Context, tx *sql.Tx, complex complexmodel.ComplexInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE complex SET
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
	LIMIT 1;`

	result, err := ca.ExecContext(
		ctx,
		tx,
		"update",
		query,
		complex.Name(),
		complex.ZipCode(),
		nullableStringValue(complex.Street()),
		complex.Number(),
		nullableStringValue(complex.Neighborhood()),
		complex.City(),
		complex.State(),
		nullableStringValue(complex.PhoneNumber()),
		complex.Sector(),
		nullableStringValue(complex.MainRegistration()),
		complex.GetPropertyType(),
		complex.ID(),
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.update.exec_error", "error", err, "complex_id", complex.ID())
		return 0, fmt.Errorf("update complex: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.update.rows_affected_error", "error", err, "complex_id", complex.ID())
		return 0, fmt.Errorf("complex rows affected: %w", err)
	}

	return affected, nil
}
