package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// alias para garantir que o swagger reconheça o tipo de erro padrão
type _ = dto.ErrorResponse

// SwitchUserRole handles POST /user/role/switch, allowing owners and realtors to toggle the active role.
//
//	@Summary	Switch the authenticated user's active role (owner ↔ realtor)
//	@Description	Allows users who have both owner and realtor roles to switch between them.
//	@Tags	User
//	@Accept	json
//	@Produce	json
//	@Success	200	{object}	dto.SwitchUserRoleResponse
//	@Failure	400	{object}	dto.ErrorResponse
//	@Failure	401	{object}	dto.ErrorResponse
//	@Failure	403	{object}	dto.ErrorResponse
//	@Failure	409	{object}	dto.ErrorResponse
//	@Failure	500	{object}	dto.ErrorResponse
//	@Router	/user/role/switch [post]
//	@Security	BearerAuth
func (uh *UserHandler) SwitchUserRole(c *gin.Context) {
	// Enrich context with request info and user
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Call service to switch user role
	tokens, err := uh.userService.SwitchUserRole(ctx)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.SwitchUserRoleResponse{
		Tokens: dto.TokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}

	c.JSON(http.StatusOK, response)
}
