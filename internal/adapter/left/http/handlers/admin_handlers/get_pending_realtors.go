package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

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
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
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
