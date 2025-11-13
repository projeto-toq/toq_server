package schedulehandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteBlockRule handles DELETE /schedules/listing/block.
//
// @Summary     Delete a recurring block rule
// @Description Removes a recurring blocking rule from a listing agenda.
// @Tags        Listing Schedules
// @Accept      json
// @Produce     json
// @Param       request body dto.ScheduleRuleDeleteRequest true "Rule deletion payload" Extensions(x-example={"ruleId":5021,"listingIdentityId":3241})
// @Success     204 "Rule deleted successfully"
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     403 {object} dto.ErrorResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /schedules/listing/block [delete]
// @Security    BearerAuth
func (h *ScheduleHandler) DeleteBlockRule(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "User context not found")
		return
	}

	var req dto.ScheduleRuleDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	input := scheduleservices.DeleteRuleInput{
		RuleID:            req.RuleID,
		ListingIdentityID: req.ListingIdentityID,
		OwnerID:           userInfo.ID,
		ActorID:           userInfo.ID,
	}

	ctx := coreutils.ContextWithLogger(baseCtx)
	if err := h.scheduleService.DeleteRule(ctx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
