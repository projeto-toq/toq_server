package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

func (uh *UserHandler) DeleteRealtorByID(c *gin.Context) {
	_ = coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	// Method not implemented in service layer yet
	httperrors.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Method not implemented yet")
}
