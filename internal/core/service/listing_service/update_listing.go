package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

func (ls *listingService) UpdateListing(ctx context.Context, input UpdateListingInput) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := ls.gsi.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.update.tx_start_error", "err", err)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.update.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	err = ls.updateListing(ctx, tx, input)
	if err != nil {
		return
	}

	err = ls.gsi.CommitTransaction(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.update.tx_commit_error", "err", err)
		return utils.InternalError("")
	}

	return
}

func (ls *listingService) updateListing(ctx context.Context, tx *sql.Tx, input UpdateListingInput) (err error) {
	logger := utils.LoggerFromContext(ctx)

	// Validate required fields
	if input.ListingIdentityID == 0 {
		return utils.BadRequest("listingIdentityId is required")
	}
	if input.VersionID == 0 {
		return utils.BadRequest("versionId is required")
	}

	// Get user ID early for ownership validation
	userID, uidErr := ls.gsi.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return uidErr
	}

	// Get listing identity to validate ownership BEFORE fetching version
	identity, err := ls.listingRepository.GetListingIdentityByID(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.update.get_identity_error", "err", err, "identity_id", input.ListingIdentityID)
		return utils.InternalError("")
	}

	// Validate ownership using identity
	if identity.UserID != userID {
		logger.Warn("unauthorized_update_attempt",
			"listing_identity_id", input.ListingIdentityID,
			"listing_version_id", input.VersionID,
			"requester_user_id", userID,
			"owner_user_id", identity.UserID)
		return utils.AuthorizationError("not the listing owner")
	}

	// Now get the specific version
	existing, err := ls.listingRepository.GetListingVersionByID(ctx, tx, input.VersionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("listing version")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.update.get_version_error", "err", err, "version_id", input.VersionID)
		return utils.InternalError("")
	}

	// Verify version belongs to the identity
	if existing.IdentityID() != input.ListingIdentityID {
		logger.Warn("version_identity_mismatch",
			"listing_identity_id", input.ListingIdentityID,
			"listing_version_id", input.VersionID,
			"version_actual_identity_id", existing.IdentityID(),
			"requester_user_id", userID)
		return utils.BadRequest("version does not belong to specified listing")
	}

	// Check if the listing is in draft status
	if existing.Status() != listingmodel.StatusDraft {
		return utils.ConflictError("listing cannot be updated unless in draft status")
	}

	//update only the allowed fields respeitando campos opcionais
	if input.Owner.IsPresent() {
		if input.Owner.IsNull() {
			existing.UnsetOwner()
		} else if selection, ok := input.Owner.Value(); ok {
			resolvedID, resolveErr := ls.resolveCatalogValue(ctx, tx, listingmodel.CatalogCategoryPropertyOwner, "owner", selection)
			if resolveErr != nil {
				return resolveErr
			}
			existing.SetOwner(listingmodel.PropertyOwner(resolvedID))
		}
	}

	if input.Features.IsPresent() {
		if input.Features.IsNull() {
			existing.SetFeatures(nil)
		} else if features, ok := input.Features.Value(); ok {
			for _, feature := range features {
				feature.SetListingVersionID(existing.ID())
			}
			existing.SetFeatures(features)
		}
	}

	if input.LandSize.IsPresent() {
		if input.LandSize.IsNull() {
			existing.UnsetLandSize()
		} else if value, ok := input.LandSize.Value(); ok {
			existing.SetLandSize(value)
		}
	}

	if input.Corner.IsPresent() {
		if input.Corner.IsNull() {
			existing.UnsetCorner()
		} else if value, ok := input.Corner.Value(); ok {
			existing.SetCorner(value)
		}
	}

	if input.NonBuildable.IsPresent() {
		if input.NonBuildable.IsNull() {
			existing.UnsetNonBuildable()
		} else if value, ok := input.NonBuildable.Value(); ok {
			existing.SetNonBuildable(value)
		}
	}

	if input.Buildable.IsPresent() {
		if input.Buildable.IsNull() {
			existing.UnsetBuildable()
		} else if value, ok := input.Buildable.Value(); ok {
			existing.SetBuildable(value)
		}
	}

	if input.Delivered.IsPresent() {
		if input.Delivered.IsNull() {
			existing.UnsetDelivered()
		} else if selection, ok := input.Delivered.Value(); ok {
			resolvedID, resolveErr := ls.resolveCatalogValue(ctx, tx, listingmodel.CatalogCategoryPropertyDelivered, "delivered", selection)
			if resolveErr != nil {
				return resolveErr
			}
			existing.SetDelivered(listingmodel.PropertyDelivered(resolvedID))
		}
	}

	if input.WhoLives.IsPresent() {
		if input.WhoLives.IsNull() {
			existing.UnsetWhoLives()
		} else if selection, ok := input.WhoLives.Value(); ok {
			resolvedID, resolveErr := ls.resolveCatalogValue(ctx, tx, listingmodel.CatalogCategoryWhoLives, "whoLives", selection)
			if resolveErr != nil {
				return resolveErr
			}
			existing.SetWhoLives(listingmodel.WhoLives(resolvedID))
		}
	}

	if input.Title.IsPresent() {
		if input.Title.IsNull() {
			existing.UnsetTitle()
		} else if value, ok := input.Title.Value(); ok {
			trimmed := strings.TrimSpace(value)
			if trimmed == "" {
				return utils.ValidationError("title", "Title must not be empty when provided.")
			}
			if len([]rune(trimmed)) > 255 {
				return utils.ValidationError("title", "Title must not exceed 255 characters.")
			}
			existing.SetTitle(trimmed)
		}
	}

	if input.Description.IsPresent() {
		if input.Description.IsNull() {
			existing.UnsetDescription()
		} else if value, ok := input.Description.Value(); ok {
			existing.SetDescription(value)
		}
	}

	if input.Transaction.IsPresent() {
		if input.Transaction.IsNull() {
			existing.UnsetTransaction()
		} else if selection, ok := input.Transaction.Value(); ok {
			resolvedID, resolveErr := ls.resolveCatalogValue(ctx, tx, listingmodel.CatalogCategoryTransactionType, "transaction", selection)
			if resolveErr != nil {
				return resolveErr
			}
			existing.SetTransaction(listingmodel.TransactionType(resolvedID))
		}
	}

	if input.SellNet.IsPresent() {
		if input.SellNet.IsNull() {
			existing.UnsetSellNet()
		} else if value, ok := input.SellNet.Value(); ok {
			existing.SetSellNet(value)
		}
	}

	if input.RentNet.IsPresent() {
		if input.RentNet.IsNull() {
			existing.UnsetRentNet()
		} else if value, ok := input.RentNet.Value(); ok {
			existing.SetRentNet(value)
		}
	}

	if input.Condominium.IsPresent() {
		if input.Condominium.IsNull() {
			existing.UnsetCondominium()
		} else if value, ok := input.Condominium.Value(); ok {
			existing.SetCondominium(value)
		}
	}

	if input.AnnualTax.IsPresent() {
		if input.AnnualTax.IsNull() {
			existing.UnsetAnnualTax()
		} else if value, ok := input.AnnualTax.Value(); ok {
			existing.SetAnnualTax(value)
		}
	}

	if input.MonthlyTax.IsPresent() {
		if input.MonthlyTax.IsNull() {
			existing.UnsetMonthlyTax()
		} else if value, ok := input.MonthlyTax.Value(); ok {
			existing.SetMonthlyTax(value)
		}
	}

	if input.AnnualGroundRent.IsPresent() {
		if input.AnnualGroundRent.IsNull() {
			existing.UnsetAnnualGroundRent()
		} else if value, ok := input.AnnualGroundRent.Value(); ok {
			existing.SetAnnualGroundRent(value)
		}
	}

	if input.MonthlyGroundRent.IsPresent() {
		if input.MonthlyGroundRent.IsNull() {
			existing.UnsetMonthlyGroundRent()
		} else if value, ok := input.MonthlyGroundRent.Value(); ok {
			existing.SetMonthlyGroundRent(value)
		}
	}

	if input.Exchange.IsPresent() {
		if input.Exchange.IsNull() {
			existing.UnsetExchange()
		} else if value, ok := input.Exchange.Value(); ok {
			existing.SetExchange(value)
		}
	}

	if input.ExchangePercentual.IsPresent() {
		if input.ExchangePercentual.IsNull() {
			existing.UnsetExchangePercentual()
		} else if value, ok := input.ExchangePercentual.Value(); ok {
			existing.SetExchangePercentual(value)
		}
	}

	if input.ExchangePlaces.IsPresent() {
		if input.ExchangePlaces.IsNull() {
			existing.SetExchangePlaces(nil)
		} else if places, ok := input.ExchangePlaces.Value(); ok {
			for _, place := range places {
				place.SetListingVersionID(existing.ID())
			}
			existing.SetExchangePlaces(places)
		}
	}

	if input.Installment.IsPresent() {
		if input.Installment.IsNull() {
			existing.UnsetInstallment()
		} else if selection, ok := input.Installment.Value(); ok {
			resolvedID, resolveErr := ls.resolveCatalogValue(ctx, tx, listingmodel.CatalogCategoryInstallmentPlan, "installment", selection)
			if resolveErr != nil {
				return resolveErr
			}
			existing.SetInstallment(listingmodel.InstallmentPlan(resolvedID))
		}
	}

	if input.Financing.IsPresent() {
		if input.Financing.IsNull() {
			existing.UnsetFinancing()
		} else if value, ok := input.Financing.Value(); ok {
			existing.SetFinancing(value)
		}
	}

	if input.FinancingBlockers.IsPresent() {
		if input.FinancingBlockers.IsNull() {
			existing.SetFinancingBlockers(nil)
		} else if selections, ok := input.FinancingBlockers.Value(); ok {
			blockers := make([]listingmodel.FinancingBlockerInterface, 0, len(selections))
			for _, selection := range selections {
				resolvedID, resolveErr := ls.resolveCatalogValue(ctx, tx, listingmodel.CatalogCategoryFinancingBlocker, "financingBlockers", selection)
				if resolveErr != nil {
					return resolveErr
				}
				blocker := listingmodel.NewFinancingBlocker()
				blocker.SetBlocker(listingmodel.FinancingBlocker(resolvedID))
				blocker.SetListingVersionID(existing.ID())
				blockers = append(blockers, blocker)
			}
			existing.SetFinancingBlockers(blockers)
		}
	}

	if input.Guarantees.IsPresent() {
		if input.Guarantees.IsNull() {
			existing.SetGuarantees(nil)
		} else if guaranteesUpdate, ok := input.Guarantees.Value(); ok {
			guarantees := make([]listingmodel.GuaranteeInterface, 0, len(guaranteesUpdate))
			for _, update := range guaranteesUpdate {
				resolvedID, resolveErr := ls.resolveCatalogValue(ctx, tx, listingmodel.CatalogCategoryGuaranteeType, "guarantees", update.Selection)
				if resolveErr != nil {
					return resolveErr
				}
				guarantee := listingmodel.NewGuarantee()
				guarantee.SetPriority(update.Priority)
				guarantee.SetGuarantee(listingmodel.GuaranteeType(resolvedID))
				guarantee.SetListingVersionID(existing.ID())
				guarantees = append(guarantees, guarantee)
			}
			existing.SetGuarantees(guarantees)
		}
	}

	if input.Visit.IsPresent() {
		if input.Visit.IsNull() {
			existing.UnsetVisit()
		} else if selection, ok := input.Visit.Value(); ok {
			resolvedID, resolveErr := ls.resolveCatalogValue(ctx, tx, listingmodel.CatalogCategoryVisitType, "visit", selection)
			if resolveErr != nil {
				return resolveErr
			}
			existing.SetVisit(listingmodel.VisitType(resolvedID))
		}
	}

	if input.TenantName.IsPresent() {
		if input.TenantName.IsNull() {
			existing.UnsetTenantName()
		} else if value, ok := input.TenantName.Value(); ok {
			existing.SetTenantName(value)
		}
	}

	if input.TenantPhone.IsPresent() {
		if input.TenantPhone.IsNull() {
			existing.UnsetTenantPhone()
		} else if value, ok := input.TenantPhone.Value(); ok {
			trimmed := strings.TrimSpace(value)
			if trimmed == "" {
				return utils.ValidationError("tenantPhone", "Tenant phone must not be empty when provided.")
			}
			normalizedPhone, phoneErr := validators.NormalizeToE164(trimmed)
			if phoneErr != nil {
				return utils.ValidationError("tenantPhone", "Tenant phone must be in valid E.164 format (e.g., +5511912345678).")
			}
			existing.SetTenantPhone(normalizedPhone)
		}
	}

	if input.TenantEmail.IsPresent() {
		if input.TenantEmail.IsNull() {
			existing.UnsetTenantEmail()
		} else if value, ok := input.TenantEmail.Value(); ok {
			existing.SetTenantEmail(value)
		}
	}

	if input.Accompanying.IsPresent() {
		if input.Accompanying.IsNull() {
			existing.UnsetAccompanying()
		} else if selection, ok := input.Accompanying.Value(); ok {
			resolvedID, resolveErr := ls.resolveCatalogValue(ctx, tx, listingmodel.CatalogCategoryAccompanyingType, "accompanying", selection)
			if resolveErr != nil {
				return resolveErr
			}
			existing.SetAccompanying(listingmodel.AccompanyingType(resolvedID))
		}
	}

	if input.CompletionForecast.IsPresent() {
		if input.CompletionForecast.IsNull() {
			existing.UnsetCompletionForecast()
		} else if value, ok := input.CompletionForecast.Value(); ok {
			existing.SetCompletionForecast(strings.TrimSpace(value))
		}
	}

	if input.LandBlock.IsPresent() {
		if input.LandBlock.IsNull() {
			existing.UnsetLandBlock()
		} else if value, ok := input.LandBlock.Value(); ok {
			existing.SetLandBlock(strings.TrimSpace(value))
		}
	}

	if input.LandLot.IsPresent() {
		if input.LandLot.IsNull() {
			existing.UnsetLandLot()
		} else if value, ok := input.LandLot.Value(); ok {
			existing.SetLandLot(strings.TrimSpace(value))
		}
	}

	if input.LandFront.IsPresent() {
		if input.LandFront.IsNull() {
			existing.UnsetLandFront()
		} else if value, ok := input.LandFront.Value(); ok {
			existing.SetLandFront(value)
		}
	}

	if input.LandSide.IsPresent() {
		if input.LandSide.IsNull() {
			existing.UnsetLandSide()
		} else if value, ok := input.LandSide.Value(); ok {
			existing.SetLandSide(value)
		}
	}

	if input.LandBack.IsPresent() {
		if input.LandBack.IsNull() {
			existing.UnsetLandBack()
		} else if value, ok := input.LandBack.Value(); ok {
			existing.SetLandBack(value)
		}
	}

	if input.LandTerrainType.IsPresent() {
		if input.LandTerrainType.IsNull() {
			existing.UnsetLandTerrainType()
		} else if selection, ok := input.LandTerrainType.Value(); ok {
			resolvedID, resolveErr := ls.resolveCatalogValue(ctx, tx, listingmodel.CatalogCategoryLandTerrainType, "landTerrainType", selection)
			if resolveErr != nil {
				return resolveErr
			}
			existing.SetLandTerrainType(listingmodel.LandTerrainType(resolvedID))
		}
	}

	if input.HasKmz.IsPresent() {
		if input.HasKmz.IsNull() {
			existing.UnsetHasKmz()
		} else if value, ok := input.HasKmz.Value(); ok {
			existing.SetHasKmz(value)
		}
	}

	if input.KmzFile.IsPresent() {
		if input.KmzFile.IsNull() {
			existing.UnsetKmzFile()
		} else if value, ok := input.KmzFile.Value(); ok {
			existing.SetKmzFile(strings.TrimSpace(value))
		}
	}

	if input.BuildingFloors.IsPresent() {
		if input.BuildingFloors.IsNull() {
			existing.UnsetBuildingFloors()
		} else if value, ok := input.BuildingFloors.Value(); ok {
			existing.SetBuildingFloors(int(value))
		}
	}

	if input.UnitTower.IsPresent() {
		if input.UnitTower.IsNull() {
			existing.UnsetUnitTower()
		} else if value, ok := input.UnitTower.Value(); ok {
			existing.SetUnitTower(strings.TrimSpace(value))
		}
	}

	if input.UnitFloor.IsPresent() {
		if input.UnitFloor.IsNull() {
			existing.UnsetUnitFloor()
		} else if value, ok := input.UnitFloor.Value(); ok {
			existing.SetUnitFloor(strings.TrimSpace(fmt.Sprintf("%d", value)))
		}
	}

	if input.UnitNumber.IsPresent() {
		if input.UnitNumber.IsNull() {
			existing.UnsetUnitNumber()
		} else if value, ok := input.UnitNumber.Value(); ok {
			existing.SetUnitNumber(strings.TrimSpace(value))
		}
	}

	if input.WarehouseManufacturingArea.IsPresent() {
		if input.WarehouseManufacturingArea.IsNull() {
			existing.UnsetWarehouseManufacturingArea()
		} else if value, ok := input.WarehouseManufacturingArea.Value(); ok {
			existing.SetWarehouseManufacturingArea(value)
		}
	}

	if input.WarehouseSector.IsPresent() {
		if input.WarehouseSector.IsNull() {
			existing.UnsetWarehouseSector()
		} else if selection, ok := input.WarehouseSector.Value(); ok {
			resolvedID, resolveErr := ls.resolveCatalogValue(ctx, tx, listingmodel.CatalogCategoryWarehouseSector, "warehouseSector", selection)
			if resolveErr != nil {
				return resolveErr
			}
			existing.SetWarehouseSector(listingmodel.WarehouseSector(resolvedID))
		}
	}

	if input.WarehouseHasPrimaryCabin.IsPresent() {
		if input.WarehouseHasPrimaryCabin.IsNull() {
			existing.UnsetWarehouseHasPrimaryCabin()
		} else if value, ok := input.WarehouseHasPrimaryCabin.Value(); ok {
			existing.SetWarehouseHasPrimaryCabin(value)
		}
	}

	if input.WarehouseCabinKva.IsPresent() {
		if input.WarehouseCabinKva.IsNull() {
			existing.UnsetWarehouseCabinKva()
		} else if value, ok := input.WarehouseCabinKva.Value(); ok {
			existing.SetWarehouseCabinKva(strings.TrimSpace(fmt.Sprintf("%.2f", value)))
		}
	}

	if input.WarehouseGroundFloor.IsPresent() {
		if input.WarehouseGroundFloor.IsNull() {
			existing.UnsetWarehouseGroundFloor()
		} else if value, ok := input.WarehouseGroundFloor.Value(); ok {
			existing.SetWarehouseGroundFloor(int(value))
		}
	}

	if input.WarehouseFloorResistance.IsPresent() {
		if input.WarehouseFloorResistance.IsNull() {
			existing.UnsetWarehouseFloorResistance()
		} else if value, ok := input.WarehouseFloorResistance.Value(); ok {
			existing.SetWarehouseFloorResistance(value)
		}
	}

	if input.WarehouseZoning.IsPresent() {
		if input.WarehouseZoning.IsNull() {
			existing.UnsetWarehouseZoning()
		} else if value, ok := input.WarehouseZoning.Value(); ok {
			existing.SetWarehouseZoning(strings.TrimSpace(value))
		}
	}

	if input.WarehouseHasOfficeArea.IsPresent() {
		if input.WarehouseHasOfficeArea.IsNull() {
			existing.UnsetWarehouseHasOfficeArea()
		} else if value, ok := input.WarehouseHasOfficeArea.Value(); ok {
			existing.SetWarehouseHasOfficeArea(value)
		}
	}

	if input.WarehouseOfficeArea.IsPresent() {
		if input.WarehouseOfficeArea.IsNull() {
			existing.UnsetWarehouseOfficeArea()
		} else if value, ok := input.WarehouseOfficeArea.Value(); ok {
			existing.SetWarehouseOfficeArea(value)
		}
	}

	if input.WarehouseAdditionalFloors.IsPresent() {
		if input.WarehouseAdditionalFloors.IsNull() {
			existing.SetWarehouseAdditionalFloors(nil)
		} else if floors, ok := input.WarehouseAdditionalFloors.Value(); ok {
			for _, floor := range floors {
				floor.SetListingVersionID(existing.ID())
			}
			existing.SetWarehouseAdditionalFloors(floors)
		}
	}

	if input.StoreHasMezzanine.IsPresent() {
		if input.StoreHasMezzanine.IsNull() {
			existing.UnsetStoreHasMezzanine()
		} else if value, ok := input.StoreHasMezzanine.Value(); ok {
			existing.SetStoreHasMezzanine(value)
		}
	}

	if input.StoreMezzanineArea.IsPresent() {
		if input.StoreMezzanineArea.IsNull() {
			existing.UnsetStoreMezzanineArea()
		} else if value, ok := input.StoreMezzanineArea.Value(); ok {
			existing.SetStoreMezzanineArea(value)
		}
	}

	// Update satellite tables
	if input.Features.IsPresent() {
		err = ls.listingRepository.UpdateFeatures(ctx, tx, existing.ID(), existing.Features())
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("listing.update.update_features_error", "err", err, "version_id", existing.ID())
			return utils.InternalError("Failed to update listing features")
		}
	}

	if input.ExchangePlaces.IsPresent() {
		err = ls.listingRepository.UpdateExchangePlaces(ctx, tx, existing.ID(), existing.ExchangePlaces())
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("listing.update.update_exchange_places_error", "err", err, "version_id", existing.ID())
			return utils.InternalError("Failed to update exchange places")
		}
	}

	if input.FinancingBlockers.IsPresent() {
		err = ls.listingRepository.UpdateFinancingBlockers(ctx, tx, existing.ID(), existing.FinancingBlockers())
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("listing.update.update_financing_blockers_error", "err", err, "version_id", existing.ID())
			return utils.InternalError("Failed to update financing blockers")
		}
	}

	if input.Guarantees.IsPresent() {
		err = ls.listingRepository.UpdateGuarantees(ctx, tx, existing.ID(), existing.Guarantees())
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("listing.update.update_guarantees_error", "err", err, "version_id", existing.ID())
			return utils.InternalError("Failed to update guarantees")
		}
	}

	if input.WarehouseAdditionalFloors.IsPresent() {
		err = ls.listingRepository.UpdateWarehouseAdditionalFloors(ctx, tx, existing.ID(), existing.WarehouseAdditionalFloors())
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("listing.update.update_warehouse_additional_floors_error", "err", err, "version_id", existing.ID())
			return utils.InternalError("Failed to update warehouse additional floors")
		}
	}

	//update the listing version
	err = ls.listingRepository.UpdateListingVersion(ctx, tx, existing)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.update.update_version_error", "err", err, "version_id", input.VersionID)
		return utils.InternalError("Failed to update listing")
	}

	err = ls.gsi.CreateAudit(ctx, tx, globalmodel.TableListings, "Anúncio atualizado")
	if err != nil {
		return err
	}

	return
}
