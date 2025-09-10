package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
)

// PostPhotoUploadURL generates a pre-signed URL to upload a profile photo
//
// @Summary      Create pre-signed upload URL for profile photo
// @Description  Create a pre-signed URL to upload a profile photo to storage
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request  body      dto.GetPhotoUploadURLRequest  true  "Upload request"
// @Success      200      {object}  dto.GetPhotoUploadURLResponse
// @Failure      400      {object}  dto.ErrorResponse  "Invalid request"
// @Failure      401      {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      422      {object}  dto.ErrorResponse  "Validation error (photo type or content type)"
// @Failure      500      {object}  dto.ErrorResponse  "Internal server error"
// @Router       /user/photo/upload-url [post]
// @Security     BearerAuth
func (uh *UserHandler) PostPhotoUploadURL(c *gin.Context) {
	ctx := c.Request.Context()

	var request dto.GetPhotoUploadURLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	if request.ObjectName == "" || request.ContentType == "" {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "MISSING_FIELDS", "object_name and content_type are required")
		return
	}

	signedURL, err := uh.userService.GetPhotoUploadURL(ctx, request.ObjectName, request.ContentType)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.GetPhotoUploadURLResponse{SignedURL: signedURL})
}
