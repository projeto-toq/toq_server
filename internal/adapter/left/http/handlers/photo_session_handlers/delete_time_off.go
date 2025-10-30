package photosessionhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

// DeleteTimeOff handles the deletion of a time-off entry.
// @Summary      Delete Photographer Time-Off
// @Description  Removes an existing time-off period for the authenticated photographer, making slots available again.
// @Tags         Photographer
// @Accept       json
// @Produce      json
// @Param        input body dto.DeleteTimeOffRequest true "Delete Time-Off payload" Extensions(x-example={"timeOffId":42,"timezone":"America/Sao_Paulo"})
// @Success      200 {object} object{message=string}
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/agenda/time-off [delete]
func (h *PhotoSessionHandler) DeleteTimeOff(c *gin.Context) {
	var req dto.DeleteTimeOffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_json", "Invalid JSON body")
		return
	}

	userID, err := h.globalService.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	input := photosessionservices.DeleteTimeOffInput{
		TimeOffID:      req.TimeOffID,
		PhotographerID: uint64(userID),
		Timezone:       req.Timezone,
	}

	dErr := h.service.DeleteTimeOff(c.Request.Context(), input)
	if dErr != nil {
		http_errors.SendHTTPErrorObj(c, dErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Time-off deleted successfully"})
}
