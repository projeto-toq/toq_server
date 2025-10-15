package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	// "time"

	listingconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/converters"
	listingentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/entity"
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

	conditions := []string{"l.deleted = 0"}
	args := make([]any, 0)

	if filter.Status != nil {
		conditions = append(conditions, "l.status = ?")
		args = append(args, int(*filter.Status))
	}
	if filter.Code != nil {
		conditions = append(conditions, "l.code = ?")
		args = append(args, int64(*filter.Code))
	}
	if filter.Title != "" {
		conditions = append(conditions, "COALESCE(l.description, '') LIKE ?")
		args = append(args, filter.Title)
	}
	if filter.ZipCode != "" {
		conditions = append(conditions, "l.zip_code LIKE ?")
		args = append(args, filter.ZipCode)
	}
	if filter.City != "" {
		conditions = append(conditions, "l.city LIKE ?")
		args = append(args, filter.City)
	}
	if filter.Neighborhood != "" {
		conditions = append(conditions, "l.neighborhood LIKE ?")
		args = append(args, filter.Neighborhood)
	}
	if filter.UserID != nil {
		conditions = append(conditions, "l.user_id = ?")
		args = append(args, *filter.UserID)
	}
	if filter.CreatedFrom != nil {
		conditions = append(conditions, "l.created_at >= ?")
		args = append(args, *filter.CreatedFrom)
	}
	if filter.CreatedTo != nil {
		conditions = append(conditions, "l.created_at <= ?")
		args = append(args, *filter.CreatedTo)
	}
	if filter.MinSellPrice != nil {
		conditions = append(conditions, "COALESCE(l.sell_net, 0) >= ?")
		args = append(args, *filter.MinSellPrice)
	}
	if filter.MaxSellPrice != nil {
		conditions = append(conditions, "COALESCE(l.sell_net, 0) <= ?")
		args = append(args, *filter.MaxSellPrice)
	}
	if filter.MinRentPrice != nil {
		conditions = append(conditions, "COALESCE(l.rent_net, 0) >= ?")
		args = append(args, *filter.MinRentPrice)
	}
	if filter.MaxRentPrice != nil {
		conditions = append(conditions, "COALESCE(l.rent_net, 0) <= ?")
		args = append(args, *filter.MaxRentPrice)
	}
	if filter.MinLandSize != nil {
		conditions = append(conditions, "COALESCE(l.land_size, 0) >= ?")
		args = append(args, *filter.MinLandSize)
	}
	if filter.MaxLandSize != nil {
		conditions = append(conditions, "COALESCE(l.land_size, 0) <= ?")
		args = append(args, *filter.MaxLandSize)
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	baseSelect := `SELECT
		l.id,
		l.user_id,
		l.code,
		l.version,
		l.status,
		l.zip_code,
		l.street,
		l.number,
		l.complement,
		l.neighborhood,
		l.city,
		l.state,
		l.type,
		l.owner,
		l.land_size,
		l.corner,
		l.non_buildable,
		l.buildable,
		l.delivered,
		l.who_lives,
		l.description,
		l.transaction,
		l.sell_net,
		l.rent_net,
		l.condominium,
		l.annual_tax,
		l.annual_ground_rent,
		l.exchange,
		l.exchange_perc,
		l.installment,
		l.financing,
	l.visit,
	l.tenant_name,
	l.tenant_email,
	l.tenant_phone,
	l.accompanying,
	l.deleted,
	// l.created_at, // não existe no database
	// l.updated_at  // não existe no database
	FROM listings l`

	listQuery := baseSelect + " " + whereClause + " ORDER BY l.id DESC LIMIT ? OFFSET ?"
	listArgs := append([]any{}, args...)
	offset := (filter.Page - 1) * filter.Limit
	listArgs = append(listArgs, filter.Limit, offset)

	rows, readErr := tx.QueryContext(ctx, listQuery, listArgs...)
	if readErr != nil {
		utils.SetSpanError(ctx, readErr)
		logger.Error("mysql.listing.list.read_error", "error", readErr)
		return listingrepository.ListListingsResult{}, fmt.Errorf("list listings read: %w", readErr)
	}
	defer rows.Close()

	result := listingrepository.ListListingsResult{}

	for rows.Next() {
		entity := listingentity.ListingEntity{}
		scanErr := rows.Scan(
			&entity.ID,
			&entity.UserID,
			&entity.Code,
			&entity.Version,
			&entity.Status,
			&entity.ZipCode,
			&entity.Street,
			&entity.Number,
			&entity.Complement,
			&entity.Neighborhood,
			&entity.City,
			&entity.State,
			&entity.ListingType,
			&entity.Owner,
			&entity.LandSize,
			&entity.Corner,
			&entity.NonBuildable,
			&entity.Buildable,
			&entity.Delivered,
			&entity.WhoLives,
			&entity.Description,
			&entity.Transaction,
			&entity.SellNet,
			&entity.RentNet,
			&entity.Condominium,
			&entity.AnnualTax,
			&entity.AnnualGroundRent,
			&entity.Exchange,
			&entity.ExchangePercentual,
			&entity.Installment,
			&entity.Financing,
			&entity.Visit,
			&entity.TenantName,
			&entity.TenantEmail,
			&entity.TenantPhone,
			&entity.Accompanying,
			&entity.Deleted,
			// &entity.CreatedAt,
			// &entity.UpdatedAt,
		)
		if scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.listing.list.scan_error", "error", scanErr)
			return listingrepository.ListListingsResult{}, fmt.Errorf("scan listing row: %w", scanErr)
		}

		listing := listingconverters.ListingEntityToDomain(entity)
		if listing != nil {
			// var createdAt *time.Time
			// if entity.CreatedAt.Valid {
			// 	ts := entity.CreatedAt.Time
			// 	createdAt = &ts
			// }
			// var updatedAt *time.Time
			// if entity.UpdatedAt.Valid {
			// 	ts := entity.UpdatedAt.Time
			// 	updatedAt = &ts
			// }

			result.Records = append(result.Records, listingrepository.ListingRecord{
				Listing: listing,
				// CreatedAt: createdAt,
				// UpdatedAt: updatedAt,
			})
		}
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.list.rows_error", "error", err)
		return listingrepository.ListListingsResult{}, fmt.Errorf("iterate listing rows: %w", err)
	}

	countQuery := "SELECT COUNT(*) FROM listings l " + whereClause
	var total int64
	countErr := tx.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if countErr != nil {
		utils.SetSpanError(ctx, countErr)
		logger.Error("mysql.listing.list.count_error", "error", countErr)
		return listingrepository.ListListingsResult{}, fmt.Errorf("count listings: %w", countErr)
	}
	result.Total = total

	return result, nil
}
