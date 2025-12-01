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

// UpdateMedia updates metadata or properties of an existing media asset.
// @Summary Update media asset
// @Description Updates the title or metadata of a specific media asset.
// @Tags Listings Media
// @Accept json
// @Produce json
// @Param request body dto.UpdateMediaInput true "Update Request"
// @Success 200 {object} map[string]string "Asset updated"
// @Failure 400 {object} httpdto.ErrorResponse "Validation Error"
// @Failure 404 {object} httpdto.ErrorResponse "Asset Not Found"
// @Failure 500 {object} httpdto.ErrorResponse "Internal Server Error"
// @Router /api/v2/listings/media/update [post]
func (h *MediaProcessingHandler) UpdateMedia(c *gin.Context) {
	var input dto.UpdateMediaInput
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

	if err := h.service.UpdateMedia(c.Request.Context(), input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "asset_updated"})
}
