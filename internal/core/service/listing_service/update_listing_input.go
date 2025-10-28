package listingservices

import (
	"strings"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateListingInput encapsula os campos opcionais para atualização de um listing.
type UpdateListingInput struct {
	ID                 int64
	Owner              coreutils.Optional[CatalogSelection]
	Features           coreutils.Optional[[]listingmodel.FeatureInterface]
	LandSize           coreutils.Optional[float64]
	Corner             coreutils.Optional[bool]
	NonBuildable       coreutils.Optional[float64]
	Buildable          coreutils.Optional[float64]
	Delivered          coreutils.Optional[CatalogSelection]
	WhoLives           coreutils.Optional[CatalogSelection]
	Title              coreutils.Optional[string]
	Description        coreutils.Optional[string]
	Transaction        coreutils.Optional[CatalogSelection]
	SellNet            coreutils.Optional[float64]
	RentNet            coreutils.Optional[float64]
	Condominium        coreutils.Optional[float64]
	AnnualTax          coreutils.Optional[float64]
	AnnualGroundRent   coreutils.Optional[float64]
	Exchange           coreutils.Optional[bool]
	ExchangePercentual coreutils.Optional[float64]
	ExchangePlaces     coreutils.Optional[[]listingmodel.ExchangePlaceInterface]
	Installment        coreutils.Optional[CatalogSelection]
	Financing          coreutils.Optional[bool]
	FinancingBlockers  coreutils.Optional[[]CatalogSelection]
	Guarantees         coreutils.Optional[[]GuaranteeUpdate]
	Visit              coreutils.Optional[CatalogSelection]
	TenantName         coreutils.Optional[string]
	TenantEmail        coreutils.Optional[string]
	TenantPhone        coreutils.Optional[string]
	Accompanying       coreutils.Optional[CatalogSelection]
}

// CatalogSelection armazena a seleção de catálogo enviada no payload.
type CatalogSelection struct {
	id   *uint8
	slug string
}

// NewCatalogSelectionFromID cria uma seleção baseada em ID.
func NewCatalogSelectionFromID(id uint8) CatalogSelection {
	return CatalogSelection{id: &id}
}

// NewCatalogSelectionFromSlug cria uma seleção baseada em slug.
func NewCatalogSelectionFromSlug(slug string) CatalogSelection {
	return CatalogSelection{slug: strings.TrimSpace(slug)}
}

// HasID indica se a seleção possui ID configurado.
func (s CatalogSelection) HasID() bool {
	return s.id != nil && *s.id > 0
}

// HasSlug indica se a seleção possui slug configurado.
func (s CatalogSelection) HasSlug() bool {
	return strings.TrimSpace(s.slug) != ""
}

// IDValue retorna o ID quando presente.
func (s CatalogSelection) IDValue() uint8 {
	if s.id == nil {
		return 0
	}
	return *s.id
}

// SlugValue retorna o slug normalizado.
func (s CatalogSelection) SlugValue() string {
	return strings.TrimSpace(strings.ToLower(s.slug))
}

// WithResolvedID atribui o ID resolvido, preservando o slug informado.
func (s CatalogSelection) WithResolvedID(id uint8) CatalogSelection {
	s.id = &id
	return s
}

// GuaranteeUpdate agrupa prioridade e seleção de garantia vinda do payload.
type GuaranteeUpdate struct {
	Priority  uint8
	Selection CatalogSelection
}
