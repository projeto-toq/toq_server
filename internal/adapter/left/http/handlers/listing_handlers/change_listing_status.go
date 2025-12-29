package listinghandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/projeto-toq/toq_server/internal/adapter/left/http/utils"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ChangeListingStatus publishes or suspends listings owned by the requester.
//
// @Summary     Update listing publication status
// @Description Allows the listing owner to publish (READY → PUBLISHED) or suspend (PUBLISHED/UNDER_OFFER/UNDER_NEGOTIATION → READY) the active version.
// @Tags        Listings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body dto.ChangeListingStatusRequest true "Listing status change payload"
// @Success     200 {object} dto.ChangeListingStatusResponse "Transition applied"
// @Failure     400 {object} dto.ErrorResponse "Validation error"
// @Failure     401 {object} dto.ErrorResponse "Unauthorized"
// @Failure     403 {object} dto.ErrorResponse "Forbidden"
// @Failure     404 {object} dto.ErrorResponse "Listing not found"
// @Failure     409 {object} dto.ErrorResponse "Conflict"
// @Failure     500 {object} dto.ErrorResponse "Internal server error"
// @Router      /listings/status [post]
func (lh *ListingHandler) ChangeListingStatus(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	userInfo, err := coreutils.GetUserInfoFromGinContext(c)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	var request dto.ChangeListingStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	action := listingservices.ListingStatusAction(strings.ToUpper(request.Action))
	input := listingservices.ChangeListingStatusInput{
		ListingIdentityID: request.ListingIdentityID,
		Action:            action,
		RequesterUserID:   int64(userInfo.ID),
	}

	output, err := lh.listingService.ChangeListingStatus(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ChangeListingStatusResponse{
		ListingIdentityID: output.ListingIdentityID,
		ActiveVersionID:   output.ActiveVersionID,
		PreviousStatus:    output.PreviousStatus.String(),
		NewStatus:         output.NewStatus.String(),
	})
}
