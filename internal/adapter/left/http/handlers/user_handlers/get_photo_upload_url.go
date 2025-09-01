package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
)

func (uh *UserHandler) GetPhotoUploadURL(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request body
	var request dto.GetPhotoUploadURLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Validate required fields
	if request.ObjectName == "" || request.ContentType == "" {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "MISSING_FIELDS", "object_name and content_type are required")
		return
	}

	// Call service to get photo upload URL
	signedURL, err := uh.userService.GetPhotoUploadURL(ctx, request.ObjectName, request.ContentType)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.GetPhotoUploadURLResponse{
		SignedURL: signedURL,
	}

	c.JSON(http.StatusOK, response)
}
