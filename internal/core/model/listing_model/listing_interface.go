package listingmodel

import globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"

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
	Title() string
	SetTitle(title string)
	HasTitle() bool
	UnsetTitle()
	ListingType() globalmodel.PropertyType
	SetListingType(listingType globalmodel.PropertyType)
	Owner() PropertyOwner
	SetOwner(owner PropertyOwner)
	HasOwner() bool
	UnsetOwner()
	Features() []FeatureInterface
	SetFeatures(features []FeatureInterface)
	LandSize() float64
	SetLandSize(landsize float64)
	HasLandSize() bool
	UnsetLandSize()
	Corner() bool
	SetCorner(corner bool)
	HasCorner() bool
	UnsetCorner()
	NonBuildable() float64
	SetNonBuildable(nonBuildable float64)
	HasNonBuildable() bool
	UnsetNonBuildable()
	Buildable() float64
	SetBuildable(buildable float64)
	HasBuildable() bool
	UnsetBuildable()
	Delivered() PropertyDelivered
	SetDelivered(delivered PropertyDelivered)
	HasDelivered() bool
	UnsetDelivered()
	WhoLives() WhoLives
	SetWhoLives(whoLives WhoLives)
	HasWhoLives() bool
	UnsetWhoLives()
	Description() string
	SetDescription(description string)
	HasDescription() bool
	UnsetDescription()
	Transaction() TransactionType
	SetTransaction(transaction TransactionType)
	HasTransaction() bool
	UnsetTransaction()
	SellNet() float64
	SetSellNet(sellNet float64)
	HasSellNet() bool
	UnsetSellNet()
	RentNet() float64
	SetRentNet(rentNet float64)
	HasRentNet() bool
	UnsetRentNet()
	Condominium() float64
	SetCondominium(condominium float64)
	HasCondominium() bool
	UnsetCondominium()
	AnnualTax() float64
	SetAnnualTax(annualTax float64)
	HasAnnualTax() bool
	UnsetAnnualTax()
	MonthlyTax() float64
	SetMonthlyTax(monthlyTax float64)
	HasMonthlyTax() bool
	UnsetMonthlyTax()
	AnnualGroundRent() float64
	SetAnnualGroundRent(annualGroundRent float64)
	HasAnnualGroundRent() bool
	UnsetAnnualGroundRent()
	MonthlyGroundRent() float64
	SetMonthlyGroundRent(monthlyGroundRent float64)
	HasMonthlyGroundRent() bool
	UnsetMonthlyGroundRent()
	Exchange() bool
	SetExchange(exchange bool)
	HasExchange() bool
	UnsetExchange()
	ExchangePercentual() float64
	SetExchangePercentual(exchangePercentual float64)
	HasExchangePercentual() bool
	UnsetExchangePercentual()
	ExchangePlaces() []ExchangePlaceInterface
	SetExchangePlaces(exchangePlaces []ExchangePlaceInterface)
	Installment() InstallmentPlan
	SetInstallment(installment InstallmentPlan)
	HasInstallment() bool
	UnsetInstallment()
	Financing() bool
	SetFinancing(financing bool)
	HasFinancing() bool
	UnsetFinancing()
	FinancingBlockers() []FinancingBlockerInterface
	SetFinancingBlockers(financingBlockers []FinancingBlockerInterface)
	Guarantees() []GuaranteeInterface
	SetGuarantees(guarantees []GuaranteeInterface)
	Visit() VisitType
	SetVisit(visit VisitType)
	HasVisit() bool
	UnsetVisit()
	TenantName() string
	SetTenantName(tenantName string)
	HasTenantName() bool
	UnsetTenantName()
	TenantEmail() string
	SetTenantEmail(tenantEmail string)
	HasTenantEmail() bool
	UnsetTenantEmail()
	TenantPhone() string
	SetTenantPhone(tenantPhone string)
	HasTenantPhone() bool
	UnsetTenantPhone()
	Accompanying() AccompanyingType
	SetAccompanying(accompanying AccompanyingType)
	HasAccompanying() bool
	UnsetAccompanying()
	Deleted() bool
	SetDeleted(deleted bool)
	HasDeleted() bool
	UnsetDeleted()
}

func NewListing() ListingInterface {
	return &listing{}
}
