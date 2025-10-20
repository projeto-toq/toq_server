package mysqlcomplexadapter

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
)

func (ca *ComplexAdapter) CreateComplex(ctx context.Context, tx *sql.Tx, complex complexmodel.ComplexInterface) (int64, error) {
	query := `INSERT INTO complex (
		name,
		zip_code,
		street,
		number,
		neighborhood,
		city,
		state,
		reception_phone,
		sector,
		main_registration,
		type
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	return ca.Create(
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
	)
}
