package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	listingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/converters"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) ListListingVersions(ctx context.Context, tx *sql.Tx, filter listingrepository.ListListingVersionsFilter) ([]listingrepository.ListingVersionSummary, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if filter.ListingIdentityID == 0 {
		return nil, fmt.Errorf("listing identity id is required")
	}

	conditions := []string{"lv.listing_identity_id = ?"}
	args := []any{filter.ListingIdentityID}

	if !filter.IncludeDeleted {
		conditions = append(conditions, "lv.deleted = 0")
	}

	query := fmt.Sprintf(`SELECT
%s
FROM listing_versions lv
INNER JOIN listing_identities li ON li.id = lv.listing_identity_id
WHERE %s
ORDER BY lv.version DESC`, listingSelectColumns, strings.Join(conditions, " AND "))

	rows, queryErr := la.QueryContext(ctx, tx, "select", query, args...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.listing.list_versions.query_error", "error", queryErr)
		return nil, fmt.Errorf("query listing versions: %w", queryErr)
	}
	defer rows.Close()

	summaries := make([]listingrepository.ListingVersionSummary, 0)

	for rows.Next() {
		entity, scanErr := scanListingEntity(rows)
		if scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.listing.list_versions.scan_error", "error", scanErr)
			return nil, fmt.Errorf("scan listing version row: %w", scanErr)
		}

		listing := listingconverters.ListingEntityToDomain(entity)
		if listing == nil {
			continue
		}

		activeVersion := listing.ActiveVersion()
		isActive := entity.ActiveVersionID.Valid && entity.ActiveVersionID.Int64 == entity.ID

		summaries = append(summaries, listingrepository.ListingVersionSummary{
			Version:  activeVersion,
			IsActive: isActive,
		})
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.list_versions.rows_error", "error", err)
		return nil, fmt.Errorf("iterate listing versions rows: %w", err)
	}

	return summaries, nil
}
