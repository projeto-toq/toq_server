package mediaprocessinghandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	httpdto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/projeto-toq/toq_server/internal/adapter/left/http/utils"
	"github.com/projeto-toq/toq_server/internal/core/derrors"
	dto "github.com/projeto-toq/toq_server/internal/core/domain/dto"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ensure swagger error response type is referenced for imports
var _ = httpdto.ErrorResponse{}

// DeleteProjectMedia removes a project media asset (plan/render) for OffPlanHouse listings.
//
// @Summary     Delete project media asset
// @Description Deletes a specific project asset (PROJECT_DOC/PROJECT_RENDER) for a listing in project upload flow.
// @Tags        Listings Media
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body deleteProjectMediaRequest true "Delete payload"
// @Success     204 {string} string "Asset deleted"
// @Failure     400 {object} httpdto.ErrorResponse "Invalid payload"
// @Failure     401 {object} httpdto.ErrorResponse "Unauthorized"
// @Failure     403 {object} httpdto.ErrorResponse "Forbidden"
// @Failure     404 {object} httpdto.ErrorResponse "Asset not found"
// @Failure     500 {object} httpdto.ErrorResponse "Internal server error"
// @Router      /listings/project-media [delete]
func (h *MediaProcessingHandler) DeleteProjectMedia(c *gin.Context) {
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

	var request deleteProjectMediaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	normalized := strings.ToUpper(request.AssetType)
	if normalized != string(mediaprocessingmodel.MediaAssetTypeProjectDoc) && normalized != string(mediaprocessingmodel.MediaAssetTypeProjectRender) {
		httperrors.SendHTTPErrorObj(c, derrors.Validation("unsupported assetType", map[string]any{"assetType": request.AssetType}))
		return
	}

	input := dto.DeleteMediaInput{
		ListingIdentityID: int64(request.ListingIdentityID),
		AssetType:         mediaprocessingmodel.MediaAssetType(normalized),
		Sequence:          request.Sequence,
		RequestedBy:       uint64(userInfo.ID),
	}

	if err := h.service.DeleteMedia(ctx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// deleteProjectMediaRequest is the HTTP payload for deleting a project asset.
type deleteProjectMediaRequest struct {
	ListingIdentityID uint64 `json:"listingIdentityId" binding:"required,min=1" example:"1024"`
	AssetType         string `json:"assetType" binding:"required" example:"PROJECT_DOC"`
	Sequence          uint8  `json:"sequence" binding:"required,min=1" example:"1"`
}
