package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetPendingRealtors handles GET /admin/users/creci/pending
//
//	@Summary      List realtors pending manual validation
//	@Description  Returns id, nickname, fullName, nationalID, creciNumber, creciValidity, creciState
//	@Tags         Admin
//	@Produce      json
//	@Param        page   query  int  false  "Page number" default(1) example(1)
//	@Param        limit  query  int  false  "Page size" default(20) example(20)
//	@Success      200  {object}  dto.AdminGetPendingRealtorsResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/users/creci/pending [get]
func (h *AdminHandler) GetPendingRealtors(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminGetPendingRealtorsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	result, err := h.userService.ListPendingRealtors(ctx, req.Page, req.Limit)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminGetPendingRealtorsResponse{
		Realtors: make([]dto.AdminPendingRealtor, 0, len(result.Realtors)),
		Pagination: dto.PaginationResponse{
			Page:       result.Page,
			Limit:      result.Limit,
			Total:      result.Total,
			TotalPages: computeTotalPages(result.Total, result.Limit),
		},
	}

	for _, r := range result.Realtors {
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
