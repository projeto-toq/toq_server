package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

// StartListingInput carries the data required to create a new listing.
type StartListingInput struct {
	ZipCode      string
	Number       string
	City         string
	State        string
	Street       string
	Neighborhood *string
	Complement   *string
	PropertyType globalmodel.PropertyType
}

func (ls *listingService) StartListing(ctx context.Context, input StartListingInput) (listing listingmodel.ListingInterface, err error) {
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
		logger.Error("listing.start.tx_start_error", "err", txErr)
		return listing, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.start.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	listing, err = ls.startListing(ctx, tx, input)
	if err != nil {
		return
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.start.tx_commit_error", "err", cmErr)
		return listing, utils.InternalError("")
	}

	return
}

func (ls *listingService) startListing(ctx context.Context, tx *sql.Tx, input StartListingInput) (listing listingmodel.ListingInterface, err error) {
	zipCode := strings.TrimSpace(input.ZipCode)
	normalizedZip, normErr := validators.NormalizeCEP(zipCode)
	if normErr != nil {
		return nil, utils.ValidationError("zipCode", "Zip code must contain exactly 8 digits without separators.")
	}
	zipCode = normalizedZip
	number := strings.TrimSpace(input.Number)
	propertyType := input.PropertyType

	exist := true
	//check if the zipCode and number there is not already a listing
	_, err = ls.listingRepository.GetListingByZipNumber(ctx, tx, zipCode, number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exist = false
		} else {
			utils.SetSpanError(ctx, err)
			return nil, utils.InternalError("")
		}
	}

	if exist {
		return nil, utils.ConflictError("Listing already exists for this address")
	}

	//get the propertyTypes allowed on the zipCode and number
	allowedTypes, err := ls.csi.GetOptions(ctx, zipCode, number)
	if err != nil {
		return
	}

	//check if the propertyType is allowed
	allowed, err := ls.isPropertyTypeAllowed(ctx, allowedTypes, propertyType)
	if err != nil {
		return
	}
	if !allowed {
		logger := utils.LoggerFromContext(ctx)
		logger.Warn("listing.start.not_allowed_property_type", "zip", zipCode, "number", number, "property_type", propertyType)
		return nil, utils.BadRequest("Property type not allowed for this area")
	}

	//create a new code for the listing in the sequence
	code, err := ls.listingRepository.GetListingCode(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	//get the user doing the request
	userID, uidErr := ls.gsi.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return nil, uidErr
	}

	listing = listingmodel.NewListing()
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

	// if !equalsNormalized(cepStreet, input.Street) || !equalsNormalized(cepCity, input.City) || !equalsNormalized(cepState, input.State) {
	// 	logger := utils.LoggerFromContext(ctx)
	// 	logger.Warn(
	// 		"listing.start.address_mismatch",
	// 		"zip", zipCode,
	// 		"number", number,
	// 		"cep_street", cepStreet,
	// 		"cep_city", cepCity,
	// 		"cep_state", cepState,
	// 		"request_street", strings.TrimSpace(input.Street),
	// 		"request_city", strings.TrimSpace(input.City),
	// 		"request_state", strings.TrimSpace(input.State),
	// 	)
	// 	return nil, utils.BadRequest("Provided address does not match zip code information")
	// }

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

	//create the listing
	err = ls.listingRepository.CreateListing(ctx, tx, listing)
	if err != nil {
		utils.SetSpanError(ctx, err)
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
		logger.Warn("listing.start.invalid_property_type_format")
		return false, utils.BadRequest("propertyType must be a single type")
	}
	alloweds := ls.DecodePropertyTypes(ctx, allowedTypes)
	if slices.Contains(alloweds, requested[0]) {
		return true, nil
	}

	return false, nil
}

// func equalsNormalized(origin, candidate string) bool {
// 	return strings.EqualFold(strings.TrimSpace(origin), strings.TrimSpace(candidate))
// }
