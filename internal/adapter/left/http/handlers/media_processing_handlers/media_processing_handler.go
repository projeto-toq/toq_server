package mediaprocessinghandlers

import (
	"log/slog"

	mediaprocessinghandlerport "github.com/projeto-toq/toq_server/internal/core/port/left/http/mediaprocessinghandler"
	mediaprocessingservice "github.com/projeto-toq/toq_server/internal/core/service/media_processing_service"
)

// MediaProcessingHandler implementa MediaProcessingHandlerPort.
type MediaProcessingHandler struct {
	service mediaprocessingservice.MediaProcessingServiceInterface
	logger  *slog.Logger
}

// NewMediaProcessingHandler cria uma nova instância do handler.
// @Summary Factory para MediaProcessingHandler
func NewMediaProcessingHandler(
	service mediaprocessingservice.MediaProcessingServiceInterface,
	logger *slog.Logger,
) mediaprocessinghandlerport.MediaProcessingHandlerPort {
	return &MediaProcessingHandler{
		service: service,
		logger:  logger,
	}
}

var _ mediaprocessinghandlerport.MediaProcessingHandlerPort = (*MediaProcessingHandler)(nil)
