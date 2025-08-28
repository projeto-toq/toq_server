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

// GetAllListingsHandler representa o handler responsável pela obtenção de todas as listagens do usuário
// Este handler encapsula toda a lógica HTTP relacionada à busca de listagens próprias do usuário autenticado
type GetAllListingsHandler struct {
	service listingservices.ListingServiceInterface
}

// NewGetAllListingsHandler cria uma nova instância do handler de listagens do usuário
//
// Parâmetros:
//   - service: Interface do serviço de listings para operações de negócio
//
// Retorna:
//   - *GetAllListingsHandler: Nova instância do handler
func NewGetAllListingsHandler(service listingservices.ListingServiceInterface) *GetAllListingsHandler {
	return &GetAllListingsHandler{
		service: service,
	}
}

// GetAllListings obtém todas as listagens de propriedades do usuário autenticado
//
// Este endpoint retorna uma lista completa de todas as propriedades que o usuário
// atual tem listadas na plataforma, incluindo informações detalhadas de cada listagem.
//
// @Summary Get all listings
// @Description Get all listings for the authenticated user
// @Tags Listings
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} httpmodels.GetAllListingsResponse
// @Failure 401 {object} httpmodels.ErrorResponse
// @Failure 500 {object} httpmodels.ErrorResponse
// @Router /api/v2/listings [get]
func (h *GetAllListingsHandler) GetAllListings(c *gin.Context) {
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

	// Busca todas as listagens do usuário através do serviço
	listings, err := h.service.GetAllListingsByUser(ctx, userInfo.ID)
	if err != nil {
		slog.Error("Error getting all listings", "error", err, "user_id", userInfo.ID)
		c.JSON(http.StatusInternalServerError, httpmodels.ErrorResponse{
			Error:   "internal_error",
			Code:    http.StatusInternalServerError,
			Message: "Failed to get listings",
		})
		return
	}

	// Conversão dos modelos de domínio para modelos HTTP
	httpListings := make([]httpmodels.Listing, 0, len(listings))
	for _, listing := range listings {
		httpListings = append(httpListings, shared.ConvertDomainListingToHTTP(listing))
	}

	// Preparação e envio da resposta
	response := httpmodels.GetAllListingsResponse{
		Listings: httpListings,
	}

	c.JSON(http.StatusOK, response)
}
