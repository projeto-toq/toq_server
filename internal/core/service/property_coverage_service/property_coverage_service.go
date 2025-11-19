package propertycoverageservice

import (
	"context"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	propertycoveragerepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/property_coverage_repository"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
)

// PropertyCoverageServiceInterface exposes the available operations for property coverage lookups.
type PropertyCoverageServiceInterface interface {
	ResolvePropertyTypes(ctx context.Context, input propertycoveragemodel.ResolvePropertyTypesInput) (propertycoveragemodel.ResolvePropertyTypesOutput, error)
	ListComplexes(ctx context.Context, input ListComplexesInput) ([]propertycoveragemodel.ManagedComplexInterface, error)
	GetComplexDetail(ctx context.Context, input GetComplexDetailInput) (propertycoveragemodel.ManagedComplexInterface, error)
	CreateComplex(ctx context.Context, input CreateComplexInput) (propertycoveragemodel.ManagedComplexInterface, error)
	UpdateComplex(ctx context.Context, input UpdateComplexInput) (propertycoveragemodel.ManagedComplexInterface, error)
	DeleteComplex(ctx context.Context, input DeleteComplexInput) error
	CreateComplexTower(ctx context.Context, input CreateComplexTowerInput) (propertycoveragemodel.VerticalComplexTowerInterface, error)
	UpdateComplexTower(ctx context.Context, input UpdateComplexTowerInput) (propertycoveragemodel.VerticalComplexTowerInterface, error)
	DeleteComplexTower(ctx context.Context, towerID int64) error
	ListComplexTowers(ctx context.Context, input ListComplexTowersInput) ([]propertycoveragemodel.VerticalComplexTowerInterface, error)
	GetComplexTowerDetail(ctx context.Context, towerID int64) (propertycoveragemodel.VerticalComplexTowerInterface, error)
	CreateComplexSize(ctx context.Context, input CreateComplexSizeInput) (propertycoveragemodel.VerticalComplexSizeInterface, error)
	UpdateComplexSize(ctx context.Context, input UpdateComplexSizeInput) (propertycoveragemodel.VerticalComplexSizeInterface, error)
	DeleteComplexSize(ctx context.Context, sizeID int64) error
	ListComplexSizes(ctx context.Context, input ListComplexSizesInput) ([]propertycoveragemodel.VerticalComplexSizeInterface, error)
	GetComplexSizeDetail(ctx context.Context, sizeID int64) (propertycoveragemodel.VerticalComplexSizeInterface, error)
	CreateComplexZipCode(ctx context.Context, input CreateComplexZipCodeInput) (propertycoveragemodel.HorizontalComplexZipCodeInterface, error)
	UpdateComplexZipCode(ctx context.Context, input UpdateComplexZipCodeInput) (propertycoveragemodel.HorizontalComplexZipCodeInterface, error)
	DeleteComplexZipCode(ctx context.Context, zipCodeID int64) error
	ListComplexZipCodes(ctx context.Context, input ListComplexZipCodesInput) ([]propertycoveragemodel.HorizontalComplexZipCodeInterface, error)
	GetComplexZipCodeDetail(ctx context.Context, zipCodeID int64) (propertycoveragemodel.HorizontalComplexZipCodeInterface, error)
	GetComplexByAddress(ctx context.Context, input GetComplexByAddressInput) (propertycoveragemodel.ManagedComplexInterface, error)
}

type propertyCoverageService struct {
	repository    propertycoveragerepository.RepositoryInterface
	globalService globalservice.GlobalServiceInterface
}

// NewPropertyCoverageService wires the repository and global service dependencies.
func NewPropertyCoverageService(repo propertycoveragerepository.RepositoryInterface, globalSvc globalservice.GlobalServiceInterface) PropertyCoverageServiceInterface {
	return &propertyCoverageService{
		repository:    repo,
		globalService: globalSvc,
	}
}
