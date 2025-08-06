package listingmodel

import (
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

type ListingInterface interface {
	ID() int64
	SetID(id int64)
	UserID() int64
	SetUserID(userID int64)
	Code() uint32
	SetCode(code uint32)
	Version() uint8
	SetVersion(version uint8)
	Status() ListingStatus
	SetStatus(status ListingStatus)
	ZipCode() string
	SetZipCode(zipCode string)
	Street() string
	SetStreet(street string)
	Number() string
	SetNumber(number string)
	Complement() string
	SetComplement(complement string)
	Neighborhood() string
	SetNeighborhood(neighborhood string)
	City() string
	SetCity(city string)
	State() string
	SetState(state string)
	ListingType() globalmodel.PropertyType
	SetListingType(listingType globalmodel.PropertyType)
	Owner() PropertyOwner
	SetOwner(owner PropertyOwner)
	Features() []FeatureInterface
	SetFeatures(features []FeatureInterface)
	LandSize() float64
	SetLandSize(landsize float64)
	Corner() bool
	SetCorner(corner bool)
	NonBuildable() float64
	SetNonBuildable(nonBuildable float64)
	Buildable() float64
	SetBuildable(buildable float64)
	Delivered() PropertyDelivered
	SetDelivered(delivered PropertyDelivered)
	WhoLives() WhoLives
	SetWhoLives(whoLives WhoLives)
	Description() string
	SetDescription(description string)
	Transaction() TransactionType
	SetTransaction(transaction TransactionType)
	SellNet() float64
	SetSellNet(sellNet float64)
	RentNet() float64
	SetRentNet(rentNet float64)
	Condominium() float64
	SetCondominium(condominium float64)
	AnnualTax() float64
	SetAnnualTax(annualTax float64)
	AnnualGroundRent() float64
	SetAnnualGroundRent(annualGroundRent float64)
	Exchange() bool
	SetExchange(exchange bool)
	ExchangePercentual() float64
	SetExchangePercentual(exchangePercentual float64)
	ExchangePlaces() []ExchangePlaceInterface
	SetExchangePlaces(exchangePlaces []ExchangePlaceInterface)
	Installment() InstallmentPlan
	SetInstallment(installment InstallmentPlan)
	Financing() bool
	SetFinancing(financing bool)
	FinancingBlockers() []FinancingBlockerInterface
	SetFinancingBlockers(financingBlockers []FinancingBlockerInterface)
	Guarantees() []GuaranteeInterface
	SetGuarantees(guarantees []GuaranteeInterface)
	Visit() VisitType
	SetVisit(visit VisitType)
	TenantName() string
	SetTenantName(tenantName string)
	TenantEmail() string
	SetTenantEmail(tenantEmail string)
	TenantPhone() string
	SetTenantPhone(tenantPhone string)
	Accompanying() AccompanyingType
	SetAccompanying(accompanying AccompanyingType)
	Deleted() bool
	SetDeleted(deleted bool)

	ToSQLNullString(input string) sql.NullString
	ToSQLNullInt(input any) sql.NullInt64
	ToSQLNullFloat64(value float64) sql.NullFloat64
}

func NewListing() ListingInterface {
	return &listing{}
}
