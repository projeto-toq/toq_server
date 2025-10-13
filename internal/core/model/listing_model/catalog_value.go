package listingmodel

// CatalogValueInterface representa um valor de catálogo utilizado em vários campos de Listing.
type CatalogValueInterface interface {
	ID() uint8
	SetID(uint8)
	Category() string
	SetCategory(string)
	Slug() string
	SetSlug(string)
	Label() string
	SetLabel(string)
	Description() *string
	SetDescription(*string)
	IsActive() bool
	SetIsActive(bool)
}

type catalogValue struct {
	id          uint8
	category    string
	slug        string
	label       string
	description *string
	isActive    bool
}

// NewCatalogValue inicializa uma nova instância de CatalogValueInterface.
func NewCatalogValue() CatalogValueInterface {
	return &catalogValue{}
}

func (cv *catalogValue) ID() uint8 {
	return cv.id
}

func (cv *catalogValue) SetID(id uint8) {
	cv.id = id
}

func (cv *catalogValue) Category() string {
	return cv.category
}

func (cv *catalogValue) SetCategory(category string) {
	cv.category = category
}

func (cv *catalogValue) Slug() string {
	return cv.slug
}

func (cv *catalogValue) SetSlug(slug string) {
	cv.slug = slug
}

func (cv *catalogValue) Label() string {
	return cv.label
}

func (cv *catalogValue) SetLabel(label string) {
	cv.label = label
}

func (cv *catalogValue) Description() *string {
	return cv.description
}

func (cv *catalogValue) SetDescription(description *string) {
	cv.description = description
}

func (cv *catalogValue) IsActive() bool {
	return cv.isActive
}

func (cv *catalogValue) SetIsActive(active bool) {
	cv.isActive = active
}

// AllowedCatalogCategories retorna todas as categorias de catálogo suportadas.
func AllowedCatalogCategories() []string {
	return []string{
		CatalogCategoryPropertyOwner,
		CatalogCategoryPropertyDelivered,
		CatalogCategoryWhoLives,
		CatalogCategoryTransactionType,
		CatalogCategoryInstallmentPlan,
		CatalogCategoryFinancingBlocker,
		CatalogCategoryVisitType,
		CatalogCategoryAccompanyingType,
		CatalogCategoryGuaranteeType,
	}
}

// IsValidCatalogCategory indica se a categoria enviada é suportada pelo domínio.
func IsValidCatalogCategory(category string) bool {
	for _, allowed := range AllowedCatalogCategories() {
		if allowed == category {
			return true
		}
	}
	return false
}
