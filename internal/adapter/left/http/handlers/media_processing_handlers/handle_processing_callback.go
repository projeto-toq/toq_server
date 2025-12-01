package mediaprocessinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	domaindto "github.com/projeto-toq/toq_server/internal/core/domain/dto"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// HandleProcessingCallback handles the callback from the media processing pipeline
func (h *MediaProcessingHandler) HandleProcessingCallback(c *gin.Context) {
	baseCtx := utils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := utils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	var req mediaprocessingmodel.MediaProcessingCallback
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results := make([]domaindto.ProcessingResult, 0, len(req.Outputs))
	for _, output := range req.Outputs {
		status := "PROCESSED"
		errorMsg := ""
		if output.ErrorCode != "" {
			status = "FAILED"
			errorMsg = output.ErrorMessage
		}

		results = append(results, domaindto.ProcessingResult{
			RawKey:       output.RawKey,
			Status:       status,
			ProcessedKey: output.ProcessedKey,
			ThumbnailKey: output.ThumbnailKey,
			Metadata:     output.Outputs,
			Error:        errorMsg,
		})
	}

	input := domaindto.HandleProcessingCallbackInput{
		JobID:   req.JobID,
		Status:  string(req.Status),
		Results: results,
	}

	if _, err := h.service.HandleProcessingCallback(ctx, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
