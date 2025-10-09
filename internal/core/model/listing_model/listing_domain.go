package listingmodel

import (
	"database/sql"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

type listing struct {
	id                 int64
	user_id            int64
	code               uint32
	version            uint8
	status             ListingStatus
	zip_code           string
	street             string
	number             string
	complement         string
	neighborhood       string
	city               string
	state              string
	listingType        globalmodel.PropertyType
	owner              PropertyOwner
	features           []FeatureInterface
	landSize           float64
	corner             bool
	nonBuildable       float64
	buildable          float64
	delivered          PropertyDelivered
	whoLives           WhoLives
	description        string
	transaction        TransactionType
	sellNet            float64
	rentNet            float64
	condominium        float64
	annualTax          float64
	annualGroundRent   float64
	exchange           bool
	exchangePercentual float64
	exchangePlaces     []ExchangePlaceInterface
	installment        InstallmentPlan
	financing          bool
	financingBlocker   []FinancingBlockerInterface
	guarantees         []GuaranteeInterface
	visit              VisitType
	tenantName         string
	tenantEmail        string
	tenantPhone        string
	accompanying       AccompanyingType
	deleted            bool
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
}

func (l *listing) Corner() bool {
	return l.corner
}

func (l *listing) SetCorner(corner bool) {
	l.corner = corner
}

func (l *listing) NonBuildable() float64 {
	return l.nonBuildable
}

func (l *listing) SetNonBuildable(nonBuildable float64) {
	l.nonBuildable = nonBuildable
}

func (l *listing) Buildable() float64 {
	return l.buildable
}

func (l *listing) SetBuildable(buildable float64) {
	l.buildable = buildable
}

func (l *listing) Delivered() PropertyDelivered {
	return l.delivered
}

func (l *listing) SetDelivered(delivered PropertyDelivered) {
	l.delivered = delivered
}

func (l *listing) WhoLives() WhoLives {
	return l.whoLives
}

func (l *listing) SetWhoLives(whoLives WhoLives) {
	l.whoLives = whoLives
}

func (l *listing) Description() string {
	return l.description
}

func (l *listing) SetDescription(description string) {
	l.description = description
}

func (l *listing) Transaction() TransactionType {
	return l.transaction
}

func (l *listing) SetTransaction(transaction TransactionType) {
	l.transaction = transaction
}

func (l *listing) SellNet() float64 {
	return l.sellNet
}

func (l *listing) SetSellNet(sellNet float64) {
	l.sellNet = sellNet
}

func (l *listing) RentNet() float64 {
	return l.rentNet
}

func (l *listing) SetRentNet(rentNet float64) {
	l.rentNet = rentNet
}

func (l *listing) Condominium() float64 {
	return l.condominium
}

func (l *listing) SetCondominium(condominium float64) {
	l.condominium = condominium
}

func (l *listing) AnnualTax() float64 {
	return l.annualTax
}

func (l *listing) SetAnnualTax(annualTax float64) {
	l.annualTax = annualTax
}

func (l *listing) AnnualGroundRent() float64 {
	return l.annualGroundRent
}

func (l *listing) SetAnnualGroundRent(annualGroundRent float64) {
	l.annualGroundRent = annualGroundRent
}

func (l *listing) Exchange() bool {
	return l.exchange
}

func (l *listing) SetExchange(exchange bool) {
	l.exchange = exchange
}

func (l *listing) ExchangePercentual() float64 {
	return l.exchangePercentual
}

func (l *listing) SetExchangePercentual(exchangePercentual float64) {
	l.exchangePercentual = exchangePercentual
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
}

func (l *listing) Financing() bool {
	return l.financing
}

func (l *listing) SetFinancing(financing bool) {
	l.financing = financing
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
}

func (l *listing) TenantName() string {
	return l.tenantName
}

func (l *listing) SetTenantName(tenantName string) {
	l.tenantName = tenantName
}

func (l *listing) TenantEmail() string {
	return l.tenantEmail
}

func (l *listing) SetTenantEmail(tenantEmail string) {
	l.tenantEmail = tenantEmail
}

func (l *listing) TenantPhone() string {
	return l.tenantPhone
}

func (l *listing) SetTenantPhone(tenantPhone string) {
	l.tenantPhone = tenantPhone
}

func (l *listing) Accompanying() AccompanyingType {
	return l.accompanying
}

func (l *listing) SetAccompanying(accompanying AccompanyingType) {
	l.accompanying = accompanying
}

func (l *listing) Deleted() bool {
	return l.deleted
}

func (l *listing) SetDeleted(deleted bool) {
	l.deleted = deleted
}

func (l *listing) ToSQLNullString(input string) sql.NullString {
	if input == "" {
		return sql.NullString{String: input, Valid: false}
	}
	return sql.NullString{String: input, Valid: true}
}

func (l *listing) ToSQLNullInt(input any) sql.NullInt64 {
	switch val := input.(type) {
	case bool:
		if val {
			return sql.NullInt64{Int64: 1, Valid: false}
		} else {
			return sql.NullInt64{Int64: 0, Valid: true}
		}
	case int64:
		return sql.NullInt64{Int64: val, Valid: true}
	case uint32:
		return sql.NullInt64{Int64: int64(val), Valid: true}
	case uint16:
		return sql.NullInt64{Int64: int64(val), Valid: true}
	case uint8:
		return sql.NullInt64{Int64: int64(val), Valid: true}
	case uint:
		return sql.NullInt64{Int64: int64(val), Valid: true}
	default:
		return sql.NullInt64{Int64: 0, Valid: false}
	}
}

func (l *listing) ToSQLNullFloat64(value float64) sql.NullFloat64 {
	if value == 0 {
		return sql.NullFloat64{Float64: 0, Valid: false}
	}
	return sql.NullFloat64{Float64: value, Valid: true}
}
