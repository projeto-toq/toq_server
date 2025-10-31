package photosessionhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateTimeOff handles the creation of a new time-off entry for a photographer.
// @Summary      Create Photographer Time-Off
// @Description  Registers a new time-off period for the authenticated photographer, blocking their agenda.
// @Tags         Photographer
// @Accept       json
// @Produce      json
// @Param        input body dto.CreateTimeOffRequest true "Time-Off payload" Extensions(x-example={"startDate":"2025-07-05T09:00:00-03:00","endDate":"2025-07-05T18:00:00-03:00","reason":"Participação em evento"})
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

	startDate, err := coreutils.ParseRFC3339Relaxed("startDate", req.StartDate)
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	endDate, err := coreutils.ParseRFC3339Relaxed("endDate", req.EndDate)
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	loc := coreutils.DetermineRangeLocation(startDate, endDate, nil)
	startDate = coreutils.ConvertToLocation(startDate, loc)
	endDate = coreutils.ConvertToLocation(endDate, loc)

	userID, err := h.globalService.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	input := photosessionservices.TimeOffInput{
		PhotographerID: uint64(userID),
		StartDate:      startDate,
		EndDate:        endDate,
		Reason:         req.Reason,
		Location:       loc,
	}

	id, dErr := h.service.CreateTimeOff(c.Request.Context(), input)
	if dErr != nil {
		http_errors.SendHTTPErrorObj(c, dErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Time-off created successfully", "timeOffId": id})
}
