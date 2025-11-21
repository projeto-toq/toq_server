package mysqllistingadapter

import (
	"context"
	"database/sql"
	"strings"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// CheckDuplicity checks if a listing exists based on dynamic criteria.
// It returns true if a duplicate is found, false otherwise.
func (la *ListingAdapter) CheckDuplicity(ctx context.Context, tx *sql.Tx, criteria listingmodel.DuplicityCriteria) (bool, error) {
	var queryBuilder strings.Builder
	var args []interface{}

	// Base query: Active listings (not deleted, status not in Closed/Expired/Archived)
	// Status codes: Closed=13, Expired=15, Archived=16
	// Note: StatusDraft (1) is NOW INCLUDED in the check.
	queryBuilder.WriteString(`
		SELECT 1 
		FROM listing_versions lv 
		JOIN listing_identities li ON lv.listing_identity_id = li.id 
		WHERE lv.zip_code = ? 
		  AND lv.number = ? 
		  AND lv.deleted = 0 
		  AND li.deleted = 0 
		  AND lv.status NOT IN (13, 15, 16)
	`)
	args = append(args, criteria.ZipCode, criteria.Number)

	// Dynamic clauses based on non-nil criteria fields
	if criteria.UnitTower != nil {
		queryBuilder.WriteString(" AND lv.unit_tower = ?")
		args = append(args, *criteria.UnitTower)
	}

	if criteria.UnitFloor != nil {
		queryBuilder.WriteString(" AND lv.unit_floor = ?")
		args = append(args, *criteria.UnitFloor)
	}

	if criteria.UnitNumber != nil {
		queryBuilder.WriteString(" AND lv.unit_number = ?")
		args = append(args, *criteria.UnitNumber)
	}

	if criteria.LandBlock != nil {
		queryBuilder.WriteString(" AND lv.land_block = ?")
		args = append(args, *criteria.LandBlock)
	}

	if criteria.LandLot != nil {
		queryBuilder.WriteString(" AND lv.land_lot = ?")
		args = append(args, *criteria.LandLot)
	}

	queryBuilder.WriteString(" LIMIT 1")

	var exists int
	err := la.QueryRowContext(ctx, tx, "check_duplicity", queryBuilder.String(), args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
