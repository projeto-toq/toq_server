package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) RejectInvitation(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Get user information from Gin context (set by AuthMiddleware)
	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	// Call service to reject invitation
	if err := uh.userService.RejectInvitation(ctx, userInfo.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.RejectInvitationResponse{
		Message: "Invitation rejected successfully",
	}

	c.JSON(http.StatusOK, response)
}
