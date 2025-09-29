package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/converters"
	dto "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
)

type AdminHandler struct {
	userService userservices.UserServiceInterface
}

func NewAdminHandlerAdapter(userService userservices.UserServiceInterface) *AdminHandler {
	return &AdminHandler{userService: userService}
}

// GetPendingRealtors handles GET /admin/user/pending
//
//	@Summary      List realtors pending manual validation
//	@Description  Returns id, nickname, fullName, nationalID, creciNumber, creciValidity, creciState
//	@Tags         Admin
//	@Produce      json
//	@Success      200  {object}  dto.AdminGetPendingRealtorsResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/user/pending [get]
func (h *AdminHandler) GetPendingRealtors(c *gin.Context) {
	ctx := c.Request.Context()
	// Use service to list by StatusPendingManual
	realtors, err := h.userService.GetCrecisToValidateByStatus(ctx, permissionmodel.StatusPendingManual)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}
	resp := dto.AdminGetPendingRealtorsResponse{Realtors: make([]dto.AdminPendingRealtor, 0, len(realtors))}
	for _, r := range realtors {
		resp.Realtors = append(resp.Realtors, dto.AdminPendingRealtor{
			ID:            r.GetID(),
			NickName:      r.GetNickName(),
			FullName:      r.GetFullName(),
			NationalID:    r.GetNationalID(),
			CreciNumber:   r.GetCreciNumber(),
			CreciValidity: httpconv.FormatDate(r.GetCreciValidity()),
			CreciState:    r.GetCreciState(),
		})
	}
	c.JSON(http.StatusOK, resp)
}

// PostAdminGetUser handles POST /admin/user
//
//	@Summary      Get full user by ID
//	@Tags         Admin
//	@Accept       json
//	@Produce      json
//	@Param        request  body  dto.AdminGetUserRequest  true  "User ID"
//	@Success      200  {object}  dto.AdminGetUserResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      404  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/user [post]
func (h *AdminHandler) PostAdminGetUser(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.AdminGetUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	user, err := h.userService.GetUserByID(ctx, req.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}
	// Reuse existing converter for public user DTO
	profile := httpconv.ToGetProfileResponse(user)
	c.JSON(http.StatusOK, dto.AdminGetUserResponse{User: profile.User})
}

// PostAdminApproveUser handles POST /admin/user/approve
//
//	@Summary      Approve or refuse realtor status manually
//	@Description  Status must be one of the allowed enum values. On success, sends FCM notification.
//	@Tags         Admin
//	@Accept       json
//	@Produce      json
//	@Param        request  body  dto.AdminApproveUserRequest  true  "User ID and target status (int)"
//	@Success      200  {object}  dto.AdminApproveUserResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      409  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/user/approve [post]
func (h *AdminHandler) PostAdminApproveUser(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.AdminApproveUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	// Validate status as enum
	target := permissionmodel.UserRoleStatus(req.Status)
	if err := h.userService.ApproveCreciManual(ctx, req.ID, target); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.AdminApproveUserResponse{Message: "Status updated"})
}
