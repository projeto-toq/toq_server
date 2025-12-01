package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpdto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
)

var _ = httpdto.ErrorResponse{}

// CompleteMedia finalizes the media processing workflow.
// @Summary Complete media processing
// @Description Consolidates media, generates ZIP, and advances listing status.
// @Tags Listings Media
// @Accept json
// @Produce json
// @Param request body dto.CompleteMediaInput true "Completion Request"
// @Success 200 {object} map[string]string "Processing completed"
// @Failure 400 {object} httpdto.ErrorResponse "Validation Error"
// @Failure 404 {object} httpdto.ErrorResponse "Listing Not Found"
// @Failure 500 {object} httpdto.ErrorResponse "Internal Server Error"
// @Router /api/v2/listings/media/uploads/complete [post]
func (lh *ListingHandler) CompleteMedia(c *gin.Context) {
	var input dto.CompleteMediaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
		return
	}

	// TODO: Implement CompleteMedia logic in service layer first
	// This handler is a placeholder for the route definition
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "CompleteMedia service method not implemented yet")
}
