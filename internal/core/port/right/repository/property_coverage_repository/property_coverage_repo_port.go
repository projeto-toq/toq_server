package propertycoveragerepository

import (
	"context"
	"database/sql"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

// RepositoryInterface defines the persistence contract for property coverage lookups.
type RepositoryInterface interface {
	GetVerticalCoverage(ctx context.Context, tx *sql.Tx, zipCode, number string) (propertycoveragemodel.CoverageInterface, error)
	GetHorizontalCoverage(ctx context.Context, tx *sql.Tx, zipCode string) (propertycoveragemodel.CoverageInterface, error)
	GetNoComplexCoverage(ctx context.Context, tx *sql.Tx, zipCode string) (propertycoveragemodel.CoverageInterface, error)
	GetVerticalComplexByZipNumber(ctx context.Context, tx *sql.Tx, zipCode, number string) (propertycoveragemodel.ManagedComplexInterface, error)
	GetHorizontalComplexByZip(ctx context.Context, tx *sql.Tx, zipCode string) (propertycoveragemodel.ManagedComplexInterface, error)

	ListManagedComplexes(ctx context.Context, tx *sql.Tx, params ListManagedComplexesParams) ([]propertycoveragemodel.ManagedComplexInterface, error)
	GetManagedComplex(ctx context.Context, tx *sql.Tx, id int64, kind propertycoveragemodel.CoverageKind) (propertycoveragemodel.ManagedComplexInterface, error)
	CreateManagedComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error)
	UpdateManagedComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error)
	DeleteManagedComplex(ctx context.Context, tx *sql.Tx, id int64, kind propertycoveragemodel.CoverageKind) (int64, error)

	CreateVerticalComplexTower(ctx context.Context, tx *sql.Tx, tower propertycoveragemodel.VerticalComplexTowerInterface) (int64, error)
	UpdateVerticalComplexTower(ctx context.Context, tx *sql.Tx, tower propertycoveragemodel.VerticalComplexTowerInterface) (int64, error)
	DeleteVerticalComplexTower(ctx context.Context, tx *sql.Tx, id int64) (int64, error)
	GetVerticalComplexTower(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.VerticalComplexTowerInterface, error)
	ListVerticalComplexTowers(ctx context.Context, tx *sql.Tx, params ListVerticalComplexTowersParams) ([]propertycoveragemodel.VerticalComplexTowerInterface, error)

	CreateVerticalComplexSize(ctx context.Context, tx *sql.Tx, size propertycoveragemodel.VerticalComplexSizeInterface) (int64, error)
	UpdateVerticalComplexSize(ctx context.Context, tx *sql.Tx, size propertycoveragemodel.VerticalComplexSizeInterface) (int64, error)
	DeleteVerticalComplexSize(ctx context.Context, tx *sql.Tx, id int64) (int64, error)
	GetVerticalComplexSize(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.VerticalComplexSizeInterface, error)
	ListVerticalComplexSizes(ctx context.Context, tx *sql.Tx, params ListVerticalComplexSizesParams) ([]propertycoveragemodel.VerticalComplexSizeInterface, error)

	CreateHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, zip propertycoveragemodel.HorizontalComplexZipCodeInterface) (int64, error)
	UpdateHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, zip propertycoveragemodel.HorizontalComplexZipCodeInterface) (int64, error)
	DeleteHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, id int64) (int64, error)
	GetHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.HorizontalComplexZipCodeInterface, error)
	ListHorizontalComplexZipCodes(ctx context.Context, tx *sql.Tx, params ListHorizontalComplexZipCodesParams) ([]propertycoveragemodel.HorizontalComplexZipCodeInterface, error)
}
