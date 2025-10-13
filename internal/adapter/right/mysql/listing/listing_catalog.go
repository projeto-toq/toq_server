package mysqllistingadapter

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) ListCatalogValues(ctx context.Context, tx *sql.Tx, category string, includeInactive bool) ([]listingmodel.CatalogValueInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, category, slug, label, description, is_active FROM listing_catalog_values WHERE category = ?`
	args := []any{category}
	if !includeInactive {
		query += ` AND is_active = 1`
	}
	query += ` ORDER BY id`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.list.prepare_error", "error", err, "category", category)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.list.query_error", "error", err, "category", category)
		return nil, err
	}
	defer rows.Close()

	values := make([]listingmodel.CatalogValueInterface, 0)
	for rows.Next() {
		value, scanErr := scanCatalogValue(rows)
		if scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.listing.catalog.list.scan_error", "error", scanErr, "category", category)
			return nil, scanErr
		}
		values = append(values, value)
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.list.rows_error", "error", err, "category", category)
		return nil, err
	}

	return values, nil
}

func (la *ListingAdapter) GetCatalogValueByID(ctx context.Context, tx *sql.Tx, category string, id uint8) (listingmodel.CatalogValueInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, category, slug, label, description, is_active FROM listing_catalog_values WHERE category = ? AND id = ?`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.get_by_id.prepare_error", "error", err, "category", category, "id", id)
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, category, id)
	value, scanErr := scanCatalogValueRow(row)
	if scanErr != nil {
		if scanErr != sql.ErrNoRows {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.listing.catalog.get_by_id.scan_error", "error", scanErr, "category", category, "id", id)
		}
		return nil, scanErr
	}
	return value, nil
}

func (la *ListingAdapter) GetCatalogValueBySlug(ctx context.Context, tx *sql.Tx, category, slug string) (listingmodel.CatalogValueInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, category, slug, label, description, is_active FROM listing_catalog_values WHERE category = ? AND slug = ?`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.get_by_slug.prepare_error", "error", err, "category", category, "slug", slug)
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, category, slug)
	value, scanErr := scanCatalogValueRow(row)
	if scanErr != nil {
		if scanErr != sql.ErrNoRows {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.listing.catalog.get_by_slug.scan_error", "error", scanErr, "category", category, "slug", slug)
		}
		return nil, scanErr
	}
	return value, nil
}

func (la *ListingAdapter) GetNextCatalogValueID(ctx context.Context, tx *sql.Tx, category string) (uint8, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT IFNULL(MAX(id), 0) + 1 FROM listing_catalog_values WHERE category = ?`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.next_id.prepare_error", "error", err, "category", category)
		return 0, err
	}
	defer stmt.Close()

	var nextID uint8
	if err := stmt.QueryRowContext(ctx, category).Scan(&nextID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.next_id.scan_error", "error", err, "category", category)
		return 0, err
	}
	return nextID, nil
}

func (la *ListingAdapter) CreateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO listing_catalog_values (category, id, slug, label, description, is_active) VALUES (?, ?, ?, ?, ?, ?)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.create.prepare_error", "error", err, "category", value.Category(), "id", value.ID())
		return err
	}
	defer stmt.Close()

	var description sql.NullString
	if desc := value.Description(); desc != nil {
		description.Valid = true
		description.String = *desc
	}

	if _, err := stmt.ExecContext(ctx, value.Category(), value.ID(), value.Slug(), value.Label(), description, value.IsActive()); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.create.exec_error", "error", err, "category", value.Category(), "id", value.ID())
		return err
	}

	return nil
}

func (la *ListingAdapter) UpdateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE listing_catalog_values SET slug = ?, label = ?, description = ?, is_active = ? WHERE category = ? AND id = ?`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.update.prepare_error", "error", err, "category", value.Category(), "id", value.ID())
		return err
	}
	defer stmt.Close()

	var description sql.NullString
	if desc := value.Description(); desc != nil {
		description.Valid = true
		description.String = *desc
	}

	result, err := stmt.ExecContext(ctx, value.Slug(), value.Label(), description, value.IsActive(), value.Category(), value.ID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.update.exec_error", "error", err, "category", value.Category(), "id", value.ID())
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.update.rows_affected_error", "error", err, "category", value.Category(), "id", value.ID())
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (la *ListingAdapter) SoftDeleteCatalogValue(ctx context.Context, tx *sql.Tx, category string, id uint8) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE listing_catalog_values SET is_active = 0 WHERE category = ? AND id = ?`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.delete.prepare_error", "error", err, "category", category, "id", id)
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, category, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.delete.exec_error", "error", err, "category", category, "id", id)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.delete.rows_affected_error", "error", err, "category", category, "id", id)
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func scanCatalogValue(rows *sql.Rows) (listingmodel.CatalogValueInterface, error) {
	var (
		id          uint8
		category    string
		slug        string
		label       string
		description sql.NullString
		isActive    bool
	)

	if err := rows.Scan(&id, &category, &slug, &label, &description, &isActive); err != nil {
		return nil, err
	}

	return buildCatalogValue(id, category, slug, label, description, isActive), nil
}

func scanCatalogValueRow(row *sql.Row) (listingmodel.CatalogValueInterface, error) {
	var (
		id          uint8
		category    string
		slug        string
		label       string
		description sql.NullString
		isActive    bool
	)

	if err := row.Scan(&id, &category, &slug, &label, &description, &isActive); err != nil {
		return nil, err
	}

	return buildCatalogValue(id, category, slug, label, description, isActive), nil
}

func buildCatalogValue(id uint8, category, slug, label string, description sql.NullString, isActive bool) listingmodel.CatalogValueInterface {
	value := listingmodel.NewCatalogValue()
	value.SetID(id)
	value.SetCategory(category)
	value.SetSlug(slug)
	value.SetLabel(label)
	if description.Valid {
		desc := description.String
		value.SetDescription(&desc)
	} else {
		value.SetDescription(nil)
	}
	value.SetIsActive(isActive)
	return value
}
