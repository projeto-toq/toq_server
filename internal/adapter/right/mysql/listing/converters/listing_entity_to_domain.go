package listingconverters

import (
	listingentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/entity"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

func ListingEntityToDomain(e listingentity.ListingEntity) (listing listingmodel.ListingInterface) {
	listing = listingmodel.NewListing()

	listing.SetID(e.ID)
	listing.SetIdentityID(e.ListingIdentityID)
	listing.SetUUID(e.ListingUUID)
	if e.ActiveVersionID.Valid {
		listing.SetActiveVersionID(e.ActiveVersionID.Int64)
	}
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
	if e.Complex.Valid {
		listing.SetComplex(e.Complex.String)
	} else {
		listing.UnsetComplex()
	}
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

	// New property-specific fields
	if e.CompletionForecast.Valid {
		listing.SetCompletionForecast(e.CompletionForecast.String)
	} else {
		listing.UnsetCompletionForecast()
	}
	if e.LandBlock.Valid {
		listing.SetLandBlock(e.LandBlock.String)
	} else {
		listing.UnsetLandBlock()
	}
	if e.LandLot.Valid {
		listing.SetLandLot(e.LandLot.String)
	} else {
		listing.UnsetLandLot()
	}
	if e.LandFront.Valid {
		listing.SetLandFront(e.LandFront.Float64)
	} else {
		listing.UnsetLandFront()
	}
	if e.LandSide.Valid {
		listing.SetLandSide(e.LandSide.Float64)
	} else {
		listing.UnsetLandSide()
	}
	if e.LandBack.Valid {
		listing.SetLandBack(e.LandBack.Float64)
	} else {
		listing.UnsetLandBack()
	}
	if e.LandTerrainType.Valid {
		listing.SetLandTerrainType(listingmodel.LandTerrainType(uint8(e.LandTerrainType.Int16)))
	} else {
		listing.UnsetLandTerrainType()
	}
	if e.HasKmz.Valid {
		listing.SetHasKmz(e.HasKmz.Int16 == 1)
	} else {
		listing.UnsetHasKmz()
	}
	if e.KmzFile.Valid {
		listing.SetKmzFile(e.KmzFile.String)
	} else {
		listing.UnsetKmzFile()
	}
	if e.BuildingFloors.Valid {
		listing.SetBuildingFloors(int(e.BuildingFloors.Int16))
	} else {
		listing.UnsetBuildingFloors()
	}
	if e.UnitTower.Valid {
		listing.SetUnitTower(e.UnitTower.String)
	} else {
		listing.UnsetUnitTower()
	}
	if e.UnitFloor.Valid {
		listing.SetUnitFloor(e.UnitFloor.String)
	} else {
		listing.UnsetUnitFloor()
	}
	if e.UnitNumber.Valid {
		listing.SetUnitNumber(e.UnitNumber.String)
	} else {
		listing.UnsetUnitNumber()
	}
	if e.WarehouseManufacturingArea.Valid {
		listing.SetWarehouseManufacturingArea(e.WarehouseManufacturingArea.Float64)
	} else {
		listing.UnsetWarehouseManufacturingArea()
	}
	if e.WarehouseSector.Valid {
		listing.SetWarehouseSector(listingmodel.WarehouseSector(uint8(e.WarehouseSector.Int16)))
	} else {
		listing.UnsetWarehouseSector()
	}
	if e.WarehouseHasPrimaryCabin.Valid {
		listing.SetWarehouseHasPrimaryCabin(e.WarehouseHasPrimaryCabin.Int16 == 1)
	} else {
		listing.UnsetWarehouseHasPrimaryCabin()
	}
	if e.WarehouseCabinKva.Valid {
		listing.SetWarehouseCabinKva(e.WarehouseCabinKva.String)
	} else {
		listing.UnsetWarehouseCabinKva()
	}
	if e.WarehouseGroundFloor.Valid {
		listing.SetWarehouseGroundFloor(int(e.WarehouseGroundFloor.Int16))
	} else {
		listing.UnsetWarehouseGroundFloor()
	}
	if e.WarehouseFloorResistance.Valid {
		listing.SetWarehouseFloorResistance(e.WarehouseFloorResistance.Float64)
	} else {
		listing.UnsetWarehouseFloorResistance()
	}
	if e.WarehouseZoning.Valid {
		listing.SetWarehouseZoning(e.WarehouseZoning.String)
	} else {
		listing.UnsetWarehouseZoning()
	}
	if e.WarehouseHasOfficeArea.Valid {
		listing.SetWarehouseHasOfficeArea(e.WarehouseHasOfficeArea.Int16 == 1)
	} else {
		listing.UnsetWarehouseHasOfficeArea()
	}
	if e.WarehouseOfficeArea.Valid {
		listing.SetWarehouseOfficeArea(e.WarehouseOfficeArea.Float64)
	} else {
		listing.UnsetWarehouseOfficeArea()
	}
	if e.StoreHasMezzanine.Valid {
		listing.SetStoreHasMezzanine(e.StoreHasMezzanine.Int16 == 1)
	} else {
		listing.UnsetStoreHasMezzanine()
	}
	if e.StoreMezzanineArea.Valid {
		listing.SetStoreMezzanineArea(e.StoreMezzanineArea.Float64)
	} else {
		listing.UnsetStoreMezzanineArea()
	}
	listing.SetWarehouseAdditionalFloors(e.WarehouseAdditionalFloorsToDomain())

	return
}
