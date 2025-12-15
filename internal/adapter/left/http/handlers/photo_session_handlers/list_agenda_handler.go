package photosessionhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

type _ = dto.ErrorResponse

const (
	defaultAgendaPage = 1
	defaultAgendaSize = 20
	maxAgendaSize     = 100
)

// ListAgenda handles the retrieval of the photographer's agenda.
// @Summary      List Photographer Agenda
// @Description  Retrieves the photographer's agenda, including available and blocked slots, within a given date range. Optional entryType filter to retrieve only specific types of entries. Photo session entries include the photoSessionId field for approval workflows.
// @Tags         Photographer
// @Produce      json
// @Param        startDate query string true "Start date in RFC3339 format"
// @Param        endDate   query string true "End date in RFC3339 format"
// @Param        page      query int    false "Page number (defaults to 1)"
// @Param        size      query int    false "Page size (defaults to 20, max 100)"
// @Param        timezone  query string false "Timezone identifier" default(America/Sao_Paulo)
// @Param        entryType query string false "Filter by entry type" Enums(PHOTO_SESSION, BLOCK, TIME_OFF, HOLIDAY)
// @Param        sortField query string false "Sort field" Enums(startDate, endDate, entryType) default(startDate)
// @Param        sortOrder query string false "Sort direction" Enums(asc, desc) default(asc)
// @Success      200 {object} photosessionservices.ListAgendaOutput
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/agenda [get]
func (h *PhotoSessionHandler) ListAgenda(c *gin.Context) {
	var query dto.ListAgendaQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_query", "Invalid query parameters: "+err.Error())
		return
	}

	startDate, err := coreutils.ParseRFC3339Relaxed("startDate", query.StartDate)
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	endDate, err := coreutils.ParseRFC3339Relaxed("endDate", query.EndDate)
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	loc := coreutils.DetermineRangeLocation(startDate, endDate, nil)
	startDate = coreutils.ConvertToLocation(startDate, loc)
	endDate = coreutils.ConvertToLocation(endDate, loc)

	page := query.Page
	if page <= 0 {
		page = defaultAgendaPage
	}

	size := query.Size
	if size <= 0 {
		size = defaultAgendaSize
	}
	if size > maxAgendaSize {
		size = maxAgendaSize
	}

	userID, dErr := h.globalService.GetUserIDFromContext(c.Request.Context())
	if dErr != nil {
		http_errors.SendHTTPErrorObj(c, dErr)
		return
	}

	var entryTypeFilter *photosessionmodel.AgendaEntryType
	if query.EntryType != nil {
		entryType := photosessionmodel.AgendaEntryType(*query.EntryType)
		entryTypeFilter = &entryType
	}

	input := photosessionservices.ListAgendaInput{
		PhotographerID: uint64(userID),
		StartDate:      startDate,
		EndDate:        endDate,
		Page:           page,
		Size:           size,
		Location:       loc,
		EntryType:      entryTypeFilter,
		SortField:      photosessionservices.AgendaSortField(query.SortField),
		SortOrder:      photosessionservices.AgendaSortOrder(query.SortOrder),
	}

	output, dErr := h.service.ListAgenda(c.Request.Context(), input)
	if dErr != nil {
		http_errors.SendHTTPErrorObj(c, dErr)
		return
	}

	c.JSON(http.StatusOK, output)
}
