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
}
