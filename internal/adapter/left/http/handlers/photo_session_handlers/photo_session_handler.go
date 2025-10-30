package photosessionhandlers

import (
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

// PhotoSessionHandler handles HTTP requests for photographer agenda management.
type PhotoSessionHandler struct {
	service       photosessionservices.PhotoSessionServiceInterface
	globalService globalservice.GlobalServiceInterface
}

// NewPhotoSessionHandler creates a new handler with its dependencies.
func NewPhotoSessionHandler(service photosessionservices.PhotoSessionServiceInterface, globalService globalservice.GlobalServiceInterface) *PhotoSessionHandler {
	return &PhotoSessionHandler{
		service:       service,
		globalService: globalService,
	}
}
