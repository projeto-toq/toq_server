package mediaprocessinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/listing_handlers/converters"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/projeto-toq/toq_server/internal/adapter/left/http/utils"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// RequestProjectUploadURLs issues signed URLs for project media uploads (plans/renders).
//
// @Summary     Request project media upload URLs
// @Description Generates signed PUT URLs for OffPlanHouse project assets (PROJECT_DOC/PROJECT_RENDER) while the listing is in StatusPendingPlanLoading. Enforces payload validation and returns upload instructions with required headers.
// @Tags        Listings Media
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body dto.RequestProjectUploadURLsRequest true "Project media manifest"
// @Success     201 {object} dto.RequestProjectUploadURLsResponse "Signed URLs issued"
// @Failure     400 {object} dto.ErrorResponse "Invalid request payload"
// @Failure     401 {object} dto.ErrorResponse "Unauthorized"
// @Failure     403 {object} dto.ErrorResponse "Forbidden"
// @Failure     404 {object} dto.ErrorResponse "Listing not found or invalid status"
// @Failure     409 {object} dto.ErrorResponse "Open batch exists"
// @Failure     500 {object} dto.ErrorResponse "Internal server error"
// @Router      /listings/project-media/uploads [post]
func (h *MediaProcessingHandler) RequestProjectUploadURLs(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	var request dto.RequestProjectUploadURLsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	input, convErr := converters.DTOToRequestProjectUploadURLsInput(request)
	if convErr != nil {
		httperrors.SendHTTPErrorObj(c, convErr)
		return
	}

	if userID, ok := coreutils.GetUserIDFromContext(ctx); ok {
		input.RequestedBy = userID
	}

	output, err := h.service.RequestProjectUploadURLs(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	files := make([]dto.RequestUploadInstructionResponse, 0, len(output.Files))
	for _, f := range output.Files {
		files = append(files, dto.RequestUploadInstructionResponse{
			AssetType: f.AssetType,
			UploadURL: f.UploadURL,
			Method:    f.Method,
			Headers:   f.Headers,
			ObjectKey: f.ObjectKey,
			Sequence:  f.Sequence,
			Title:     f.Title,
		})
	}

	response := dto.RequestProjectUploadURLsResponse{
		ListingIdentityID:   uint64(output.ListingIdentityID),
		UploadURLTTLSeconds: output.UploadURLTTLSeconds,
		Files:               files,
	}

	c.JSON(http.StatusCreated, response)
}
