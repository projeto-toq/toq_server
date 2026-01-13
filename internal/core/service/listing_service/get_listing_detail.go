package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
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
	OwnerDetail       *ListingOwnerDetail
	Owner             *CatalogValueDetail
	Delivered         *CatalogValueDetail
	WhoLives          *CatalogValueDetail
	Transaction       *CatalogValueDetail
	Installment       *CatalogValueDetail
	Visit             *CatalogValueDetail
	Accompanying      *CatalogValueDetail
	LandTerrainType   *CatalogValueDetail
	WarehouseSector   *CatalogValueDetail
	FinancingBlockers []FinancingBlockerDetail
	Guarantees        []GuaranteeDetail
	PhotoSessionID    *uint64
	Performance       ListingPerformanceMetrics
	FavoritesCount    int64
	IsFavorite        bool
}

// ListingOwnerDetail exposes owner profile metadata enriched with engagement metrics.
type ListingOwnerDetail struct {
	ID                int64
	FullName          string
	MemberSinceMonths int
	PhotoURL          string
	Metrics           OwnerEngagementMetrics
}

// OwnerEngagementMetrics contains response KPI placeholders sourced from owner metrics repository.
type OwnerEngagementMetrics struct {
	VisitAverageSeconds    sql.NullInt64
	ProposalAverageSeconds sql.NullInt64
}

// ListingPerformanceMetrics describes property engagement metrics (shares, views, favorites).
// TODO(listing-metrics): Populate these counters once the analytics provider is connected.
type ListingPerformanceMetrics struct {
	Shares    int64
	Views     int64
	Favorites int64
}

// GetListingDetail retrieves comprehensive details of a listing by its identity ID
//
// This method returns the ACTIVE version of a listing (referenced by listing_identities.active_version_id).
// It orchestrates the following operations:
//  1. Validates listingIdentityId parameter (must be > 0)
//  2. Starts read-only transaction for data consistency
//  3. Fetches listing identity record (listing_identities table)
//  4. Validates ownership (identity.user_id == authenticated user_id)
//  5. Fetches active version via identity.active_version_id
//  6. Enriches catalog values (owner, delivered, whoLives, transaction, etc.)
//  7. Fetches related entities (features, guarantees, exchange places, financing blockers)
//  8. Queries active photo session booking (if exists)
//  9. Lists all versions to populate draft metadata
//
// The operation is read-only and uses a read-only transaction for performance.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging. Must contain request metadata and authenticated user.
//   - listingIdentityId: Unique identifier of the listing identity (listing_identities.id). Must be > 0.
//
// Returns:
//   - detail: ListingDetailOutput with enriched listing data, catalog values, and related entities
//   - err: Domain error with appropriate HTTP status code:
//   - 400 (Bad Request) if listingIdentityId <= 0
//   - 403 (Forbidden) if requester is not the listing owner
//   - 404 (Not Found) if listing identity does not exist
//   - 500 (Internal) for infrastructure failures (DB, cache, transaction errors)
//
// Business Rules:
//   - Only the listing owner (identity.user_id == authenticated user_id) can access details
//   - Returns the ACTIVE version (referenced by identity.active_version_id)
//   - Draft version metadata is included in response if it exists
//   - Catalog values are enriched with slug and label for frontend display
//   - Photo session booking ID is included if active booking exists
//
// Side Effects:
//   - Logs warning if requester is not the listing owner (audit trail)
//   - Logs warning if catalog values are not found (data integrity issue)
//
// Example:
//
//	detail, err := svc.GetListingDetail(ctx, 1024)
//	if err != nil {
//	    // Handle error (already logged by service)
//	    return err
//	}
//	// detail.Listing contains active version data
//	// detail.PhotoSessionID contains booking ID (if exists)
func (ls *listingService) GetListingDetail(ctx context.Context, listingIdentityId int64) (ListingDetailOutput, error) {
	var output ListingDetailOutput

	// Validate listingIdentityId parameter (business rule)
	if listingIdentityId <= 0 {
		return output, utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}

	// Initialize tracing for distributed observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		ctx = utils.ContextWithLogger(ctx)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("listing.detail.tracer_error", "stage", "tracer_init", "err", err)
		return output, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Failed to start tracer for listing detail", map[string]any{
			"stage": "tracer_init",
		})
	}
	defer spanEnd()

	// Ensure logger propagation with request_id and trace_id
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Start read-only transaction (performance optimization for SELECT-only operations)
	tx, txErr := ls.gsi.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.detail.tx_start_error", "stage", "tx_ro_start", "err", txErr)
		return output, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Failed to start read-only transaction", map[string]any{
			"stage": "tx_ro_start",
		})
	}
	defer func() {
		// Always rollback read-only transactions (no side effects)
		_ = ls.gsi.RollbackTransaction(ctx, tx)
	}()

	// Fetch listing identity record to validate existence and ownership
	identity, identityErr := ls.listingRepository.GetListingIdentityByID(ctx, tx, listingIdentityId)
	if identityErr != nil {
		if errors.Is(identityErr, sql.ErrNoRows) {
			// Listing identity not found (404)
			return output, utils.NotFoundError("Listing")
		}
		// Infrastructure error (500)
		utils.SetSpanError(ctx, identityErr)
		logger.Error("listing.detail.get_identity_error", "stage", "identity_lookup", "err", identityErr, "listing_identity_id", listingIdentityId)
		return output, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Failed to load listing identity", map[string]any{
			"stage":             "identity_lookup",
			"listingIdentityId": listingIdentityId,
		})
	}

	// Validate ownership before returning sensitive data
	userID, uidErr := ls.gsi.GetUserIDFromContext(ctx)
	if uidErr != nil {
		// Context missing user_id (should never happen after AuthMiddleware)
		return output, uidErr
	}

	if identity.UserID != userID {
		// Requester is not the owner (audit trail + 403 Forbidden)
		logger.Warn("unauthorized_detail_access_attempt",
			"listing_identity_id", listingIdentityId,
			"requester_user_id", userID,
			"owner_user_id", identity.UserID)
		return output, utils.AuthorizationError("not authorized to access this listing")
	}

	// Fetch active version via identity.active_version_id
	if !identity.ActiveVersionID.Valid || identity.ActiveVersionID.Int64 == 0 {
		// Edge case: identity exists but no active version (should not happen in production)
		logger.Error("listing.detail.no_active_version", "stage", "active_version_missing", "listing_identity_id", listingIdentityId)
		return output, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Listing has no active version", map[string]any{
			"stage":             "active_version_missing",
			"listingIdentityId": listingIdentityId,
		})
	}

	activeVersionID := identity.ActiveVersionID.Int64

	// Fetch active version FULLY ENRICHED (includes Features, Guarantees, ExchangePlaces, etc.)
	// Note: GetListingVersionByID internally calls GetListingByQuery which enriches all satellite tables
	listing, repoErr := ls.listingRepository.GetListingVersionByID(ctx, tx, activeVersionID)
	if repoErr != nil {
		if errors.Is(repoErr, sql.ErrNoRows) {
			// Active version referenced but not found (data integrity issue)
			logger.Error("listing.detail.active_version_not_found", "stage", "active_version_lookup", "listing_identity_id", listingIdentityId, "active_version_id", activeVersionID)
			return output, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Active listing version not found", map[string]any{
				"stage":             "active_version_lookup",
				"listingIdentityId": listingIdentityId,
				"activeVersionId":   activeVersionID,
			})
		}
		// Infrastructure error
		utils.SetSpanError(ctx, repoErr)
		logger.Error("listing.detail.get_listing_error", "stage", "listing_version_fetch", "err", repoErr, "listing_identity_id", listingIdentityId, "active_version_id", activeVersionID)
		return output, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Failed to load listing version", map[string]any{
			"stage":             "listing_version_fetch",
			"listingIdentityId": listingIdentityId,
			"activeVersionId":   activeVersionID,
		})
	}

	// Set listing identity metadata (UUID, identity ID)
	listing.SetIdentityID(listingIdentityId)
	listing.SetUUID(identity.UUID)
	listing.SetActiveVersionID(activeVersionID)

	ownerDetail, ownerErr := ls.buildOwnerDetail(ctx, tx, identity.UserID)
	if ownerErr != nil {
		return output, ownerErr
	}
	output.OwnerDetail = ownerDetail

	// Populate favorites engagement metrics and requester flag
	favCounts, favErr := ls.favoriteRepo.CountByListingIdentities(ctx, tx, []int64{listingIdentityId})
	if favErr != nil {
		utils.SetSpanError(ctx, favErr)
		logger.Error("listing.detail.fav_count_error", "stage", "favorites_count", "err", favErr, "listing_identity_id", listingIdentityId)
		return output, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Failed to load listing favorites", map[string]any{
			"stage":             "favorites_count",
			"listingIdentityId": listingIdentityId,
		})
	}

	viewsCount, viewsErr := ls.viewRepo.GetCount(ctx, tx, listingIdentityId)
	if viewsErr != nil {
		utils.SetSpanError(ctx, viewsErr)
		logger.Error("listing.detail.views_count_error", "stage", "views_fetch", "err", viewsErr, "listing_identity_id", listingIdentityId)
		return output, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Failed to load listing views", map[string]any{
			"stage":             "views_fetch",
			"listingIdentityId": listingIdentityId,
		})
	}

	output.Performance = ListingPerformanceMetrics{Favorites: favCounts[listingIdentityId], Views: viewsCount}
	output.FavoritesCount = favCounts[listingIdentityId]

	if favFlags, flagErr := ls.favoriteRepo.GetUserFlags(ctx, tx, []int64{listingIdentityId}, userID); flagErr == nil {
		output.IsFavorite = favFlags[listingIdentityId]
	} else {
		utils.SetSpanError(ctx, flagErr)
		logger.Warn("listing.detail.fav_flags_error", "err", flagErr, "listing_identity_id", listingIdentityId)
	}

	// Fetch draft version (if exists) for metadata
	// Note: Avoid fetching all versions; only draft is needed
	draftVersion, draftErr := ls.listingRepository.GetDraftVersionByListingIdentityID(ctx, tx, listingIdentityId)
	if draftErr != nil && !errors.Is(draftErr, sql.ErrNoRows) {
		// Log warning but do not fail request (draft is optional)
		logger.Warn("listing.detail.get_draft_warning", "listing_identity_id", listingIdentityId, "err", draftErr)
	} else if draftErr == nil && draftVersion != nil {
		listing.SetDraftVersion(draftVersion)
	} else {
		listing.ClearDraftVersion()
	}

	output.Listing = listing

	// Fetch active photo session booking (optional, may not exist)
	booking, bookingErr := ls.photoSessionSvc.GetActiveBookingByListingIdentityID(ctx, tx, listingIdentityId)
	if bookingErr != nil && !errors.Is(bookingErr, sql.ErrNoRows) {
		// Log warning if not ErrNoRows (absence of booking is expected)
		logger.Warn("listing.detail.get_active_booking_warning", "listing_identity_id", listingIdentityId, "err", bookingErr)
		// Do not return error; just omit photoSessionId field
	} else if bookingErr == nil && booking != nil {
		bookingID := booking.ID()
		output.PhotoSessionID = &bookingID
	}

	// Cache for catalog values to avoid duplicate queries
	cache := make(map[string]*CatalogValueDetail)

	// Enrich features with base feature metadata (description, name)
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
			logger.Error("listing.detail.get_features_metadata_error", "err", ferr, "listing_identity_id", listingIdentityId)
			return output, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Failed to load feature metadata", map[string]any{
				"stage":             "features_metadata",
				"listingIdentityId": listingIdentityId,
			})
		}

		featureDetails := make([]FeatureDetail, 0, len(listingFeatures))
		for _, feature := range listingFeatures {
			featureID := feature.FeatureID()
			metadata, ok := featureMap[featureID]
			if !ok {
				// Feature ID referenced but not found in base_features (data integrity issue)
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

	// Enrich catalog values with slug and label for frontend display
	ownerCatalogDetail, derr := ls.fetchCatalogValueDetail(ctx, tx, listingmodel.CatalogCategoryPropertyOwner, uint8(listing.Owner()), cache)
	if derr != nil {
		return output, derr
	}
	output.Owner = ownerCatalogDetail

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

	landTerrainTypeDetail, derr := ls.fetchCatalogValueDetail(ctx, tx, listingmodel.CatalogCategoryLandTerrainType, uint8(listing.LandTerrainType()), cache)
	if derr != nil {
		return output, derr
	}
	output.LandTerrainType = landTerrainTypeDetail

	warehouseSectorDetail, derr := ls.fetchCatalogValueDetail(ctx, tx, listingmodel.CatalogCategoryWarehouseSector, uint8(listing.WarehouseSector()), cache)
	if derr != nil {
		return output, derr
	}
	output.WarehouseSector = warehouseSectorDetail

	// Enrich financing blockers with catalog labels
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

	// Enrich guarantees with catalog labels
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

	updatedViews, viewIncErr := ls.registerListingView(ctx, listingIdentityId)
	if viewIncErr != nil {
		return output, viewIncErr
	}
	output.Performance.Views = updatedViews

	return output, nil
}

// buildOwnerDetail fetches owner profile and response metrics, computing member since months.
func (ls *listingService) buildOwnerDetail(ctx context.Context, tx *sql.Tx, ownerID int64) (*ListingOwnerDetail, error) {
	owner, err := ls.userRepository.GetUserByID(ctx, tx, ownerID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("listing.detail.owner_not_found", "stage", "owner_lookup", "owner_id", ownerID, "err", err)
		} else {
			logger.Error("listing.detail.owner_fetch_error", "stage", "owner_lookup", "owner_id", ownerID, "err", err)
		}
		return nil, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Failed to load owner information", map[string]any{
			"stage":   "owner_lookup",
			"ownerId": ownerID,
		})
	}

	metrics, metricsErr := ls.ownerMetricsRepo.GetByOwnerID(ctx, tx, ownerID)
	switch {
	case metricsErr == nil:
		// OK
	case errors.Is(metricsErr, sql.ErrNoRows):
		metrics = usermodel.NewOwnerResponseMetrics()
	default:
		utils.SetSpanError(ctx, metricsErr)
		logger := utils.LoggerFromContext(ctx)
		logger.Warn("listing.detail.owner_metrics_error", "owner_id", ownerID, "err", metricsErr)
		metrics = usermodel.NewOwnerResponseMetrics()
	}

	memberSinceMonths := monthsSince(owner.GetCreatedAt(), time.Now().UTC())

	return &ListingOwnerDetail{
		ID:                owner.GetID(),
		FullName:          owner.GetFullName(),
		MemberSinceMonths: memberSinceMonths,
		Metrics: OwnerEngagementMetrics{
			VisitAverageSeconds:    metrics.VisitAverageSeconds(),
			ProposalAverageSeconds: metrics.ProposalAverageSeconds(),
		},
	}, nil
}

// fetchCatalogValueDetail retrieves and caches catalog value metadata
//
// This private helper reduces redundant database queries by caching catalog values.
// Returns nil for numeric value 0 (indicates no value selected).
//
// Parameters:
//   - ctx: current context
//   - tx: database transaction
//   - category: catalog category (e.g., "property_owner", "transaction_type")
//   - numeric: catalog numeric value (tinyint)
//   - cache: in-memory cache map to avoid duplicate queries
//
// Returns:
//   - *CatalogValueDetail: enriched catalog value with slug and label
//   - error: infrastructure error if database query fails
func (ls *listingService) fetchCatalogValueDetail(
	ctx context.Context,
	tx *sql.Tx,
	category string,
	numeric uint8,
	cache map[string]*CatalogValueDetail,
) (*CatalogValueDetail, error) {
	// Numeric value 0 indicates no value selected (return nil)
	if numeric == 0 {
		return nil, nil
	}

	// Check cache to avoid redundant database queries
	key := fmt.Sprintf("%s:%d", category, numeric)
	if cached, ok := cache[key]; ok {
		return cached, nil
	}

	// Query catalog value from database
	value, err := ls.listingRepository.GetCatalogValueByNumeric(ctx, tx, category, numeric)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Catalog value not found (data integrity issue, log warning)
			logger := utils.LoggerFromContext(ctx)
			logger.Warn("listing.detail.catalog_not_found", "category", category, "numeric", numeric)
			detail := &CatalogValueDetail{NumericValue: numeric}
			cache[key] = detail
			return detail, nil
		}
		// Infrastructure error
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("listing.detail.catalog_error", "stage", "catalog_lookup", "category", category, "numeric", numeric, "err", err)
		return nil, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Failed to load catalog value", map[string]any{
			"stage":    "catalog_lookup",
			"category": category,
			"numeric":  numeric,
		})
	}

	// Cache result for future use in this request
	detail := &CatalogValueDetail{
		NumericValue: value.NumericValue(),
		Slug:         value.Slug(),
		Label:        value.Label(),
	}
	cache[key] = detail

	return detail, nil
}
