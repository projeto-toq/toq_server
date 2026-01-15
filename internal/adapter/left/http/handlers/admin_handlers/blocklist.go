package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListBlocklist handles GET /admin/blocklist
//
//	@Summary	List access-token blocklist entries
//	@Tags		Admin Blocklist
//	@Produce	json
//	@Param		page		query	int	false	"Page number" default(1)
//	@Param		pageSize	query	int	false	"Page size (max 500)" default(100)
//	@Success	200	{object}	dto.AdminBlocklistListResponse
//	@Failure	400	{object} map[string]any
//	@Failure	401	{object} map[string]any
//	@Failure	403	{object} map[string]any
//	@Failure	500	{object} map[string]any
//	@Router		/admin/blocklist [get]
func (h *AdminHandler) ListBlocklist(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	if h.tokenBlocklist == nil {
		httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Blocklist adapter not configured"))
		return
	}

	var req dto.AdminBlocklistListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 100
	}

	items, err := h.tokenBlocklist.List(ctx, page, pageSize)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Failed to list blocklist"))
		return
	}
	total, errCount := h.tokenBlocklist.Count(ctx)
	if errCount != nil {
		total = int64(len(items))
	}

	resp := dto.AdminBlocklistListResponse{
		Items: make([]dto.AdminBlocklistItemResponse, 0, len(items)),
		Pagination: dto.PaginationResponse{
			Page:       int(page),
			Limit:      int(pageSize),
			Total:      total,
			TotalPages: computeTotalPages(total, int(pageSize)),
		},
	}

	for _, it := range items {
		resp.Items = append(resp.Items, dto.AdminBlocklistItemResponse{JTI: it.JTI, ExpiresAt: it.ExpiresAt})
	}

	c.JSON(http.StatusOK, resp)
}

// AddBlocklist handles POST /admin/blocklist
//
//	@Summary	Add an access-token JTI to the blocklist
//	@Tags		Admin Blocklist
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminBlocklistAddRequest	true	"Blocklist add payload"
//	@Success	200	{object} map[string]any
//	@Failure	400	{object} map[string]any
//	@Failure	401	{object} map[string]any
//	@Failure	403	{object} map[string]any
//	@Failure	500	{object} map[string]any
//	@Router		/admin/blocklist [post]
func (h *AdminHandler) AddBlocklist(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	if h.tokenBlocklist == nil {
		httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Blocklist adapter not configured"))
		return
	}

	var req dto.AdminBlocklistAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	if req.TTL <= 0 {
		httperrors.SendHTTPErrorObj(c, coreutils.BadRequest("ttl must be positive"))
		return
	}

	if err := h.tokenBlocklist.Add(ctx, req.JTI, req.TTL); err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Failed to add blocklist entry"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "JTI added to blocklist"})
}

// DeleteBlocklist handles DELETE /admin/blocklist
//
//	@Summary	Remove an access-token JTI from the blocklist
//	@Tags		Admin Blocklist
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.AdminBlocklistDeleteRequest	true	"Blocklist delete payload"
//	@Success	200	{object} map[string]any
//	@Failure	400	{object} map[string]any
//	@Failure	401	{object} map[string]any
//	@Failure	403	{object} map[string]any
//	@Failure	500	{object} map[string]any
//	@Router		/admin/blocklist [delete]
func (h *AdminHandler) DeleteBlocklist(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	if h.tokenBlocklist == nil {
		httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Blocklist adapter not configured"))
		return
	}

	var req dto.AdminBlocklistDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	if err := h.tokenBlocklist.Delete(ctx, req.JTI); err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Failed to delete blocklist entry"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "JTI removed from blocklist"})
}
