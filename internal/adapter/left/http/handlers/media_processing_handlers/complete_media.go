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

// CompleteMedia finalizes the upload process
// @Summary Complete media upload
// @Description Marks the upload batch as complete and triggers finalization steps.
// @Tags Listings Media
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CompleteMediaInput true "Completion Request"
// @Success 200 {object} map[string]string "Upload completed"
// @Failure 400 {object} httpdto.ErrorResponse "Validation Error"
// @Failure 401 {object} httpdto.ErrorResponse "Unauthorized"
// @Failure 500 {object} httpdto.ErrorResponse "Internal Server Error"
// @Router /listings/media/uploads/complete [post]
func (h *MediaProcessingHandler) CompleteMedia(c *gin.Context) {
	baseCtx := utils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := utils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	userInfo, err := utils.GetUserInfoFromGinContext(c)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User info not found in context")
		return
	}

	var input dto.CompleteMediaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
		return
	}

	input.RequestedBy = uint64(userInfo.ID)

	if err := h.service.CompleteMedia(ctx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "upload_completed"})
}
