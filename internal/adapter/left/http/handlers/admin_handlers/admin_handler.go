package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

type AdminHandler struct {
	userService    userservices.UserServiceInterface
	listingService listingservices.ListingServiceInterface
}

func NewAdminHandlerAdapter(userService userservices.UserServiceInterface, listingService listingservices.ListingServiceInterface) *AdminHandler {
	return &AdminHandler{
		userService:    userService,
		listingService: listingService,
	}
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
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
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
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
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
	c.JSON(http.StatusOK, dto.AdminGetUserResponse(profile))
}

// PostAdminApproveUser handles POST /admin/user/approve
//
//	@Summary      Approve or refuse realtor status manually
//	@Description  Status must be one of the allowed enum values. On success, sends FCM notification.
//	@Tags         Admin
//	@Accept       json
//	@Produce      json
//	@Param        request  body  dto.AdminApproveUserRequest  true  "User ID and target status (enum: 0=active,10=refused_image,11=refused_document,12=refused_data)"
//	@Success      200  {object}  dto.AdminApproveUserResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      409  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/user/approve [post]
func (h *AdminHandler) PostAdminApproveUser(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminApproveUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	// Validate status as enum
	target, statusErr := req.ToStatus()
	if statusErr != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("status", statusErr.Error()))
		return
	}
	if err := h.userService.ApproveCreciManual(ctx, req.ID, target); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.AdminApproveUserResponse{Message: "Status updated"})
}

// PostAdminCreciDownloadURL handles POST /admin/user/creci-download-url
//
//	@Summary      Get signed download URLs for CRECI documents
//	@Description  Returns signed URLs (selfie/front/back) for a realtor user, valid for a limited time
//	@Tags         Admin
//	@Accept       json
//	@Produce      json
//	@Param        request  body  dto.AdminCreciDownloadURLRequest  true  "User ID"
//	@Success      200  {object}  dto.AdminCreciDownloadURLResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      404  {object}  map[string]any
//	@Failure      409  {object}  map[string]any
//	@Failure      422  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/user/creci-download-url [post]
func (h *AdminHandler) PostAdminCreciDownloadURL(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminCreciDownloadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	urls, err := h.userService.GetCreciDownloadURLs(ctx, req.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminCreciDownloadURLResponse{
		URLs: dto.AdminCreciDocumentURLs{
			Selfie: urls.Selfie,
			Front:  urls.Front,
			Back:   urls.Back,
		},
		ExpiresInMinutes: urls.ExpiresInMinutes,
	}

	c.JSON(http.StatusOK, resp)
}

// ListListingCatalogValues handles GET /admin/listing/catalog
//
//	@Summary	Listar valores de catálogo de listings
//	@Tags		Admin
//	@Produce	json
//	@Param		category		query	string	true	"Categoria do catálogo"
//	@Param		includeInactive	query	bool	false	"Retornar valores inativos"
//	@Success	200	{object}	dto.ListingCatalogValuesResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/listing/catalog [get]
func (h *AdminHandler) ListListingCatalogValues(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var query dto.AdminListingCatalogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	values, err := h.listingService.ListCatalogValues(ctx, query.Category, query.IncludeInactive)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToListingCatalogValuesResponse(values))
}

// CreateListingCatalogValue handles POST /admin/listing/catalog
//
//	@Summary	Criar valor de catálogo de listings
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.ListingCatalogCreateRequest	true	"Payload de criação"
//	@Success	201	{object}	dto.ListingCatalogValueResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	409	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/listing/catalog [post]
func (h *AdminHandler) CreateListingCatalogValue(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.ListingCatalogCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := listingservices.CreateCatalogValueInput{
		Category:    req.Category,
		Slug:        req.Slug,
		Label:       req.Label,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	value, err := h.listingService.CreateCatalogValue(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusCreated, httpconv.ToListingCatalogValueResponse(value))
}

// UpdateListingCatalogValue handles PUT /admin/listing/catalog
//
//	@Summary	Atualizar valor de catálogo de listings
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.ListingCatalogUpdateRequest	true	"Campos para atualização parcial"
//	@Success	200	{object}	dto.ListingCatalogValueResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	409	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/listing/catalog [put]
func (h *AdminHandler) UpdateListingCatalogValue(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.ListingCatalogUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := listingservices.UpdateCatalogValueInput{
		Category:    req.Category,
		ID:          req.ID,
		Slug:        req.Slug,
		Label:       req.Label,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	value, err := h.listingService.UpdateCatalogValue(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, httpconv.ToListingCatalogValueResponse(value))
}

// DeleteListingCatalogValue handles DELETE /admin/listing/catalog
//
//	@Summary	Desativar valor de catálogo de listings
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.ListingCatalogDeleteRequest	true	"Dados para desativação"
//	@Success	204	"Valor desativado com sucesso"
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/listing/catalog [delete]
func (h *AdminHandler) DeleteListingCatalogValue(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.ListingCatalogDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	if err := h.listingService.DeleteCatalogValue(ctx, req.Category, req.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
