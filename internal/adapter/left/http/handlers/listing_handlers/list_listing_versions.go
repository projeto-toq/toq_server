package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListListingVersions retorna as versões de um listing agrupadas por identidade.
//
// @Summary    List listing versions
// @Description    Returns all versions attached to a listing identity, indicating which one is active.
// @Tags       Listings
// @Produce    json
// @Param      listingIdentityId query int  true  "Listing identity identifier" example(1024)
// @Param      includeDeleted    query bool false "Include soft-deleted versions" default(false)
// @Success    200 {object} dto.ListListingVersionsResponse
// @Failure    400 {object} dto.ErrorResponse "Validation error"
// @Failure    401 {object} dto.ErrorResponse "Unauthorized"
// @Failure    403 {object} dto.ErrorResponse "Forbidden"
// @Failure    404 {object} dto.ErrorResponse "Listing not found"
// @Failure    500 {object} dto.ErrorResponse "Internal error"
// @Router     /listings/versions [get]
// @Security   BearerAuth
func (lh *ListingHandler) ListListingVersions(c *gin.Context) {
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

	var request dto.ListListingVersionsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	input := listingservices.ListListingVersionsInput{
		ListingIdentityID: request.ListingIdentityID,
		IncludeDeleted:    request.IncludeDeleted,
	}

	output, serviceErr := lh.listingService.ListListingVersions(ctx, input)
	if serviceErr != nil {
		httperrors.SendHTTPErrorObj(c, serviceErr)
		return
	}

	response := converters.ListingVersionsToDTO(output)
	c.JSON(http.StatusOK, response)
}
