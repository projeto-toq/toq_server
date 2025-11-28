package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	mediaprocessingservice "github.com/projeto-toq/toq_server/internal/core/service/media_processing_service"
)

// HandleProcessingCallback handles the callback from the media processing pipeline
func (h *ListingHandler) HandleProcessingCallback(c *gin.Context) {
	var req mediaprocessingmodel.MediaProcessingCallback
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := mediaprocessingservice.HandleProcessingCallbackInput{
		Callback: req,
	}

	_, err := h.mediaProcessingService.HandleProcessingCallback(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
