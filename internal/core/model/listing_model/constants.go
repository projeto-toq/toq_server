package listingmodel

type ListingStatus uint8

const (
	StatusDraft ListingStatus = iota + 1
	StatusAwaitingPhoto
	StatusAwaitingApproval
	StatusPublished
)

func (s ListingStatus) String() string {
	switch s {
	case StatusDraft:
		return "Draft"
	case StatusAwaitingPhoto:
		return "Awaiting Photo"
	case StatusAwaitingApproval:
		return "Awaiting Approval"
	case StatusPublished:
		return "Published"
	default:
		return "Unknown"
	}
}

type PropertyOwner uint8
type PropertyDelivered uint8
type WhoLives uint8
type TransactionType uint8
type InstallmentPlan uint8
type FinancingBlocker uint8
type VisitType uint8
type AccompanyingType uint8
type GuaranteeType uint8

const (
	CatalogCategoryPropertyOwner     = "property_owner"
	CatalogCategoryPropertyDelivered = "property_delivered"
	CatalogCategoryWhoLives          = "who_lives"
	CatalogCategoryTransactionType   = "transaction_type"
	CatalogCategoryInstallmentPlan   = "installment_plan"
	CatalogCategoryFinancingBlocker  = "financing_blocker"
	CatalogCategoryVisitType         = "visit_type"
	CatalogCategoryAccompanyingType  = "accompanying_type"
	CatalogCategoryGuaranteeType     = "guarantee_type"
)
