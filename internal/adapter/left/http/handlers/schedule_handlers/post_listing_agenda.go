package schedulehandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostListingAgenda handles POST /schedules/listing/detail.
//
// @Summary		List agenda entries for a listing
// @Description	Returns all agenda entries for a specific listing owned by the authenticated user.
// @Tags		Schedules
// @Accept		json
// @Produce	json
// @Param		request	body	dto.ListingAgendaDetailRequest	true	"Agenda detail filter"
// @Success	200	{object}	dto.ListingAgendaDetailResponse
// @Failure	400	{object}	dto.ErrorResponse
// @Failure	401	{object}	dto.ErrorResponse
// @Failure	403	{object}	dto.ErrorResponse
// @Failure	404	{object}	dto.ErrorResponse
// @Failure	500	{object}	dto.ErrorResponse
// @Router		/schedules/listing/detail [post]
// @Security	BearerAuth
func (h *ScheduleHandler) PostListingAgenda(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "User context not found")
		return
	}

	var req dto.ListingAgendaDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	rangeFilter, err := parseScheduleRange(req.Range)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	pagination := sanitizeSchedulePagination(req.Pagination)

	filter := schedulemodel.AgendaDetailFilter{
		OwnerID:    userInfo.ID,
		ListingID:  req.ListingID,
		Range:      rangeFilter,
		Pagination: pagination,
	}

	ctx = coreutils.ContextWithLogger(ctx)
	result, serviceErr := h.scheduleService.ListAgendaEntries(ctx, filter)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	page, limit := schedulePaginationValues(pagination)
	response := converters.ScheduleEntriesToDTO(result.Entries, page, limit, result.Total)

	c.JSON(http.StatusOK, response)
}
