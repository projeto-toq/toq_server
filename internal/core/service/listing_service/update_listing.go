package listingservices

import (
	"context"
	"database/sql"
	"errors"
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

	// Determine which ID to use (prefer VersionID over deprecated ID)
	versionID := input.VersionID
	if versionID == 0 {
		versionID = input.ID
	}
	if versionID == 0 {
		return utils.BadRequest("versionId is required")
	}

	// Get the listing version
	existing, err := ls.listingRepository.GetListingVersionByID(ctx, tx, versionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("listing version")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.update.get_version_error", "err", err, "version_id", versionID)
		return utils.InternalError("")
	}

	// Check if the listing is in draft status
	if existing.Status() != listingmodel.StatusDraft {
		return utils.ConflictError("listing cannot be updated unless in draft status")
	}

	// Check if the user is the owner of the listing
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

	//update the listing version
	err = ls.listingRepository.UpdateListingVersion(ctx, tx, existing)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.update.update_version_error", "err", err, "version_id", versionID)
		return utils.InternalError("Failed to update listing")
	}

	err = ls.gsi.CreateAudit(ctx, tx, globalmodel.TableListings, "Anúncio atualizado")
	if err != nil {
		return err
	}

	return
}
