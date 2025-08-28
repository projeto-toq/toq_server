package listing

import (
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
	httpmodels "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/models"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
)

// GetOptionsHandler representa o handler responsável pela obtenção das opções de tipo de propriedade
// Este handler encapsula toda a lógica HTTP relacionada à busca de tipos de imóveis por localização
type GetOptionsHandler struct {
	service listingservices.ListingServiceInterface
}

// NewGetOptionsHandler cria uma nova instância do handler de opções de propriedade
//
// Parâmetros:
//   - service: Interface do serviço de listings para operações de negócio
//
// Retorna:
//   - *GetOptionsHandler: Nova instância do handler
func NewGetOptionsHandler(service listingservices.ListingServiceInterface) *GetOptionsHandler {
	return &GetOptionsHandler{
		service: service,
	}
}

// GetOptions obtém as opções de tipo de propriedade disponíveis para uma determinada localização
//
// Este endpoint retorna os tipos de propriedade disponíveis baseados no CEP e número
// fornecidos, permitindo que o usuário selecione o tipo adequado para sua listagem.
//
// @Summary Get property options
// @Description Get available property type options for a given location
// @Tags Listings
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body httpmodels.GetOptionsRequest true "Get options request"
// @Success 200 {object} httpmodels.GetOptionsResponse
// @Failure 400 {object} httpmodels.ErrorResponse
// @Failure 401 {object} httpmodels.ErrorResponse
// @Failure 500 {object} httpmodels.ErrorResponse
// @Router /api/v2/listings/options [post]
func (h *GetOptionsHandler) GetOptions(c *gin.Context) {
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

	// Binding e validação do corpo da requisição
	var req httpmodels.GetOptionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Error binding JSON for GetOptions", "error", err)
		c.JSON(http.StatusBadRequest, httpmodels.ErrorResponse{
			Error:   "bad_request",
			Code:    http.StatusBadRequest,
			Message: "Invalid request format",
		})
		return
	}

	// Busca as opções de tipo através do serviço
	types, err := h.service.GetOptions(c.Request.Context(), req.ZipCode, req.Number)
	if err != nil {
		slog.Error("Error getting options", "error", err)
		c.JSON(http.StatusInternalServerError, httpmodels.ErrorResponse{
			Error:   "internal_error",
			Code:    http.StatusInternalServerError,
			Message: "Failed to get options",
		})
		return
	}

	// Preparação e envio da resposta
	response := httpmodels.GetOptionsResponse{
		Types: types,
	}

	c.JSON(http.StatusOK, response)
}
