package globalhandlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
)

// CSPHandler exposes endpoints to manage Content Security Policy configuration.
type CSPHandler struct {
	globalService globalservice.GlobalServiceInterface
}

// NewCSPHandlerAdapter creates a new handler instance.
func NewCSPHandlerAdapter(globalService globalservice.GlobalServiceInterface) *CSPHandler {
	return &CSPHandler{globalService: globalService}
}

// GetCSPPolicy handles GET /admin/security/csp
//
//	@Summary      Get the current Content Security Policy
//	@Description  Returns the active CSP directives and version used by the platform
//	@Tags         Security
//	@Produce      json
//	@Security     BearerAuth
//	@Success      200  {object}  dto.CSPPolicyResponse
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      404  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/security/csp [get]
func (h *CSPHandler) GetCSPPolicy(c *gin.Context) {
	ctx := c.Request.Context()
	policy, err := h.globalService.GetCSPPolicy(ctx)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.CSPPolicyResponse{
		Version:    policy.Version,
		Directives: policy.Directives,
	})
}

// UpdateCSPPolicy handles PUT /admin/security/csp
//
//	@Summary      Update the Content Security Policy
//	@Description  Replaces the active CSP directives using optimistic concurrency. Pass version 0 to create a new policy.
//	@Tags         Security
//	@Accept       json
//	@Produce      json
//	@Security     BearerAuth
//	@Param        request  body  dto.UpdateCSPPolicyRequest  true  "Content Security Policy payload"
//	@Success      200  {object}  dto.CSPPolicyResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      409  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/security/csp [put]
func (h *CSPHandler) UpdateCSPPolicy(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.UpdateCSPPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	if len(req.Directives) == 0 {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(errEmptyDirectives))
		return
	}

	updated, err := h.globalService.UpdateCSPPolicy(ctx, req.Version, req.Directives)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	slog.Info("security.csp.updated", "version", updated.Version)

	c.JSON(http.StatusOK, dto.CSPPolicyResponse{
		Version:    updated.Version,
		Directives: updated.Directives,
	})
}

var errEmptyDirectives = errors.New("directives cannot be empty")
