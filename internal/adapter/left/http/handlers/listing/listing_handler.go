package listing

import (
	"log/slog"

	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
)

// ListingHandler representa o handler DEPRECIADO de listagens
//
// ATENÇÃO: Este handler foi refatorado e dividido em handlers especializados seguindo Clean Architecture.
// Utilize os novos handlers específicos em vez deste:
//   - StartListingHandler: Para criação de novas listagens
//   - GetBaseFeaturesHandler: Para obter características base
//   - GetOptionsHandler: Para obter opções de propriedade
//   - GetAllListingsHandler: Para listar propriedades do usuário
//   - GetAllOffersHandler: Para listar ofertas do usuário
//   - GetAllVisitsHandler: Para listar visitas do usuário
//
// Este arquivo será removido em versões futuras.
type ListingHandler struct {
	service listingservices.ListingServiceInterface
}

// NewListingHandler cria uma nova instância do handler DEPRECIADO
//
// DEPRECIADO: Use os novos handlers especializados em vez deste.
func NewListingHandler(service listingservices.ListingServiceInterface) *ListingHandler {
	slog.Warn("ListingHandler está depreciado. Use os handlers especializados (StartListingHandler, GetBaseFeaturesHandler, etc.)")
	return &ListingHandler{
		service: service,
	}
}
