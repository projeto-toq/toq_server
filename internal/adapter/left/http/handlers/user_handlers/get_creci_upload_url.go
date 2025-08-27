package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetCreciUploadURL(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request body
	var request dto.GetCreciUploadURLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Validate required fields
	if request.DocumentType == "" || request.ContentType == "" {
		utils.SendHTTPError(c, http.StatusBadRequest, "MISSING_FIELDS", "document_type and content_type are required")
		return
	}

	// Call service to get CRECI upload URL
	signedURL, err := uh.userService.GetCreciUploadURL(ctx, request.DocumentType, request.ContentType)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "GET_CRECI_UPLOAD_URL_FAILED", "Failed to get CRECI upload URL")
		return
	}

	// Prepare response
	response := dto.GetCreciUploadURLResponse{
		SignedURL: signedURL,
	}

	c.JSON(http.StatusOK, response)
}
