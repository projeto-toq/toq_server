package listinghandlers

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
// @Param request body dto.ProcessMediaInput true "Processing Request"
// @Success 202 {object} map[string]string "Processing started"
// @Failure 400 {object} httpdto.ErrorResponse "Validation Error"
// @Failure 404 {object} httpdto.ErrorResponse "Listing Not Found"
// @Failure 500 {object} httpdto.ErrorResponse "Internal Server Error"
// @Router /api/v2/listings/media/uploads/process [post]
func (lh *ListingHandler) ProcessMedia(c *gin.Context) {
	var input dto.ProcessMediaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
		return
	}

	// Extract user context
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context")
		return
	}
	input.RequestedBy = userID

	if err := lh.mediaProcessingService.ProcessMedia(c.Request.Context(), input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "processing_started"})
}
