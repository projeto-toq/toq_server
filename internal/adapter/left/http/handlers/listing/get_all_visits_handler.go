package listing

import (
	"context"
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/shared"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
	httpmodels "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/models"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
)

// GetAllVisitsHandler representa o handler responsável pela obtenção de todas as visitas do usuário
// Este handler encapsula toda a lógica HTTP relacionada à busca de visitas agendadas ou realizadas pelo usuário
type GetAllVisitsHandler struct {
	service listingservices.ListingServiceInterface
}

// NewGetAllVisitsHandler cria uma nova instância do handler de visitas do usuário
//
// Parâmetros:
//   - service: Interface do serviço de listings para operações de negócio
//
// Retorna:
//   - *GetAllVisitsHandler: Nova instância do handler
func NewGetAllVisitsHandler(service listingservices.ListingServiceInterface) *GetAllVisitsHandler {
	return &GetAllVisitsHandler{
		service: service,
	}
}

// GetAllVisits obtém todas as visitas relacionadas ao usuário autenticado
//
// Este endpoint retorna uma lista completa de todas as visitas que o usuário
// agendou ou que foram agendadas para suas propriedades, incluindo detalhes sobre
// datas, horários e status das visitas.
//
// @Summary Get all visits
// @Description Get all visits for the authenticated user
// @Tags Listings
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} httpmodels.GetAllVisitsResponse
// @Failure 401 {object} httpmodels.ErrorResponse
// @Failure 500 {object} httpmodels.ErrorResponse
// @Router /api/v2/visits [get]
func (h *GetAllVisitsHandler) GetAllVisits(c *gin.Context) {
	// Verificação de autenticação e obtenção das informações do usuário
	userInfo, exists := middlewares.GetUserInfoFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, httpmodels.ErrorResponse{
			Error:   "unauthorized",
			Code:    http.StatusUnauthorized,
			Message: "User not authenticated",
		})
		return
	}

	// Cria contexto com informações do usuário para o serviço
	// O contexto da requisição HTTP não contém as informações do Gin
	ctx := context.WithValue(c.Request.Context(), globalmodel.TokenKey, userInfo)

	// Busca todas as visitas do usuário através do serviço
	visits, err := h.service.GetAllVisitsByUser(ctx, userInfo.ID)
	if err != nil {
		slog.Error("Error getting all visits", "error", err, "user_id", userInfo.ID)
		c.JSON(http.StatusInternalServerError, httpmodels.ErrorResponse{
			Error:   "internal_error",
			Code:    http.StatusInternalServerError,
			Message: "Failed to get visits",
		})
		return
	}

	// Conversão dos modelos de domínio para modelos HTTP
	httpVisits := make([]httpmodels.Visit, 0, len(visits))
	for _, visit := range visits {
		httpVisits = append(httpVisits, shared.ConvertDomainVisitToHTTP(visit))
	}

	// Preparação e envio da resposta
	response := httpmodels.GetAllVisitsResponse{
		Visits: httpVisits,
	}

	c.JSON(http.StatusOK, response)
}
