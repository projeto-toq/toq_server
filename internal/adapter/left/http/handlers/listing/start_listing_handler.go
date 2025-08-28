package listing

import (
	"context"
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
	httpmodels "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/models"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
)

// StartListingHandler representa o handler responsável pelo início de anúncios
// Este handler encapsula toda a lógica HTTP relacionada à criação de novos anúncios
type StartListingHandler struct {
	service listingservices.ListingServiceInterface
}

// NewStartListingHandler cria uma nova instância do handler de início de anúncios
// Parâmetros:
//   - service: Interface do serviço de anúncios que contém a lógica de negócio
//
// Retorna:
//   - *StartListingHandler: Nova instância do handler
func NewStartListingHandler(service listingservices.ListingServiceInterface) *StartListingHandler {
	return &StartListingHandler{
		service: service,
	}
}

// StartListing manipula requisições HTTP para iniciar um novo anúncio
// Este endpoint requer autenticação e permite que usuários criem anúncios imobiliários
// Aplica validações de negócio e verifica permissões antes de criar o anúncio
//
// @Summary Iniciar novo anúncio
// @Description Cria um novo anúncio imobiliário no sistema
// @Tags Listings
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body httpmodels.StartListingRequest true "Dados do anúncio a ser criado"
// @Success 201 {object} httpmodels.StartListingResponse "Anúncio criado com sucesso"
// @Failure 400 {object} httpmodels.ErrorResponse "Dados de entrada inválidos"
// @Failure 401 {object} httpmodels.ErrorResponse "Usuário não autenticado"
// @Failure 500 {object} httpmodels.ErrorResponse "Erro interno do servidor"
// @Router /api/v2/listings [post]
func (h *StartListingHandler) StartListing(c *gin.Context) {
	// Obtém informações do usuário autenticado a partir do contexto
	// O middleware de autenticação injeta essas informações no contexto
	userInfo, exists := middlewares.GetUserInfoFromContext(c)
	if !exists {
		slog.Error("Usuário não encontrado no contexto de autenticação")
		c.JSON(http.StatusUnauthorized, httpmodels.ErrorResponse{
			Error:   "unauthorized",
			Code:    http.StatusUnauthorized,
			Message: "Usuário não autenticado",
		})
		return
	}

	var req httpmodels.StartListingRequest

	// Extrai e valida os dados JSON da requisição
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Erro ao fazer bind do JSON para StartListing", "error", err)
		c.JSON(http.StatusBadRequest, httpmodels.ErrorResponse{
			Error:   "bad_request",
			Code:    http.StatusBadRequest,
			Message: "Formato de requisição inválido",
		})
		return
	}

	// Cria contexto com informações do usuário para o serviço
	// O contexto da requisição HTTP não contém as informações do Gin
	ctx := context.WithValue(c.Request.Context(), globalmodel.TokenKey, userInfo)

	// Chama o serviço para criar o anúncio
	// O serviço aplica validações de negócio específicas para anúncios
	// Verifica se o usuário tem permissão para criar anúncios
	// Valida dados do imóvel e informações de localização
	listing, err := h.service.StartListing(
		ctx,
		req.ZipCode,
		req.Number,
		globalmodel.PropertyType(req.PropertyType),
	)
	if err != nil {
		slog.Error("Erro ao criar anúncio", "error", err, "userID", userInfo.ID)
		c.JSON(http.StatusInternalServerError, httpmodels.ErrorResponse{
			Error:   "internal_error",
			Code:    http.StatusInternalServerError,
			Message: "Falha ao criar anúncio",
		})
		return
	}

	// Prepara resposta de sucesso com ID do anúncio criado
	// O ID permite que o cliente acompanhe o status do anúncio
	response := httpmodels.StartListingResponse{
		ID: listing.ID(),
	}

	// Retorna resposta de sucesso com status 201 Created
	c.JSON(http.StatusCreated, response)
}
