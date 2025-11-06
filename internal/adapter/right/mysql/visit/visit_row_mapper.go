package mysqlvisitadapter

import "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/entities"

// rowScanner defines the interface for scanning database rows
//
// This interface abstracts both sql.Row and sql.Rows, allowing scanVisitEntity
// to work with QueryRowContext (single row) and QueryContext (multiple rows).
//
// Implementations:
//   - *sql.Row: From QueryRowContext
//   - *sql.Rows: From QueryContext (within iteration loop)
type rowScanner interface {
	Scan(dest ...any) error
}

// scanVisitEntity scans a database row into a VisitEntity struct
//
// This function handles the mapping from database columns to struct fields,
// maintaining strict column order matching the SELECT queries in repository methods.
//
// Column Order (MUST match all SELECT queries):
//  1. id              → VisitEntity.ID
//  2. listing_id      → VisitEntity.ListingID
//  3. owner_id        → VisitEntity.OwnerID
//  4. realtor_id      → VisitEntity.RealtorID
//  5. scheduled_start → VisitEntity.ScheduledStart
//  6. scheduled_end   → VisitEntity.ScheduledEnd
//  7. status          → VisitEntity.Status
//  8. cancel_reason   → VisitEntity.CancelReason (sql.NullString)
//  9. notes           → VisitEntity.Notes (sql.NullString)
//
// 10. created_by      → VisitEntity.CreatedBy
// 11. updated_by      → VisitEntity.UpdatedBy (sql.NullInt64)
//
// Parameters:
//   - scanner: rowScanner interface (sql.Row or sql.Rows)
//
// Returns:
//   - entity: VisitEntity with all fields populated from database
//   - error: Scan errors (type mismatch, null constraint violation, etc.)
//
// Error Scenarios:
//   - sql.ErrNoRows: No row available (only from sql.Row, not sql.Rows)
//   - Type mismatch: Database type incompatible with struct field
//   - Column count mismatch: SELECT query columns ≠ Scan arguments
//
// Important:
//   - Column order MUST be synchronized with all SELECT queries
//   - Changing column order requires updating ALL queries in this adapter
//   - Used by: GetVisitByID (single row), ListVisits (multiple rows)
func scanVisitEntity(scanner rowScanner) (entities.VisitEntity, error) {
	var visit entities.VisitEntity

	// Scan all 11 columns in exact order matching SELECT queries
	if err := scanner.Scan(
		&visit.ID,
		&visit.ListingID,
		&visit.OwnerID,
		&visit.RealtorID,
		&visit.ScheduledStart,
		&visit.ScheduledEnd,
		&visit.Status,
		&visit.CancelReason, // sql.NullString handles NULL
		&visit.Notes,        // sql.NullString handles NULL
		&visit.CreatedBy,
		&visit.UpdatedBy, // sql.NullInt64 handles NULL
	); err != nil {
		return entities.VisitEntity{}, err
	}

	return visit, nil
}
