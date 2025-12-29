package mediaprocessinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/projeto-toq/toq_server/internal/adapter/left/http/utils"
	domaindto "github.com/projeto-toq/toq_server/internal/core/domain/dto"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ApproveListingMedia allows the listing owner to approve or reject processed assets.
//
// @Summary     Approve or reject listing media
// @Description Owner-only endpoint that validates the listing status (PENDING_OWNER_APPROVAL) and applies the requested decision.
// @Tags        Listings Media
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body dto.ListingMediaApprovalRequest true "Owner decision payload"
// @Success     200 {object} dto.ListingMediaApprovalResponse "Decision applied"
// @Failure     400 {object} dto.ErrorResponse "Listing not awaiting owner approval"
// @Failure     401 {object} dto.ErrorResponse "Unauthorized"
// @Failure     403 {object} dto.ErrorResponse "Only the owner can approve/reject"
// @Failure     404 {object} dto.ErrorResponse "Listing not found"
// @Failure     500 {object} dto.ErrorResponse "Internal server error"
// @Router      /listings/media/approve [post]
func (h *MediaProcessingHandler) ApproveListingMedia(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	userInfo, err := coreutils.GetUserInfoFromGinContext(c)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User info not found in context")
		return
	}

	var request dto.ListingMediaApprovalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	input := domaindto.ListingMediaApprovalInput{
		ListingIdentityID: int64(request.ListingIdentityID),
		Approve:           request.Approve,
		RequestedBy:       uint64(userInfo.ID),
	}

	output, err := h.service.HandleOwnerMediaApproval(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	decision := "rejected"
	if request.Approve {
		decision = "approved"
	}

	c.JSON(http.StatusOK, dto.ListingMediaApprovalResponse{
		ListingIdentityID: uint64(output.ListingIdentityID),
		Decision:          decision,
		NewStatus:         output.NewStatus.String(),
	})
}
