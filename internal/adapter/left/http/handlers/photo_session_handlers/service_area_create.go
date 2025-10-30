package photosessionhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

// CreateServiceArea handles the creation of a new service area for the authenticated photographer.
// @Summary      Create Photographer Service Area
// @Description  Creates a new service area for the authenticated photographer.
// @Tags         Photographer
// @Accept       json
// @Produce      json
// @Param        payload body dto.PhotographerServiceAreaRequest true "Service area payload"
// @Success      201 {object} dto.PhotographerServiceAreaResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/service-area [post]
func (h *PhotoSessionHandler) CreateServiceArea(c *gin.Context) {
	var request dto.PhotographerServiceAreaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_payload", "Invalid request payload: "+err.Error())
		return
	}

	userID, err := h.globalService.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	input := photosessionservices.CreateServiceAreaInput{
		PhotographerID: uint64(userID),
		City:           request.City,
		State:          request.State,
	}

	result, svcErr := h.service.CreateServiceArea(c.Request.Context(), input)
	if svcErr != nil {
		http_errors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.JSON(http.StatusCreated, dto.PhotographerServiceAreaResponse{
		ID:    result.Area.ID(),
		City:  result.Area.City(),
		State: result.Area.State(),
	})
}
