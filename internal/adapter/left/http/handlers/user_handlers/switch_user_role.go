package userhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// alias para garantir que o swagger reconheça o tipo de erro padrão
type _ = dto.ErrorResponse

// SwitchUserRole processa POST /user/role/switch permitindo owners e realtors alternarem a role ativa.
//
//	@Summary	Altera o role ativo do usuário autenticado (owner ↔ realtor)
//	@Description	Permite que usuários com roles owner e realtor alternem entre as funções atribuídas.
//	@Tags	User
//	@Accept	json
//	@Produce	json
//	@Param	request	body	dto.SwitchUserRoleRequest	true	"Slug do role desejado"
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

	// Parse request body
	var request dto.SwitchUserRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	request.RoleSlug = strings.ToLower(request.RoleSlug)
	switch permissionmodel.RoleSlug(request.RoleSlug) {
	case permissionmodel.RoleSlugOwner, permissionmodel.RoleSlugRealtor:
		// válido
	default:
		httperrors.SendHTTPErrorObj(c, coreutils.BadRequest("Only owners or realtors can switch roles"))
		return
	}

	// Call service to switch user role
	tokens, err := uh.userService.SwitchUserRole(ctx, permissionmodel.RoleSlug(request.RoleSlug))
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
