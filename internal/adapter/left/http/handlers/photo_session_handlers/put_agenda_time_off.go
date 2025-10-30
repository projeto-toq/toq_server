package photosessionhandlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

// UpdateTimeOff handles PUT /photographer/agenda/time-off.
//
// @Summary      Update photographer time-off
// @Description  Updates an existing time-off period for the authenticated photographer and refreshes agenda slots.
// @Tags         Photographer
// @Accept       json
// @Produce      json
// @Param        input body dto.UpdateTimeOffRequest true "Update Time-Off payload" Extensions(x-example={"timeOffId":42,"startDate":"2025-07-05T10:00:00-03:00","endDate":"2025-07-05T12:00:00-03:00","reason":"Consulta médica","timezone":"America/Sao_Paulo"})
// @Success      200 {object} dto.PhotographerTimeOffResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/agenda/time-off [put]
func (h *PhotoSessionHandler) UpdateTimeOff(c *gin.Context) {
	var req dto.UpdateTimeOffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_json", "Invalid JSON body")
		return
	}

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_start_date", "Invalid startDate format, use RFC3339")
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_end_date", "Invalid endDate format, use RFC3339")
		return
	}

	userID, err := h.globalService.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	result, serviceErr := h.service.UpdateTimeOff(c.Request.Context(), photosessionservices.UpdateTimeOffInput{
		TimeOffID:      req.TimeOffID,
		PhotographerID: uint64(userID),
		StartDate:      startDate,
		EndDate:        endDate,
		Reason:         req.Reason,
		Timezone:       req.Timezone,
	})
	if serviceErr != nil {
		http_errors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	c.JSON(http.StatusOK, converters.TimeOffResultToDTO(result))
}
