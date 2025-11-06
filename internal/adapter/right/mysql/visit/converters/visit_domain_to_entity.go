package converters

import (
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/entities"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// ToVisitEntity converts a domain VisitInterface to a database VisitEntity
//
// This converter handles the translation from clean domain types to database-specific
// types (sql.Null*), preparing data for database insertion/update operations.
//
// Conversion Rules:
//   - string → sql.NullString (Valid=true if non-empty)
//   - int64 → sql.NullInt64 (Valid=true if value present)
//   - VisitStatus enum → string (stored as ENUM in database)
//
// Parameters:
//   - model: VisitInterface from core layer with all required fields populated
//
// Returns:
//   - entity: VisitEntity ready for database operations (INSERT/UPDATE)
//
// Important:
//   - ID may be 0 for new records (populated by AUTO_INCREMENT)
//   - Empty strings are converted to NULL for optional fields
//   - Optional fields use (value, ok) pattern to check presence
func ToVisitEntity(model listingmodel.VisitInterface) entities.VisitEntity {
	// Map mandatory fields directly
	entity := entities.VisitEntity{
		ID:             model.ID(),
		ListingID:      model.ListingID(),
		OwnerID:        model.OwnerID(),
		RealtorID:      model.RealtorID(),
		ScheduledStart: model.ScheduledStart(),
		ScheduledEnd:   model.ScheduledEnd(),
		Status:         string(model.Status()),
		CreatedBy:      model.CreatedBy(),
	}

	// Map optional fields - convert to sql.Null* with Valid based on value presence
	if value, ok := model.CancelReason(); ok {
		entity.CancelReason = sql.NullString{String: value, Valid: true}
	}

	if value, ok := model.Notes(); ok {
		entity.Notes = sql.NullString{String: value, Valid: true}
	}

	if value, ok := model.UpdatedBy(); ok {
		entity.UpdatedBy = sql.NullInt64{Int64: value, Valid: true}
	}

	return entity
}
