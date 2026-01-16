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

// CompleteProjectMedia finalizes project media uploads (plans/renders) for OffPlanHouse listings.
//
// @Summary     Complete project media upload
// @Description Confirms project media uploads for a listing in StatusPendingPlanLoading and triggers finalization.
// @Tags        Listings Media
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body dto.CompleteProjectMediaRequest true "Completion payload"
// @Success     200 {object} map[string]string "Project media completed"
// @Failure     400 {object} dto.ErrorResponse "Invalid payload"
// @Failure     401 {object} dto.ErrorResponse "Unauthorized"
// @Failure     403 {object} dto.ErrorResponse "Forbidden"
// @Failure     404 {object} dto.ErrorResponse "Listing not found or invalid status"
// @Failure     500 {object} dto.ErrorResponse "Internal server error"
// @Router      /listings/project-media/complete [post]
func (h *MediaProcessingHandler) CompleteProjectMedia(c *gin.Context) {
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

	var request dto.CompleteProjectMediaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	input := converters.DTOToCompleteProjectMediaInput(request)
	input.RequestedBy = uint64(userInfo.ID)

	if err := h.service.CompleteProjectMedia(ctx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "project_media_completed"})
}
