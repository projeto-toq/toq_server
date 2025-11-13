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

// ListListings aplica filtros com paginação para retornar listings ao painel admin.
func (la *ListingAdapter) ListListings(ctx context.Context, tx *sql.Tx, filter listingrepository.ListListingsFilter) (listingrepository.ListListingsResult, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return listingrepository.ListListingsResult{}, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	conditions := []string{"lv.deleted = 0", "li.deleted = 0"}
	args := make([]any, 0)

	if filter.Status != nil {
		conditions = append(conditions, "lv.status = ?")
		args = append(args, int(*filter.Status))
	}
	if filter.Code != nil {
		conditions = append(conditions, "lv.code = ?")
		args = append(args, int64(*filter.Code))
	}
	if filter.Title != "" {
		conditions = append(conditions, "(COALESCE(lv.title, '') LIKE ? OR COALESCE(lv.description, '') LIKE ?)")
		args = append(args, filter.Title, filter.Title)
	}
	if filter.ZipCode != "" {
		conditions = append(conditions, "lv.zip_code LIKE ?")
		args = append(args, filter.ZipCode)
	}
	if filter.City != "" {
		conditions = append(conditions, "lv.city LIKE ?")
		args = append(args, filter.City)
	}
	if filter.Neighborhood != "" {
		conditions = append(conditions, "lv.neighborhood LIKE ?")
		args = append(args, filter.Neighborhood)
	}
	if filter.UserID != nil {
		conditions = append(conditions, "lv.user_id = ?")
		args = append(args, *filter.UserID)
	}

	if filter.MinSellPrice != nil {
		conditions = append(conditions, "COALESCE(lv.sell_net, 0) >= ?")
		args = append(args, *filter.MinSellPrice)
	}
	if filter.MaxSellPrice != nil {
		conditions = append(conditions, "COALESCE(lv.sell_net, 0) <= ?")
		args = append(args, *filter.MaxSellPrice)
	}
	if filter.MinRentPrice != nil {
		conditions = append(conditions, "COALESCE(lv.rent_net, 0) >= ?")
		args = append(args, *filter.MinRentPrice)
	}
	if filter.MaxRentPrice != nil {
		conditions = append(conditions, "COALESCE(lv.rent_net, 0) <= ?")
		args = append(args, *filter.MaxRentPrice)
	}
	if filter.MinLandSize != nil {
		conditions = append(conditions, "COALESCE(lv.land_size, 0) >= ?")
		args = append(args, *filter.MinLandSize)
	}
	if filter.MaxLandSize != nil {
		conditions = append(conditions, "COALESCE(lv.land_size, 0) <= ?")
		args = append(args, *filter.MaxLandSize)
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	baseSelect := fmt.Sprintf(`SELECT
%s
FROM listing_versions lv
INNER JOIN listing_identities li ON li.id = lv.listing_identity_id`, listingSelectColumns)

	listQuery := baseSelect + " " + whereClause + " ORDER BY l.id DESC LIMIT ? OFFSET ?"
	listArgs := append([]any{}, args...)
	offset := (filter.Page - 1) * filter.Limit
	listArgs = append(listArgs, filter.Limit, offset)

	rows, queryErr := la.QueryContext(ctx, tx, "select", listQuery, listArgs...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.listing.list.query_error", "error", queryErr)
		return listingrepository.ListListingsResult{}, fmt.Errorf("list listings query: %w", queryErr)
	}
	defer rows.Close()

	result := listingrepository.ListListingsResult{}

	for rows.Next() {
		entity, scanErr := scanListingEntity(rows)
		if scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.listing.list.scan_error", "error", scanErr)
			return listingrepository.ListListingsResult{}, fmt.Errorf("scan listing row: %w", scanErr)
		}

		listing := listingconverters.ListingEntityToDomain(entity)
		if listing != nil {
			result.Records = append(result.Records, listingrepository.ListingRecord{
				Listing: listing,
			})
		}
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.list.rows_error", "error", err)
		return listingrepository.ListListingsResult{}, fmt.Errorf("iterate listing rows: %w", err)
	}

	countQuery := "SELECT COUNT(*) FROM listing_versions lv INNER JOIN listing_identities li ON li.id = lv.listing_identity_id " + whereClause
	var total int64
	if countErr := la.QueryRowContext(ctx, tx, "select", countQuery, args...).Scan(&total); countErr != nil {
		utils.SetSpanError(ctx, countErr)
		logger.Error("mysql.listing.list.count_error", "error", countErr)
		return listingrepository.ListListingsResult{}, fmt.Errorf("count listings: %w", countErr)
	}
	result.Total = total

	return result, nil
}
