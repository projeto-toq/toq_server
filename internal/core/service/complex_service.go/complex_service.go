package complexservices

import (
	"context"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	complexrepository "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/complex_repository"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
)

type complexService struct {
	complexRepository complexrepository.ComplexRepoPortInterface
	gsi               globalservice.GlobalServiceInterface
}

func NewComplexService(
	cr complexrepository.ComplexRepoPortInterface,
	gsi globalservice.GlobalServiceInterface,

) ComplexServiceInterface {
	return &complexService{
		complexRepository: cr,
		gsi:               gsi,
	}
}

type ComplexServiceInterface interface {
	GetOptions(ctx context.Context, zipCode string, number string) (propertyTypes globalmodel.PropertyType, err error)
}
