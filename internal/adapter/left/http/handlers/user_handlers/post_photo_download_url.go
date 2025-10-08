package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// PostPhotoDownloadURL generates a pre-signed URL to download a profile photo variant
//
// @Summary      Create pre-signed download URL for profile photo
// @Description  Create a pre-signed URL to download a specific profile photo variant from storage
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request  body      dto.GetPhotoDownloadURLRequest  true  "Download request"
// @Success      200      {object}  dto.GetPhotoDownloadURLResponse
// @Failure      400      {object}  dto.ErrorResponse  "Invalid request"
// @Failure      401      {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      422      {object}  dto.ErrorResponse  "Validation error (variant)"
// @Failure      500      {object}  dto.ErrorResponse  "Internal server error"
// @Router       /user/photo/download-url [post]
// @Security     BearerAuth
func (uh *UserHandler) PostPhotoDownloadURL(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var request dto.GetPhotoDownloadURLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	if request.Variant == "" {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "MISSING_FIELDS", "variant is required")
		return
	}

	signedURL, err := uh.userService.GetPhotoDownloadURL(ctx, request.Variant)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.GetPhotoDownloadURLResponse{SignedURL: signedURL})
}
