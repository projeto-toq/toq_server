package listingservices

import (
	"context"
	"database/sql"
	"errors"

	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ls *listingService) StartListing(ctx context.Context, zipCode string, number string, propertyType globalmodel.PropertyType) (listing listingmodel.ListingInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, txErr := ls.gsi.StartTransaction(ctx)
	if txErr != nil {
		slog.Error("listing.start.tx_start_error", "err", txErr)
		return listing, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("listing.start.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	listing, err = ls.startListing(ctx, tx, zipCode, number, propertyType)
	if err != nil {
		return
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		slog.Error("listing.start.tx_commit_error", "err", cmErr)
		return listing, utils.InternalError("Failed to commit transaction")
	}

	return
}

func (ls *listingService) startListing(ctx context.Context, tx *sql.Tx, zipCode string, number string, propertyType globalmodel.PropertyType) (listing listingmodel.ListingInterface, err error) {

	exist := true
	//check if the zipCode and number there is not already a listing
	_, err = ls.listingRepository.GetListingByZipNumber(ctx, tx, zipCode, number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exist = false
		} else {
			return
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
		slog.Warn("listing.start.not_allowed_property_type", "zip", zipCode, "number", number, "property_type", propertyType)
		return nil, utils.BadRequest("Property type not allowed for this area")
	}

	//create a new code for the listing in the sequence
	code, err := ls.listingRepository.GetListingCode(ctx, tx)
	if err != nil {
		return
	}

	//get the user doing the request
	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	listing = listingmodel.NewListing()
	listing.SetUserID(infos.ID)
	listing.SetCode(code)
	listing.SetVersion(1)
	listing.SetStatus(listingmodel.StatusDraft)
	listing.SetZipCode(zipCode)
	listing.SetNumber(number)
	listing.SetListingType(propertyType)

	//recover the adress from the zipCode and number
	address, err := ls.gsi.GetCEP(ctx, zipCode)
	if err != nil {
		return
	}

	listing.SetStreet(address.GetStreet())
	listing.SetComplement(address.GetComplement())
	listing.SetNeighborhood(address.GetNeighborhood())
	listing.SetCity(address.GetCity())
	listing.SetState(address.GetState())

	//create the listing
	err = ls.listingRepository.CreateListing(ctx, tx, listing)
	if err != nil {
		return
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
		slog.Warn("listing.start.invalid_property_type_format")
		return false, utils.BadRequest("propertyType must be a single type")
	}
	alloweds := ls.DecodePropertyTypes(ctx, allowedTypes)
	for _, allowedType := range alloweds {
		if allowedType == requested[0] {
			return true, nil
		}
	}

	return false, nil
}
