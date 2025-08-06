package listingconverters

import (
	listingentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing/entity"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
)

func ListingEntityToDomain(e listingentity.ListingEntity) (listing listingmodel.ListingInterface) {
	listing = listingmodel.NewListing()

	listing.SetID(e.ID)
	listing.SetUserID(e.UserID)
	listing.SetCode(e.Code)
	listing.SetVersion(e.Version)
	listing.SetStatus(listingmodel.ListingStatus(e.Status))
	listing.SetZipCode(e.ZipCode)
	listing.SetStreet(e.ToString(e.Street))

	listing.SetNumber(e.Number)
	listing.SetComplement(e.ToString(e.Complement))
	listing.SetNeighborhood(e.ToString(e.Neighborhood))
	listing.SetCity(e.ToString(e.City))
	listing.SetState(e.ToString(e.State))
	listing.SetListingType(globalmodel.PropertyType(e.ListingType))
	listing.SetOwner(listingmodel.PropertyOwner(e.ToUint8(e.Owner)))
	listing.SetFeatures(e.FeaturesToDomain())
	listing.SetLandSize(e.ToFloat64(e.LandSize))
	if e.Corner.Valid {
		listing.SetCorner(e.Corner.Int16 == 1)
	} else {
		listing.SetCorner(false)
	}
	listing.SetNonBuildable(e.ToFloat64(e.NonBuildable))
	listing.SetBuildable(e.ToFloat64(e.Buildable))
	listing.SetDelivered(listingmodel.PropertyDelivered(e.ToUint8(e.Delivered)))
	listing.SetWhoLives(listingmodel.WhoLives(e.ToUint8(e.WhoLives)))
	listing.SetDescription(e.ToString(e.Description))
	listing.SetTransaction(listingmodel.TransactionType(e.ToUint8(e.Transaction)))
	listing.SetSellNet(e.ToFloat64(e.SellNet))
	listing.SetRentNet(e.ToFloat64(e.RentNet))
	listing.SetCondominium(e.ToFloat64(e.Condominium))
	listing.SetAnnualTax(e.ToFloat64(e.AnnualTax))
	listing.SetAnnualGroundRent(e.ToFloat64(e.AnnualGroundRent))
	if e.Exchange.Valid {
		listing.SetExchange(e.Exchange.Int16 == 1)
	} else {
		listing.SetExchange(false)
	}
	listing.SetExchangePercentual(e.ToFloat64(e.ExchangePercentual))
	listing.SetExchangePlaces(e.ExchangePlacesToDomain())
	listing.SetInstallment(listingmodel.InstallmentPlan(e.ToUint8(e.Installment)))
	if e.Financing.Valid {
		listing.SetFinancing(e.Financing.Int16 == 1)
	} else {
		listing.SetFinancing(false)
	}
	listing.SetFinancingBlockers(e.FinancingBlockersToDomain())
	listing.SetGuarantees(e.GuaranteesToDomain())
	listing.SetVisit(listingmodel.VisitType(e.ToUint8(e.Visit)))
	listing.SetTenantName(e.ToString(e.TenantName))
	listing.SetTenantEmail(e.ToString(e.TenantEmail))
	listing.SetTenantPhone(e.ToString(e.TenantPhone))
	listing.SetAccompanying(listingmodel.AccompanyingType(e.ToUint8(e.Accompanying)))
	if e.Deleted.Valid {
		listing.SetDeleted(e.Deleted.Int16 == 1)
	} else {
		listing.SetDeleted(false)
	}

	return
}
