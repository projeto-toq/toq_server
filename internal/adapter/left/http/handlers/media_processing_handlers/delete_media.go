package mediaprocessinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteMedia removes a media asset.
// @Summary Delete media asset
// @Description Removes a specific media asset from the listing.
// @Tags Media Processing
// @Accept json
// @Produce json
// @Param request body dto.DeleteMediaInput true "Delete Request"
// @Success 200 {object} map[string]string "Asset deleted"
// @Failure 400 {object} httpdto.ErrorResponse "Validation Error"
// @Failure 404 {object} httpdto.ErrorResponse "Asset Not Found"
// @Failure 500 {object} httpdto.ErrorResponse "Internal Server Error"
// @Router /api/v2/media/delete [delete]
func (h *MediaProcessingHandler) DeleteMedia(c *gin.Context) {
	var input dto.DeleteMediaInput
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

	if err := h.service.DeleteMedia(c.Request.Context(), input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "asset_deleted"})
}
