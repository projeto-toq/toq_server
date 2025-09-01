package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
)

func (uh *UserHandler) GetDocumentsUploadURL(c *gin.Context) {
	// Parse request body
	var request dto.GetDocumentsUploadURLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// TODO: This method is not implemented in the service layer yet
	// For now, return a placeholder response
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "GetDocumentsUploadURL service method not implemented yet")
}
