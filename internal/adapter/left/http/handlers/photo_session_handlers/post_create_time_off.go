package photosessionhandlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

// CreateTimeOff handles the creation of a new time-off entry for a photographer.
// @Summary      Create Photographer Time-Off
// @Description  Registers a new time-off period for the authenticated photographer, blocking their agenda.
// @Tags         Photographer
// @Accept       json
// @Produce      json
// @Param        input body dto.CreateTimeOffRequest true "Time-Off payload" Extensions(x-example={"startDate":"2025-07-05T09:00:00-03:00","endDate":"2025-07-05T18:00:00-03:00","reason":"Participação em evento","timezone":"America/Sao_Paulo","holidayCalendarId":1,"horizonMonths":3,"workdayStartHour":8,"workdayEndHour":19})
// @Success      201 {object} object{message=string,timeOffId=int}
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/agenda/time-off [post]
func (h *PhotoSessionHandler) CreateTimeOff(c *gin.Context) {
	var req dto.CreateTimeOffRequest
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

	input := photosessionservices.TimeOffInput{
		PhotographerID:    uint64(userID),
		StartDate:         startDate,
		EndDate:           endDate,
		Reason:            req.Reason,
		Timezone:          req.Timezone,
		HolidayCalendarID: req.HolidayCalendarID,
		HorizonMonths:     req.HorizonMonths,
		WorkdayStartHour:  req.WorkdayStartHour,
		WorkdayEndHour:    req.WorkdayEndHour,
	}

	id, dErr := h.service.CreateTimeOff(c.Request.Context(), input)
	if dErr != nil {
		http_errors.SendHTTPErrorObj(c, dErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Time-off created successfully", "timeOffId": id})
}
