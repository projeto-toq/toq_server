package photosessionhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

// DeleteServiceArea removes a service area owned by the authenticated photographer.
// @Summary      Delete Photographer Service Area
// @Description  Deletes a service area owned by the authenticated photographer.
// @Tags         Photographer
// @Accept       json
// @Produce      json
// @Param        payload body dto.PhotographerServiceAreaDeleteRequest true "Service area identifier"
// @Success      204 "No Content"
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/service-area [delete]
func (h *PhotoSessionHandler) DeleteServiceArea(c *gin.Context) {
	var request dto.PhotographerServiceAreaDeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_payload", "Invalid request payload: "+err.Error())
		return
	}

	userID, err := h.globalService.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	input := photosessionservices.DeleteServiceAreaInput{
		PhotographerID: uint64(userID),
		ServiceAreaID:  request.ServiceAreaID,
	}

	if svcErr := h.service.DeleteServiceArea(c.Request.Context(), input); svcErr != nil {
		http_errors.SendHTTPErrorObj(c, svcErr)
		return
	}

	c.Status(http.StatusNoContent)
}
