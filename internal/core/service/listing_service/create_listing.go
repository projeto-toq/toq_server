package listingservices

import (
	"context"
	"database/sql"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	propertycoverageservice "github.com/projeto-toq/toq_server/internal/core/service/property_coverage_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

// CreateListingInput carries the data required to create a new listing.
type CreateListingInput struct {
	ZipCode      string
	Number       string
	City         string
	State        string
	Street       string
	Neighborhood *string
	Complement   *string
	PropertyType globalmodel.PropertyType
	Complex      *string
	UnitTower    *string
	UnitFloor    *int16
	UnitNumber   *string
	LandBlock    *string
	LandLot      *string
}

func (ls *listingService) CreateListing(ctx context.Context, input CreateListingInput) (listing listingmodel.ListingInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return listing, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := ls.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.create.tx_start_error", "err", txErr)
		return listing, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.create.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	listing, err = ls.createListing(ctx, tx, input)
	if err != nil {
		return
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.create.tx_commit_error", "err", cmErr)
		return listing, utils.InternalError("")
	}

	return
}

func (ls *listingService) createListing(ctx context.Context, tx *sql.Tx, input CreateListingInput) (listing listingmodel.ListingInterface, err error) {
	logger := utils.LoggerFromContext(ctx)

	// Get the user doing the request
	userID, uidErr := ls.gsi.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return nil, uidErr
	}

	// 1. Normalize ZipCode
	zipCode := strings.TrimSpace(input.ZipCode)
	normalizedZip, normErr := validators.NormalizeCEP(zipCode)
	if normErr != nil {
		return nil, utils.ValidationError("zipCode", "Zip code must contain exactly 8 digits without separators.")
	}
	zipCode = normalizedZip
	number := strings.TrimSpace(input.Number)

	// 2. Fetch Address from CEP Adapter (Source of Truth)
	cepAddress, err := ls.gsi.GetCEP(ctx, zipCode)
	if err != nil {
		return nil, err
	}

	// 3. Get Complex by Address (Coverage & Rules)
	complex, err := ls.propertyCoverage.GetComplexByAddress(ctx, propertycoverageservice.GetComplexByAddressInput{
		ZipCode: zipCode,
		Number:  number,
	})
	if err != nil {
		return nil, err
	}

	// 4. Validate Property Type Allowance
	allowed, err := ls.isPropertyTypeAllowed(ctx, complex.PropertyTypes(), input.PropertyType)
	if err != nil {
		return nil, err
	}
	if !allowed {
		logger.Warn("listing.create.not_allowed_property_type", "zip", zipCode, "number", number, "property_type", input.PropertyType)
		return nil, utils.BadRequest("Property type not allowed for this area/complex")
	}

	// 5. Validate Complex Specific Rules (Towers, Floors)
	// Only if it is a managed complex (not standalone)
	if complex.Kind() != propertycoveragemodel.CoverageKindStandalone {
		if err := ls.validateComplexRules(ctx, complex, input); err != nil {
			return nil, err
		}
	}

	// 6. Check Duplicity
	criteria := ls.buildDuplicityCriteria(zipCode, number, input)
	exists, err := ls.listingRepository.CheckDuplicity(ctx, tx, criteria)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.create.check_duplicity_error", "err", err, "criteria", criteria)
		return nil, utils.InternalError("")
	}
	if exists {
		return nil, utils.ConflictError("Listing already exists for this address/unit")
	}

	// 7. Create Listing Entity
	code, err := ls.listingRepository.GetListingCode(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	listing = listingmodel.NewListing()
	listing.SetUUID(uuid.NewString())
	listing.SetUserID(userID)
	listing.SetCode(code)
	listing.SetVersion(1)
	listing.SetStatus(listingmodel.StatusDraft)
	listing.SetZipCode(zipCode)
	listing.SetNumber(number)
	listing.SetListingType(input.PropertyType)

	// Overwrite address with CEP data
	listing.SetStreet(strings.TrimSpace(cepAddress.GetStreet()))
	listing.SetCity(strings.TrimSpace(cepAddress.GetCity()))
	listing.SetState(strings.TrimSpace(cepAddress.GetState()))

	// Use input for Neighborhood/Complement if provided, else CEP
	neighborhood := strings.TrimSpace(cepAddress.GetNeighborhood())
	if input.Neighborhood != nil {
		neighborhood = strings.TrimSpace(*input.Neighborhood)
	}
	listing.SetNeighborhood(neighborhood)

	complement := strings.TrimSpace(cepAddress.GetComplement())
	if input.Complement != nil {
		complement = strings.TrimSpace(*input.Complement)
	}
	listing.SetComplement(complement)

	// Set Unit/Land details from criteria
	if criteria.UnitTower != nil {
		listing.SetUnitTower(*criteria.UnitTower)
	}
	if criteria.UnitFloor != nil {
		listing.SetUnitFloor(*criteria.UnitFloor)
	}
	if criteria.UnitNumber != nil {
		listing.SetUnitNumber(*criteria.UnitNumber)
	}
	if criteria.LandBlock != nil {
		listing.SetLandBlock(*criteria.LandBlock)
	}
	if criteria.LandLot != nil {
		listing.SetLandLot(*criteria.LandLot)
	}

	listing.SetDeleted(false)

	// 8. Persist Identity
	if identityErr := ls.listingRepository.CreateListingIdentity(ctx, tx, listing); identityErr != nil {
		utils.SetSpanError(ctx, identityErr)
		logger.Error("listing.create.create_identity_error", "err", identityErr)
		return nil, utils.InternalError("")
	}

	// 9. Create Version
	activeVersion := listing.ActiveVersion()
	activeVersion.SetListingIdentityID(listing.IdentityID())
	activeVersion.SetListingUUID(listing.UUID())

	// Determine Complex Name: System > User Input
	complexName := strings.TrimSpace(complex.Name())
	if complexName == "" && input.Complex != nil {
		complexName = strings.TrimSpace(*input.Complex)
	}
	if complexName != "" {
		activeVersion.SetComplex(complexName)
	}

	if versionErr := ls.listingRepository.CreateListingVersion(ctx, tx, activeVersion); versionErr != nil {
		utils.SetSpanError(ctx, versionErr)
		logger.Error("listing.create.create_version_error", "err", versionErr, "identity_id", listing.IdentityID())
		return nil, utils.InternalError("")
	}

	listing.SetActiveVersionID(activeVersion.ID())

	if setErr := ls.listingRepository.SetListingActiveVersion(ctx, tx, listing.IdentityID(), listing.ActiveVersionID()); setErr != nil {
		utils.SetSpanError(ctx, setErr)
		logger.Error("listing.create.set_active_error", "err", setErr, "identity_id", listing.IdentityID(), "version_id", listing.ActiveVersionID())
		return nil, utils.InternalError("")
	}

	err = ls.gsi.CreateAudit(ctx, tx, globalmodel.TableListings, "Anúncio criado")
	if err != nil {
		return
	}

	return
}

// validateComplexRules checks if the input unit details match the complex configuration.
// If the complex configuration is incomplete (e.g., no towers defined), it accepts the user input.
func (ls *listingService) validateComplexRules(ctx context.Context, complex propertycoveragemodel.ManagedComplexInterface, input CreateListingInput) error {
	// Only validate for types that require unit details
	if input.PropertyType != globalmodel.Apartment && input.PropertyType != globalmodel.CommercialFloor {
		return nil
	}

	towers := complex.Towers()

	// Fallback: If the complex has no towers defined in the system, we cannot validate.
	// In this case, we trust the user input and return success.
	if len(towers) == 0 {
		return nil
	}

	// If the complex HAS towers defined, the user MUST provide a valid one.
	if input.UnitTower == nil {
		return utils.ValidationError("unit_tower", "Tower is required for this complex")
	}
	inputTower := strings.TrimSpace(*input.UnitTower)

	var foundTower propertycoveragemodel.VerticalComplexTowerInterface
	for _, t := range towers {
		if strings.EqualFold(t.Tower(), inputTower) {
			foundTower = t
			break
		}
	}

	if foundTower == nil {
		// Strict validation: Tower must exist if towers are mapped.
		return utils.ValidationError("unit_tower", "Tower not found in this complex")
	}

	// Validate Floor if provided and if the tower has a limit
	if input.UnitFloor != nil {
		floor := int(*input.UnitFloor)
		// Only validate max floor if the system has this data
		if foundTower.Floors() != nil && *foundTower.Floors() > 0 && floor > *foundTower.Floors() {
			return utils.ValidationError("unit_floor", "Floor exceeds the tower limit")
		}
	}

	return nil
}

func (ls *listingService) isPropertyTypeAllowed(ctx context.Context, allowedTypes globalmodel.PropertyType, propertyType globalmodel.PropertyType) (allow bool, err error) {
	requested := ls.DecodePropertyTypes(ctx, propertyType)
	if len(requested) != 1 {
		logger := utils.LoggerFromContext(ctx)
		logger.Warn("listing.create.invalid_property_type_format")
		return false, utils.BadRequest("propertyType must be a single type")
	}
	alloweds := ls.DecodePropertyTypes(ctx, allowedTypes)
	if slices.Contains(alloweds, requested[0]) {
		return true, nil
	}

	return false, nil
}

// buildDuplicityCriteria constructs the criteria ensuring only relevant fields are set based on PropertyType.
// It strictly follows the rules defined in duplicity_criteria.md.
func (ls *listingService) buildDuplicityCriteria(zipCode, number string, input CreateListingInput) listingmodel.DuplicityCriteria {
	criteria := listingmodel.DuplicityCriteria{
		ZipCode: zipCode,
		Number:  number,
	}

	// Helper to treat empty strings as nil
	toPtr := func(s *string) *string {
		if s == nil {
			return nil
		}
		trimmed := strings.TrimSpace(*s)
		if trimmed == "" {
			return nil
		}
		return &trimmed
	}

	switch input.PropertyType {
	case globalmodel.Apartment, globalmodel.Suite:
		// Checks: UnitTower, UnitFloor, UnitNumber
		criteria.UnitTower = toPtr(input.UnitTower)
		if input.UnitFloor != nil {
			s := strconv.Itoa(int(*input.UnitFloor))
			criteria.UnitFloor = &s
		}
		criteria.UnitNumber = toPtr(input.UnitNumber)

	case globalmodel.CommercialStore:
		// Checks: UnitNumber only
		criteria.UnitNumber = toPtr(input.UnitNumber)

	case globalmodel.CommercialFloor:
		// Checks: UnitTower, UnitFloor (Ignores UnitNumber)
		criteria.UnitTower = toPtr(input.UnitTower)
		if input.UnitFloor != nil {
			s := strconv.Itoa(int(*input.UnitFloor))
			criteria.UnitFloor = &s
		}

	case globalmodel.ResidencialLand:
		// Checks: LandBlock, LandLot (Ignores UnitNumber)
		criteria.LandBlock = toPtr(input.LandBlock)
		criteria.LandLot = toPtr(input.LandLot)

	case globalmodel.House, globalmodel.OffPlanHouse, globalmodel.CommercialLand, globalmodel.Building, globalmodel.Warehouse:
		// Checks: ZipCode + Number ONLY.
		// No additional fields are checked for these types.

	default:
		// Fallback: If a new type is added but not mapped, we default to strict Zip+Number check
		// to avoid accidental duplicates, or we could log a warning.
		// For now, no extra criteria added.
	}

	return criteria
}
