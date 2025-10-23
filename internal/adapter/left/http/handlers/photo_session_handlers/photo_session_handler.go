package photosessionhandlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

type _ = dto.ErrorResponse

// PhotoSessionHandler handles HTTP requests for photographer agenda management.
type PhotoSessionHandler struct {
	service       photosessionservices.PhotoSessionServiceInterface
	globalService globalservice.GlobalServiceInterface
}

// NewPhotoSessionHandler creates a new handler with its dependencies.
func NewPhotoSessionHandler(service photosessionservices.PhotoSessionServiceInterface, globalService globalservice.GlobalServiceInterface) *PhotoSessionHandler {
	return &PhotoSessionHandler{
		service:       service,
		globalService: globalService,
	}
}

// CreateTimeOff handles the creation of a new time-off entry for a photographer.
// @Summary      Create Photographer Time-Off
// @Description  Registers a new time-off period for the authenticated photographer, blocking their agenda.
// @Tags         Photo Session
// @Accept       json
// @Produce      json
// @Param        input body dto.CreateTimeOffRequest true "Time-Off payload"
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

// DeleteTimeOff handles the deletion of a time-off entry.
// @Summary      Delete Photographer Time-Off
// @Description  Removes an existing time-off period for the authenticated photographer, making slots available again.
// @Tags         Photo Session
// @Accept       json
// @Produce      json
// @Param        input body dto.DeleteTimeOffRequest true "Delete Time-Off payload"
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
		TimeOffID:         req.TimeOffID,
		PhotographerID:    uint64(userID),
		Timezone:          req.Timezone,
		HolidayCalendarID: req.HolidayCalendarID,
		HorizonMonths:     req.HorizonMonths,
		WorkdayStartHour:  req.WorkdayStartHour,
		WorkdayEndHour:    req.WorkdayEndHour,
	}

	dErr := h.service.DeleteTimeOff(c.Request.Context(), input)
	if dErr != nil {
		http_errors.SendHTTPErrorObj(c, dErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Time-off deleted successfully"})
}
