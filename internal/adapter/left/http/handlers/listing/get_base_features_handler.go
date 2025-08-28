package listing

import (
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
	httpmodels "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/models"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
)

// GetBaseFeaturesHandler representa o handler responsável pela obtenção das características base dos imóveis
// Este handler encapsula toda a lógica HTTP relacionada à busca de features/características disponíveis
type GetBaseFeaturesHandler struct {
	service listingservices.ListingServiceInterface
}

// NewGetBaseFeaturesHandler cria uma nova instância do handler de características base
//
// Parâmetros:
//   - service: Interface do serviço de listings para operações de negócio
//
// Retorna:
//   - *GetBaseFeaturesHandler: Nova instância do handler
func NewGetBaseFeaturesHandler(service listingservices.ListingServiceInterface) *GetBaseFeaturesHandler {
	return &GetBaseFeaturesHandler{
		service: service,
	}
}

// GetBaseFeatures obtém as características base disponíveis para listagens de imóveis
//
// Este endpoint retorna todas as características/features que podem ser associadas
// a uma propriedade durante o processo de criação de listagem.
//
// @Summary Get base features
// @Description Get available base features for property listings
// @Tags Listings
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} httpmodels.GetBaseFeaturesResponse
// @Failure 401 {object} httpmodels.ErrorResponse
// @Failure 500 {object} httpmodels.ErrorResponse
// @Router /api/v2/listings/features [get]
func (h *GetBaseFeaturesHandler) GetBaseFeatures(c *gin.Context) {
	// Verificação de autenticação do usuário
	_, exists := middlewares.GetUserInfoFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, httpmodels.ErrorResponse{
			Error:   "unauthorized",
			Code:    http.StatusUnauthorized,
			Message: "User not authenticated",
		})
		return
	}

	// Busca as características base através do serviço
	features, err := h.service.GetBaseFeatures(c.Request.Context())
	if err != nil {
		slog.Error("Error getting base features", "error", err)
		c.JSON(http.StatusInternalServerError, httpmodels.ErrorResponse{
			Error:   "internal_error",
			Code:    http.StatusInternalServerError,
			Message: "Failed to get base features",
		})
		return
	}

	// Conversão dos modelos de domínio para modelos HTTP
	httpFeatures := make([]httpmodels.BaseFeature, 0, len(features))
	for _, feature := range features {
		httpFeatures = append(httpFeatures, httpmodels.BaseFeature{
			ID:          feature.ID(),
			Feature:     feature.Feature(),
			Description: feature.Description(),
		})
	}

	// Preparação e envio da resposta
	response := httpmodels.GetBaseFeaturesResponse{
		Features: httpFeatures,
	}

	c.JSON(http.StatusOK, response)
}
