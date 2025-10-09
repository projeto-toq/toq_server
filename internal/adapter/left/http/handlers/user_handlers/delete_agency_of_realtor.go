package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

func (uh *UserHandler) DeleteAgencyOfRealtor(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Get user information from context (set by middleware)
	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		// Se chegar aqui, é erro de pipeline (middleware deveria ter setado)
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	// Call service to delete agency of realtor
	if err := uh.userService.DeleteAgencyOfRealtor(ctx, userInfo.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.DeleteAgencyOfRealtorResponse{
		Message: "Agency of realtor deleted successfully",
	}

	c.JSON(http.StatusOK, response)
}
