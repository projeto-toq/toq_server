package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
)

// PostCreciUploadURL issues a pre-signed URL to upload CRECI documents (realtor-only)
//
// @Summary      Create pre-signed upload URL for CRECI documents
// @Description  Create a pre-signed URL to upload a CRECI document (selfie/front/back)
// @Tags         Realtor
// @Accept       json
// @Produce      json
// @Param        request  body      dto.GetCreciUploadURLRequest  true  "Upload request"
// @Success      200      {object}  dto.GetCreciUploadURLResponse
// @Failure      400      {object}  dto.ErrorResponse  "Invalid request"
// @Failure      401      {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      403      {object}  dto.ErrorResponse  "Forbidden"
// @Failure      422      {object}  dto.ErrorResponse  "Validation error (documentType or contentType)"
// @Failure      500      {object}  dto.ErrorResponse  "Internal server error"
// @Router       /realtor/creci/upload-url [post]
// @Security     BearerAuth
func (uh *UserHandler) PostCreciUploadURL(c *gin.Context) {
	ctx := c.Request.Context()

	var request dto.GetCreciUploadURLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	if request.DocumentType == "" || request.ContentType == "" {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "MISSING_FIELDS", "documentType and contentType are required")
		return
	}

	signedURL, err := uh.userService.GetCreciUploadURL(ctx, request.DocumentType, request.ContentType)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.GetCreciUploadURLResponse{SignedURL: signedURL})
}
