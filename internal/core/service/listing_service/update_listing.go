package listingservices

import (
	"context"
	"database/sql"
	"errors"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
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

	exist := true
	//check if exists the listing
	existing, err := ls.listingRepository.GetListingByID(ctx, tx, input.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exist = false
		} else {
			utils.SetSpanError(ctx, err)
			return utils.InternalError("")
		}
	}

	if !exist {
		return utils.NotFoundError("listing")
	}

	//check if the listing is in draft status
	if existing.Status() != listingmodel.StatusDraft {
		return utils.ConflictError("listing cannot be updated unless in draft status")
	}

	//check if the user is the owner of the listing
	userID, uidErr := ls.gsi.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return uidErr
	}
	if existing.UserID() != userID {
		return utils.AuthorizationError("not the listing owner")
	}

	//update only the allowed fields respeitando campos opcionais
	if input.Owner.IsPresent() {
		if input.Owner.IsNull() {
			existing.SetOwner(0)
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
				feature.SetListingID(existing.ID())
			}
			existing.SetFeatures(features)
		}
	}

	if input.LandSize.IsPresent() {
		if input.LandSize.IsNull() {
			existing.SetLandSize(0)
		} else if value, ok := input.LandSize.Value(); ok {
			existing.SetLandSize(value)
		}
	}

	if input.Corner.IsPresent() {
		if input.Corner.IsNull() {
			existing.SetCorner(false)
		} else if value, ok := input.Corner.Value(); ok {
			existing.SetCorner(value)
		}
	}

	if input.NonBuildable.IsPresent() {
		if input.NonBuildable.IsNull() {
			existing.SetNonBuildable(0)
		} else if value, ok := input.NonBuildable.Value(); ok {
			existing.SetNonBuildable(value)
		}
	}

	if input.Buildable.IsPresent() {
		if input.Buildable.IsNull() {
			existing.SetBuildable(0)
		} else if value, ok := input.Buildable.Value(); ok {
			existing.SetBuildable(value)
		}
	}

	if input.Delivered.IsPresent() {
		if input.Delivered.IsNull() {
			existing.SetDelivered(0)
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
			existing.SetWhoLives(0)
		} else if selection, ok := input.WhoLives.Value(); ok {
			resolvedID, resolveErr := ls.resolveCatalogValue(ctx, tx, listingmodel.CatalogCategoryWhoLives, "whoLives", selection)
			if resolveErr != nil {
				return resolveErr
			}
			existing.SetWhoLives(listingmodel.WhoLives(resolvedID))
		}
	}

	if input.Description.IsPresent() {
		if input.Description.IsNull() {
			existing.SetDescription("")
		} else if value, ok := input.Description.Value(); ok {
			existing.SetDescription(value)
		}
	}

	if input.Transaction.IsPresent() {
		if input.Transaction.IsNull() {
			existing.SetTransaction(0)
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
			existing.SetSellNet(0)
		} else if value, ok := input.SellNet.Value(); ok {
			existing.SetSellNet(value)
		}
	}

	if input.RentNet.IsPresent() {
		if input.RentNet.IsNull() {
			existing.SetRentNet(0)
		} else if value, ok := input.RentNet.Value(); ok {
			existing.SetRentNet(value)
		}
	}

	if input.Condominium.IsPresent() {
		if input.Condominium.IsNull() {
			existing.SetCondominium(0)
		} else if value, ok := input.Condominium.Value(); ok {
			existing.SetCondominium(value)
		}
	}

	if input.AnnualTax.IsPresent() {
		if input.AnnualTax.IsNull() {
			existing.SetAnnualTax(0)
		} else if value, ok := input.AnnualTax.Value(); ok {
			existing.SetAnnualTax(value)
		}
	}

	if input.AnnualGroundRent.IsPresent() {
		if input.AnnualGroundRent.IsNull() {
			existing.SetAnnualGroundRent(0)
		} else if value, ok := input.AnnualGroundRent.Value(); ok {
			existing.SetAnnualGroundRent(value)
		}
	}

	if input.Exchange.IsPresent() {
		if input.Exchange.IsNull() {
			existing.SetExchange(false)
		} else if value, ok := input.Exchange.Value(); ok {
			existing.SetExchange(value)
		}
	}

	if input.ExchangePercentual.IsPresent() {
		if input.ExchangePercentual.IsNull() {
			existing.SetExchangePercentual(0)
		} else if value, ok := input.ExchangePercentual.Value(); ok {
			existing.SetExchangePercentual(value)
		}
	}

	if input.ExchangePlaces.IsPresent() {
		if input.ExchangePlaces.IsNull() {
			existing.SetExchangePlaces(nil)
		} else if places, ok := input.ExchangePlaces.Value(); ok {
			for _, place := range places {
				place.SetListingID(existing.ID())
			}
			existing.SetExchangePlaces(places)
		}
	}

	if input.Installment.IsPresent() {
		if input.Installment.IsNull() {
			existing.SetInstallment(0)
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
			existing.SetFinancing(false)
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
				blocker.SetListingID(existing.ID())
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
				guarantee.SetListingID(existing.ID())
				guarantees = append(guarantees, guarantee)
			}
			existing.SetGuarantees(guarantees)
		}
	}

	if input.Visit.IsPresent() {
		if input.Visit.IsNull() {
			existing.SetVisit(0)
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
			existing.SetTenantName("")
		} else if value, ok := input.TenantName.Value(); ok {
			existing.SetTenantName(value)
		}
	}

	if input.TenantPhone.IsPresent() {
		if input.TenantPhone.IsNull() {
			existing.SetTenantPhone("")
		} else if value, ok := input.TenantPhone.Value(); ok {
			existing.SetTenantPhone(value)
		}
	}

	if input.TenantEmail.IsPresent() {
		if input.TenantEmail.IsNull() {
			existing.SetTenantEmail("")
		} else if value, ok := input.TenantEmail.Value(); ok {
			existing.SetTenantEmail(value)
		}
	}

	if input.Accompanying.IsPresent() {
		if input.Accompanying.IsNull() {
			existing.SetAccompanying(0)
		} else if selection, ok := input.Accompanying.Value(); ok {
			resolvedID, resolveErr := ls.resolveCatalogValue(ctx, tx, listingmodel.CatalogCategoryAccompanyingType, "accompanying", selection)
			if resolveErr != nil {
				return resolveErr
			}
			existing.SetAccompanying(listingmodel.AccompanyingType(resolvedID))
		}
	}

	//update the listing
	err = ls.listingRepository.UpdateListing(ctx, tx, existing)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	err = ls.gsi.CreateAudit(ctx, tx, globalmodel.TableListings, "Anúncio atualizado")
	if err != nil {
		return err
	}

	return
}
