package photosessionhandlers

import (
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

const (
	defaultServiceAreaPage = 1
	defaultServiceAreaSize = 20
	maxServiceAreaSize     = 100
)

// ListServiceAreas handles the retrieval of service areas for the authenticated photographer.
// @Summary      List Photographer Service Areas
// @Description  Lists the service areas configured by the authenticated photographer with optional filters.
// @Tags         Photographer
// @Produce      json
// @Param        city  query string false "City filter"
// @Param        state query string false "State filter"
// @Param        page  query int    false "Page number"
// @Param        size  query int    false "Page size"
// @Success      200 {object} dto.PhotographerServiceAreaListResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/service-area [get]
func (h *PhotoSessionHandler) ListServiceAreas(c *gin.Context) {
	var query dto.PhotographerServiceAreaListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_query", "Invalid query parameters: "+err.Error())
		return
	}

	page := query.Page
	if page <= 0 {
		page = defaultServiceAreaPage
	}

	size := query.Size
	if size <= 0 {
		size = defaultServiceAreaSize
	}
	if size > maxServiceAreaSize {
		size = maxServiceAreaSize
	}

	userID, err := h.globalService.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	var cityFilter *string
	if trimmed := strings.TrimSpace(query.City); trimmed != "" {
		cityFilter = &trimmed
	}

	var stateFilter *string
	if trimmed := strings.TrimSpace(query.State); trimmed != "" {
		stateFilter = &trimmed
	}

	input := photosessionservices.ListServiceAreasInput{
		PhotographerID: uint64(userID),
		City:           cityFilter,
		State:          stateFilter,
		Page:           page,
		Size:           size,
	}

	output, svcErr := h.service.ListServiceAreas(c.Request.Context(), input)
	if svcErr != nil {
		http_errors.SendHTTPErrorObj(c, svcErr)
		return
	}

	responses := make([]dto.PhotographerServiceAreaResponse, 0, len(output.Areas))
	for _, area := range output.Areas {
		responses = append(responses, dto.PhotographerServiceAreaResponse{
			ID:    area.ID(),
			City:  area.City(),
			State: area.State(),
		})
	}

	totalPages := 0
	if output.Size > 0 {
		totalPages = int(math.Ceil(float64(output.Total) / float64(output.Size)))
	}

	c.JSON(http.StatusOK, dto.PhotographerServiceAreaListResponse{
		ServiceAreas: responses,
		Pagination: dto.PaginationResponse{
			Page:       output.Page,
			Limit:      output.Size,
			Total:      output.Total,
			TotalPages: totalPages,
		},
	})
}
