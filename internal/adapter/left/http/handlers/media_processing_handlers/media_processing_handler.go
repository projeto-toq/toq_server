package mediaprocessinghandlers

import (
	"log/slog"

	mediaprocessinghandlerport "github.com/projeto-toq/toq_server/internal/core/port/left/http/mediaprocessinghandler"
	mediaprocessingcallbackport "github.com/projeto-toq/toq_server/internal/core/port/right/functions/mediaprocessingcallback"
	mediaprocessingservice "github.com/projeto-toq/toq_server/internal/core/service/media_processing_service"
)

// MediaProcessingHandler implementa MediaProcessingHandlerPort.
type MediaProcessingHandler struct {
	service           mediaprocessingservice.MediaProcessingServiceInterface
	logger            *slog.Logger
	callbackValidator mediaprocessingcallbackport.CallbackPortInterface
}

// NewMediaProcessingHandler cria uma nova instância do handler.
// @Summary Factory para MediaProcessingHandler
func NewMediaProcessingHandler(
	service mediaprocessingservice.MediaProcessingServiceInterface,
	logger *slog.Logger,
	callbackValidator mediaprocessingcallbackport.CallbackPortInterface,
) mediaprocessinghandlerport.MediaProcessingHandlerPort {
	return &MediaProcessingHandler{
		service:           service,
		logger:            logger,
		callbackValidator: callbackValidator,
	}
}

var _ mediaprocessinghandlerport.MediaProcessingHandlerPort = (*MediaProcessingHandler)(nil)
