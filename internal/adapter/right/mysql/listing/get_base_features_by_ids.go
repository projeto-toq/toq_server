package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	listingentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/entity"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetBaseFeaturesByIDs(
	ctx context.Context,
	tx *sql.Tx,
	ids []int64,
) (map[int64]listingmodel.BaseFeatureInterface, error) {
	result := make(map[int64]listingmodel.BaseFeatureInterface)

	if len(ids) == 0 {
		return result, nil
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	placeholders := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		placeholders = append(placeholders, "?")
		args = append(args, id)
	}

	query := fmt.Sprintf(
		"SELECT id, feature, description, priority FROM base_features WHERE id IN (%s)",
		strings.Join(placeholders, ","),
	)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_base_features_by_ids.prepare_error", "error", err)
		return nil, fmt.Errorf("prepare get base features by ids: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_base_features_by_ids.query_error", "error", err)
		return nil, fmt.Errorf("query get base features by ids: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		entity := listingentity.EntityBaseFeature{}
		if err := rows.Scan(&entity.ID, &entity.Feature, &entity.Description, &entity.Priority); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.get_base_features_by_ids.scan_error", "error", err)
			return nil, fmt.Errorf("scan base feature by ids: %w", err)
		}

		feature := listingmodel.NewBaseFeature()
		feature.SetID(entity.ID)
		feature.SetFeature(entity.Feature)
		feature.SetDescription(entity.Description)
		feature.SetPriority(entity.Priority)

		result[entity.ID] = feature
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_base_features_by_ids.rows_error", "error", err)
		return nil, fmt.Errorf("rows iteration for base features by ids: %w", err)
	}

	return result, nil
}
