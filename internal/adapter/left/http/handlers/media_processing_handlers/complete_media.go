package mediaprocessinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CompleteMedia finalizes the media processing workflow.
// @Summary Complete media processing
// @Description Consolidates media, generates ZIP, and advances listing status.
// @Tags Media Processing
// @Accept json
// @Produce json
// @Param request body dto.CompleteMediaInput true "Completion Request"
// @Success 200 {object} map[string]string "Processing completed"
// @Failure 400 {object} httpdto.ErrorResponse "Validation Error"
// @Failure 404 {object} httpdto.ErrorResponse "Listing Not Found"
// @Failure 500 {object} httpdto.ErrorResponse "Internal Server Error"
// @Router /api/v2/media/complete [post]
func (h *MediaProcessingHandler) CompleteMedia(c *gin.Context) {
	var input dto.CompleteMediaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
		return
	}

	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context")
		return
	}
	input.RequestedBy = userID

	if err := h.service.CompleteMedia(c.Request.Context(), input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "processing_completed"})
}
