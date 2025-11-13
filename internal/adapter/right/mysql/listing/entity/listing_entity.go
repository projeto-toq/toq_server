package listingentity

import (
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

type ListingEntity struct {
	ID                 int64
	UserID             int64
	Code               uint32
	Version            uint8
	Status             uint8
	ZipCode            string
	Street             sql.NullString
	Number             string
	Complement         sql.NullString
	Neighborhood       sql.NullString
	City               sql.NullString
	State              sql.NullString
	Title              sql.NullString
	ListingType        uint8
	Owner              sql.NullInt16
	Features           []EntityFeature
	LandSize           sql.NullFloat64
	Corner             sql.NullInt16
	NonBuildable       sql.NullFloat64
	Buildable          sql.NullFloat64
	Delivered          sql.NullInt16
	WhoLives           sql.NullInt16
	Description        sql.NullString
	Transaction        sql.NullInt16
	SellNet            sql.NullFloat64
	RentNet            sql.NullFloat64
	Condominium        sql.NullFloat64
	AnnualTax          sql.NullFloat64
	MonthlyTax         sql.NullFloat64
	AnnualGroundRent   sql.NullFloat64
	MonthlyGroundRent  sql.NullFloat64
	Exchange           sql.NullInt16
	ExchangePercentual sql.NullFloat64
	ExchangePlaces     []EntityExchangePlace
	Installment        sql.NullInt16
	Financing          sql.NullInt16
	FinancingBlocker   []EntityFinancingBlocker
	Guarantees         []EntityGuarantee
	Visit              sql.NullInt16
	TenantName         sql.NullString
	TenantEmail        sql.NullString
	TenantPhone        sql.NullString
	Accompanying       sql.NullInt16
	Deleted            sql.NullInt16
}

func (e *ListingEntity) ToString(entity sql.NullString) string {
	if entity.Valid {
		return entity.String
	}
	return ""
}

func (e *ListingEntity) ToUint8(entity sql.NullInt16) uint8 {
	if entity.Valid {
		return uint8(entity.Int16)
	}
	return 0
}

func (e *ListingEntity) ToFloat64(entity sql.NullFloat64) float64 {
	if entity.Valid {
		return entity.Float64
	}
	return 0
}

func (e *ListingEntity) FeaturesToDomain() (features []listingmodel.FeatureInterface) {
	for _, entity := range e.Features {
		feature := listingmodel.NewFeature()
		feature.SetID(entity.ID)
		feature.SetListingID(entity.ListingID)
		feature.SetFeatureID(entity.FeatureID)
		feature.SetQuantity(entity.Quantity)
		features = append(features, feature)
	}
	return
}

func (e *ListingEntity) ExchangePlacesToDomain() (exchangePlaces []listingmodel.ExchangePlaceInterface) {
	for _, entity := range e.ExchangePlaces {
		exchangePlace := listingmodel.NewExchangePlace()
		exchangePlace.SetID(entity.ID)
		exchangePlace.SetListingID(entity.ListingID)
		exchangePlace.SetNeighborhood(entity.Neighborhood)
		exchangePlace.SetCity(entity.City)
		exchangePlace.SetState(entity.State)
		exchangePlaces = append(exchangePlaces, exchangePlace)
	}
	return
}

func (e *ListingEntity) FinancingBlockersToDomain() (financingBlockers []listingmodel.FinancingBlockerInterface) {
	for _, entity := range e.FinancingBlocker {
		financingBlocker := listingmodel.NewFinancingBlocker()
		financingBlocker.SetID(entity.ID)
		financingBlocker.SetListingID(entity.ListingID)
		financingBlocker.SetBlocker(listingmodel.FinancingBlocker(entity.Blocker))
		financingBlockers = append(financingBlockers, financingBlocker)
	}
	return
}

func (e *ListingEntity) GuaranteesToDomain() (guarantees []listingmodel.GuaranteeInterface) {
	for _, entity := range e.Guarantees {
		guarantee := listingmodel.NewGuarantee()
		guarantee.SetID(entity.ID)
		guarantee.SetListingID(entity.ListingID)
		guarantee.SetPriority(entity.Priority)
		guarantee.SetGuarantee(listingmodel.GuaranteeType(entity.Guarantee))
		guarantees = append(guarantees, guarantee)
	}
	return
}
