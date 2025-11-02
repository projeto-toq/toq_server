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

	query := `SELECT id, category, numeric_value, slug, label, description, is_active FROM listing_catalog_values WHERE category = ?`
	args := []any{category}
	if !includeInactive {
		query += ` AND is_active = 1`
	}
	query += ` ORDER BY numeric_value`

	rows, queryErr := la.QueryContext(ctx, tx, "select", query, args...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.listing.catalog.list.query_error", "error", queryErr, "category", category)
		return nil, queryErr
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

	query := `SELECT id, category, numeric_value, slug, label, description, is_active FROM listing_catalog_values WHERE category = ? AND id = ?`

	row := la.QueryRowContext(ctx, tx, "select", query, category, id)
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

	query := `SELECT id, category, numeric_value, slug, label, description, is_active FROM listing_catalog_values WHERE category = ? AND slug = ?`

	row := la.QueryRowContext(ctx, tx, "select", query, category, slug)
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

func (la *ListingAdapter) GetCatalogValueByNumeric(ctx context.Context, tx *sql.Tx, category string, numericValue uint8) (listingmodel.CatalogValueInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, category, numeric_value, slug, label, description, is_active FROM listing_catalog_values WHERE category = ? AND numeric_value = ?`

	row := la.QueryRowContext(ctx, tx, "select", query, category, numericValue)
	value, scanErr := scanCatalogValueRow(row)
	if scanErr != nil {
		if scanErr != sql.ErrNoRows {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.listing.catalog.get_by_numeric.scan_error", "error", scanErr, "category", category, "numeric_value", numericValue)
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

	row := la.QueryRowContext(ctx, tx, "select", query, category)

	var nextID uint8
	if err := row.Scan(&nextID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.next_id.scan_error", "error", err, "category", category)
		return 0, err
	}
	return nextID, nil
}

func (la *ListingAdapter) GetNextCatalogNumericValue(ctx context.Context, tx *sql.Tx, category string) (uint8, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT IFNULL(MAX(numeric_value), 0) + 1 FROM listing_catalog_values WHERE category = ?`

	row := la.QueryRowContext(ctx, tx, "select", query, category)

	var nextNumeric uint8
	if err := row.Scan(&nextNumeric); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.catalog.next_numeric.scan_error", "error", err, "category", category)
		return 0, err
	}

	return nextNumeric, nil
}

func (la *ListingAdapter) CreateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO listing_catalog_values (category, id, numeric_value, slug, label, description, is_active) VALUES (?, ?, ?, ?, ?, ?, ?)`

	var description sql.NullString
	if desc := value.Description(); desc != nil {
		description.Valid = true
		description.String = *desc
	}

	_, execErr := la.ExecContext(ctx, tx, "insert", query, value.Category(), value.ID(), value.NumericValue(), value.Slug(), value.Label(), description, value.IsActive())
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.catalog.create.exec_error", "error", execErr, "category", value.Category(), "id", value.ID())
		return execErr
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

	var description sql.NullString
	if desc := value.Description(); desc != nil {
		description.Valid = true
		description.String = *desc
	}

	result, execErr := la.ExecContext(ctx, tx, "update", query, value.Slug(), value.Label(), description, value.IsActive(), value.Category(), value.ID())
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.catalog.update.exec_error", "error", execErr, "category", value.Category(), "id", value.ID())
		return execErr
	}

	rows, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.listing.catalog.update.rows_affected_error", "error", rowsErr, "category", value.Category(), "id", value.ID())
		return rowsErr
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

	result, execErr := la.ExecContext(ctx, tx, "update", query, category, id)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.catalog.delete.exec_error", "error", execErr, "category", category, "id", id)
		return execErr
	}

	rows, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.listing.catalog.delete.rows_affected_error", "error", rowsErr, "category", category, "id", id)
		return rowsErr
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func scanCatalogValue(rows *sql.Rows) (listingmodel.CatalogValueInterface, error) {
	var (
		id           uint8
		category     string
		numericValue uint8
		slug         string
		label        string
		description  sql.NullString
		isActive     bool
	)

	if err := rows.Scan(&id, &category, &numericValue, &slug, &label, &description, &isActive); err != nil {
		return nil, err
	}

	return buildCatalogValue(id, numericValue, category, slug, label, description, isActive), nil
}

func scanCatalogValueRow(row *sql.Row) (listingmodel.CatalogValueInterface, error) {
	var (
		id           uint8
		category     string
		numericValue uint8
		slug         string
		label        string
		description  sql.NullString
		isActive     bool
	)

	if err := row.Scan(&id, &category, &numericValue, &slug, &label, &description, &isActive); err != nil {
		return nil, err
	}

	return buildCatalogValue(id, numericValue, category, slug, label, description, isActive), nil
}

func buildCatalogValue(id, numericValue uint8, category, slug, label string, description sql.NullString, isActive bool) listingmodel.CatalogValueInterface {
	value := listingmodel.NewCatalogValue()
	value.SetID(id)
	value.SetNumericValue(numericValue)
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
