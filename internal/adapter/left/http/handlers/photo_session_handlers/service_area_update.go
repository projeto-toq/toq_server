package photosessionhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

// UpdateServiceArea updates the city and state of an existing service area.
// @Summary      Update Photographer Service Area
// @Description  Updates the city and state of an existing service area.
// @Tags         Photographer
// @Accept       json
// @Produce      json
// @Param        payload body dto.PhotographerServiceAreaUpdateRequest true "Service area payload"
// @Success      200 {object} dto.PhotographerServiceAreaResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/service-area [put]
func (h *PhotoSessionHandler) UpdateServiceArea(c *gin.Context) {
	var request dto.PhotographerServiceAreaUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_payload", "Invalid request payload: "+err.Error())
		return
	}

	userID, err := h.globalService.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	input := photosessionservices.UpdateServiceAreaInput{
		PhotographerID: uint64(userID),
		ServiceAreaID:  request.ServiceAreaID,
		City:           request.City,
		State:          request.State,
	}

	result, svcErr := h.service.UpdateServiceArea(c.Request.Context(), input)
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
