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
