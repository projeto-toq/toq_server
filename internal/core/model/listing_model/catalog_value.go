package listingmodel

// CatalogValueInterface represents a catalog value used across multiple listing fields.
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

// NewCatalogValue creates a new CatalogValueInterface instance.
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

// AllowedCatalogCategories returns all supported catalog categories.
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

// IsValidCatalogCategory indicates whether the provided category is supported by the domain.
func IsValidCatalogCategory(category string) bool {
	for _, allowed := range AllowedCatalogCategories() {
		if allowed == category {
			return true
		}
	}
	return false
}
