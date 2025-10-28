package photosessionhandlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
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

// DeleteTimeOff handles the deletion of a time-off entry.
// @Summary      Delete Photographer Time-Off
// @Description  Removes an existing time-off period for the authenticated photographer, making slots available again.
// @Tags         Photographer
// @Accept       json
// @Produce      json
// @Param        input body dto.DeleteTimeOffRequest true "Delete Time-Off payload" Extensions(x-example={"timeOffId":42,"timezone":"America/Sao_Paulo","holidayCalendarId":1,"horizonMonths":3,"workdayStartHour":8,"workdayEndHour":19})
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

// ListTimeOff handles GET /photographer/agenda/time-off.
//
// @Summary      List photographer time-off
// @Description  Lists time-off periods for the authenticated photographer within a date range.
// @Tags         Photographer
// @Produce      json
// @Param        rangeFrom query string true "Range start (RFC3339)" example("2025-07-01T00:00:00Z")
// @Param        rangeTo query string true "Range end (RFC3339)" example("2025-07-31T23:59:59Z")
// @Param        page query int false "Page number" minimum(1) example(1)
// @Param        size query int false "Items per page" minimum(1) example(20)
// @Param        timezone query string false "Preferred timezone" example("America/Sao_Paulo")
// @Success      200 {object} dto.ListPhotographerTimeOffResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/agenda/time-off [get]
func (h *PhotoSessionHandler) ListTimeOff(c *gin.Context) {
	var query dto.ListTimeOffQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_query", "Invalid query parameters")
		return
	}

	rangeFrom, err := time.Parse(time.RFC3339, query.RangeFrom)
	if err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_range_from", "Invalid rangeFrom format, use RFC3339")
		return
	}

	rangeTo, err := time.Parse(time.RFC3339, query.RangeTo)
	if err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_range_to", "Invalid rangeTo format, use RFC3339")
		return
	}

	userID, err := h.globalService.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	input := photosessionservices.ListTimeOffInput{
		PhotographerID: uint64(userID),
		RangeFrom:      rangeFrom,
		RangeTo:        rangeTo,
		Page:           query.Page,
		Size:           query.Size,
		Timezone:       query.Timezone,
	}

	result, serviceErr := h.service.ListTimeOff(c.Request.Context(), input)
	if serviceErr != nil {
		http_errors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	c.JSON(http.StatusOK, converters.ListTimeOffOutputToDTO(result))
}

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

// UpdateTimeOff handles PUT /photographer/agenda/time-off.
//
// @Summary      Update photographer time-off
// @Description  Updates an existing time-off period for the authenticated photographer and refreshes agenda slots.
// @Tags         Photographer
// @Accept       json
// @Produce      json
// @Param        input body dto.UpdateTimeOffRequest true "Update Time-Off payload" Extensions(x-example={"timeOffId":42,"startDate":"2025-07-05T10:00:00-03:00","endDate":"2025-07-05T12:00:00-03:00","reason":"Consulta médica","timezone":"America/Sao_Paulo","holidayCalendarId":1,"horizonMonths":3,"workdayStartHour":8,"workdayEndHour":19})
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
		TimeOffID:         req.TimeOffID,
		PhotographerID:    uint64(userID),
		StartDate:         startDate,
		EndDate:           endDate,
		Reason:            req.Reason,
		Timezone:          req.Timezone,
		HolidayCalendarID: req.HolidayCalendarID,
		HorizonMonths:     req.HorizonMonths,
		WorkdayStartHour:  req.WorkdayStartHour,
		WorkdayEndHour:    req.WorkdayEndHour,
	})
	if serviceErr != nil {
		http_errors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	c.JSON(http.StatusOK, converters.TimeOffResultToDTO(result))
}
