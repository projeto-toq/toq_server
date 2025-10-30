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
