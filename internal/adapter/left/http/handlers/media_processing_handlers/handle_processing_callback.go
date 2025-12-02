package mediaprocessinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpdto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
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

	// Validate shared secret when configured.
	if h.callbackValidator != nil {
		signature := c.GetHeader("X-Toq-Signature")
		if err := h.callbackValidator.ValidateSharedSecret(ctx, signature); err != nil {
			httperrors.SendHTTPErrorObj(c, err)
			return
		}
	} else if h.logger != nil {
		h.logger.Warn("handler.media.callback.validator_missing")
	}

	request, err := httpdto.BindMediaProcessingCallbackRequest(c.Request)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_CALLBACK", err.Error())
		return
	}

	input, err := toHandleProcessingCallbackInput(request)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_CALLBACK", err.Error())
		return
	}

	if h.logger != nil {
		h.logger.Info("handler.media.callback.forward", "job_id", input.JobID, "status", input.Status, "provider", input.Provider)
	}

	if _, err := h.service.HandleProcessingCallback(ctx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
