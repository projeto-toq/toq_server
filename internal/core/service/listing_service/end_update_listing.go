package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const fallbackListingTimezone = "America/Sao_Paulo"

var stateTimezoneLookup = map[string]string{
	"AC": "America/Rio_Branco",
	"AL": "America/Maceio",
	"AM": "America/Manaus",
	"AP": "America/Belem",
	"BA": "America/Bahia",
	"CE": "America/Fortaleza",
	"DF": "America/Sao_Paulo",
	"ES": "America/Sao_Paulo",
	"GO": "America/Sao_Paulo",
	"MA": "America/Fortaleza",
	"MG": "America/Sao_Paulo",
	"MS": "America/Campo_Grande",
	"MT": "America/Cuiaba",
	"PA": "America/Belem",
	"PB": "America/Fortaleza",
	"PE": "America/Recife",
	"PI": "America/Fortaleza",
	"PR": "America/Sao_Paulo",
	"RJ": "America/Sao_Paulo",
	"RN": "America/Fortaleza",
	"RO": "America/Porto_Velho",
	"RR": "America/Boa_Vista",
	"RS": "America/Sao_Paulo",
	"SC": "America/Sao_Paulo",
	"SE": "America/Maceio",
	"SP": "America/Sao_Paulo",
	"TO": "America/Araguaina",
}

func (ls *listingService) EndUpdateListing(ctx context.Context, input EndUpdateListingInput) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	listingUUID := strings.TrimSpace(input.ListingUUID)
	if listingUUID == "" {
		return utils.ValidationError("listingUuid", "listingUuid is required")
	}
	if _, parseErr := uuid.Parse(listingUUID); parseErr != nil {
		return utils.ValidationError("listingUuid", "listingUuid must be a valid UUID")
	}

	var listingVersionID int64

	tx, err := ls.gsi.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.end_update.tx_start_error", "err", err, "listing_uuid", listingUUID)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.end_update.tx_rollback_error", "err", rbErr, "listing_uuid", listingUUID, "listing_version_id", listingVersionID)
			}
		}
	}()

	identity, identityErr := ls.listingRepository.GetListingIdentityByUUID(ctx, tx, listingUUID)
	if identityErr != nil {
		if errors.Is(identityErr, sql.ErrNoRows) {
			return utils.NotFoundError("listing")
		}
		utils.SetSpanError(ctx, identityErr)
		logger.Error("listing.end_update.identity_fetch_error", "err", identityErr, "listing_uuid", listingUUID)
		return utils.InternalError("")
	}

	if identity.Deleted {
		return utils.BadRequest("Listing is not available")
	}

	listingVersionID = identity.ActiveVersionID.Int64
	if !identity.ActiveVersionID.Valid || listingVersionID <= 0 {
		logger.Warn("listing.end_update.missing_active_version", "listing_uuid", listingUUID, "listing_identity_id", identity.ID)
		return utils.ConflictError("Listing does not have an active version to end update")
	}

	snapshot, err := ls.listingRepository.GetListingForEndUpdate(ctx, tx, listingVersionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.end_update.fetch_error", "err", err, "listing_uuid", listingUUID, "listing_version_id", listingVersionID)
		return utils.InternalError("")
	}

	if snapshot.Status != listingmodel.StatusDraft {
		return utils.ConflictError("Listing must be in draft status to end update")
	}

	userID, uidErr := ls.gsi.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return uidErr
	}

	if snapshot.UserID != userID {
		return utils.AuthorizationError("Only listing owner can end update")
	}

	if verr := ls.validateListingBeforeEndUpdate(ctx, tx, snapshot); verr != nil {
		return verr
	}

	updateErr := ls.listingRepository.UpdateListingStatus(ctx, tx, listingVersionID, listingmodel.StatusPendingAvailability, listingmodel.StatusDraft)
	if updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return utils.ConflictError("Listing status changed while finishing update")
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("listing.end_update.update_status_error", "err", updateErr, "listing_uuid", listingUUID, "listing_version_id", listingVersionID)
		return utils.InternalError("")
	}

	timezone := resolveListingTimezone(snapshot)
	version := int64(snapshot.Version)
	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		userID,
		auditmodel.AuditTarget{Type: auditmodel.TargetListingIdentity, ID: identity.ID, Version: &version},
		auditmodel.OperationStatusChange,
		map[string]any{
			"listing_identity_id": identity.ID,
			"listing_version_id":  listingVersionID,
			"version":             snapshot.Version,
			"status_from":         listingmodel.StatusDraft.String(),
			"status_to":           listingmodel.StatusPendingAvailability.String(),
			"actor_role":          string(permissionmodel.RoleSlugOwner),
			"timezone":            timezone,
		},
	)

	if auditErr := ls.auditService.RecordChange(ctx, tx, auditRecord); auditErr != nil {
		utils.SetSpanError(ctx, auditErr)
		logger.Error("listing.end_update.audit_error", "err", auditErr, "listing_uuid", listingUUID, "listing_identity_id", identity.ID, "listing_version_id", listingVersionID)
		return auditErr
	}
	agendaInput := scheduleservices.CreateDefaultAgendaInput{
		ListingIdentityID: identity.ID,
		OwnerID:           userID,
		Timezone:          timezone,
		ActorID:           userID,
	}
	if _, agendaErr := ls.scheduleService.CreateDefaultAgendaWithTx(ctx, tx, agendaInput); agendaErr != nil {
		utils.SetSpanError(ctx, agendaErr)
		logger.Error("listing.end_update.create_default_agenda_error", "err", agendaErr, "listing_uuid", listingUUID, "listing_identity_id", identity.ID)
		return agendaErr
	}

	if err = ls.gsi.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.end_update.tx_commit_error", "err", err, "listing_uuid", listingUUID, "listing_version_id", listingVersionID)
		return utils.InternalError("")
	}

	logger.Info("listing.end_update.completed", "listing_uuid", listingUUID, "listing_version_id", listingVersionID, "new_status", listingmodel.StatusPendingPhotoScheduling.String())

	return nil
}

// validateListingBeforeEndUpdate validates all required fields before ending listing update.
// This function orchestrates validation in layers: basic fields, transaction-specific rules,
// property type conditionals, and detailed property-specific validations.
//
// Validation follows the documented business rules in procedimento_de_criação_de_novo_anuncio.md
// section 4.5 (Regras de Validação do Promote).
//
// Parameters:
//   - ctx: context for logging
//   - tx: database transaction for catalog lookups
//   - data: complete listing snapshot with all fields
//
// Returns error with appropriate HTTP status code:
//   - 400 Bad Request: missing required field or invalid value
//   - 500 Internal Server Error: infrastructure failure during catalog lookups
func (ls *listingService) validateListingBeforeEndUpdate(ctx context.Context, tx *sql.Tx, data listingrepository.ListingEndUpdateData) error {
	logger := utils.LoggerFromContext(ctx)

	// ========== LAYER 1: Basic Required Fields (all property types) ==========
	if err := ls.validateBasicFields(data); err != nil {
		return err
	}

	// ========== LAYER 2: Transaction-Specific Validations ==========
	if err := ls.validateTransactionRules(ctx, tx, data, logger); err != nil {
		return err
	}

	// ========== LAYER 3: Property Type Conditional Validations ==========
	if err := ls.validatePropertyTypeConditionals(ctx, data); err != nil {
		return err
	}

	// ========== LAYER 4: Tenant-Specific Validations ==========
	if err := ls.validateTenantFields(ctx, tx, data, logger); err != nil {
		return err
	}

	// ========== LAYER 5: Property-Specific Detailed Validations ==========
	if err := ls.validatePropertySpecificFields(ctx, data); err != nil {
		return err
	}

	return nil
}

// validateBasicFields validates mandatory fields required for all listing types.
// These fields must be present regardless of property type or transaction type.
func (ls *listingService) validateBasicFields(data listingrepository.ListingEndUpdateData) error {
	if data.Code == 0 {
		return utils.BadRequest("Listing code is required")
	}
	if data.Version == 0 {
		return utils.BadRequest("Listing version is required")
	}
	if strings.TrimSpace(data.ZipCode) == "" {
		return utils.BadRequest("Zip code is required")
	}
	if !data.Street.Valid || strings.TrimSpace(data.Street.String) == "" {
		return utils.BadRequest("Street is required")
	}
	if !data.Number.Valid || strings.TrimSpace(data.Number.String) == "" {
		return utils.BadRequest("Number is required")
	}
	if !data.City.Valid || strings.TrimSpace(data.City.String) == "" {
		return utils.BadRequest("City is required")
	}
	if !data.State.Valid || strings.TrimSpace(data.State.String) == "" {
		return utils.BadRequest("State is required")
	}
	if !data.Title.Valid || strings.TrimSpace(data.Title.String) == "" {
		return utils.BadRequest("Title is required")
	}
	if data.ListingType == 0 {
		return utils.BadRequest("Property type is required")
	}
	if !data.Owner.Valid {
		return utils.BadRequest("Property owner is required")
	}
	if !data.Buildable.Valid {
		return utils.BadRequest("Buildable size is required")
	}
	if !data.Delivered.Valid {
		return utils.BadRequest("Delivered status is required")
	}
	if !data.WhoLives.Valid {
		return utils.BadRequest("Who lives information is required")
	}
	if !data.Description.Valid || strings.TrimSpace(data.Description.String) == "" {
		return utils.BadRequest("Description is required")
	}
	if !data.Transaction.Valid {
		return utils.BadRequest("Transaction type is required")
	}
	if !data.Visit.Valid {
		return utils.BadRequest("Visit type is required")
	}
	if !data.Accompanying.Valid {
		return utils.BadRequest("Accompanying type is required")
	}

	// Validate IPTU (property tax): AT LEAST ONE must be provided, but NOT BOTH
	// Business rule: Frontend chooses to send either annual OR monthly, never both
	hasAnnualTax := data.AnnualTax.Valid
	hasMonthlyTax := data.MonthlyTax.Valid

	if !hasAnnualTax && !hasMonthlyTax {
		return utils.BadRequest("IPTU is required: provide either annual_tax or monthly_tax")
	}

	if hasAnnualTax && hasMonthlyTax {
		return utils.BadRequest("IPTU conflict: cannot provide both annual_tax and monthly_tax simultaneously")
	}

	// Validate Laudêmio (ground rent): BOTH are optional, but NOT BOTH simultaneously
	// Business rule: Laudêmio may not exist for all properties
	hasAnnualGroundRent := data.AnnualGroundRent.Valid
	hasMonthlyGroundRent := data.MonthlyGroundRent.Valid

	if hasAnnualGroundRent && hasMonthlyGroundRent {
		return utils.BadRequest("Laudêmio conflict: cannot provide both annual_ground_rent and monthly_ground_rent simultaneously")
	}

	// Note: Features validation moved to validatePropertySpecificFields() (LAYER 5)
	// Features are MANDATORY only for residential types: Apartment, House, OffPlanHouse
	// Other property types (commercial, land, warehouse) can proceed without features

	return nil
}

// validateTransactionRules validates fields required based on transaction type (sale, rent, both).
//
// For Sale transactions:
//   - saleNet (sell value)
//   - exchange flag and related fields if enabled
//   - financing flag and blockers if disabled
//
// For Rent transactions:
//   - rentNet (rent value)
//   - guarantees (at least one)
func (ls *listingService) validateTransactionRules(ctx context.Context, tx *sql.Tx, data listingrepository.ListingEndUpdateData, logger *slog.Logger) error {
	txnValue := uint8(data.Transaction.Int16)
	txnCatalog, err := ls.listingRepository.GetCatalogValueByNumeric(ctx, tx, listingmodel.CatalogCategoryTransactionType, txnValue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.BadRequest("Transaction type is invalid")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.end_update.transaction_catalog_error", "err", err, "listing_id", data.ListingID, "transaction_id", txnValue)
		return utils.InternalError("")
	}

	slug := strings.ToLower(strings.TrimSpace(txnCatalog.Slug()))
	needsSaleValidation := slug == "sale" || slug == "both"
	needsRentValidation := slug == "rent" || slug == "both"

	if needsSaleValidation {
		if !data.SaleNet.Valid {
			return utils.BadRequest("Sale net value is required")
		}
		if !data.Exchange.Valid {
			return utils.BadRequest("Exchange flag is required")
		}
		if data.Exchange.Valid && data.Exchange.Int16 == 1 {
			if !data.ExchangePercentual.Valid {
				return utils.BadRequest("Exchange percentual is required when exchange is enabled")
			}
			if data.ExchangePlacesCount == 0 {
				return utils.BadRequest("Exchange places are required when exchange is enabled")
			}
		}
		if !data.Financing.Valid {
			return utils.BadRequest("Financing flag is required")
		}
		if data.Financing.Int16 == 0 && data.FinancingBlockersCount == 0 {
			return utils.BadRequest("Financing blockers are required when financing is disabled")
		}
	}

	if needsRentValidation {
		if !data.RentNet.Valid {
			return utils.BadRequest("Rent net value is required")
		}
		if data.GuaranteesCount == 0 {
			return utils.BadRequest("Guarantees are required for rent transactions")
		}
	}

	return nil
}

// validatePropertyTypeConditionals validates fields required based on broad property categories.
//
// Categories:
//   - Condominium-based: Apartment, CommercialFloor (require condominium value)
//   - Land-based: House, OffPlanHouse, ResidencialLand, CommercialLand (require land size and corner flag)
func (ls *listingService) validatePropertyTypeConditionals(ctx context.Context, data listingrepository.ListingEndUpdateData) error {
	propertyOptions := ls.DecodePropertyTypes(ctx, data.ListingType)
	if len(propertyOptions) == 0 {
		return utils.BadRequest("Property type is invalid")
	}

	needsCondominium := false
	needsLandData := false

	for _, option := range propertyOptions {
		switch option.Code {
		case int64(globalmodel.Apartment), int64(globalmodel.CommercialFloor):
			needsCondominium = true
		case int64(globalmodel.House), int64(globalmodel.OffPlanHouse),
			int64(globalmodel.ResidencialLand), int64(globalmodel.CommercialLand):
			needsLandData = true
		}
	}

	if needsCondominium && !data.Condominium.Valid {
		return utils.BadRequest("Condominium value is required for the selected property type")
	}

	if needsLandData {
		if !data.LandSize.Valid {
			return utils.BadRequest("Land size is required for the selected property type")
		}
		if !data.Corner.Valid {
			return utils.BadRequest("Corner information is required for the selected property type")
		}
	}

	return nil
}

// validateTenantFields validates tenant-specific fields when whoLives = "tenant".
// Requires tenant name, phone, and email for proper contact management.
func (ls *listingService) validateTenantFields(ctx context.Context, tx *sql.Tx, data listingrepository.ListingEndUpdateData, logger *slog.Logger) error {
	if !data.WhoLives.Valid {
		return nil
	}

	whoLivesValue := uint8(data.WhoLives.Int16)
	whoLivesCatalog, catalogErr := ls.listingRepository.GetCatalogValueByNumeric(ctx, tx, listingmodel.CatalogCategoryWhoLives, whoLivesValue)
	if catalogErr != nil {
		if errors.Is(catalogErr, sql.ErrNoRows) {
			return utils.BadRequest("Who lives value is invalid")
		}
		utils.SetSpanError(ctx, catalogErr)
		logger.Error("listing.end_update.wholives_catalog_error", "err", catalogErr, "listing_id", data.ListingID, "who_lives_id", whoLivesValue)
		return utils.InternalError("")
	}

	if strings.ToLower(strings.TrimSpace(whoLivesCatalog.Slug())) == "tenant" {
		if !data.TenantName.Valid || strings.TrimSpace(data.TenantName.String) == "" {
			return utils.BadRequest("Tenant name is required when tenant lives in the property")
		}
		if !data.TenantPhone.Valid || strings.TrimSpace(data.TenantPhone.String) == "" {
			return utils.BadRequest("Tenant phone is required when tenant lives in the property")
		}
		if !data.TenantEmail.Valid || strings.TrimSpace(data.TenantEmail.String) == "" {
			return utils.BadRequest("Tenant email is required when tenant lives in the property")
		}
	}

	return nil
}

// validatePropertySpecificFields validates detailed fields specific to each property type.
// Each property type has its own validation rules based on business requirements.
//
// Property types validated:
//   - Apartment (1): features (MANDATORY) + unit fields
//   - House (16): features (MANDATORY)
//   - OffPlanHouse (32): features (MANDATORY) + completion forecast
//   - ResidencialLand, CommercialLand (64, 128): land block, lot, terrain type, KMZ
//   - CommercialFloor (4): unit tower, floor, number
//   - CommercialStore (2): unit fields + mezzanine flag and area
//   - Warehouse (512): manufacturing area, sector, cabin, floor specs, zoning, office area
func (ls *listingService) validatePropertySpecificFields(ctx context.Context, data listingrepository.ListingEndUpdateData) error {
	propertyOptions := ls.DecodePropertyTypes(ctx, data.ListingType)

	for _, option := range propertyOptions {
		switch option.Code {
		case int64(globalmodel.Apartment):
			// Apartment requires both features validation (residential) and unit validation
			if err := ls.validateResidentialFeatures(data, option.Code); err != nil {
				return err
			}
			if err := ls.validateUnit(data); err != nil {
				return err
			}

		case int64(globalmodel.House), int64(globalmodel.OffPlanHouse):
			// House types require features validation only
			if err := ls.validateResidentialFeatures(data, option.Code); err != nil {
				return err
			}
			// OffPlanHouse specifically requires completion forecast
			if option.Code == int64(globalmodel.OffPlanHouse) {
				if !data.CompletionForecast.Valid || strings.TrimSpace(data.CompletionForecast.String) == "" {
					return utils.BadRequest("Completion forecast (YYYY-MM) is required for Off Plan House")
				}
			}

		case int64(globalmodel.ResidencialLand), int64(globalmodel.CommercialLand):
			if err := ls.validateLand(data, option.Code); err != nil {
				return err
			}

		case int64(globalmodel.CommercialFloor):
			if err := ls.validateUnit(data); err != nil {
				return err
			}

		case int64(globalmodel.CommercialStore):
			// CommercialStore requires both unit validation and mezzanine validation
			if err := ls.validateUnit(data); err != nil {
				return err
			}
			if err := ls.validateCommercialStore(data); err != nil {
				return err
			}

		case int64(globalmodel.Warehouse):
			if err := ls.validateWarehouse(data); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateLand validates fields specific to land property types.
//
// All land types (ResidencialLand=64, CommercialLand=128) require:
//   - landBlock: block identification
//   - landLot: lot number
//   - landTerrainType: terrain type catalog value
//   - hasKmz: flag indicating KMZ file presence
//
// For CommercialLand specifically:
//   - kmzFile: required when hasKmz is true
func (ls *listingService) validateLand(data listingrepository.ListingEndUpdateData, propertyCode int64) error {
	// All land types require land block
	if !data.LandBlock.Valid || strings.TrimSpace(data.LandBlock.String) == "" {
		return utils.BadRequest("Land block is required for land properties")
	}

	// All land types require lot number
	if !data.LandLot.Valid || strings.TrimSpace(data.LandLot.String) == "" {
		return utils.BadRequest("Land lot is required for land properties")
	}

	// All land types require terrain type
	if !data.LandTerrainType.Valid {
		return utils.BadRequest("Land terrain type is required for land properties")
	}

	// All land types require hasKmz flag
	if !data.HasKmz.Valid {
		return utils.BadRequest("HasKmz flag is required for land properties")
	}

	// CommercialLand (128) requires KMZ file when hasKmz is true
	if propertyCode == int64(globalmodel.CommercialLand) && data.HasKmz.Valid && data.HasKmz.Int16 == 1 {
		if !data.KmzFile.Valid || strings.TrimSpace(data.KmzFile.String) == "" {
			return utils.BadRequest("KMZ file path is required when HasKmz is true for Commercial Land")
		}
	}

	return nil
}

// validateUnit validates fields specific to unit-based property types.
//
// Applies to:
//   - Apartment (1): residential unit in multi-family building
//   - CommercialStore (2): commercial store/shop
//   - CommercialFloor (4): commercial floor/office space
//
// All require:
//   - unitTower: tower/building identification
//   - unitFloor: floor number
//   - unitNumber: unit identification number
func (ls *listingService) validateUnit(data listingrepository.ListingEndUpdateData) error {
	if !data.UnitTower.Valid || strings.TrimSpace(data.UnitTower.String) == "" {
		return utils.BadRequest("Unit tower is required for unit-based properties")
	}
	if !data.UnitFloor.Valid || strings.TrimSpace(data.UnitFloor.String) == "" {
		return utils.BadRequest("Unit floor is required for unit-based properties")
	}
	if !data.UnitNumber.Valid || strings.TrimSpace(data.UnitNumber.String) == "" {
		return utils.BadRequest("Unit number is required for unit-based properties")
	}
	return nil
}

// validateWarehouse validates fields specific to Warehouse/Industrial property type (code: 512).
//
// Warehouses are complex industrial/logistics facilities requiring detailed technical specifications:
//   - warehouseManufacturingArea: production/manufacturing area in m²
//   - warehouseSector: industrial sector catalog value (manufacturing, logistics, etc.)
//   - warehouseHasPrimaryCabin: flag for primary electrical cabin presence
//   - warehouseCabinKva: cabin power in KVA (required when primary cabin exists)
//   - warehouseGroundFloor: ground floor height in meters
//   - warehouseFloorResistance: floor load capacity in kg/m²
//   - warehouseZoning: municipal zoning classification
//   - warehouseHasOfficeArea: flag for office space presence
//   - warehouseOfficeArea: office area in m² (required when office space exists)
func (ls *listingService) validateWarehouse(data listingrepository.ListingEndUpdateData) error {
	if !data.WarehouseManufacturingArea.Valid {
		return utils.BadRequest("Warehouse manufacturing area is required for Warehouse")
	}
	if !data.WarehouseSector.Valid {
		return utils.BadRequest("Warehouse sector is required for Warehouse")
	}
	if !data.WarehouseHasPrimaryCabin.Valid {
		return utils.BadRequest("Warehouse has primary cabin flag is required for Warehouse")
	}
	if data.WarehouseHasPrimaryCabin.Valid && data.WarehouseHasPrimaryCabin.Int16 == 1 {
		if !data.WarehouseCabinKva.Valid || strings.TrimSpace(data.WarehouseCabinKva.String) == "" {
			return utils.BadRequest("Warehouse cabin KVA is required when has primary cabin is true for Warehouse")
		}
	}
	if !data.WarehouseGroundFloor.Valid {
		return utils.BadRequest("Warehouse ground floor height is required for Warehouse")
	}
	if !data.WarehouseFloorResistance.Valid {
		return utils.BadRequest("Warehouse floor resistance is required for Warehouse")
	}
	if !data.WarehouseZoning.Valid || strings.TrimSpace(data.WarehouseZoning.String) == "" {
		return utils.BadRequest("Warehouse zoning is required for Warehouse")
	}
	if !data.WarehouseHasOfficeArea.Valid {
		return utils.BadRequest("Warehouse has office area flag is required for Warehouse")
	}
	if data.WarehouseHasOfficeArea.Valid && data.WarehouseHasOfficeArea.Int16 == 1 {
		if !data.WarehouseOfficeArea.Valid {
			return utils.BadRequest("Warehouse office area is required when has office area is true for Warehouse")
		}
	}
	return nil
}

// validateCommercialStore validates fields specific to CommercialStore property type (code: 2).
//
// Commercial stores may have mezzanine space:
//   - storeHasMezzanine: flag indicating mezzanine presence
//   - storeMezzanineArea: mezzanine area in m² (required when mezzanine exists)
func (ls *listingService) validateCommercialStore(data listingrepository.ListingEndUpdateData) error {
	if !data.StoreHasMezzanine.Valid {
		return utils.BadRequest("Store has mezzanine flag is required for Commercial Store")
	}
	if data.StoreHasMezzanine.Valid && data.StoreHasMezzanine.Int16 == 1 {
		if !data.StoreMezzanineArea.Valid {
			return utils.BadRequest("Store mezzanine area is required when has mezzanine is true for Commercial Store")
		}
	}
	return nil
}

// validateResidentialFeatures validates that residential property types have required features.
//
// Features are MANDATORY for residential properties where amenities significantly impact
// property value and buyer/renter decisions. These property types require detailed
// feature information for accurate listing presentation and search filtering.
//
// Mandatory for:
//   - Apartment (1): Multi-family residential units (bedrooms, bathrooms, garage, etc.)
//   - House (16): Single-family residences (backyard, pool, gourmet area, etc.)
//   - OffPlanHouse (32): Pre-construction houses (planned features and finishes)
//
// Optional for (not validated here):
//   - Commercial properties (stores, offices, warehouses): features less relevant
//   - Land properties: no built amenities to list
//   - Buildings: aggregate property, features handled at unit level
//
// Parameters:
//   - data: Complete listing snapshot with features count
//   - propertyCode: Numeric code identifying the property type being validated
//
// Returns:
//   - error: 400 Bad Request if features are missing for residential types
//
// Business Rule Reference:
//   - docs/procedimento_de_criação_de_novo_anuncio.md - Section 4.5
func (ls *listingService) validateResidentialFeatures(data listingrepository.ListingEndUpdateData, propertyCode int64) error {
	// Residential properties must have at least one feature for listing completeness
	if data.FeaturesCount == 0 {
		// Provide clear error message indicating which property type requires features
		propertyName := getPropertyTypeName(propertyCode)
		return utils.BadRequest("Features are required for " + propertyName + " property type")
	}

	return nil
}

// getPropertyTypeName returns human-readable name for property type code.
// Used for constructing clear validation error messages.
func getPropertyTypeName(code int64) string {
	switch code {
	case int64(globalmodel.Apartment):
		return "Apartment"
	case int64(globalmodel.House):
		return "House"
	case int64(globalmodel.OffPlanHouse):
		return "Off Plan House"
	default:
		return "residential"
	}
}

func resolveListingTimezone(data listingrepository.ListingEndUpdateData) string {
	if data.State.Valid {
		state := strings.ToUpper(strings.TrimSpace(data.State.String))
		if tz, ok := stateTimezoneLookup[state]; ok && tz != "" {
			return tz
		}
	}
	return fallbackListingTimezone
}
