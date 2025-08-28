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

// GetAllOffersHandler representa o handler responsável pela obtenção de todas as ofertas do usuário
// Este handler encapsula toda a lógica HTTP relacionada à busca de ofertas feitas ou recebidas pelo usuário
type GetAllOffersHandler struct {
	service listingservices.ListingServiceInterface
}

// NewGetAllOffersHandler cria uma nova instância do handler de ofertas do usuário
//
// Parâmetros:
//   - service: Interface do serviço de listings para operações de negócio
//
// Retorna:
//   - *GetAllOffersHandler: Nova instância do handler
func NewGetAllOffersHandler(service listingservices.ListingServiceInterface) *GetAllOffersHandler {
	return &GetAllOffersHandler{
		service: service,
	}
}

// GetAllOffers obtém todas as ofertas relacionadas ao usuário autenticado
//
// Este endpoint retorna uma lista completa de todas as ofertas que o usuário
// fez ou recebeu, incluindo detalhes sobre status, valores e propriedades relacionadas.
//
// @Summary Get all offers
// @Description Get all offers for the authenticated user
// @Tags Listings
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} httpmodels.GetAllOffersResponse
// @Failure 401 {object} httpmodels.ErrorResponse
// @Failure 500 {object} httpmodels.ErrorResponse
// @Router /api/v2/offers [get]
func (h *GetAllOffersHandler) GetAllOffers(c *gin.Context) {
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

	// Busca todas as ofertas do usuário através do serviço
	offers, err := h.service.GetAllOffersByUser(ctx, userInfo.ID)
	if err != nil {
		slog.Error("Error getting all offers", "error", err, "user_id", userInfo.ID)
		c.JSON(http.StatusInternalServerError, httpmodels.ErrorResponse{
			Error:   "internal_error",
			Code:    http.StatusInternalServerError,
			Message: "Failed to get offers",
		})
		return
	}

	// Conversão dos modelos de domínio para modelos HTTP
	httpOffers := make([]httpmodels.Offer, 0, len(offers))
	for _, offer := range offers {
		httpOffers = append(httpOffers, shared.ConvertDomainOfferToHTTP(offer))
	}

	// Preparação e envio da resposta
	response := httpmodels.GetAllOffersResponse{
		Offers: httpOffers,
	}

	c.JSON(http.StatusOK, response)
}
