package listingmodel

import globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"

type listing struct {
	id                      int64
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
	annualGroundRent        float64
	annualGroundRentValid   bool
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
}

func (l *listing) ID() int64 {
	return l.id
}

func (l *listing) SetID(id int64) {
	l.id = id
}

func (l *listing) UserID() int64 {
	return l.user_id
}

func (l *listing) SetUserID(user_id int64) {
	l.user_id = user_id
}

func (l *listing) Code() uint32 {
	return l.code
}

func (l *listing) SetCode(code uint32) {
	l.code = code
}

func (l *listing) Version() uint8 {
	return l.version
}

func (l *listing) SetVersion(version uint8) {
	l.version = version
}

func (l *listing) Status() ListingStatus {
	return l.status
}

func (l *listing) SetStatus(status ListingStatus) {
	l.status = status
}

func (l *listing) ZipCode() string {
	return l.zip_code
}

func (l *listing) SetZipCode(zip_code string) {
	l.zip_code = zip_code
}

func (l *listing) Street() string {
	return l.street
}

func (l *listing) SetStreet(street string) {
	l.street = street
}

func (l *listing) Number() string {
	return l.number
}

func (l *listing) SetNumber(number string) {
	l.number = number
}

func (l *listing) Complement() string {
	return l.complement
}

func (l *listing) SetComplement(complement string) {
	l.complement = complement
}

func (l *listing) Neighborhood() string {
	return l.neighborhood
}

func (l *listing) SetNeighborhood(neighborhood string) {
	l.neighborhood = neighborhood
}

func (l *listing) City() string {
	return l.city
}

func (l *listing) SetCity(city string) {
	l.city = city
}

func (l *listing) State() string {
	return l.state
}

func (l *listing) SetState(state string) {
	l.state = state
}

func (l *listing) ListingType() globalmodel.PropertyType {
	return l.listingType
}

func (l *listing) SetListingType(listingType globalmodel.PropertyType) {
	l.listingType = listingType
}

func (l *listing) Owner() PropertyOwner {
	return l.owner
}

func (l *listing) SetOwner(owner PropertyOwner) {
	l.owner = owner
	l.ownerValid = true
}

func (l *listing) HasOwner() bool {
	return l.ownerValid
}

func (l *listing) UnsetOwner() {
	l.owner = 0
	l.ownerValid = false
}

func (l *listing) Features() []FeatureInterface {
	return l.features
}

func (l *listing) SetFeatures(features []FeatureInterface) {
	l.features = features
}

func (l *listing) LandSize() float64 {
	return l.landSize
}

func (l *listing) SetLandSize(landsize float64) {
	l.landSize = landsize
	l.landSizeValid = true
}

func (l *listing) HasLandSize() bool {
	return l.landSizeValid
}

func (l *listing) UnsetLandSize() {
	l.landSize = 0
	l.landSizeValid = false
}

func (l *listing) Corner() bool {
	return l.corner
}

func (l *listing) SetCorner(corner bool) {
	l.corner = corner
	l.cornerValid = true
}

func (l *listing) HasCorner() bool {
	return l.cornerValid
}

func (l *listing) UnsetCorner() {
	l.corner = false
	l.cornerValid = false
}

func (l *listing) NonBuildable() float64 {
	return l.nonBuildable
}

func (l *listing) SetNonBuildable(nonBuildable float64) {
	l.nonBuildable = nonBuildable
	l.nonBuildableValid = true
}

func (l *listing) HasNonBuildable() bool {
	return l.nonBuildableValid
}

func (l *listing) UnsetNonBuildable() {
	l.nonBuildable = 0
	l.nonBuildableValid = false
}

func (l *listing) Buildable() float64 {
	return l.buildable
}

func (l *listing) SetBuildable(buildable float64) {
	l.buildable = buildable
	l.buildableValid = true
}

func (l *listing) HasBuildable() bool {
	return l.buildableValid
}

func (l *listing) UnsetBuildable() {
	l.buildable = 0
	l.buildableValid = false
}

func (l *listing) Delivered() PropertyDelivered {
	return l.delivered
}

func (l *listing) SetDelivered(delivered PropertyDelivered) {
	l.delivered = delivered
	l.deliveredValid = true
}

func (l *listing) HasDelivered() bool {
	return l.deliveredValid
}

func (l *listing) UnsetDelivered() {
	l.delivered = 0
	l.deliveredValid = false
}

func (l *listing) WhoLives() WhoLives {
	return l.whoLives
}

func (l *listing) SetWhoLives(whoLives WhoLives) {
	l.whoLives = whoLives
	l.whoLivesValid = true
}

func (l *listing) HasWhoLives() bool {
	return l.whoLivesValid
}

func (l *listing) UnsetWhoLives() {
	l.whoLives = 0
	l.whoLivesValid = false
}

func (l *listing) Description() string {
	return l.description
}

func (l *listing) SetDescription(description string) {
	l.description = description
	l.descriptionValid = true
}

func (l *listing) HasDescription() bool {
	return l.descriptionValid
}

func (l *listing) UnsetDescription() {
	l.description = ""
	l.descriptionValid = false
}

func (l *listing) Transaction() TransactionType {
	return l.transaction
}

func (l *listing) SetTransaction(transaction TransactionType) {
	l.transaction = transaction
	l.transactionValid = true
}

func (l *listing) HasTransaction() bool {
	return l.transactionValid
}

func (l *listing) UnsetTransaction() {
	l.transaction = 0
	l.transactionValid = false
}

func (l *listing) SellNet() float64 {
	return l.sellNet
}

func (l *listing) SetSellNet(sellNet float64) {
	l.sellNet = sellNet
	l.sellNetValid = true
}

func (l *listing) HasSellNet() bool {
	return l.sellNetValid
}

func (l *listing) UnsetSellNet() {
	l.sellNet = 0
	l.sellNetValid = false
}

func (l *listing) RentNet() float64 {
	return l.rentNet
}

func (l *listing) SetRentNet(rentNet float64) {
	l.rentNet = rentNet
	l.rentNetValid = true
}

func (l *listing) HasRentNet() bool {
	return l.rentNetValid
}

func (l *listing) UnsetRentNet() {
	l.rentNet = 0
	l.rentNetValid = false
}

func (l *listing) Condominium() float64 {
	return l.condominium
}

func (l *listing) SetCondominium(condominium float64) {
	l.condominium = condominium
	l.condominiumValid = true
}

func (l *listing) HasCondominium() bool {
	return l.condominiumValid
}

func (l *listing) UnsetCondominium() {
	l.condominium = 0
	l.condominiumValid = false
}

func (l *listing) AnnualTax() float64 {
	return l.annualTax
}

func (l *listing) SetAnnualTax(annualTax float64) {
	l.annualTax = annualTax
	l.annualTaxValid = true
}

func (l *listing) HasAnnualTax() bool {
	return l.annualTaxValid
}

func (l *listing) UnsetAnnualTax() {
	l.annualTax = 0
	l.annualTaxValid = false
}

func (l *listing) AnnualGroundRent() float64 {
	return l.annualGroundRent
}

func (l *listing) SetAnnualGroundRent(annualGroundRent float64) {
	l.annualGroundRent = annualGroundRent
	l.annualGroundRentValid = true
}

func (l *listing) HasAnnualGroundRent() bool {
	return l.annualGroundRentValid
}

func (l *listing) UnsetAnnualGroundRent() {
	l.annualGroundRent = 0
	l.annualGroundRentValid = false
}

func (l *listing) Exchange() bool {
	return l.exchange
}

func (l *listing) SetExchange(exchange bool) {
	l.exchange = exchange
	l.exchangeValid = true
}

func (l *listing) HasExchange() bool {
	return l.exchangeValid
}

func (l *listing) UnsetExchange() {
	l.exchange = false
	l.exchangeValid = false
}

func (l *listing) ExchangePercentual() float64 {
	return l.exchangePercentual
}

func (l *listing) SetExchangePercentual(exchangePercentual float64) {
	l.exchangePercentual = exchangePercentual
	l.exchangePercentualValid = true
}

func (l *listing) HasExchangePercentual() bool {
	return l.exchangePercentualValid
}

func (l *listing) UnsetExchangePercentual() {
	l.exchangePercentual = 0
	l.exchangePercentualValid = false
}

func (l *listing) ExchangePlaces() []ExchangePlaceInterface {
	return l.exchangePlaces
}

func (l *listing) SetExchangePlaces(exchangePlaces []ExchangePlaceInterface) {
	l.exchangePlaces = exchangePlaces
}

func (l *listing) Installment() InstallmentPlan {
	return l.installment
}

func (l *listing) SetInstallment(installment InstallmentPlan) {
	l.installment = installment
	l.installmentValid = true
}

func (l *listing) HasInstallment() bool {
	return l.installmentValid
}

func (l *listing) UnsetInstallment() {
	l.installment = 0
	l.installmentValid = false
}

func (l *listing) Financing() bool {
	return l.financing
}

func (l *listing) SetFinancing(financing bool) {
	l.financing = financing
	l.financingValid = true
}

func (l *listing) HasFinancing() bool {
	return l.financingValid
}

func (l *listing) UnsetFinancing() {
	l.financing = false
	l.financingValid = false
}

func (l *listing) FinancingBlockers() []FinancingBlockerInterface {
	return l.financingBlocker
}

func (l *listing) SetFinancingBlockers(financingBlocker []FinancingBlockerInterface) {
	l.financingBlocker = financingBlocker
}

func (l *listing) Guarantees() []GuaranteeInterface {
	return l.guarantees
}

func (l *listing) SetGuarantees(guarantees []GuaranteeInterface) {
	l.guarantees = guarantees
}

func (l *listing) Visit() VisitType {
	return l.visit
}

func (l *listing) SetVisit(visit VisitType) {
	l.visit = visit
	l.visitValid = true
}

func (l *listing) HasVisit() bool {
	return l.visitValid
}

func (l *listing) UnsetVisit() {
	l.visit = 0
	l.visitValid = false
}

func (l *listing) TenantName() string {
	return l.tenantName
}

func (l *listing) SetTenantName(tenantName string) {
	l.tenantName = tenantName
	l.tenantNameValid = true
}

func (l *listing) HasTenantName() bool {
	return l.tenantNameValid
}

func (l *listing) UnsetTenantName() {
	l.tenantName = ""
	l.tenantNameValid = false
}

func (l *listing) TenantEmail() string {
	return l.tenantEmail
}

func (l *listing) SetTenantEmail(tenantEmail string) {
	l.tenantEmail = tenantEmail
	l.tenantEmailValid = true
}

func (l *listing) HasTenantEmail() bool {
	return l.tenantEmailValid
}

func (l *listing) UnsetTenantEmail() {
	l.tenantEmail = ""
	l.tenantEmailValid = false
}

func (l *listing) TenantPhone() string {
	return l.tenantPhone
}

func (l *listing) SetTenantPhone(tenantPhone string) {
	l.tenantPhone = tenantPhone
	l.tenantPhoneValid = true
}

func (l *listing) HasTenantPhone() bool {
	return l.tenantPhoneValid
}

func (l *listing) UnsetTenantPhone() {
	l.tenantPhone = ""
	l.tenantPhoneValid = false
}

func (l *listing) Accompanying() AccompanyingType {
	return l.accompanying
}

func (l *listing) SetAccompanying(accompanying AccompanyingType) {
	l.accompanying = accompanying
	l.accompanyingValid = true
}

func (l *listing) HasAccompanying() bool {
	return l.accompanyingValid
}

func (l *listing) UnsetAccompanying() {
	l.accompanying = 0
	l.accompanyingValid = false
}

func (l *listing) Deleted() bool {
	return l.deleted
}

func (l *listing) SetDeleted(deleted bool) {
	l.deleted = deleted
	l.deletedValid = true
}

func (l *listing) HasDeleted() bool {
	return l.deletedValid
}

func (l *listing) UnsetDeleted() {
	l.deleted = false
	l.deletedValid = false
}
