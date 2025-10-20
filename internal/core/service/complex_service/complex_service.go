package complexservices

import (
	"context"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	complexrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/complex_repository"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
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
	CreateComplex(ctx context.Context, input CreateComplexInput) (complexmodel.ComplexInterface, error)
	UpdateComplex(ctx context.Context, input UpdateComplexInput) (complexmodel.ComplexInterface, error)
	DeleteComplex(ctx context.Context, id int64) error
	ListComplexes(ctx context.Context, filter ListComplexesInput) ([]complexmodel.ComplexInterface, error)
	GetComplexDetail(ctx context.Context, id int64) (complexmodel.ComplexInterface, error)
	CreateComplexTower(ctx context.Context, input CreateComplexTowerInput) (complexmodel.ComplexTowerInterface, error)
	UpdateComplexTower(ctx context.Context, input UpdateComplexTowerInput) (complexmodel.ComplexTowerInterface, error)
	DeleteComplexTower(ctx context.Context, id int64) error
	ListComplexTowers(ctx context.Context, filter ListComplexTowersInput) ([]complexmodel.ComplexTowerInterface, error)
	CreateComplexSize(ctx context.Context, input CreateComplexSizeInput) (complexmodel.ComplexSizeInterface, error)
	UpdateComplexSize(ctx context.Context, input UpdateComplexSizeInput) (complexmodel.ComplexSizeInterface, error)
	DeleteComplexSize(ctx context.Context, id int64) error
	ListComplexSizes(ctx context.Context, filter ListComplexSizesInput) ([]complexmodel.ComplexSizeInterface, error)
	CreateComplexZipCode(ctx context.Context, input CreateComplexZipCodeInput) (complexmodel.ComplexZipCodeInterface, error)
	UpdateComplexZipCode(ctx context.Context, input UpdateComplexZipCodeInput) (complexmodel.ComplexZipCodeInterface, error)
	DeleteComplexZipCode(ctx context.Context, id int64) error
	ListComplexZipCodes(ctx context.Context, filter ListComplexZipCodesInput) ([]complexmodel.ComplexZipCodeInterface, error)
	ListSizesByAddress(ctx context.Context, input ListSizesByAddressInput) ([]complexmodel.ComplexSizeInterface, error)
}
