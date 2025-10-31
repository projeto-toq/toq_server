package photosessionhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
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

	rangeFrom, err := coreutils.ParseRFC3339Relaxed("rangeFrom", query.RangeFrom)
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	rangeTo, err := coreutils.ParseRFC3339Relaxed("rangeTo", query.RangeTo)
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	loc := coreutils.DetermineRangeLocation(rangeFrom, rangeTo, nil)
	rangeFrom = coreutils.ConvertToLocation(rangeFrom, loc)
	rangeTo = coreutils.ConvertToLocation(rangeTo, loc)

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
		Location:       loc,
	}

	result, serviceErr := h.service.ListTimeOff(c.Request.Context(), input)
	if serviceErr != nil {
		http_errors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	c.JSON(http.StatusOK, converters.ListTimeOffOutputToDTO(result))
}
