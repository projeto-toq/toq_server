package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const fallbackListingTimezone = "America/Sao_Paulo"

var stateTimezoneLookup = map[string]string{
	"AC": "America/Rio_Branco",
	"AL": "America/Maceio",
	"AM": "America/Manaus",
	"AP": "America/Belem",
	"BA": "America/Bahia",
	"CE": "America/Fortaleza",
	"DF": "America/Sao_Paulo",
	"ES": "America/Sao_Paulo",
	"GO": "America/Sao_Paulo",
	"MA": "America/Fortaleza",
	"MG": "America/Sao_Paulo",
	"MS": "America/Campo_Grande",
	"MT": "America/Cuiaba",
	"PA": "America/Belem",
	"PB": "America/Fortaleza",
	"PE": "America/Recife",
	"PI": "America/Fortaleza",
	"PR": "America/Sao_Paulo",
	"RJ": "America/Sao_Paulo",
	"RN": "America/Fortaleza",
	"RO": "America/Porto_Velho",
	"RR": "America/Boa_Vista",
	"RS": "America/Sao_Paulo",
	"SC": "America/Sao_Paulo",
	"SE": "America/Maceio",
	"SP": "America/Sao_Paulo",
	"TO": "America/Araguaina",
}

func (ls *listingService) EndUpdateListing(ctx context.Context, input EndUpdateListingInput) (err error) {
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
		logger.Error("listing.end_update.tx_start_error", "err", err, "listing_id", input.ListingID)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.end_update.tx_rollback_error", "err", rbErr, "listing_id", input.ListingID)
			}
		}
	}()

	snapshot, err := ls.listingRepository.GetListingForEndUpdate(ctx, tx, input.ListingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.end_update.fetch_error", "err", err, "listing_id", input.ListingID)
		return utils.InternalError("")
	}

	if snapshot.Status != listingmodel.StatusDraft {
		return utils.ConflictError("Listing must be in draft status to end update")
	}

	userID, uidErr := ls.gsi.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return uidErr
	}

	if snapshot.UserID != userID {
		return utils.AuthorizationError("Only listing owner can end update")
	}

	if verr := ls.validateListingBeforeEndUpdate(ctx, tx, snapshot); verr != nil {
		return verr
	}

	updateErr := ls.listingRepository.UpdateListingStatus(ctx, tx, input.ListingID, listingmodel.StatusPendingAvailability, listingmodel.StatusDraft)
	if updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return utils.ConflictError("Listing status changed while finishing update")
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("listing.end_update.update_status_error", "err", updateErr, "listing_id", input.ListingID)
		return utils.InternalError("")
	}

	if auditErr := ls.gsi.CreateAudit(ctx, tx, globalmodel.TableListings, "Anúncio finalizado (end-update)"); auditErr != nil {
		return auditErr
	}

	timezone := resolveListingTimezone(snapshot)
	agendaInput := scheduleservices.CreateDefaultAgendaInput{
		ListingID: input.ListingID,
		OwnerID:   userID,
		Timezone:  timezone,
		ActorID:   userID,
	}
	if _, agendaErr := ls.scheduleService.CreateDefaultAgendaWithTx(ctx, tx, agendaInput); agendaErr != nil {
		utils.SetSpanError(ctx, agendaErr)
		logger.Error("listing.end_update.create_default_agenda_error", "err", agendaErr, "listing_id", input.ListingID)
		return agendaErr
	}

	if err = ls.gsi.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.end_update.tx_commit_error", "err", err, "listing_id", input.ListingID)
		return utils.InternalError("")
	}

	logger.Info("listing.end_update.completed", "listing_id", input.ListingID, "new_status", listingmodel.StatusPendingPhotoScheduling.String())

	return nil
}

func (ls *listingService) validateListingBeforeEndUpdate(ctx context.Context, tx *sql.Tx, data listingrepository.ListingEndUpdateData) error {
	logger := utils.LoggerFromContext(ctx)

	// Campos básicos obrigatórios para qualquer anúncio em draft.
	if data.Code == 0 {
		return utils.BadRequest("Listing code is required")
	}
	if data.Version == 0 {
		return utils.BadRequest("Listing version is required")
	}
	if strings.TrimSpace(data.ZipCode) == "" {
		return utils.BadRequest("Zip code is required")
	}
	if !data.Street.Valid || strings.TrimSpace(data.Street.String) == "" {
		return utils.BadRequest("Street is required")
	}
	if !data.Number.Valid || strings.TrimSpace(data.Number.String) == "" {
		return utils.BadRequest("Number is required")
	}
	if !data.City.Valid || strings.TrimSpace(data.City.String) == "" {
		return utils.BadRequest("City is required")
	}
	if !data.State.Valid || strings.TrimSpace(data.State.String) == "" {
		return utils.BadRequest("State is required")
	}
	if !data.Title.Valid || strings.TrimSpace(data.Title.String) == "" {
		return utils.BadRequest("Title is required")
	}
	if data.ListingType == 0 {
		return utils.BadRequest("Property type is required")
	}
	if !data.Owner.Valid {
		return utils.BadRequest("Property owner is required")
	}
	if !data.Buildable.Valid {
		return utils.BadRequest("Buildable size is required")
	}
	if !data.Delivered.Valid {
		return utils.BadRequest("Delivered status is required")
	}
	if !data.WhoLives.Valid {
		return utils.BadRequest("Who lives information is required")
	}
	if !data.Description.Valid || strings.TrimSpace(data.Description.String) == "" {
		return utils.BadRequest("Description is required")
	}
	if !data.Transaction.Valid {
		return utils.BadRequest("Transaction type is required")
	}
	if !data.Visit.Valid {
		return utils.BadRequest("Visit type is required")
	}
	if !data.Accompanying.Valid {
		return utils.BadRequest("Accompanying type is required")
	}
	if !data.AnnualTax.Valid {
		return utils.BadRequest("Annual tax is required")
	}
	if data.FeaturesCount == 0 {
		return utils.BadRequest("Listing must include features")
	}

	// Regras condicionais para o tipo de transação.
	txnValue := uint8(data.Transaction.Int16)
	txnCatalog, err := ls.listingRepository.GetCatalogValueByNumeric(ctx, tx, listingmodel.CatalogCategoryTransactionType, txnValue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.BadRequest("Transaction type is invalid")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.end_update.transaction_catalog_error", "err", err, "listing_id", data.ListingID, "transaction_id", txnValue)
		return utils.InternalError("")
	}

	slug := strings.ToLower(strings.TrimSpace(txnCatalog.Slug()))
	needsSaleValidation := slug == "sale" || slug == "both"
	needsRentValidation := slug == "rent" || slug == "both"

	if needsSaleValidation {
		// Quando a transação envolve venda, validamos preço líquido, permuta e barreiras de financiamento.
		if !data.SaleNet.Valid {
			return utils.BadRequest("Sale net value is required")
		}
		if !data.Exchange.Valid {
			return utils.BadRequest("Exchange flag is required")
		}
		if data.Exchange.Valid && data.Exchange.Int16 == 1 {
			if !data.ExchangePercentual.Valid {
				return utils.BadRequest("Exchange percentual is required when exchange is enabled")
			}
			if data.ExchangePlacesCount == 0 {
				return utils.BadRequest("Exchange places are required when exchange is enabled")
			}
		}
		if !data.Financing.Valid {
			return utils.BadRequest("Financing flag is required")
		}
		if data.Financing.Int16 == 0 && data.FinancingBlockersCount == 0 {
			return utils.BadRequest("Financing blockers are required when financing is disabled")
		}
	}

	if needsRentValidation {
		// Nas locações exigimos valor líquido e garantias cadastradas para prosseguir.
		if !data.RentNet.Valid {
			return utils.BadRequest("Rent net value is required")
		}
		if data.GuaranteesCount == 0 {
			return utils.BadRequest("Guarantees are required for rent transactions")
		}
	}

	// Regras específicas por tipo de imóvel.
	propertyOptions := ls.DecodePropertyTypes(ctx, data.ListingType)
	if len(propertyOptions) == 0 {
		return utils.BadRequest("Property type is invalid")
	}
	// Cada option representa um bit ativo na máscara do tipo; usamos isso para derivar validações adicionais.
	needsCondominium := false
	needsLandData := false
	for _, option := range propertyOptions {
		switch option.Code {
		case 1, 4:
			needsCondominium = true
		case 16, 32, 64, 128:
			needsLandData = true
		}
	}

	if needsCondominium && !data.Condominium.Valid {
		return utils.BadRequest("Condominium value is required for the selected property type")
	}

	if needsLandData {
		if !data.LandSize.Valid {
			return utils.BadRequest("Land size is required for the selected property type")
		}
		if !data.Corner.Valid {
			return utils.BadRequest("Corner information is required for the selected property type")
		}
	}

	// Regras adicionais quando quem mora é inquilino.
	if data.WhoLives.Valid {
		whoLivesValue := uint8(data.WhoLives.Int16)
		whoLivesCatalog, catalogErr := ls.listingRepository.GetCatalogValueByNumeric(ctx, tx, listingmodel.CatalogCategoryWhoLives, whoLivesValue)
		if catalogErr != nil {
			if errors.Is(catalogErr, sql.ErrNoRows) {
				return utils.BadRequest("Who lives value is invalid")
			}
			utils.SetSpanError(ctx, catalogErr)
			logger.Error("listing.end_update.wholives_catalog_error", "err", catalogErr, "listing_id", data.ListingID, "who_lives_id", whoLivesValue)
			return utils.InternalError("")
		}

		if strings.ToLower(strings.TrimSpace(whoLivesCatalog.Slug())) == "tenant" {
			if !data.TenantName.Valid || strings.TrimSpace(data.TenantName.String) == "" {
				return utils.BadRequest("Tenant name is required when tenant lives in the property")
			}
			if !data.TenantPhone.Valid || strings.TrimSpace(data.TenantPhone.String) == "" {
				return utils.BadRequest("Tenant phone is required when tenant lives in the property")
			}
			if !data.TenantEmail.Valid || strings.TrimSpace(data.TenantEmail.String) == "" {
				return utils.BadRequest("Tenant email is required when tenant lives in the property")
			}
		}
	}

	return nil
}

func resolveListingTimezone(data listingrepository.ListingEndUpdateData) string {
	if data.State.Valid {
		state := strings.ToUpper(strings.TrimSpace(data.State.String))
		if tz, ok := stateTimezoneLookup[state]; ok && tz != "" {
			return tz
		}
	}
	return fallbackListingTimezone
}
