package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PromoteListingVersion promove uma versão draft para ativa após passar pelas validações necessárias.
//
// @Summary    Promote listing version
// @Description    Promotes a draft listing version to become the active version for its identity, preserving historical records.
// @Tags       Listings
// @Accept     json
// @Produce    json
// @Param      request body dto.PromoteListingVersionRequest true "Listing version identifier" Extensions(x-example={"versionId":12345})
// @Success    200 {object} dto.PromoteListingVersionResponse
// @Failure    400 {object} dto.ErrorResponse "Validation error"
// @Failure    401 {object} dto.ErrorResponse "Unauthorized"
// @Failure    403 {object} dto.ErrorResponse "Forbidden"
// @Failure    404 {object} dto.ErrorResponse "Listing version not found"
// @Failure    409 {object} dto.ErrorResponse "Conflict"
// @Failure    500 {object} dto.ErrorResponse "Internal error"
// @Router     /listings/versions/promote [post]
// @Security   BearerAuth
func (lh *ListingHandler) PromoteListingVersion(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	var request dto.PromoteListingVersionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	input := listingservices.PromoteListingVersionInput{VersionID: request.VersionID}
	if err := lh.listingService.PromoteListingVersion(ctx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.PromoteListingVersionResponse{
		Success: true,
		Message: "Listing version promoted",
	})
}
