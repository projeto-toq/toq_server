package schedulehandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostCreateBlockRule handles POST /schedules/listing/block.
//
// @Summary     Create recurring block rules
// @Description Creates recurring blocking rules for a listing agenda.
// @Tags        Listing Schedules
// @Accept      json
// @Produce     json
// @Param       request body dto.ScheduleRuleRequest true "Rule creation payload" Extensions(x-example={"listingId":3241,"weekDays":["MONDAY","TUESDAY"],"rangeStart":"00:00","rangeEnd":"07:59","active":true,"timezone":"America/Sao_Paulo"})
// @Success     200 {object} dto.ScheduleRulesResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     403 {object} dto.ErrorResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     409 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /schedules/listing/block [post]
// @Security    BearerAuth
func (h *ScheduleHandler) PostCreateBlockRule(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "User context not found")
		return
	}

	var req dto.ScheduleRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	weekdays, err := parseScheduleWeekdays(req.WeekDays)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	startMinute, err := parseScheduleRuleMinutes("rangeStart", req.RangeStart)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	endMinute, err := parseScheduleRuleMinutes("rangeEnd", req.RangeEnd)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	input := scheduleservices.CreateRuleInput{
		ListingID: req.ListingID,
		OwnerID:   userInfo.ID,
		Weekdays:  weekdays,
		Range: scheduleservices.RuleTimeRange{
			StartMinute: startMinute,
			EndMinute:   endMinute,
		},
		Active:   req.Active,
		Timezone: req.Timezone,
		ActorID:  userInfo.ID,
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	result, serviceErr := h.scheduleService.CreateRules(ctx, input)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	response := converters.ScheduleRulesMutationToDTO(result)
	c.JSON(http.StatusOK, response)
}
