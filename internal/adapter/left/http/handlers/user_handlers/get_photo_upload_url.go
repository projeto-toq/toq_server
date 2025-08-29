package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetPhotoUploadURL(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request body
	var request dto.GetPhotoUploadURLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Validate required fields
	if request.ObjectName == "" || request.ContentType == "" {
		utils.SendHTTPError(c, http.StatusBadRequest, "MISSING_FIELDS", "object_name and content_type are required")
		return
	}

	// Call service to get photo upload URL
	signedURL, err := uh.userService.GetPhotoUploadURL(ctx, request.ObjectName, request.ContentType)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "GET_PHOTO_UPLOAD_URL_FAILED", "Failed to get photo upload URL")
		return
	}

	// Prepare response
	response := dto.GetPhotoUploadURLResponse{
		SignedURL: signedURL,
	}

	c.JSON(http.StatusOK, response)
}
