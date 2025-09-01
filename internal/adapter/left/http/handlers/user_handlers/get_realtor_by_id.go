package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
)

func (uh *UserHandler) GetRealtorByID(c *gin.Context) {
	// Method not implemented in service layer yet
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Method not implemented yet")
}
