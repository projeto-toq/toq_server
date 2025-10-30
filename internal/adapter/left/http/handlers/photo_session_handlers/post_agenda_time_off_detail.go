package photosessionhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

// GetTimeOffDetail handles POST /photographer/agenda/time-off/detail.
//
// @Summary      Get photographer time-off detail
// @Description  Retrieves a specific time-off entry for the authenticated photographer.
// @Tags         Photographer
// @Accept       json
// @Produce      json
// @Param        input body dto.TimeOffDetailRequest true "Time-Off detail payload" Extensions(x-example={"timeOffId":42,"timezone":"America/Sao_Paulo"})
// @Success      200 {object} dto.PhotographerTimeOffResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/agenda/time-off/detail [post]
func (h *PhotoSessionHandler) GetTimeOffDetail(c *gin.Context) {
	var req dto.TimeOffDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_json", "Invalid JSON body")
		return
	}

	userID, err := h.globalService.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	result, serviceErr := h.service.GetTimeOffDetail(c.Request.Context(), photosessionservices.TimeOffDetailInput{
		TimeOffID:      req.TimeOffID,
		PhotographerID: uint64(userID),
		Timezone:       req.Timezone,
	})
	if serviceErr != nil {
		http_errors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	c.JSON(http.StatusOK, converters.TimeOffResultToDTO(result))
}
