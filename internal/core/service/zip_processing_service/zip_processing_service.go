package zipprocessingservice

import (
	lambdaport "github.com/projeto-toq/toq_server/internal/core/port/left/lambda_port"
	storageport "github.com/projeto-toq/toq_server/internal/core/port/right/storage"
)

type zipProcessingService struct {
	s3Adapter storageport.ListingMediaStoragePort
}

// NewZipProcessingService creates a new instance of ZipProcessingService
func NewZipProcessingService(s3Adapter storageport.ListingMediaStoragePort) lambdaport.ZipProcessingServiceInterface {
	return &zipProcessingService{
		s3Adapter: s3Adapter,
	}
}
