package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
)

// alias to ensure swagger can link error response
type _ = dto.ErrorResponse

// GetUserStatus handles GET /user/status
//
//	@Summary      Get active user role status
//	@Description  Returns the status of the currently active user role. Absence of an active role is treated as internal inconsistency.
//	@Tags         User
//	@Accept       json
//	@Produce      json
//	@Success      200 {object} dto.UserStatusResponse "Active role status (status campo inteiro enum ver comentários DTO)"
//	@Failure      400 {object} dto.ErrorResponse "Bad request"
//	@Failure      401 {object} dto.ErrorResponse "Unauthorized"
//	@Failure      500 {object} dto.ErrorResponse "Internal error"
//	@Router       /user/status [get]
//	@Security     BearerAuth
func (uh *UserHandler) GetUserStatus(c *gin.Context) {
	// Handlers não devem criar spans; middleware de telemetria já cria.
	ctx := c.Request.Context()

	status, serr := uh.userService.GetActiveRoleStatus(ctx)
	if serr != nil {
		httperrors.SendHTTPErrorObj(c, serr)
		return
	}

	c.JSON(http.StatusOK, dto.UserStatusResponse{Data: dto.UserStatusData{Status: int(status)}})
}
