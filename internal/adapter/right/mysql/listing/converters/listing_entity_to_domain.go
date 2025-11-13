package listingconverters

import (
	listingentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/entity"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
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
	if e.Title.Valid {
		listing.SetTitle(e.Title.String)
	} else {
		listing.UnsetTitle()
	}
	listing.SetListingType(globalmodel.PropertyType(e.ListingType))
	if e.Owner.Valid {
		listing.SetOwner(listingmodel.PropertyOwner(uint8(e.Owner.Int16)))
	} else {
		listing.UnsetOwner()
	}
	listing.SetFeatures(e.FeaturesToDomain())
	if e.LandSize.Valid {
		listing.SetLandSize(e.LandSize.Float64)
	} else {
		listing.UnsetLandSize()
	}
	if e.Corner.Valid {
		listing.SetCorner(e.Corner.Int16 == 1)
	} else {
		listing.UnsetCorner()
	}
	if e.NonBuildable.Valid {
		listing.SetNonBuildable(e.NonBuildable.Float64)
	} else {
		listing.UnsetNonBuildable()
	}
	if e.Buildable.Valid {
		listing.SetBuildable(e.Buildable.Float64)
	} else {
		listing.UnsetBuildable()
	}
	if e.Delivered.Valid {
		listing.SetDelivered(listingmodel.PropertyDelivered(uint8(e.Delivered.Int16)))
	} else {
		listing.UnsetDelivered()
	}
	if e.WhoLives.Valid {
		listing.SetWhoLives(listingmodel.WhoLives(uint8(e.WhoLives.Int16)))
	} else {
		listing.UnsetWhoLives()
	}
	if e.Description.Valid {
		listing.SetDescription(e.Description.String)
	} else {
		listing.UnsetDescription()
	}
	if e.Transaction.Valid {
		listing.SetTransaction(listingmodel.TransactionType(uint8(e.Transaction.Int16)))
	} else {
		listing.UnsetTransaction()
	}
	if e.SellNet.Valid {
		listing.SetSellNet(e.SellNet.Float64)
	} else {
		listing.UnsetSellNet()
	}
	if e.RentNet.Valid {
		listing.SetRentNet(e.RentNet.Float64)
	} else {
		listing.UnsetRentNet()
	}
	if e.Condominium.Valid {
		listing.SetCondominium(e.Condominium.Float64)
	} else {
		listing.UnsetCondominium()
	}
	if e.AnnualTax.Valid {
		listing.SetAnnualTax(e.AnnualTax.Float64)
	} else {
		listing.UnsetAnnualTax()
	}
	if e.MonthlyTax.Valid {
		listing.SetMonthlyTax(e.MonthlyTax.Float64)
	} else {
		listing.UnsetMonthlyTax()
	}
	if e.AnnualGroundRent.Valid {
		listing.SetAnnualGroundRent(e.AnnualGroundRent.Float64)
	} else {
		listing.UnsetAnnualGroundRent()
	}
	if e.MonthlyGroundRent.Valid {
		listing.SetMonthlyGroundRent(e.MonthlyGroundRent.Float64)
	} else {
		listing.UnsetMonthlyGroundRent()
	}
	if e.Exchange.Valid {
		listing.SetExchange(e.Exchange.Int16 == 1)
	} else {
		listing.UnsetExchange()
	}
	if e.ExchangePercentual.Valid {
		listing.SetExchangePercentual(e.ExchangePercentual.Float64)
	} else {
		listing.UnsetExchangePercentual()
	}
	listing.SetExchangePlaces(e.ExchangePlacesToDomain())
	if e.Installment.Valid {
		listing.SetInstallment(listingmodel.InstallmentPlan(uint8(e.Installment.Int16)))
	} else {
		listing.UnsetInstallment()
	}
	if e.Financing.Valid {
		listing.SetFinancing(e.Financing.Int16 == 1)
	} else {
		listing.UnsetFinancing()
	}
	listing.SetFinancingBlockers(e.FinancingBlockersToDomain())
	listing.SetGuarantees(e.GuaranteesToDomain())
	if e.Visit.Valid {
		listing.SetVisit(listingmodel.VisitType(uint8(e.Visit.Int16)))
	} else {
		listing.UnsetVisit()
	}
	if e.TenantName.Valid {
		listing.SetTenantName(e.TenantName.String)
	} else {
		listing.UnsetTenantName()
	}
	if e.TenantEmail.Valid {
		listing.SetTenantEmail(e.TenantEmail.String)
	} else {
		listing.UnsetTenantEmail()
	}
	if e.TenantPhone.Valid {
		listing.SetTenantPhone(e.TenantPhone.String)
	} else {
		listing.UnsetTenantPhone()
	}
	if e.Accompanying.Valid {
		listing.SetAccompanying(listingmodel.AccompanyingType(uint8(e.Accompanying.Int16)))
	} else {
		listing.UnsetAccompanying()
	}
	if e.Deleted.Valid {
		listing.SetDeleted(e.Deleted.Int16 == 1)
	} else {
		listing.UnsetDeleted()
	}

	return
}
