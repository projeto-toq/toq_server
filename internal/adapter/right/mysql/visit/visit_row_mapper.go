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
//  1. id                   → VisitEntity.ID
//  2. listing_identity_id  → VisitEntity.ListingIdentityID
//  3. listing_version      → VisitEntity.ListingVersion
//  4. user_id              → VisitEntity.RequesterUserID
//  5. owner_user_id        → VisitEntity.OwnerUserID
//  6. scheduled_start      → VisitEntity.ScheduledStart
//  7. scheduled_end        → VisitEntity.ScheduledEnd
//  8. duration_minutes     → VisitEntity.DurationMinutes
//  9. status               → VisitEntity.Status
//
// 10. type                 → VisitEntity.Type
// 11. source               → VisitEntity.Source (sql.NullString)
// 12. realtor_notes        → VisitEntity.RealtorNotes (sql.NullString)
// 13. owner_notes          → VisitEntity.OwnerNotes (sql.NullString)
// 14. rejection_reason     → VisitEntity.RejectionReason (sql.NullString)
// 15. cancel_reason        → VisitEntity.CancelReason (sql.NullString)
// 16. first_owner_action_at→ VisitEntity.FirstOwnerActionAt (sql.NullTime)
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

	// Scan all columns in exact order matching SELECT queries
	if err := scanner.Scan(
		&visit.ID,
		&visit.ListingIdentityID,
		&visit.ListingVersion,
		&visit.RequesterUserID,
		&visit.OwnerUserID,
		&visit.ScheduledStart,
		&visit.ScheduledEnd,
		&visit.DurationMinutes,
		&visit.Status,
		&visit.Type,
		&visit.Source,
		&visit.RealtorNotes,
		&visit.OwnerNotes,
		&visit.RejectionReason,
		&visit.CancelReason,
		&visit.FirstOwnerActionAt,
	); err != nil {
		return entities.VisitEntity{}, err
	}

	return visit, nil
}
