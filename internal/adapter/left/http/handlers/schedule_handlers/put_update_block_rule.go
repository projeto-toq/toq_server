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

// PutUpdateBlockRule handles PUT /schedules/listing/block.
//
// @Summary     Update a recurring block rule
// @Description Updates an existing recurring blocking rule for a listing agenda.
// @Tags        Listing Schedules
// @Accept      json
// @Produce     json
// @Param       request body dto.ScheduleRuleUpdateRequest true "Rule update payload" Extensions(x-example={"ruleId":5021,"listingIdentityId":3241,"weekDays":["MONDAY"],"rangeStart":"10:00","rangeEnd":"18:00","active":true,"timezone":"America/Sao_Paulo"})
// @Success     200 {object} dto.ScheduleRulesResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     403 {object} dto.ErrorResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     409 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /schedules/listing/block [put]
// @Security    BearerAuth
func (h *ScheduleHandler) PutUpdateBlockRule(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "User context not found")
		return
	}

	var req dto.ScheduleRuleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	weekday, err := parseSingleScheduleWeekday(req.WeekDays)
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

	input := scheduleservices.UpdateRuleInput{
		RuleID:            req.RuleID,
		ListingIdentityID: req.ListingIdentityID,
		OwnerID:           userInfo.ID,
		Weekday:           weekday,
		Range: scheduleservices.RuleTimeRange{
			StartMinute: startMinute,
			EndMinute:   endMinute,
		},
		Active:   req.Active,
		Timezone: req.Timezone,
		ActorID:  userInfo.ID,
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	rule, serviceErr := h.scheduleService.UpdateRule(ctx, input)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	response := dto.ScheduleRulesResponse{
		ListingIdentityID: req.ListingIdentityID,
		Timezone:          req.Timezone,
		Rules:             []dto.ScheduleRuleResponse{converters.ScheduleRuleToDTO(rule)},
	}
	c.JSON(http.StatusOK, response)
}
