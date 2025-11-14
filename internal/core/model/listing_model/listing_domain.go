package listingmodel

import globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"

type listingVersion struct {
	id                      int64
	listingIdentityID       int64
	listingUUID             string
	user_id                 int64
	code                    uint32
	version                 uint8
	status                  ListingStatus
	zip_code                string
	street                  string
	number                  string
	complement              string
	neighborhood            string
	city                    string
	state                   string
	title                   string
	titleValid              bool
	listingType             globalmodel.PropertyType
	owner                   PropertyOwner
	ownerValid              bool
	features                []FeatureInterface
	landSize                float64
	landSizeValid           bool
	corner                  bool
	cornerValid             bool
	nonBuildable            float64
	nonBuildableValid       bool
	buildable               float64
	buildableValid          bool
	delivered               PropertyDelivered
	deliveredValid          bool
	whoLives                WhoLives
	whoLivesValid           bool
	description             string
	descriptionValid        bool
	transaction             TransactionType
	transactionValid        bool
	sellNet                 float64
	sellNetValid            bool
	rentNet                 float64
	rentNetValid            bool
	condominium             float64
	condominiumValid        bool
	annualTax               float64
	annualTaxValid          bool
	monthlyTax              float64
	monthlyTaxValid         bool
	annualGroundRent        float64
	annualGroundRentValid   bool
	monthlyGroundRent       float64
	monthlyGroundRentValid  bool
	exchange                bool
	exchangeValid           bool
	exchangePercentual      float64
	exchangePercentualValid bool
	exchangePlaces          []ExchangePlaceInterface
	installment             InstallmentPlan
	installmentValid        bool
	financing               bool
	financingValid          bool
	financingBlocker        []FinancingBlockerInterface
	guarantees              []GuaranteeInterface
	visit                   VisitType
	visitValid              bool
	tenantName              string
	tenantNameValid         bool
	tenantEmail             string
	tenantEmailValid        bool
	tenantPhone             string
	tenantPhoneValid        bool
	accompanying            AccompanyingType
	accompanyingValid       bool
	deleted                 bool
	deletedValid            bool
	// New property-specific fields
	completionForecast              string
	completionForecastValid         bool
	landBlock                       string
	landBlockValid                  bool
	landLot                         string
	landLotValid                    bool
	landFront                       float64
	landFrontValid                  bool
	landSide                        float64
	landSideValid                   bool
	landBack                        float64
	landBackValid                   bool
	landTerrainType                 LandTerrainType
	landTerrainTypeValid            bool
	hasKmz                          bool
	hasKmzValid                     bool
	kmzFile                         string
	kmzFileValid                    bool
	buildingFloors                  int
	buildingFloorsValid             bool
	unitTower                       string
	unitTowerValid                  bool
	unitFloor                       string
	unitFloorValid                  bool
	unitNumber                      string
	unitNumberValid                 bool
	warehouseManufacturingArea      float64
	warehouseManufacturingAreaValid bool
	warehouseSector                 WarehouseSector
	warehouseSectorValid            bool
	warehouseHasPrimaryCabin        bool
	warehouseHasPrimaryCabinValid   bool
	warehouseCabinKva               string
	warehouseCabinKvaValid          bool
	warehouseGroundFloor            int
	warehouseGroundFloorValid       bool
	warehouseFloorResistance        float64
	warehouseFloorResistanceValid   bool
	warehouseZoning                 string
	warehouseZoningValid            bool
	warehouseHasOfficeArea          bool
	warehouseHasOfficeAreaValid     bool
	warehouseOfficeArea             float64
	warehouseOfficeAreaValid        bool
	storeHasMezzanine               bool
	storeHasMezzanineValid          bool
	storeMezzanineArea              float64
	storeMezzanineAreaValid         bool
	warehouseAdditionalFloors       []WarehouseAdditionalFloorInterface
}

func (l *listingVersion) ID() int64 {
	return l.id
}

func (l *listingVersion) SetID(id int64) {
	l.id = id
}

func (l *listingVersion) ListingIdentityID() int64 {
	return l.listingIdentityID
}

func (l *listingVersion) SetListingIdentityID(listingIdentityID int64) {
	l.listingIdentityID = listingIdentityID
}

func (l *listingVersion) ListingUUID() string {
	return l.listingUUID
}

func (l *listingVersion) SetListingUUID(listingUUID string) {
	l.listingUUID = listingUUID
}

func (l *listingVersion) UserID() int64 {
	return l.user_id
}

func (l *listingVersion) SetUserID(user_id int64) {
	l.user_id = user_id
}

func (l *listingVersion) Code() uint32 {
	return l.code
}

func (l *listingVersion) SetCode(code uint32) {
	l.code = code
}

func (l *listingVersion) Version() uint8 {
	return l.version
}

func (l *listingVersion) SetVersion(version uint8) {
	l.version = version
}

func (l *listingVersion) Status() ListingStatus {
	return l.status
}

func (l *listingVersion) SetStatus(status ListingStatus) {
	l.status = status
}

func (l *listingVersion) ZipCode() string {
	return l.zip_code
}

func (l *listingVersion) SetZipCode(zip_code string) {
	l.zip_code = zip_code
}

func (l *listingVersion) Street() string {
	return l.street
}

func (l *listingVersion) SetStreet(street string) {
	l.street = street
}

func (l *listingVersion) Number() string {
	return l.number
}

func (l *listingVersion) SetNumber(number string) {
	l.number = number
}

func (l *listingVersion) Complement() string {
	return l.complement
}

func (l *listingVersion) SetComplement(complement string) {
	l.complement = complement
}

func (l *listingVersion) Neighborhood() string {
	return l.neighborhood
}

func (l *listingVersion) SetNeighborhood(neighborhood string) {
	l.neighborhood = neighborhood
}

func (l *listingVersion) City() string {
	return l.city
}

func (l *listingVersion) SetCity(city string) {
	l.city = city
}

func (l *listingVersion) State() string {
	return l.state
}

func (l *listingVersion) SetState(state string) {
	l.state = state
}

func (l *listingVersion) Title() string {
	return l.title
}

func (l *listingVersion) SetTitle(title string) {
	l.title = title
	l.titleValid = true
}

func (l *listingVersion) HasTitle() bool {
	return l.titleValid
}

func (l *listingVersion) UnsetTitle() {
	l.title = ""
	l.titleValid = false
}

func (l *listingVersion) ListingType() globalmodel.PropertyType {
	return l.listingType
}

func (l *listingVersion) SetListingType(listingType globalmodel.PropertyType) {
	l.listingType = listingType
}

func (l *listingVersion) Owner() PropertyOwner {
	return l.owner
}

func (l *listingVersion) SetOwner(owner PropertyOwner) {
	l.owner = owner
	l.ownerValid = true
}

func (l *listingVersion) HasOwner() bool {
	return l.ownerValid
}

func (l *listingVersion) UnsetOwner() {
	l.owner = 0
	l.ownerValid = false
}

func (l *listingVersion) Features() []FeatureInterface {
	return l.features
}

func (l *listingVersion) SetFeatures(features []FeatureInterface) {
	l.features = features
}

func (l *listingVersion) LandSize() float64 {
	return l.landSize
}

func (l *listingVersion) SetLandSize(landsize float64) {
	l.landSize = landsize
	l.landSizeValid = true
}

func (l *listingVersion) HasLandSize() bool {
	return l.landSizeValid
}

func (l *listingVersion) UnsetLandSize() {
	l.landSize = 0
	l.landSizeValid = false
}

func (l *listingVersion) Corner() bool {
	return l.corner
}

func (l *listingVersion) SetCorner(corner bool) {
	l.corner = corner
	l.cornerValid = true
}

func (l *listingVersion) HasCorner() bool {
	return l.cornerValid
}

func (l *listingVersion) UnsetCorner() {
	l.corner = false
	l.cornerValid = false
}

func (l *listingVersion) NonBuildable() float64 {
	return l.nonBuildable
}

func (l *listingVersion) SetNonBuildable(nonBuildable float64) {
	l.nonBuildable = nonBuildable
	l.nonBuildableValid = true
}

func (l *listingVersion) HasNonBuildable() bool {
	return l.nonBuildableValid
}

func (l *listingVersion) UnsetNonBuildable() {
	l.nonBuildable = 0
	l.nonBuildableValid = false
}

func (l *listingVersion) Buildable() float64 {
	return l.buildable
}

func (l *listingVersion) SetBuildable(buildable float64) {
	l.buildable = buildable
	l.buildableValid = true
}

func (l *listingVersion) HasBuildable() bool {
	return l.buildableValid
}

func (l *listingVersion) UnsetBuildable() {
	l.buildable = 0
	l.buildableValid = false
}

func (l *listingVersion) Delivered() PropertyDelivered {
	return l.delivered
}

func (l *listingVersion) SetDelivered(delivered PropertyDelivered) {
	l.delivered = delivered
	l.deliveredValid = true
}

func (l *listingVersion) HasDelivered() bool {
	return l.deliveredValid
}

func (l *listingVersion) UnsetDelivered() {
	l.delivered = 0
	l.deliveredValid = false
}

func (l *listingVersion) WhoLives() WhoLives {
	return l.whoLives
}

func (l *listingVersion) SetWhoLives(whoLives WhoLives) {
	l.whoLives = whoLives
	l.whoLivesValid = true
}

func (l *listingVersion) HasWhoLives() bool {
	return l.whoLivesValid
}

func (l *listingVersion) UnsetWhoLives() {
	l.whoLives = 0
	l.whoLivesValid = false
}

func (l *listingVersion) Description() string {
	return l.description
}

func (l *listingVersion) SetDescription(description string) {
	l.description = description
	l.descriptionValid = true
}

func (l *listingVersion) HasDescription() bool {
	return l.descriptionValid
}

func (l *listingVersion) UnsetDescription() {
	l.description = ""
	l.descriptionValid = false
}

func (l *listingVersion) Transaction() TransactionType {
	return l.transaction
}

func (l *listingVersion) SetTransaction(transaction TransactionType) {
	l.transaction = transaction
	l.transactionValid = true
}

func (l *listingVersion) HasTransaction() bool {
	return l.transactionValid
}

func (l *listingVersion) UnsetTransaction() {
	l.transaction = 0
	l.transactionValid = false
}

func (l *listingVersion) SellNet() float64 {
	return l.sellNet
}

func (l *listingVersion) SetSellNet(sellNet float64) {
	l.sellNet = sellNet
	l.sellNetValid = true
}

func (l *listingVersion) HasSellNet() bool {
	return l.sellNetValid
}

func (l *listingVersion) UnsetSellNet() {
	l.sellNet = 0
	l.sellNetValid = false
}

func (l *listingVersion) RentNet() float64 {
	return l.rentNet
}

func (l *listingVersion) SetRentNet(rentNet float64) {
	l.rentNet = rentNet
	l.rentNetValid = true
}

func (l *listingVersion) HasRentNet() bool {
	return l.rentNetValid
}

func (l *listingVersion) UnsetRentNet() {
	l.rentNet = 0
	l.rentNetValid = false
}

func (l *listingVersion) Condominium() float64 {
	return l.condominium
}

func (l *listingVersion) SetCondominium(condominium float64) {
	l.condominium = condominium
	l.condominiumValid = true
}

func (l *listingVersion) HasCondominium() bool {
	return l.condominiumValid
}

func (l *listingVersion) UnsetCondominium() {
	l.condominium = 0
	l.condominiumValid = false
}

func (l *listingVersion) AnnualTax() float64 {
	return l.annualTax
}

func (l *listingVersion) SetAnnualTax(annualTax float64) {
	l.annualTax = annualTax
	l.annualTaxValid = true
}

func (l *listingVersion) HasAnnualTax() bool {
	return l.annualTaxValid
}

func (l *listingVersion) UnsetAnnualTax() {
	l.annualTax = 0
	l.annualTaxValid = false
}

func (l *listingVersion) MonthlyTax() float64 {
	return l.monthlyTax
}

func (l *listingVersion) SetMonthlyTax(monthlyTax float64) {
	l.monthlyTax = monthlyTax
	l.monthlyTaxValid = true
}

func (l *listingVersion) HasMonthlyTax() bool {
	return l.monthlyTaxValid
}

func (l *listingVersion) UnsetMonthlyTax() {
	l.monthlyTax = 0
	l.monthlyTaxValid = false
}

func (l *listingVersion) AnnualGroundRent() float64 {
	return l.annualGroundRent
}

func (l *listingVersion) SetAnnualGroundRent(annualGroundRent float64) {
	l.annualGroundRent = annualGroundRent
	l.annualGroundRentValid = true
}

func (l *listingVersion) HasAnnualGroundRent() bool {
	return l.annualGroundRentValid
}

func (l *listingVersion) UnsetAnnualGroundRent() {
	l.annualGroundRent = 0
	l.annualGroundRentValid = false
}

func (l *listingVersion) MonthlyGroundRent() float64 {
	return l.monthlyGroundRent
}

func (l *listingVersion) SetMonthlyGroundRent(monthlyGroundRent float64) {
	l.monthlyGroundRent = monthlyGroundRent
	l.monthlyGroundRentValid = true
}

func (l *listingVersion) HasMonthlyGroundRent() bool {
	return l.monthlyGroundRentValid
}

func (l *listingVersion) UnsetMonthlyGroundRent() {
	l.monthlyGroundRent = 0
	l.monthlyGroundRentValid = false
}

func (l *listingVersion) Exchange() bool {
	return l.exchange
}

func (l *listingVersion) SetExchange(exchange bool) {
	l.exchange = exchange
	l.exchangeValid = true
}

func (l *listingVersion) HasExchange() bool {
	return l.exchangeValid
}

func (l *listingVersion) UnsetExchange() {
	l.exchange = false
	l.exchangeValid = false
}

func (l *listingVersion) ExchangePercentual() float64 {
	return l.exchangePercentual
}

func (l *listingVersion) SetExchangePercentual(exchangePercentual float64) {
	l.exchangePercentual = exchangePercentual
	l.exchangePercentualValid = true
}

func (l *listingVersion) HasExchangePercentual() bool {
	return l.exchangePercentualValid
}

func (l *listingVersion) UnsetExchangePercentual() {
	l.exchangePercentual = 0
	l.exchangePercentualValid = false
}

func (l *listingVersion) ExchangePlaces() []ExchangePlaceInterface {
	return l.exchangePlaces
}

func (l *listingVersion) SetExchangePlaces(exchangePlaces []ExchangePlaceInterface) {
	l.exchangePlaces = exchangePlaces
}

func (l *listingVersion) Installment() InstallmentPlan {
	return l.installment
}

func (l *listingVersion) SetInstallment(installment InstallmentPlan) {
	l.installment = installment
	l.installmentValid = true
}

func (l *listingVersion) HasInstallment() bool {
	return l.installmentValid
}

func (l *listingVersion) UnsetInstallment() {
	l.installment = 0
	l.installmentValid = false
}

func (l *listingVersion) Financing() bool {
	return l.financing
}

func (l *listingVersion) SetFinancing(financing bool) {
	l.financing = financing
	l.financingValid = true
}

func (l *listingVersion) HasFinancing() bool {
	return l.financingValid
}

func (l *listingVersion) UnsetFinancing() {
	l.financing = false
	l.financingValid = false
}

func (l *listingVersion) FinancingBlockers() []FinancingBlockerInterface {
	return l.financingBlocker
}

func (l *listingVersion) SetFinancingBlockers(financingBlocker []FinancingBlockerInterface) {
	l.financingBlocker = financingBlocker
}

func (l *listingVersion) Guarantees() []GuaranteeInterface {
	return l.guarantees
}

func (l *listingVersion) SetGuarantees(guarantees []GuaranteeInterface) {
	l.guarantees = guarantees
}

func (l *listingVersion) Visit() VisitType {
	return l.visit
}

func (l *listingVersion) SetVisit(visit VisitType) {
	l.visit = visit
	l.visitValid = true
}

func (l *listingVersion) HasVisit() bool {
	return l.visitValid
}

func (l *listingVersion) UnsetVisit() {
	l.visit = 0
	l.visitValid = false
}

func (l *listingVersion) TenantName() string {
	return l.tenantName
}

func (l *listingVersion) SetTenantName(tenantName string) {
	l.tenantName = tenantName
	l.tenantNameValid = true
}

func (l *listingVersion) HasTenantName() bool {
	return l.tenantNameValid
}

func (l *listingVersion) UnsetTenantName() {
	l.tenantName = ""
	l.tenantNameValid = false
}

func (l *listingVersion) TenantEmail() string {
	return l.tenantEmail
}

func (l *listingVersion) SetTenantEmail(tenantEmail string) {
	l.tenantEmail = tenantEmail
	l.tenantEmailValid = true
}

func (l *listingVersion) HasTenantEmail() bool {
	return l.tenantEmailValid
}

func (l *listingVersion) UnsetTenantEmail() {
	l.tenantEmail = ""
	l.tenantEmailValid = false
}

func (l *listingVersion) TenantPhone() string {
	return l.tenantPhone
}

func (l *listingVersion) SetTenantPhone(tenantPhone string) {
	l.tenantPhone = tenantPhone
	l.tenantPhoneValid = true
}

func (l *listingVersion) HasTenantPhone() bool {
	return l.tenantPhoneValid
}

func (l *listingVersion) UnsetTenantPhone() {
	l.tenantPhone = ""
	l.tenantPhoneValid = false
}

func (l *listingVersion) Accompanying() AccompanyingType {
	return l.accompanying
}

func (l *listingVersion) SetAccompanying(accompanying AccompanyingType) {
	l.accompanying = accompanying
	l.accompanyingValid = true
}

func (l *listingVersion) HasAccompanying() bool {
	return l.accompanyingValid
}

func (l *listingVersion) UnsetAccompanying() {
	l.accompanying = 0
	l.accompanyingValid = false
}

func (l *listingVersion) Deleted() bool {
	return l.deleted
}

func (l *listingVersion) SetDeleted(deleted bool) {
	l.deleted = deleted
	l.deletedValid = true
}

func (l *listingVersion) HasDeleted() bool {
	return l.deletedValid
}

func (l *listingVersion) UnsetDeleted() {
	l.deleted = false
	l.deletedValid = false
}

// CompletionForecast represents completion forecast for properties under construction (month/year).
func (l *listingVersion) CompletionForecast() string {
	return l.completionForecast
}

func (l *listingVersion) SetCompletionForecast(completionForecast string) {
	l.completionForecast = completionForecast
	l.completionForecastValid = true
}

func (l *listingVersion) HasCompletionForecast() bool {
	return l.completionForecastValid
}

func (l *listingVersion) UnsetCompletionForecast() {
	l.completionForecast = ""
	l.completionForecastValid = false
}

// LandBlock represents the block identifier for land properties.
func (l *listingVersion) LandBlock() string {
	return l.landBlock
}

func (l *listingVersion) SetLandBlock(landBlock string) {
	l.landBlock = landBlock
	l.landBlockValid = true
}

func (l *listingVersion) HasLandBlock() bool {
	return l.landBlockValid
}

func (l *listingVersion) UnsetLandBlock() {
	l.landBlock = ""
	l.landBlockValid = false
}

// LandLot represents the lot identifier for commercial/residential land.
func (l *listingVersion) LandLot() string {
	return l.landLot
}

func (l *listingVersion) SetLandLot(landLot string) {
	l.landLot = landLot
	l.landLotValid = true
}

func (l *listingVersion) HasLandLot() bool {
	return l.landLotValid
}

func (l *listingVersion) UnsetLandLot() {
	l.landLot = ""
	l.landLotValid = false
}

// LandFront represents the front dimension in meters for commercial/residential land.
func (l *listingVersion) LandFront() float64 {
	return l.landFront
}

func (l *listingVersion) SetLandFront(landFront float64) {
	l.landFront = landFront
	l.landFrontValid = true
}

func (l *listingVersion) HasLandFront() bool {
	return l.landFrontValid
}

func (l *listingVersion) UnsetLandFront() {
	l.landFront = 0
	l.landFrontValid = false
}

// LandSide represents the side dimension in meters for commercial/residential land.
func (l *listingVersion) LandSide() float64 {
	return l.landSide
}

func (l *listingVersion) SetLandSide(landSide float64) {
	l.landSide = landSide
	l.landSideValid = true
}

func (l *listingVersion) HasLandSide() bool {
	return l.landSideValid
}

func (l *listingVersion) UnsetLandSide() {
	l.landSide = 0
	l.landSideValid = false
}

// LandBack represents the back dimension in meters for commercial/residential land.
func (l *listingVersion) LandBack() float64 {
	return l.landBack
}

func (l *listingVersion) SetLandBack(landBack float64) {
	l.landBack = landBack
	l.landBackValid = true
}

func (l *listingVersion) HasLandBack() bool {
	return l.landBackValid
}

func (l *listingVersion) UnsetLandBack() {
	l.landBack = 0
	l.landBackValid = false
}

// LandTerrainType represents the terrain type (uphill, flat, downhill, etc.).
func (l *listingVersion) LandTerrainType() LandTerrainType {
	return l.landTerrainType
}

func (l *listingVersion) SetLandTerrainType(landTerrainType LandTerrainType) {
	l.landTerrainType = landTerrainType
	l.landTerrainTypeValid = true
}

func (l *listingVersion) HasLandTerrainType() bool {
	return l.landTerrainTypeValid
}

func (l *listingVersion) UnsetLandTerrainType() {
	l.landTerrainType = 0
	l.landTerrainTypeValid = false
}

// HasKmz indicates whether a KMZ file is available for the land.
func (l *listingVersion) HasKmz() bool {
	return l.hasKmz
}

func (l *listingVersion) SetHasKmz(hasKmz bool) {
	l.hasKmz = hasKmz
	l.hasKmzValid = true
}

func (l *listingVersion) HasHasKmz() bool {
	return l.hasKmzValid
}

func (l *listingVersion) UnsetHasKmz() {
	l.hasKmz = false
	l.hasKmzValid = false
}

// KmzFile represents the KMZ file path or URL.
func (l *listingVersion) KmzFile() string {
	return l.kmzFile
}

func (l *listingVersion) SetKmzFile(kmzFile string) {
	l.kmzFile = kmzFile
	l.kmzFileValid = true
}

func (l *listingVersion) HasKmzFile() bool {
	return l.kmzFileValid
}

func (l *listingVersion) UnsetKmzFile() {
	l.kmzFile = ""
	l.kmzFileValid = false
}

// BuildingFloors represents the number of floors in a building.
func (l *listingVersion) BuildingFloors() int {
	return l.buildingFloors
}

func (l *listingVersion) SetBuildingFloors(buildingFloors int) {
	l.buildingFloors = buildingFloors
	l.buildingFloorsValid = true
}

func (l *listingVersion) HasBuildingFloors() bool {
	return l.buildingFloorsValid
}

func (l *listingVersion) UnsetBuildingFloors() {
	l.buildingFloors = 0
	l.buildingFloorsValid = false
}

// UnitTower represents the tower/block identifier for apartment/office/laje.
func (l *listingVersion) UnitTower() string {
	return l.unitTower
}

func (l *listingVersion) SetUnitTower(unitTower string) {
	l.unitTower = unitTower
	l.unitTowerValid = true
}

func (l *listingVersion) HasUnitTower() bool {
	return l.unitTowerValid
}

func (l *listingVersion) UnsetUnitTower() {
	l.unitTower = ""
	l.unitTowerValid = false
}

// UnitFloor represents the floor identifier for apartment/office/laje.
func (l *listingVersion) UnitFloor() string {
	return l.unitFloor
}

func (l *listingVersion) SetUnitFloor(unitFloor string) {
	l.unitFloor = unitFloor
	l.unitFloorValid = true
}

func (l *listingVersion) HasUnitFloor() bool {
	return l.unitFloorValid
}

func (l *listingVersion) UnsetUnitFloor() {
	l.unitFloor = ""
	l.unitFloorValid = false
}

// UnitNumber represents the unit number for apartment/office/laje.
func (l *listingVersion) UnitNumber() string {
	return l.unitNumber
}

func (l *listingVersion) SetUnitNumber(unitNumber string) {
	l.unitNumber = unitNumber
	l.unitNumberValid = true
}

func (l *listingVersion) HasUnitNumber() bool {
	return l.unitNumberValid
}

func (l *listingVersion) UnsetUnitNumber() {
	l.unitNumber = ""
	l.unitNumberValid = false
}

// WarehouseManufacturingArea represents the manufacturing area in m² for warehouse.
func (l *listingVersion) WarehouseManufacturingArea() float64 {
	return l.warehouseManufacturingArea
}

func (l *listingVersion) SetWarehouseManufacturingArea(warehouseManufacturingArea float64) {
	l.warehouseManufacturingArea = warehouseManufacturingArea
	l.warehouseManufacturingAreaValid = true
}

func (l *listingVersion) HasWarehouseManufacturingArea() bool {
	return l.warehouseManufacturingAreaValid
}

func (l *listingVersion) UnsetWarehouseManufacturingArea() {
	l.warehouseManufacturingArea = 0
	l.warehouseManufacturingAreaValid = false
}

// WarehouseSector represents the warehouse sector (manufacturing, industrial, logistics).
func (l *listingVersion) WarehouseSector() WarehouseSector {
	return l.warehouseSector
}

func (l *listingVersion) SetWarehouseSector(warehouseSector WarehouseSector) {
	l.warehouseSector = warehouseSector
	l.warehouseSectorValid = true
}

func (l *listingVersion) HasWarehouseSector() bool {
	return l.warehouseSectorValid
}

func (l *listingVersion) UnsetWarehouseSector() {
	l.warehouseSector = 0
	l.warehouseSectorValid = false
}

// WarehouseHasPrimaryCabin indicates whether warehouse has primary cabin.
func (l *listingVersion) WarehouseHasPrimaryCabin() bool {
	return l.warehouseHasPrimaryCabin
}

func (l *listingVersion) SetWarehouseHasPrimaryCabin(warehouseHasPrimaryCabin bool) {
	l.warehouseHasPrimaryCabin = warehouseHasPrimaryCabin
	l.warehouseHasPrimaryCabinValid = true
}

func (l *listingVersion) HasWarehouseHasPrimaryCabin() bool {
	return l.warehouseHasPrimaryCabinValid
}

func (l *listingVersion) UnsetWarehouseHasPrimaryCabin() {
	l.warehouseHasPrimaryCabin = false
	l.warehouseHasPrimaryCabinValid = false
}

// WarehouseCabinKva represents the primary cabin KVA specification.
func (l *listingVersion) WarehouseCabinKva() string {
	return l.warehouseCabinKva
}

func (l *listingVersion) SetWarehouseCabinKva(warehouseCabinKva string) {
	l.warehouseCabinKva = warehouseCabinKva
	l.warehouseCabinKvaValid = true
}

func (l *listingVersion) HasWarehouseCabinKva() bool {
	return l.warehouseCabinKvaValid
}

func (l *listingVersion) UnsetWarehouseCabinKva() {
	l.warehouseCabinKva = ""
	l.warehouseCabinKvaValid = false
}

// WarehouseGroundFloor represents the ground floor height in meters.
func (l *listingVersion) WarehouseGroundFloor() int {
	return l.warehouseGroundFloor
}

func (l *listingVersion) SetWarehouseGroundFloor(warehouseGroundFloor int) {
	l.warehouseGroundFloor = warehouseGroundFloor
	l.warehouseGroundFloorValid = true
}

func (l *listingVersion) HasWarehouseGroundFloor() bool {
	return l.warehouseGroundFloorValid
}

func (l *listingVersion) UnsetWarehouseGroundFloor() {
	l.warehouseGroundFloor = 0
	l.warehouseGroundFloorValid = false
}

// WarehouseFloorResistance represents the floor resistance in kg/m².
func (l *listingVersion) WarehouseFloorResistance() float64 {
	return l.warehouseFloorResistance
}

func (l *listingVersion) SetWarehouseFloorResistance(warehouseFloorResistance float64) {
	l.warehouseFloorResistance = warehouseFloorResistance
	l.warehouseFloorResistanceValid = true
}

func (l *listingVersion) HasWarehouseFloorResistance() bool {
	return l.warehouseFloorResistanceValid
}

func (l *listingVersion) UnsetWarehouseFloorResistance() {
	l.warehouseFloorResistance = 0
	l.warehouseFloorResistanceValid = false
}

// WarehouseZoning represents the zoning classification.
func (l *listingVersion) WarehouseZoning() string {
	return l.warehouseZoning
}

func (l *listingVersion) SetWarehouseZoning(warehouseZoning string) {
	l.warehouseZoning = warehouseZoning
	l.warehouseZoningValid = true
}

func (l *listingVersion) HasWarehouseZoning() bool {
	return l.warehouseZoningValid
}

func (l *listingVersion) UnsetWarehouseZoning() {
	l.warehouseZoning = ""
	l.warehouseZoningValid = false
}

// WarehouseHasOfficeArea indicates whether warehouse has office area.
func (l *listingVersion) WarehouseHasOfficeArea() bool {
	return l.warehouseHasOfficeArea
}

func (l *listingVersion) SetWarehouseHasOfficeArea(warehouseHasOfficeArea bool) {
	l.warehouseHasOfficeArea = warehouseHasOfficeArea
	l.warehouseHasOfficeAreaValid = true
}

func (l *listingVersion) HasWarehouseHasOfficeArea() bool {
	return l.warehouseHasOfficeAreaValid
}

func (l *listingVersion) UnsetWarehouseHasOfficeArea() {
	l.warehouseHasOfficeArea = false
	l.warehouseHasOfficeAreaValid = false
}

// WarehouseOfficeArea represents the office area in m².
func (l *listingVersion) WarehouseOfficeArea() float64 {
	return l.warehouseOfficeArea
}

func (l *listingVersion) SetWarehouseOfficeArea(warehouseOfficeArea float64) {
	l.warehouseOfficeArea = warehouseOfficeArea
	l.warehouseOfficeAreaValid = true
}

func (l *listingVersion) HasWarehouseOfficeArea() bool {
	return l.warehouseOfficeAreaValid
}

func (l *listingVersion) UnsetWarehouseOfficeArea() {
	l.warehouseOfficeArea = 0
	l.warehouseOfficeAreaValid = false
}

// StoreHasMezzanine indicates whether store has mezzanine.
func (l *listingVersion) StoreHasMezzanine() bool {
	return l.storeHasMezzanine
}

func (l *listingVersion) SetStoreHasMezzanine(storeHasMezzanine bool) {
	l.storeHasMezzanine = storeHasMezzanine
	l.storeHasMezzanineValid = true
}

func (l *listingVersion) HasStoreHasMezzanine() bool {
	return l.storeHasMezzanineValid
}

func (l *listingVersion) UnsetStoreHasMezzanine() {
	l.storeHasMezzanine = false
	l.storeHasMezzanineValid = false
}

// StoreMezzanineArea represents the mezzanine area in m².
func (l *listingVersion) StoreMezzanineArea() float64 {
	return l.storeMezzanineArea
}

func (l *listingVersion) SetStoreMezzanineArea(storeMezzanineArea float64) {
	l.storeMezzanineArea = storeMezzanineArea
	l.storeMezzanineAreaValid = true
}

func (l *listingVersion) HasStoreMezzanineArea() bool {
	return l.storeMezzanineAreaValid
}

func (l *listingVersion) UnsetStoreMezzanineArea() {
	l.storeMezzanineArea = 0
	l.storeMezzanineAreaValid = false
}

// WarehouseAdditionalFloors represents additional floors for warehouse properties.
func (l *listingVersion) WarehouseAdditionalFloors() []WarehouseAdditionalFloorInterface {
	return l.warehouseAdditionalFloors
}

func (l *listingVersion) SetWarehouseAdditionalFloors(warehouseAdditionalFloors []WarehouseAdditionalFloorInterface) {
	l.warehouseAdditionalFloors = warehouseAdditionalFloors
}

func (l *listingVersion) copyFrom(version ListingVersionInterface) {
	if version == nil {
		return
	}

	l.id = version.ID()
	l.listingIdentityID = version.ListingIdentityID()
	l.listingUUID = version.ListingUUID()
	l.user_id = version.UserID()
	l.code = version.Code()
	l.version = version.Version()
	l.status = version.Status()
	l.zip_code = version.ZipCode()
	l.street = version.Street()
	l.number = version.Number()
	l.complement = version.Complement()
	l.neighborhood = version.Neighborhood()
	l.city = version.City()
	l.state = version.State()

	if version.HasTitle() {
		l.title = version.Title()
		l.titleValid = true
	} else {
		l.title = ""
		l.titleValid = false
	}

	l.listingType = version.ListingType()

	if version.HasOwner() {
		l.owner = version.Owner()
		l.ownerValid = true
	} else {
		l.owner = 0
		l.ownerValid = false
	}

	l.features = cloneFeatureInterfaces(version.Features())

	if version.HasLandSize() {
		l.landSize = version.LandSize()
		l.landSizeValid = true
	} else {
		l.landSize = 0
		l.landSizeValid = false
	}

	if version.HasCorner() {
		l.corner = version.Corner()
		l.cornerValid = true
	} else {
		l.corner = false
		l.cornerValid = false
	}

	if version.HasNonBuildable() {
		l.nonBuildable = version.NonBuildable()
		l.nonBuildableValid = true
	} else {
		l.nonBuildable = 0
		l.nonBuildableValid = false
	}

	if version.HasBuildable() {
		l.buildable = version.Buildable()
		l.buildableValid = true
	} else {
		l.buildable = 0
		l.buildableValid = false
	}

	if version.HasDelivered() {
		l.delivered = version.Delivered()
		l.deliveredValid = true
	} else {
		l.delivered = 0
		l.deliveredValid = false
	}

	if version.HasWhoLives() {
		l.whoLives = version.WhoLives()
		l.whoLivesValid = true
	} else {
		l.whoLives = 0
		l.whoLivesValid = false
	}

	if version.HasDescription() {
		l.description = version.Description()
		l.descriptionValid = true
	} else {
		l.description = ""
		l.descriptionValid = false
	}

	if version.HasTransaction() {
		l.transaction = version.Transaction()
		l.transactionValid = true
	} else {
		l.transaction = 0
		l.transactionValid = false
	}

	if version.HasSellNet() {
		l.sellNet = version.SellNet()
		l.sellNetValid = true
	} else {
		l.sellNet = 0
		l.sellNetValid = false
	}

	if version.HasRentNet() {
		l.rentNet = version.RentNet()
		l.rentNetValid = true
	} else {
		l.rentNet = 0
		l.rentNetValid = false
	}

	if version.HasCondominium() {
		l.condominium = version.Condominium()
		l.condominiumValid = true
	} else {
		l.condominium = 0
		l.condominiumValid = false
	}

	if version.HasAnnualTax() {
		l.annualTax = version.AnnualTax()
		l.annualTaxValid = true
	} else {
		l.annualTax = 0
		l.annualTaxValid = false
	}

	if version.HasMonthlyTax() {
		l.monthlyTax = version.MonthlyTax()
		l.monthlyTaxValid = true
	} else {
		l.monthlyTax = 0
		l.monthlyTaxValid = false
	}

	if version.HasAnnualGroundRent() {
		l.annualGroundRent = version.AnnualGroundRent()
		l.annualGroundRentValid = true
	} else {
		l.annualGroundRent = 0
		l.annualGroundRentValid = false
	}

	if version.HasMonthlyGroundRent() {
		l.monthlyGroundRent = version.MonthlyGroundRent()
		l.monthlyGroundRentValid = true
	} else {
		l.monthlyGroundRent = 0
		l.monthlyGroundRentValid = false
	}

	if version.HasExchange() {
		l.exchange = version.Exchange()
		l.exchangeValid = true
	} else {
		l.exchange = false
		l.exchangeValid = false
	}

	if version.HasExchangePercentual() {
		l.exchangePercentual = version.ExchangePercentual()
		l.exchangePercentualValid = true
	} else {
		l.exchangePercentual = 0
		l.exchangePercentualValid = false
	}

	l.exchangePlaces = cloneExchangePlaces(version.ExchangePlaces())

	if version.HasInstallment() {
		l.installment = version.Installment()
		l.installmentValid = true
	} else {
		l.installment = 0
		l.installmentValid = false
	}

	if version.HasFinancing() {
		l.financing = version.Financing()
		l.financingValid = true
	} else {
		l.financing = false
		l.financingValid = false
	}

	l.financingBlocker = cloneFinancingBlockers(version.FinancingBlockers())
	l.guarantees = cloneGuarantees(version.Guarantees())

	if version.HasVisit() {
		l.visit = version.Visit()
		l.visitValid = true
	} else {
		l.visit = 0
		l.visitValid = false
	}

	if version.HasTenantName() {
		l.tenantName = version.TenantName()
		l.tenantNameValid = true
	} else {
		l.tenantName = ""
		l.tenantNameValid = false
	}

	if version.HasTenantEmail() {
		l.tenantEmail = version.TenantEmail()
		l.tenantEmailValid = true
	} else {
		l.tenantEmail = ""
		l.tenantEmailValid = false
	}

	if version.HasTenantPhone() {
		l.tenantPhone = version.TenantPhone()
		l.tenantPhoneValid = true
	} else {
		l.tenantPhone = ""
		l.tenantPhoneValid = false
	}

	if version.HasAccompanying() {
		l.accompanying = version.Accompanying()
		l.accompanyingValid = true
	} else {
		l.accompanying = 0
		l.accompanyingValid = false
	}

	if version.HasDeleted() {
		l.deleted = version.Deleted()
		l.deletedValid = true
	} else {
		l.deleted = false
		l.deletedValid = false
	}

	// New property-specific fields
	if version.HasCompletionForecast() {
		l.completionForecast = version.CompletionForecast()
		l.completionForecastValid = true
	} else {
		l.completionForecast = ""
		l.completionForecastValid = false
	}

	if version.HasLandBlock() {
		l.landBlock = version.LandBlock()
		l.landBlockValid = true
	} else {
		l.landBlock = ""
		l.landBlockValid = false
	}

	if version.HasLandLot() {
		l.landLot = version.LandLot()
		l.landLotValid = true
	} else {
		l.landLot = ""
		l.landLotValid = false
	}

	if version.HasLandFront() {
		l.landFront = version.LandFront()
		l.landFrontValid = true
	} else {
		l.landFront = 0
		l.landFrontValid = false
	}

	if version.HasLandSide() {
		l.landSide = version.LandSide()
		l.landSideValid = true
	} else {
		l.landSide = 0
		l.landSideValid = false
	}

	if version.HasLandBack() {
		l.landBack = version.LandBack()
		l.landBackValid = true
	} else {
		l.landBack = 0
		l.landBackValid = false
	}

	if version.HasLandTerrainType() {
		l.landTerrainType = version.LandTerrainType()
		l.landTerrainTypeValid = true
	} else {
		l.landTerrainType = 0
		l.landTerrainTypeValid = false
	}

	if version.HasHasKmz() {
		l.hasKmz = version.HasKmz()
		l.hasKmzValid = true
	} else {
		l.hasKmz = false
		l.hasKmzValid = false
	}

	if version.HasKmzFile() {
		l.kmzFile = version.KmzFile()
		l.kmzFileValid = true
	} else {
		l.kmzFile = ""
		l.kmzFileValid = false
	}

	if version.HasBuildingFloors() {
		l.buildingFloors = version.BuildingFloors()
		l.buildingFloorsValid = true
	} else {
		l.buildingFloors = 0
		l.buildingFloorsValid = false
	}

	if version.HasUnitTower() {
		l.unitTower = version.UnitTower()
		l.unitTowerValid = true
	} else {
		l.unitTower = ""
		l.unitTowerValid = false
	}

	if version.HasUnitFloor() {
		l.unitFloor = version.UnitFloor()
		l.unitFloorValid = true
	} else {
		l.unitFloor = ""
		l.unitFloorValid = false
	}

	if version.HasUnitNumber() {
		l.unitNumber = version.UnitNumber()
		l.unitNumberValid = true
	} else {
		l.unitNumber = ""
		l.unitNumberValid = false
	}

	if version.HasWarehouseManufacturingArea() {
		l.warehouseManufacturingArea = version.WarehouseManufacturingArea()
		l.warehouseManufacturingAreaValid = true
	} else {
		l.warehouseManufacturingArea = 0
		l.warehouseManufacturingAreaValid = false
	}

	if version.HasWarehouseSector() {
		l.warehouseSector = version.WarehouseSector()
		l.warehouseSectorValid = true
	} else {
		l.warehouseSector = 0
		l.warehouseSectorValid = false
	}

	if version.HasWarehouseHasPrimaryCabin() {
		l.warehouseHasPrimaryCabin = version.WarehouseHasPrimaryCabin()
		l.warehouseHasPrimaryCabinValid = true
	} else {
		l.warehouseHasPrimaryCabin = false
		l.warehouseHasPrimaryCabinValid = false
	}

	if version.HasWarehouseCabinKva() {
		l.warehouseCabinKva = version.WarehouseCabinKva()
		l.warehouseCabinKvaValid = true
	} else {
		l.warehouseCabinKva = ""
		l.warehouseCabinKvaValid = false
	}

	if version.HasWarehouseGroundFloor() {
		l.warehouseGroundFloor = version.WarehouseGroundFloor()
		l.warehouseGroundFloorValid = true
	} else {
		l.warehouseGroundFloor = 0
		l.warehouseGroundFloorValid = false
	}

	if version.HasWarehouseFloorResistance() {
		l.warehouseFloorResistance = version.WarehouseFloorResistance()
		l.warehouseFloorResistanceValid = true
	} else {
		l.warehouseFloorResistance = 0
		l.warehouseFloorResistanceValid = false
	}

	if version.HasWarehouseZoning() {
		l.warehouseZoning = version.WarehouseZoning()
		l.warehouseZoningValid = true
	} else {
		l.warehouseZoning = ""
		l.warehouseZoningValid = false
	}

	if version.HasWarehouseHasOfficeArea() {
		l.warehouseHasOfficeArea = version.WarehouseHasOfficeArea()
		l.warehouseHasOfficeAreaValid = true
	} else {
		l.warehouseHasOfficeArea = false
		l.warehouseHasOfficeAreaValid = false
	}

	if version.HasWarehouseOfficeArea() {
		l.warehouseOfficeArea = version.WarehouseOfficeArea()
		l.warehouseOfficeAreaValid = true
	} else {
		l.warehouseOfficeArea = 0
		l.warehouseOfficeAreaValid = false
	}

	if version.HasStoreHasMezzanine() {
		l.storeHasMezzanine = version.StoreHasMezzanine()
		l.storeHasMezzanineValid = true
	} else {
		l.storeHasMezzanine = false
		l.storeHasMezzanineValid = false
	}

	if version.HasStoreMezzanineArea() {
		l.storeMezzanineArea = version.StoreMezzanineArea()
		l.storeMezzanineAreaValid = true
	} else {
		l.storeMezzanineArea = 0
		l.storeMezzanineAreaValid = false
	}

	l.warehouseAdditionalFloors = cloneWarehouseAdditionalFloors(version.WarehouseAdditionalFloors())
}

type listing struct {
	*listingVersion
	identityID      int64
	uuid            string
	activeVersionID int64
	versions        []ListingVersionInterface
	draftVersion    ListingVersionInterface
}

func (l *listing) IdentityID() int64 {
	if l.identityID != 0 {
		return l.identityID
	}
	if l.listingVersion != nil {
		return l.listingVersion.ListingIdentityID()
	}
	return 0
}

func (l *listing) SetIdentityID(identityID int64) {
	l.identityID = identityID
	l.ensureListingVersion().SetListingIdentityID(identityID)
}

func (l *listing) UUID() string {
	if l.uuid != "" {
		return l.uuid
	}
	if l.listingVersion != nil {
		return l.listingVersion.ListingUUID()
	}
	return ""
}

func (l *listing) SetUUID(uuid string) {
	l.uuid = uuid
	l.ensureListingVersion().SetListingUUID(uuid)
}

func (l *listing) ActiveVersionID() int64 {
	if l.activeVersionID != 0 {
		return l.activeVersionID
	}
	if l.listingVersion != nil {
		return l.listingVersion.ID()
	}
	return 0
}

func (l *listing) SetActiveVersionID(versionID int64) {
	l.activeVersionID = versionID
	l.ensureListingVersion().SetID(versionID)
}

func (l *listing) ActiveVersion() ListingVersionInterface {
	return l.listingVersion
}

func (l *listing) SetActiveVersion(version ListingVersionInterface) {
	adopted := adoptListingVersion(version)
	l.listingVersion = adopted
	l.identityID = adopted.ListingIdentityID()
	l.uuid = adopted.ListingUUID()
	l.activeVersionID = adopted.ID()
}

func (l *listing) DraftVersion() (ListingVersionInterface, bool) {
	if l.draftVersion == nil {
		return nil, false
	}
	return l.draftVersion, true
}

func (l *listing) SetDraftVersion(version ListingVersionInterface) {
	l.draftVersion = adoptListingVersion(version)
}

func (l *listing) ClearDraftVersion() {
	l.draftVersion = nil
}

func (l *listing) Versions() []ListingVersionInterface {
	if len(l.versions) == 0 {
		return nil
	}

	versions := make([]ListingVersionInterface, len(l.versions))
	copy(versions, l.versions)
	return versions
}

func (l *listing) SetVersions(versions []ListingVersionInterface) {
	l.versions = cloneListingVersions(versions)
}

func (l *listing) AddVersion(version ListingVersionInterface) {
	l.versions = append(l.versions, adoptListingVersion(version))
}

func (l *listing) ensureListingVersion() *listingVersion {
	if l.listingVersion == nil {
		l.listingVersion = newListingVersion()
	}
	return l.listingVersion
}

func newListingVersion() *listingVersion {
	return &listingVersion{}
}

func adoptListingVersion(version ListingVersionInterface) *listingVersion {
	if version == nil {
		return newListingVersion()
	}

	if internal, ok := version.(*listingVersion); ok {
		return internal
	}

	clone := newListingVersion()
	clone.copyFrom(version)
	return clone
}

func cloneListingVersions(versions []ListingVersionInterface) []ListingVersionInterface {
	if len(versions) == 0 {
		return nil
	}

	cloned := make([]ListingVersionInterface, 0, len(versions))
	for _, version := range versions {
		cloned = append(cloned, adoptListingVersion(version))
	}
	return cloned
}

func cloneFeatureInterfaces(features []FeatureInterface) []FeatureInterface {
	if len(features) == 0 {
		return nil
	}

	cloned := make([]FeatureInterface, len(features))
	copy(cloned, features)
	return cloned
}

func cloneExchangePlaces(places []ExchangePlaceInterface) []ExchangePlaceInterface {
	if len(places) == 0 {
		return nil
	}

	cloned := make([]ExchangePlaceInterface, len(places))
	copy(cloned, places)
	return cloned
}

func cloneFinancingBlockers(blockers []FinancingBlockerInterface) []FinancingBlockerInterface {
	if len(blockers) == 0 {
		return nil
	}

	cloned := make([]FinancingBlockerInterface, len(blockers))
	copy(cloned, blockers)
	return cloned
}

func cloneGuarantees(guarantees []GuaranteeInterface) []GuaranteeInterface {
	if len(guarantees) == 0 {
		return nil
	}

	cloned := make([]GuaranteeInterface, len(guarantees))
	copy(cloned, guarantees)
	return cloned
}

func cloneWarehouseAdditionalFloors(floors []WarehouseAdditionalFloorInterface) []WarehouseAdditionalFloorInterface {
	if len(floors) == 0 {
		return nil
	}

	cloned := make([]WarehouseAdditionalFloorInterface, len(floors))
	copy(cloned, floors)
	return cloned
}

func NewListing() ListingInterface {
	activeVersion := newListingVersion()
	return &listing{
		listingVersion: activeVersion,
		versions:       []ListingVersionInterface{activeVersion},
	}
}

func NewListingVersion() ListingVersionInterface {
	return newListingVersion()
}
