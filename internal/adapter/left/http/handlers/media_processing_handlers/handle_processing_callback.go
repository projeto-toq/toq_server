package mediaprocessinghandlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	httpdto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// HandleProcessingCallback receives the asynchronous callback fired by the Step Functions pipeline and
// forwards it to the service layer. The endpoint stays unauthenticated but protected by the shared secret.
//
// @Summary     Receive async media processing callback
// @Description Validates the shared secret, parses the Step Functions callback payload, updates internal state and returns 200.
// @Tags        Listings Media
// @Accept      json
// @Produce     json
// @Param       X-Toq-Signature  header  string                                false  "HMAC signature generated in the callback lambda"
// @Param       request          body    httpdto.MediaProcessingCallbackRequest true   "Callback payload forwarded by Step Functions"
// @Success     200              {object} httpdto.APIResponse                  "Callback acknowledged"
// @Failure     400              {object} httpdto.ErrorResponse                "Invalid callback payload"
// @Failure     401              {object} httpdto.ErrorResponse                "Invalid shared secret"
// @Failure     500              {object} httpdto.ErrorResponse                "Internal error"
// @Router      /listings/media/callback [post]
func (h *MediaProcessingHandler) HandleProcessingCallback(c *gin.Context) {
	baseCtx := utils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := utils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	request, err := httpdto.BindMediaProcessingCallbackRequest(c.Request)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_CALLBACK", err.Error())
		return
	}

	if err := h.validateSharedSecret(ctx, c.GetHeader("X-Toq-Signature"), request.RawBody); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
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

	c.JSON(http.StatusOK, httpdto.SuccessResponse(gin.H{"status": "accepted"}))
}

func (h *MediaProcessingHandler) validateSharedSecret(ctx context.Context, signature string, payload []byte) error {
	if h.callbackValidator == nil {
		if h.logger != nil {
			h.logger.Warn("handler.media.callback.validator_missing")
		}
		return nil
	}
	return h.callbackValidator.ValidateSignature(ctx, signature, payload)
}
