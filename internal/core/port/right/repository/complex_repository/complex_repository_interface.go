package complexrepository

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
)

type ComplexRepoPortInterface interface {
	GetHorizontalComplex(ctx context.Context, tx *sql.Tx, zipCode string) (complex complexmodel.ComplexInterface, err error)
	GetVerticalComplex(ctx context.Context, tx *sql.Tx, zipCode string, number string) (complex complexmodel.ComplexInterface, err error)
	CreateComplex(ctx context.Context, tx *sql.Tx, complex complexmodel.ComplexInterface) (int64, error)
	UpdateComplex(ctx context.Context, tx *sql.Tx, complex complexmodel.ComplexInterface) (int64, error)
	DeleteComplex(ctx context.Context, tx *sql.Tx, id int64) (int64, error)
	GetComplexByID(ctx context.Context, tx *sql.Tx, id int64) (complexmodel.ComplexInterface, error)
	ListComplexes(ctx context.Context, tx *sql.Tx, params ListComplexesParams) ([]complexmodel.ComplexInterface, error)
	CreateComplexTower(ctx context.Context, tx *sql.Tx, tower complexmodel.ComplexTowerInterface) (int64, error)
	UpdateComplexTower(ctx context.Context, tx *sql.Tx, tower complexmodel.ComplexTowerInterface) (int64, error)
	DeleteComplexTower(ctx context.Context, tx *sql.Tx, id int64) (int64, error)
	GetComplexTowerByID(ctx context.Context, tx *sql.Tx, id int64) (complexmodel.ComplexTowerInterface, error)
	ListComplexTowers(ctx context.Context, tx *sql.Tx, params ListComplexTowersParams) ([]complexmodel.ComplexTowerInterface, error)
	CreateComplexSize(ctx context.Context, tx *sql.Tx, size complexmodel.ComplexSizeInterface) (int64, error)
	UpdateComplexSize(ctx context.Context, tx *sql.Tx, size complexmodel.ComplexSizeInterface) (int64, error)
	DeleteComplexSize(ctx context.Context, tx *sql.Tx, id int64) (int64, error)
	GetComplexSizeByID(ctx context.Context, tx *sql.Tx, id int64) (complexmodel.ComplexSizeInterface, error)
	ListComplexSizes(ctx context.Context, tx *sql.Tx, params ListComplexSizesParams) ([]complexmodel.ComplexSizeInterface, error)
	CreateComplexZipCode(ctx context.Context, tx *sql.Tx, zip complexmodel.ComplexZipCodeInterface) (int64, error)
	UpdateComplexZipCode(ctx context.Context, tx *sql.Tx, zip complexmodel.ComplexZipCodeInterface) (int64, error)
	DeleteComplexZipCode(ctx context.Context, tx *sql.Tx, id int64) (int64, error)
	GetComplexZipCodeByID(ctx context.Context, tx *sql.Tx, id int64) (complexmodel.ComplexZipCodeInterface, error)
	ListComplexZipCodes(ctx context.Context, tx *sql.Tx, params ListComplexZipCodesParams) ([]complexmodel.ComplexZipCodeInterface, error)
}
