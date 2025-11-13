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

// GetListingBlockRules handles GET /schedules/listing/block.
//
// @Summary     List recurring block rules for a listing agenda
// @Description Returns the recurring blocking rules configured for a listing owned by the authenticated user.
// @Tags        Listing Schedules
// @Produce     json
// @Param       listingIdentityId query int64  true  "Listing identity identifier" example(3241)
// @Param       weekDays  query []string false "Weekdays filter" collectionFormat(multi) example(MONDAY)
// @Success     200 {object} dto.ScheduleRulesResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     403 {object} dto.ErrorResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /schedules/listing/block [get]
// @Security    BearerAuth
func (h *ScheduleHandler) GetListingBlockRules(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "User context not found")
		return
	}

	var req dto.ListingBlockRulesQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_QUERY", "Invalid query parameters")
		return
	}

	weekdays, err := parseScheduleWeekdaysAllowEmpty(req.WeekDays)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	filter := schedulemodel.BlockRulesFilter{
		OwnerID:           userInfo.ID,
		ListingIdentityID: req.ListingIdentityID,
		Weekdays:          weekdays,
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	result, serviceErr := h.scheduleService.ListBlockRules(ctx, filter)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	c.JSON(http.StatusOK, converters.ScheduleRuleListToDTO(result))
}
