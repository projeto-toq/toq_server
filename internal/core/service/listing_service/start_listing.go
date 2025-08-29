package listingservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"golang.org/x/exp/slog"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
"errors"
)

func (ls *listingService) StartListing(ctx context.Context, zipCode string, number string, propertyType globalmodel.PropertyType) (listing listingmodel.ListingInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := ls.gsi.StartTransaction(ctx)
	if err != nil {
		return
	}

	listing, err = ls.startListing(ctx, tx, zipCode, number, propertyType)
	if err != nil {
		ls.gsi.RollbackTransaction(ctx, tx)
		return
	}

	err = ls.gsi.CommitTransaction(ctx, tx)
	if err != nil {
		ls.gsi.RollbackTransaction(ctx, tx)
		return
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
		err = utils.ErrInternalServer
		return
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
		slog.Error("PropertyType not allowed on this area", "error", "PropertyType not allowed on this area")
		err = utils.ErrInternalServer
		return
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
		slog.Error("Invalid propertyType", "error", "propertyType must be a single type")
		err = utils.ErrInternalServer
		return false, err
	}
	alloweds := ls.DecodePropertyTypes(ctx, allowedTypes)
	for _, allowedType := range alloweds {
		if allowedType == requested[0] {
			return true, nil
		}
	}

	return false, nil
}
