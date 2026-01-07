package listingmodel

import (
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

type ListingVersionInterface interface {
	ID() int64
	SetID(id int64)
	ListingIdentityID() int64
	SetListingIdentityID(listingIdentityID int64)
	ListingUUID() string
	SetListingUUID(listingUUID string)
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
	Complex() string
	SetComplex(complex string)
	HasComplex() bool
	UnsetComplex()
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
	// New property-specific fields
	CompletionForecast() string
	SetCompletionForecast(completionForecast string)
	HasCompletionForecast() bool
	UnsetCompletionForecast()
	LandBlock() string
	SetLandBlock(landBlock string)
	HasLandBlock() bool
	UnsetLandBlock()
	LandLot() string
	SetLandLot(landLot string)
	HasLandLot() bool
	UnsetLandLot()
	LandFront() float64
	SetLandFront(landFront float64)
	HasLandFront() bool
	UnsetLandFront()
	LandSide() float64
	SetLandSide(landSide float64)
	HasLandSide() bool
	UnsetLandSide()
	LandBack() float64
	SetLandBack(landBack float64)
	HasLandBack() bool
	UnsetLandBack()
	LandTerrainType() LandTerrainType
	SetLandTerrainType(landTerrainType LandTerrainType)
	HasLandTerrainType() bool
	UnsetLandTerrainType()
	HasKmz() bool
	SetHasKmz(hasKmz bool)
	HasHasKmz() bool
	UnsetHasKmz()
	KmzFile() string
	SetKmzFile(kmzFile string)
	HasKmzFile() bool
	UnsetKmzFile()
	BuildingFloors() int
	SetBuildingFloors(buildingFloors int)
	HasBuildingFloors() bool
	UnsetBuildingFloors()
	UnitTower() string
	SetUnitTower(unitTower string)
	HasUnitTower() bool
	UnsetUnitTower()
	UnitFloor() string
	SetUnitFloor(unitFloor string)
	HasUnitFloor() bool
	UnsetUnitFloor()
	UnitNumber() string
	SetUnitNumber(unitNumber string)
	HasUnitNumber() bool
	UnsetUnitNumber()
	WarehouseManufacturingArea() float64
	SetWarehouseManufacturingArea(warehouseManufacturingArea float64)
	HasWarehouseManufacturingArea() bool
	UnsetWarehouseManufacturingArea()
	WarehouseSector() WarehouseSector
	SetWarehouseSector(warehouseSector WarehouseSector)
	HasWarehouseSector() bool
	UnsetWarehouseSector()
	WarehouseHasPrimaryCabin() bool
	SetWarehouseHasPrimaryCabin(warehouseHasPrimaryCabin bool)
	HasWarehouseHasPrimaryCabin() bool
	UnsetWarehouseHasPrimaryCabin()
	WarehouseCabinKva() string
	SetWarehouseCabinKva(warehouseCabinKva string)
	HasWarehouseCabinKva() bool
	UnsetWarehouseCabinKva()
	WarehouseGroundFloor() int
	SetWarehouseGroundFloor(warehouseGroundFloor int)
	HasWarehouseGroundFloor() bool
	UnsetWarehouseGroundFloor()
	WarehouseFloorResistance() float64
	SetWarehouseFloorResistance(warehouseFloorResistance float64)
	HasWarehouseFloorResistance() bool
	UnsetWarehouseFloorResistance()
	WarehouseZoning() string
	SetWarehouseZoning(warehouseZoning string)
	HasWarehouseZoning() bool
	UnsetWarehouseZoning()
	WarehouseHasOfficeArea() bool
	SetWarehouseHasOfficeArea(warehouseHasOfficeArea bool)
	HasWarehouseHasOfficeArea() bool
	UnsetWarehouseHasOfficeArea()
	WarehouseOfficeArea() float64
	SetWarehouseOfficeArea(warehouseOfficeArea float64)
	HasWarehouseOfficeArea() bool
	UnsetWarehouseOfficeArea()
	StoreHasMezzanine() bool
	SetStoreHasMezzanine(storeHasMezzanine bool)
	HasStoreHasMezzanine() bool
	UnsetStoreHasMezzanine()
	StoreMezzanineArea() float64
	SetStoreMezzanineArea(storeMezzanineArea float64)
	HasStoreMezzanineArea() bool
	UnsetStoreMezzanineArea()
	WarehouseAdditionalFloors() []WarehouseAdditionalFloorInterface
	SetWarehouseAdditionalFloors(warehouseAdditionalFloors []WarehouseAdditionalFloorInterface)
}

type ListingInterface interface {
	ListingVersionInterface

	IdentityID() int64
	SetIdentityID(identityID int64)
	CreatedAt() time.Time
	SetCreatedAt(createdAt time.Time)
	HasCreatedAt() bool
	UnsetCreatedAt()
	PriceUpdatedAt() time.Time
	SetPriceUpdatedAt(priceUpdatedAt time.Time)
	HasPriceUpdatedAt() bool
	UnsetPriceUpdatedAt()
	UUID() string
	SetUUID(uuid string)
	ActiveVersionID() int64
	SetActiveVersionID(versionID int64)
	ActiveVersion() ListingVersionInterface
	SetActiveVersion(version ListingVersionInterface)
	DraftVersion() (ListingVersionInterface, bool)
	SetDraftVersion(version ListingVersionInterface)
	ClearDraftVersion()
	Versions() []ListingVersionInterface
	SetVersions(versions []ListingVersionInterface)
	AddVersion(version ListingVersionInterface)
}
