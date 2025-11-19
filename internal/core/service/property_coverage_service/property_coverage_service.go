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
