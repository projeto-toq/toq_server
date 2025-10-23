package photosessionhandlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

type _ = dto.ErrorResponse

// ListAgendaRequest defines the expected JSON body for the ListAgenda endpoint.
type ListAgendaRequest struct {
	StartDate string `json:"startDate" binding:"required" example:"2023-10-01T00:00:00Z"`
	EndDate   string `json:"endDate" binding:"required" example:"2023-10-31T23:59:59Z"`
}

// ListAgenda handles the retrieval of the photographer's agenda.
// @Summary      List Photographer Agenda
// @Description  Retrieves the photographer's agenda, including available and blocked slots, within a given date range.
// @Tags         Photo Session
// @Accept       json
// @Produce      json
// @Param        input body ListAgendaRequest true "Date range for the agenda"
// @Success      200 {object} photosessionservices.ListAgendaOutput
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/agenda [post]
func (h *PhotoSessionHandler) ListAgenda(c *gin.Context) {
	var req ListAgendaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_json", "Invalid JSON body: "+err.Error())
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

	userID, dErr := h.globalService.GetUserIDFromContext(c.Request.Context())
	if dErr != nil {
		http_errors.SendHTTPErrorObj(c, dErr)
		return
	}

	input := photosessionservices.ListAgendaInput{
		PhotographerID: uint64(userID),
		StartDate:      startDate,
		EndDate:        endDate,
	}

	output, dErr := h.service.ListAgenda(c.Request.Context(), input)
	if dErr != nil {
		http_errors.SendHTTPErrorObj(c, dErr)
		return
	}

	c.JSON(http.StatusOK, output)
}
