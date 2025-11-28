package lambdaport

import (
	"context"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// ZipProcessingServiceInterface defines the contract for ZIP generation service
type ZipProcessingServiceInterface interface {
	GenerateZipBundle(ctx context.Context, input mediaprocessingmodel.GenerateZipInput) (mediaprocessingmodel.ZipOutput, error)
}
