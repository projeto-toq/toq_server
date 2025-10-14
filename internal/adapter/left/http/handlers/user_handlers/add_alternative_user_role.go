package userhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

// alias para garantir que o swagger reconheça o tipo de erro padrão
type _ = dto.ErrorResponse

// AddAlternativeUserRole handles POST /user/role/alternative and creates an alternative role (owner ↔ realtor).
//
//	@Summary	Create an alternative role for the authenticated user (owner ↔ realtor)
//	@Description	Allows owners to request a realtor role (pending CRECI validation) and vice versa. CRECI fields are optional, but when provided: creciNumber must end with "-F", and creciState must be a valid Brazilian state abbreviation (2 letters).
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AddAlternativeUserRoleRequest	true	"Payload containing CRECI data"
//	@Success	200	{object}	dto.AddAlternativeUserRoleResponse
//	@Failure	400	{object}	dto.ErrorResponse
//	@Failure	401	{object}	dto.ErrorResponse
//	@Failure	403	{object}	dto.ErrorResponse
//	@Failure	409	{object}	dto.ErrorResponse
//	@Failure	500	{object}	dto.ErrorResponse
//	@Router		/user/role/alternative [post]
//	@Security	BearerAuth
func (uh *UserHandler) AddAlternativeUserRole(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Parse request body
	var request dto.AddAlternativeUserRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	userInfo, err := coreutils.GetUserInfoFromContext(ctx)
	if err != nil || userInfo.ID == 0 {
		httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError(""))
		return
	}

	// Get user to determine current role e status
	user, err := uh.userService.GetProfile(ctx)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	activeRole := user.GetActiveRole()
	if activeRole == nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ErrUserActiveRoleMissing)
		return
	}

	currentRole := coreutils.GetUserRoleSlugFromUserRole(activeRole)
	if currentRole != permissionmodel.RoleSlugOwner && currentRole != permissionmodel.RoleSlugRealtor {
		httperrors.SendHTTPErrorObj(c, coreutils.AuthorizationError("Only owners or realtors can request an alternative role"))
		return
	}

	if activeRole.GetStatus() != permissionmodel.StatusActive {
		httperrors.SendHTTPErrorObj(c, coreutils.ConflictError("Active role status must be active"))
		return
	}

	var (
		alternativeRole permissionmodel.RoleSlug
		creciArgs       []string
	)

	switch currentRole {
	case permissionmodel.RoleSlugOwner:
		alternativeRole = permissionmodel.RoleSlugRealtor

		creciNumber := request.CreciNumber
		creciState := request.CreciState
		creciValidity := strings.TrimSpace(request.CreciValidity)
		creciProvided := strings.TrimSpace(creciNumber) != "" || strings.TrimSpace(creciState) != "" || creciValidity != ""

		if creciProvided {
			normalizedNumber, err := validators.ValidateCreciNumber("request.creciNumber", creciNumber, true)
			if err != nil {
				httperrors.SendHTTPErrorObj(c, err)
				return
			}

			normalizedState, err := validators.ValidateCreciState("request.creciState", creciState, true)
			if err != nil {
				httperrors.SendHTTPErrorObj(c, err)
				return
			}

			if creciValidity == "" {
				httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("request.creciValidity", "Creci validity is required when sending CRECI data"))
				return
			}

			creciArgs = []string{normalizedNumber, normalizedState, creciValidity}
		}
	case permissionmodel.RoleSlugRealtor:
		alternativeRole = permissionmodel.RoleSlugOwner
	default:
		// fallback defensivo, embora já tenhamos retornado acima
		httperrors.SendHTTPErrorObj(c, coreutils.AuthorizationError("Unsupported role"))
		return
	}

	// Call service to add alternative role
	if err := uh.userService.AddAlternativeRole(ctx, userInfo.ID, alternativeRole, creciArgs...); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.AddAlternativeUserRoleResponse{
		Message: "Alternative user role added successfully",
	}

	c.JSON(http.StatusOK, response)
}
