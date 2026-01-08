package adminhandlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
)

// GetAdminRoutes handles GET /admin/permissions/routes
//
//	@Summary		List all registered HTTP routes
//	@Description	Retrieve a paginated list of all HTTP routes registered in the system.
//	@Description	Supports filtering by HTTP method (GET, POST, etc.) and path pattern (substring match).
//	@Description	The `action` field in the response follows the format `METHOD:PATH`, which is used for permission management.
//	@Description	Useful for API discovery, debugging, and populating frontend permission forms.
//	@Description	Requires admin-level permissions.
//	@Tags			Admin Permissions
//	@Produce		json
//	@Param			page			query		int		false	"Page number (1-indexed)"	minimum(1)	default(1)	Extensions(x-example=1)
//	@Param			limit			query		int		false	"Items per page"			minimum(1)	maximum(200)	default(50)	Extensions(x-example=50)
//	@Param			method			query		string	false	"Filter by HTTP method (case-insensitive)"	Enums(GET, POST, PUT, DELETE, PATCH)	Extensions(x-example=GET)
//	@Param			pathPattern		query		string	false	"Filter by path substring (case-insensitive)"	Extensions(x-example=admin)
//	@Success		200				{object}	dto.AdminListRoutesResponse	"Paginated list of routes"
//	@Failure		400				{object}	map[string]any				"Invalid query parameters"
//	@Failure		401				{object}	map[string]any				"Unauthorized (missing or invalid token)"
//	@Failure		403				{object}	map[string]any				"Forbidden (insufficient permissions)"
//	@Failure		500				{object}	map[string]any				"Internal server error"
//	@Router			/admin/permissions/routes [get]
//	@Security		BearerAuth
func (h *AdminHandler) GetAdminRoutes(c *gin.Context) {
	// Note: TelemetryMiddleware already provides request tracing
	// No need to create additional spans in handler layer

	// Parse and validate query parameters
	var req dto.AdminListRoutesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	// Get all routes from Gin router (introspection at request time)
	allRoutes := h.router.Routes()

	// Apply filters (method and path pattern)
	filtered := filterRoutes(allRoutes, req.Method, req.PathPattern)

	// Calculate pagination bounds
	total := int64(len(filtered))
	start, end := calculatePaginationBounds(req.Page, req.Limit, total)
	paginatedRoutes := filtered[start:end]

	// Convert Gin RouteInfo to DTO
	routes := make([]dto.RouteInfo, 0, len(paginatedRoutes))
	for _, route := range paginatedRoutes {
		routes = append(routes, dto.RouteInfo{
			Method:  route.Method,
			Path:    route.Path,
			Action:  fmt.Sprintf("%s:%s", route.Method, route.Path), // Format: METHOD:PATH
			Handler: extractHandlerName(route.Handler),
		})
	}

	// Build response with pagination metadata
	resp := dto.AdminListRoutesResponse{
		Routes: routes,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: computeTotalPages(total, req.Limit),
		},
	}

	c.JSON(http.StatusOK, resp)
}
