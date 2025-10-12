package listingservices

import (
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateListingInput encapsula os campos opcionais para atualização de um listing.
type UpdateListingInput struct {
	ID                 int64
	Owner              coreutils.Optional[listingmodel.PropertyOwner]
	Features           coreutils.Optional[[]listingmodel.FeatureInterface]
	LandSize           coreutils.Optional[float64]
	Corner             coreutils.Optional[bool]
	NonBuildable       coreutils.Optional[float64]
	Buildable          coreutils.Optional[float64]
	Delivered          coreutils.Optional[listingmodel.PropertyDelivered]
	WhoLives           coreutils.Optional[listingmodel.WhoLives]
	Description        coreutils.Optional[string]
	Transaction        coreutils.Optional[listingmodel.TransactionType]
	SellNet            coreutils.Optional[float64]
	RentNet            coreutils.Optional[float64]
	Condominium        coreutils.Optional[float64]
	AnnualTax          coreutils.Optional[float64]
	AnnualGroundRent   coreutils.Optional[float64]
	Exchange           coreutils.Optional[bool]
	ExchangePercentual coreutils.Optional[float64]
	ExchangePlaces     coreutils.Optional[[]listingmodel.ExchangePlaceInterface]
	Installment        coreutils.Optional[listingmodel.InstallmentPlan]
	Financing          coreutils.Optional[bool]
	FinancingBlockers  coreutils.Optional[[]listingmodel.FinancingBlockerInterface]
	Guarantees         coreutils.Optional[[]listingmodel.GuaranteeInterface]
	Visit              coreutils.Optional[listingmodel.VisitType]
	TenantName         coreutils.Optional[string]
	TenantEmail        coreutils.Optional[string]
	TenantPhone        coreutils.Optional[string]
	Accompanying       coreutils.Optional[listingmodel.AccompanyingType]
}
