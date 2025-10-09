package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

func (uh *UserHandler) AcceptInvitation(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Get user information from Gin context (set by AuthMiddleware)
	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	// Call service to accept invitation
	if err := uh.userService.AcceptInvitation(ctx, userInfo.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.AcceptInvitationResponse{
		Message: "Invitation accepted successfully",
	}

	c.JSON(http.StatusOK, response)
}
