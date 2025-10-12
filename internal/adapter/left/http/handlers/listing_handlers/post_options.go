package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostOptions handles getting available options for listings
//
//	@Summary		Get listing options
//	@Description	Get available property types for listings based on location
//	@Tags			Listings
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.GetOptionsRequest	true	"Location data for listing options"
//	@Success		200		{object}	dto.GetOptionsResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"Forbidden"
//	@Failure		404		{object}	dto.ErrorResponse	"Complex not found"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/listings/options [post]
//	@Security		BearerAuth
func (lh *ListingHandler) PostOptions(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	var request dto.GetOptionsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	types, err := lh.listingService.GetOptions(ctx, request.ZipCode, request.Number)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	propertyTypes := make([]dto.PropertyTypeOption, 0, len(types))
	for _, t := range types {
		propertyTypes = append(propertyTypes, dto.PropertyTypeOption{
			PropertyType: int(t.Code),
			Name:         t.Label,
		})
	}

	c.JSON(http.StatusOK, dto.GetOptionsResponse{
		PropertyTypes: propertyTypes,
	})
}
