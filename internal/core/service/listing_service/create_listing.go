package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"

	"github.com/google/uuid"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
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

	zipCode := strings.TrimSpace(input.ZipCode)
	normalizedZip, normErr := validators.NormalizeCEP(zipCode)
	if normErr != nil {
		return nil, utils.ValidationError("zipCode", "Zip code must contain exactly 8 digits without separators.")
	}
	zipCode = normalizedZip
	number := strings.TrimSpace(input.Number)
	propertyType := input.PropertyType

	exist := true
	// // Check if user already has an active listing
	// hasActiveListing, checkErr := ls.listingRepository.CheckActiveListingExists(ctx, tx, userID)
	// if checkErr != nil {
	// 	utils.SetSpanError(ctx, checkErr)
	// 	logger.Error("listing.create.check_active_error", "err", checkErr, "user_id", userID)
	// 	return nil, utils.InternalError("")
	// }
	// if hasActiveListing {
	// 	return nil, utils.ConflictError("User already has an active listing")
	// }

	// Check if listing already exists for this address
	_, err = ls.listingRepository.GetListingVersionByAddress(ctx, tx, zipCode, number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exist = false
		} else {
			utils.SetSpanError(ctx, err)
			logger.Error("listing.create.check_address_error", "err", err, "zip", zipCode, "number", number)
			return nil, utils.InternalError("")
		}
	}

	if exist {
		return nil, utils.ConflictError("Listing already exists for this address")
	}

	coverage, err := ls.propertyCoverage.ResolvePropertyTypes(ctx, propertycoveragemodel.ResolvePropertyTypesInput{
		ZipCode: zipCode,
		Number:  number,
	})
	if err != nil {
		return nil, err
	}

	//check if the propertyType is allowed
	allowed, err := ls.isPropertyTypeAllowed(ctx, coverage.PropertyTypes, propertyType)
	if err != nil {
		return
	}
	if !allowed {
		logger := utils.LoggerFromContext(ctx)
		logger.Warn("listing.create.not_allowed_property_type", "zip", zipCode, "number", number, "property_type", propertyType)
		return nil, utils.BadRequest("Property type not allowed for this area")
	}

	//create a new code for the listing in the sequence
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
	listing.SetListingType(propertyType)

	//recover the address from the zipCode and number
	address, err := ls.gsi.GetCEP(ctx, zipCode)
	if err != nil {
		return
	}

	cepStreet := strings.TrimSpace(address.GetStreet())
	cepCity := strings.TrimSpace(address.GetCity())
	cepState := strings.TrimSpace(address.GetState())
	cepNeighborhood := strings.TrimSpace(address.GetNeighborhood())
	cepComplement := strings.TrimSpace(address.GetComplement())

	neighborhood := cepNeighborhood
	if input.Neighborhood != nil {
		neighborhood = strings.TrimSpace(*input.Neighborhood)
	}

	complement := cepComplement
	if input.Complement != nil {
		complement = strings.TrimSpace(*input.Complement)
	}

	listing.SetStreet(cepStreet)
	listing.SetComplement(complement)
	listing.SetNeighborhood(neighborhood)
	listing.SetCity(cepCity)
	listing.SetState(cepState)
	listing.SetDeleted(false)

	if identityErr := ls.listingRepository.CreateListingIdentity(ctx, tx, listing); identityErr != nil {
		utils.SetSpanError(ctx, identityErr)
		logger.Error("listing.create.create_identity_error", "err", identityErr)
		return nil, utils.InternalError("")
	}

	activeVersion := listing.ActiveVersion()
	activeVersion.SetListingIdentityID(listing.IdentityID())
	activeVersion.SetListingUUID(listing.UUID())
	if coverage.ComplexName != "" {
		activeVersion.SetComplex(coverage.ComplexName)
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
