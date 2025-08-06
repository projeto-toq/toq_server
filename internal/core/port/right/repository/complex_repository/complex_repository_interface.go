package complexrepository

import (
	"context"
	"database/sql"

	complexmodel "github.com/giulio-alfieri/toq_server/internal/core/model/complex_model"
)

type ComplexRepoPortInterface interface {
	GetHorizontalComplex(ctx context.Context, tx *sql.Tx, zipCode string) (complex complexmodel.ComplexInterface, err error)
	GetVerticalComplex(ctx context.Context, tx *sql.Tx, zipCode string, number string) (complex complexmodel.ComplexInterface, err error)
}
