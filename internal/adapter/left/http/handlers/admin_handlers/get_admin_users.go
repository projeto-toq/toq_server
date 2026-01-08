package adminhandlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetAdminUsers handles GET /admin/users
//
//	@Summary      List admin-facing users with filters
//	@Tags         Admin Users
//	@Produce      json
//	@Param        page              query  int    false  "Page number" default(1) Extensions(x-example=1)
//	@Param        limit             query  int    false  "Page size" default(20) Extensions(x-example=20)
//	@Param        roleName          query  string false  "Filter by role name (supports '*' wildcard)" Extensions(x-example="*manager*")
//	@Param        roleSlug          query  string false  "Filter by role slug (supports '*' wildcard)" Extensions(x-example="*admin*")
//	@Param        roleStatus        query  int    false  "Filter by role status enum" Extensions(x-example=1)
//	@Param        isSystemRole      query  bool   false  "Filter by system role flag" Extensions(x-example=true)
//	@Param        fullName          query  string false  "Filter by full name (supports '*' wildcard)" Extensions(x-example="*Silva*")
//	@Param        cpf               query  string false  "Filter by CPF (supports '*' wildcard)" Extensions(x-example="123*456")
//	@Param        email             query  string false  "Filter by email (supports '*' wildcard)" Extensions(x-example="*toq.com*")
//	@Param        phoneNumber       query  string false  "Filter by phone (supports '*' wildcard)" Extensions(x-example="*119*")
//	@Param        deleted           query  bool   false  "Filter by deletion flag" Extensions(x-example=false)
//	@Param        idFrom            query  int    false  "Filter by user ID (from)" Extensions(x-example=100)
//	@Param        idTo              query  int    false  "Filter by user ID (to)" Extensions(x-example=250)
//	@Param        bornAtFrom        query  string false  "Filter by birth date from (RFC3339 or YYYY-MM-DD)" Extensions(x-example="1990-01-01")
//	@Param        bornAtTo          query  string false  "Filter by birth date to (RFC3339 or YYYY-MM-DD)" Extensions(x-example="2000-12-31")
//	@Param        lastActivityFrom  query  string false  "Filter by last activity from (RFC3339 or YYYY-MM-DD)" Extensions(x-example="2025-01-01T00:00:00Z")
//	@Param        lastActivityTo    query  string false  "Filter by last activity to (RFC3339 or YYYY-MM-DD)" Extensions(x-example="2025-01-31T23:59:59Z")
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

	var roleStatus *globalmodel.UserRoleStatus
	if req.RoleStatus != nil {
		status := globalmodel.UserRoleStatus(*req.RoleStatus)
		roleStatus = &status
	}

	bornAtFrom, err := parseOptionalISOTime(req.BornAtFrom)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("bornAtFrom", err.Error()))
		return
	}
	bornAtTo, err := parseOptionalISOTime(req.BornAtTo)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("bornAtTo", err.Error()))
		return
	}
	lastActivityFrom, err := parseOptionalISOTime(req.LastActivityFrom)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("lastActivityFrom", err.Error()))
		return
	}
	lastActivityTo, err := parseOptionalISOTime(req.LastActivityTo)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("lastActivityTo", err.Error()))
		return
	}

	input := userservices.ListUsersInput{
		Page:             req.Page,
		Limit:            req.Limit,
		RoleName:         strings.TrimSpace(req.RoleName),
		RoleSlug:         strings.TrimSpace(req.RoleSlug),
		RoleStatus:       roleStatus,
		IsSystemRole:     req.IsSystemRole,
		FullName:         strings.TrimSpace(req.FullName),
		CPF:              strings.TrimSpace(req.CPF),
		Email:            strings.TrimSpace(req.Email),
		PhoneNumber:      strings.TrimSpace(req.PhoneNumber),
		Deleted:          req.Deleted,
		IDFrom:           req.IDFrom,
		IDTo:             req.IDTo,
		BornAtFrom:       bornAtFrom,
		BornAtTo:         bornAtTo,
		LastActivityFrom: lastActivityFrom,
		LastActivityTo:   lastActivityTo,
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

func parseOptionalISOTime(value string) (*time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}
	layouts := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, trimmed); err == nil {
			result := parsed
			return &result, nil
		}
	}
	return nil, fmt.Errorf("invalid datetime format, expected RFC3339 timestamp or YYYY-MM-DD")
}
