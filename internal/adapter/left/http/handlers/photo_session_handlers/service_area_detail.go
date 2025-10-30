package photosessionhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

// GetServiceArea retrieves a service area owned by the authenticated photographer.
// @Summary      Get Photographer Service Area
// @Description  Retrieves a service area owned by the authenticated photographer.
// @Tags         Photographer
// @Accept       json
// @Produce      json
// @Param        payload body dto.PhotographerServiceAreaDetailRequest true "Service area identifier"
// @Success      200 {object} dto.PhotographerServiceAreaResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/service-area/detail [post]
func (h *PhotoSessionHandler) GetServiceArea(c *gin.Context) {
	var request dto.PhotographerServiceAreaDetailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_payload", "Invalid request payload: "+err.Error())
		return
	}

	userID, err := h.globalService.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	input := photosessionservices.ServiceAreaDetailInput{
		PhotographerID: uint64(userID),
		ServiceAreaID:  request.ServiceAreaID,
	}

	result, svcErr := h.service.GetServiceArea(c.Request.Context(), input)
	if svcErr != nil {
		http_errors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, dto.PhotographerServiceAreaResponse{
		ID:    result.Area.ID(),
		City:  result.Area.City(),
		State: result.Area.State(),
	})
}
