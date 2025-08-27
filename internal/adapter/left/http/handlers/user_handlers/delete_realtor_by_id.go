package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) DeleteRealtorByID(c *gin.Context) {
	// Method not implemented in service layer yet
	utils.SendHTTPError(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Method not implemented yet")
}
