package mysqlcomplexadapter

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
)

func (ca *ComplexAdapter) UpdateComplex(ctx context.Context, tx *sql.Tx, complex complexmodel.ComplexInterface) (int64, error) {
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

	return ca.Update(
		ctx,
		tx,
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
}
