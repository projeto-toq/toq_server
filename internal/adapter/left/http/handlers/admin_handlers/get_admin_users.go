package adminhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetAdminUsers handles GET /admin/users
//
//	@Summary      List admin-facing users with filters
//	@Tags         Admin
//	@Produce      json
//	@Param        page          query  int    false  "Page number" default(1)
//	@Param        limit         query  int    false  "Page size" default(20)
//	@Param        roleName      query  string false  "Filter by role name"
//	@Param        roleSlug      query  string false  "Filter by role slug"
//	@Param        roleStatus    query  int    false  "Filter by role status enum"
//	@Param        isSystemRole  query  bool   false  "Filter by system role flag"
//	@Param        fullName      query  string false  "Filter by full name"
//	@Param        cpf           query  string false  "Filter by CPF"
//	@Param        email         query  string false  "Filter by email"
//	@Param        phoneNumber   query  string false  "Filter by phone"
//	@Param        deleted       query  bool   false  "Filter by deletion flag"
//	@Success      200  {object}  dto.AdminListUsersResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/users [get]
func (h *AdminHandler) GetAdminUsers(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	var roleStatus *permissionmodel.UserRoleStatus
	if req.RoleStatus != nil {
		status := permissionmodel.UserRoleStatus(*req.RoleStatus)
		roleStatus = &status
	}

	input := userservices.ListUsersInput{
		Page:         req.Page,
		Limit:        req.Limit,
		RoleName:     strings.TrimSpace(req.RoleName),
		RoleSlug:     strings.TrimSpace(req.RoleSlug),
		RoleStatus:   roleStatus,
		IsSystemRole: req.IsSystemRole,
		FullName:     strings.TrimSpace(req.FullName),
		CPF:          strings.TrimSpace(req.CPF),
		Email:        strings.TrimSpace(req.Email),
		PhoneNumber:  strings.TrimSpace(req.PhoneNumber),
		Deleted:      req.Deleted,
	}

	result, err := h.userService.ListUsers(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminListUsersResponse{
		Users: make([]dto.AdminUserSummary, 0, len(result.Users)),
		Pagination: dto.PaginationResponse{
			Page:       result.Page,
			Limit:      result.Limit,
			Total:      result.Total,
			TotalPages: computeTotalPages(result.Total, result.Limit),
		},
	}

	for _, user := range result.Users {
		resp.Users = append(resp.Users, toAdminUserSummary(user))
	}

	c.JSON(http.StatusOK, resp)
}
