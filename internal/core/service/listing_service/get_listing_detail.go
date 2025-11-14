package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CatalogValueDetail agrega informações úteis de um valor de catálogo relacionado ao listing.
type CatalogValueDetail struct {
	NumericValue uint8
	Slug         string
	Label        string
}

// FinancingBlockerDetail combina o registro de bloqueio com seu valor de catálogo.
type FinancingBlockerDetail struct {
	Item    listingmodel.FinancingBlockerInterface
	Catalog *CatalogValueDetail
}

// FeatureDetail agrega informações de catálogo para uma feature associada ao listing.
type FeatureDetail struct {
	Feature     string
	Description string
	Quantity    uint8
}

// GuaranteeDetail combina a garantia com a entrada correspondente no catálogo.
type GuaranteeDetail struct {
	Item    listingmodel.GuaranteeInterface
	Catalog *CatalogValueDetail
}

// ListingDetailOutput encapsula o listing e metadados associados.
type ListingDetailOutput struct {
	Listing           listingmodel.ListingInterface
	Features          []FeatureDetail
	Owner             *CatalogValueDetail
	Delivered         *CatalogValueDetail
	WhoLives          *CatalogValueDetail
	Transaction       *CatalogValueDetail
	Installment       *CatalogValueDetail
	Visit             *CatalogValueDetail
	Accompanying      *CatalogValueDetail
	FinancingBlockers []FinancingBlockerDetail
	Guarantees        []GuaranteeDetail
	PhotoSessionID    *uint64
}

// GetListingDetail retorna todos os dados de um listing específico.
func (ls *listingService) GetListingDetail(ctx context.Context, listingID int64) (ListingDetailOutput, error) {
	var output ListingDetailOutput

	if listingID <= 0 {
		return output, utils.ValidationError("listingId", "listingId must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return output, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := ls.gsi.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.detail.tx_start_error", "err", txErr)
		return output, utils.InternalError("")
	}
	defer func() {
		_ = ls.gsi.RollbackTransaction(ctx, tx)
	}()

	listing, repoErr := ls.listingRepository.GetListingVersionByID(ctx, tx, listingID)
	if repoErr != nil {
		if errors.Is(repoErr, sql.ErrNoRows) {
			return output, utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, repoErr)
		logger.Error("listing.detail.get_listing_error", "err", repoErr, "listing_id", listingID)
		return output, utils.InternalError("")
	}

	// Validate ownership
	userID, uidErr := ls.gsi.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return output, uidErr
	}

	if listing.UserID() != userID {
		logger.Warn("unauthorized_detail_access_attempt",
			"listing_id", listingID,
			"listing_identity_id", listing.IdentityID(),
			"requester_user_id", userID,
			"owner_user_id", listing.UserID())
		return output, utils.AuthorizationError("not authorized to access this listing")
	}

	versionSummaries, listErr := ls.listingRepository.ListListingVersions(ctx, tx, listingrepository.ListListingVersionsFilter{
		ListingIdentityID: listing.IdentityID(),
		IncludeDeleted:    false,
	})
	if listErr != nil {
		utils.SetSpanError(ctx, listErr)
		logger.Error("listing.detail.list_versions_error", "err", listErr, "listing_identity_id", listing.IdentityID())
		return output, utils.InternalError("")
	}

	if len(versionSummaries) > 0 {
		versions := make([]listingmodel.ListingVersionInterface, 0, len(versionSummaries))
		var draftVersion listingmodel.ListingVersionInterface

		for _, summary := range versionSummaries {
			version := summary.Version
			if version == nil {
				continue
			}

			versions = append(versions, version)
			if summary.IsActive {
				listing.SetActiveVersion(version)
			}
			if !summary.IsActive && version.Status() == listingmodel.StatusDraft {
				draftVersion = version
			}
		}

		listing.SetVersions(versions)
		if draftVersion != nil {
			listing.SetDraftVersion(draftVersion)
		} else {
			listing.ClearDraftVersion()
		}
	} else {
		listing.ClearDraftVersion()
	}

	output.Listing = listing

	// Buscar booking ativo de photo session se existir
	booking, bookingErr := ls.photoSessionSvc.GetActiveBookingByListingID(ctx, tx, listingID)
	if bookingErr != nil && !errors.Is(bookingErr, sql.ErrNoRows) {
		// Apenas loga warning se não for ErrNoRows (ausência de booking é esperado)
		logger.Warn("listing.detail.get_active_booking_warning", "listing_id", listingID, "err", bookingErr)
		// Não retorna erro; apenas não preenche o campo
	} else if bookingErr == nil && booking != nil {
		bookingID := booking.ID()
		output.PhotoSessionID = &bookingID
	}

	cache := make(map[string]*CatalogValueDetail)

	if listingFeatures := listing.Features(); len(listingFeatures) > 0 {
		ids := make([]int64, 0, len(listingFeatures))
		seen := make(map[int64]struct{}, len(listingFeatures))
		for _, feature := range listingFeatures {
			featureID := feature.FeatureID()
			if featureID == 0 {
				continue
			}
			if _, ok := seen[featureID]; !ok {
				seen[featureID] = struct{}{}
				ids = append(ids, featureID)
			}
		}

		featureMap, ferr := ls.listingRepository.GetBaseFeaturesByIDs(ctx, tx, ids)
		if ferr != nil {
			utils.SetSpanError(ctx, ferr)
			logger.Error("listing.detail.get_features_metadata_error", "err", ferr, "listing_id", listingID)
			return output, utils.InternalError("")
		}

		featureDetails := make([]FeatureDetail, 0, len(listingFeatures))
		for _, feature := range listingFeatures {
			featureID := feature.FeatureID()
			metadata, ok := featureMap[featureID]
			if !ok {
				logger.Warn("listing.detail.base_feature_not_found", "feature_id", featureID)
				featureDetails = append(featureDetails, FeatureDetail{
					Feature:     "",
					Description: "",
					Quantity:    feature.Quantity(),
				})
				continue
			}

			featureDetails = append(featureDetails, FeatureDetail{
				Feature:     metadata.Feature(),
				Description: metadata.Description(),
				Quantity:    feature.Quantity(),
			})
		}

		output.Features = featureDetails
	}

	ownerDetail, derr := ls.fetchCatalogValueDetail(ctx, tx, listingmodel.CatalogCategoryPropertyOwner, uint8(listing.Owner()), cache)
	if derr != nil {
		return output, derr
	}
	output.Owner = ownerDetail

	deliveredDetail, derr := ls.fetchCatalogValueDetail(ctx, tx, listingmodel.CatalogCategoryPropertyDelivered, uint8(listing.Delivered()), cache)
	if derr != nil {
		return output, derr
	}
	output.Delivered = deliveredDetail

	whoLivesDetail, derr := ls.fetchCatalogValueDetail(ctx, tx, listingmodel.CatalogCategoryWhoLives, uint8(listing.WhoLives()), cache)
	if derr != nil {
		return output, derr
	}
	output.WhoLives = whoLivesDetail

	transactionDetail, derr := ls.fetchCatalogValueDetail(ctx, tx, listingmodel.CatalogCategoryTransactionType, uint8(listing.Transaction()), cache)
	if derr != nil {
		return output, derr
	}
	output.Transaction = transactionDetail

	installmentDetail, derr := ls.fetchCatalogValueDetail(ctx, tx, listingmodel.CatalogCategoryInstallmentPlan, uint8(listing.Installment()), cache)
	if derr != nil {
		return output, derr
	}
	output.Installment = installmentDetail

	visitDetail, derr := ls.fetchCatalogValueDetail(ctx, tx, listingmodel.CatalogCategoryVisitType, uint8(listing.Visit()), cache)
	if derr != nil {
		return output, derr
	}
	output.Visit = visitDetail

	accompanyingDetail, derr := ls.fetchCatalogValueDetail(ctx, tx, listingmodel.CatalogCategoryAccompanyingType, uint8(listing.Accompanying()), cache)
	if derr != nil {
		return output, derr
	}
	output.Accompanying = accompanyingDetail

	if blockers := listing.FinancingBlockers(); len(blockers) > 0 {
		details := make([]FinancingBlockerDetail, 0, len(blockers))
		for _, blocker := range blockers {
			catalog, ferr := ls.fetchCatalogValueDetail(ctx, tx, listingmodel.CatalogCategoryFinancingBlocker, uint8(blocker.Blocker()), cache)
			if ferr != nil {
				return output, ferr
			}
			details = append(details, FinancingBlockerDetail{Item: blocker, Catalog: catalog})
		}
		output.FinancingBlockers = details
	}

	if guarantees := listing.Guarantees(); len(guarantees) > 0 {
		details := make([]GuaranteeDetail, 0, len(guarantees))
		for _, guarantee := range guarantees {
			catalog, gerr := ls.fetchCatalogValueDetail(ctx, tx, listingmodel.CatalogCategoryGuaranteeType, uint8(guarantee.Guarantee()), cache)
			if gerr != nil {
				return output, gerr
			}
			details = append(details, GuaranteeDetail{Item: guarantee, Catalog: catalog})
		}
		output.Guarantees = details
	}

	return output, nil
}

func (ls *listingService) fetchCatalogValueDetail(
	ctx context.Context,
	tx *sql.Tx,
	category string,
	numeric uint8,
	cache map[string]*CatalogValueDetail,
) (*CatalogValueDetail, error) {
	if numeric == 0 {
		return nil, nil
	}

	key := fmt.Sprintf("%s:%d", category, numeric)
	if cached, ok := cache[key]; ok {
		return cached, nil
	}

	value, err := ls.listingRepository.GetCatalogValueByNumeric(ctx, tx, category, numeric)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger := utils.LoggerFromContext(ctx)
			logger.Warn("listing.detail.catalog_not_found", "category", category, "numeric", numeric)
			detail := &CatalogValueDetail{NumericValue: numeric}
			cache[key] = detail
			return detail, nil
		}
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	detail := &CatalogValueDetail{
		NumericValue: value.NumericValue(),
		Slug:         value.Slug(),
		Label:        value.Label(),
	}
	cache[key] = detail

	return detail, nil
}
