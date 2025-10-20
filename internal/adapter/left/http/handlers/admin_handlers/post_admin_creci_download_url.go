package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostAdminCreciDownloadURL handles POST /admin/users/creci/download-url
//
//	@Summary      Get signed download URLs for CRECI documents
//	@Description  Returns signed URLs (selfie/front/back) for a realtor user, valid for a limited time
//	@Tags         Admin
//	@Accept       json
//	@Produce      json
//	@Param        request  body  dto.AdminCreciDownloadURLRequest  true  "User ID"
//	@Success      200  {object}  dto.AdminCreciDownloadURLResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      404  {object}  map[string]any
//	@Failure      409  {object}  map[string]any
//	@Failure      422  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/users/creci/download-url [post]
func (h *AdminHandler) PostAdminCreciDownloadURL(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminCreciDownloadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	urls, err := h.userService.GetCreciDownloadURLs(ctx, req.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminCreciDownloadURLResponse{
		URLs: dto.AdminCreciDocumentURLs{
			Selfie: urls.Selfie,
			Front:  urls.Front,
			Back:   urls.Back,
		},
		ExpiresInMinutes: urls.ExpiresInMinutes,
	}

	c.JSON(http.StatusOK, resp)
}
