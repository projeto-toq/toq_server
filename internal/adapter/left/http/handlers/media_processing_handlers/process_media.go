package mediaprocessinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpdto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

var _ = httpdto.ErrorResponse{}

// ProcessMedia triggers the async processing pipeline for pending assets.
// @Summary Trigger media processing
// @Description Starts the asynchronous processing pipeline for assets in PENDING_UPLOAD status.
// @Tags Listings Media
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ProcessMediaInput true "Processing Request"
// @Success 202 {object} map[string]string "Processing started"
// @Failure 400 {object} httpdto.ErrorResponse "Validation Error"
// @Failure 401 {object} httpdto.ErrorResponse "Unauthorized"
// @Failure 404 {object} httpdto.ErrorResponse "Listing Not Found"
// @Failure 500 {object} httpdto.ErrorResponse "Internal Server Error"
// @Router /listings/media/uploads/process [post]
func (h *MediaProcessingHandler) ProcessMedia(c *gin.Context) {
	// 1. Observability & Context Setup
	baseCtx := utils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := utils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// 2. Extract User Info
	userInfo, err := utils.GetUserInfoFromGinContext(c)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User info not found in context")
		return
	}

	var input dto.ProcessMediaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
		return
	}

	input.RequestedBy = uint64(userInfo.ID)

	// 3. Service Call
	if err := h.service.ProcessMedia(ctx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "processing_started"})
}
