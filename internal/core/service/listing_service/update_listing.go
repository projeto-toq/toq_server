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

func (ls *listingService) UpdateListing(ctx context.Context, listing listingmodel.ListingInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	tx, err := ls.gsi.StartTransaction(ctx)
	if err != nil {
		slog.Error("listing.update.tx_start_error", "err", err)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("listing.update.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	err = ls.updateListing(ctx, tx, listing)
	if err != nil {
		return
	}

	err = ls.gsi.CommitTransaction(ctx, tx)
	if err != nil {
		slog.Error("listing.update.tx_commit_error", "err", err)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (ls *listingService) updateListing(ctx context.Context, tx *sql.Tx, listing listingmodel.ListingInterface) (err error) {

	exist := true
	//check if exists the listing
	existing, err := ls.listingRepository.GetListingByID(ctx, tx, listing.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exist = false
		} else {
			return
		}
	}

	if !exist {
		err = utils.ErrInternalServer
		return
	}

	//check if the listing is in draft status
	if existing.Status() != listingmodel.StatusDraft {
		err = utils.ErrInternalServer
		return
	}

	//check if the user is the owner of the listing
	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)
	if existing.UserID() != infos.ID {
		err = utils.ErrInternalServer
		return
	}

	//update only the allowed fields
	existing.SetOwner(listing.Owner())
	existing.SetFeatures(listing.Features())
	existing.SetLandSize(listing.LandSize())
	existing.SetCorner(listing.Corner())
	existing.SetNonBuildable(listing.NonBuildable())
	existing.SetBuildable(listing.Buildable())
	existing.SetDelivered(listing.Delivered())
	existing.SetWhoLives(listing.WhoLives())
	existing.SetDescription(listing.Description())
	existing.SetTransaction(listing.Transaction())
	existing.SetSellNet(listing.SellNet())
	existing.SetRentNet(listing.RentNet())
	existing.SetCondominium(listing.Condominium())
	existing.SetCondominium(listing.Condominium())
	existing.SetAnnualTax(listing.AnnualTax())
	existing.SetAnnualGroundRent(listing.AnnualGroundRent())
	existing.SetExchange(listing.Exchange())
	existing.SetExchangePercentual(listing.ExchangePercentual())
	existing.SetExchangePlaces(listing.ExchangePlaces())
	existing.SetInstallment(listing.Installment())
	existing.SetFinancing(listing.Financing())
	existing.SetFinancingBlockers(listing.FinancingBlockers())
	existing.SetGuarantees(listing.Guarantees())
	existing.SetVisit(listing.Visit())
	existing.SetTenantName(listing.TenantName())
	existing.SetTenantPhone(listing.TenantPhone())
	existing.SetTenantEmail(listing.TenantEmail())
	existing.SetAccompanying(listing.Accompanying())

	//update the listing
	err = ls.listingRepository.UpdateListing(ctx, tx, existing)
	if err != nil {
		return
	}

	err = ls.gsi.CreateAudit(ctx, tx, globalmodel.TableListings, "Anúncio atualizado")
	if err != nil {
		return
	}

	return
}
