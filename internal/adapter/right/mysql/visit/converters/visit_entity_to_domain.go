package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/entities"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// ToVisitModel converts a database VisitEntity to a domain VisitInterface
//
// This converter handles the translation from database-specific types (sql.Null*)
// to clean domain types, ensuring the core layer remains decoupled from database concerns.
//
// Conversion Rules:
//   - sql.NullString → string (empty string if NULL)
//   - sql.NullInt64 → int64 (0 if NULL, checked via Valid flag)
//   - status string → VisitStatus enum type
//
// Parameters:
//   - e: VisitEntity from database query with all fields populated
//
// Returns:
//   - visit: VisitInterface with all fields converted to domain types
//
// NULL Field Handling:
//   - CancelReason: Only set if Valid=true (optional field)
//   - Notes: Only set if Valid=true (optional field)
//   - UpdatedBy: Only set if Valid=true (optional audit field)
func ToVisitModel(e entities.VisitEntity) listingmodel.VisitInterface {
	visit := listingmodel.NewVisit()

	// Map mandatory fields (NOT NULL in schema)
	visit.SetID(e.ID)
	visit.SetListingIdentityID(e.ListingIdentityID)
	visit.SetListingVersion(e.ListingVersion)
	visit.SetRequesterUserID(e.RequesterUserID)
	visit.SetOwnerUserID(e.OwnerUserID)
	visit.SetScheduledStart(e.ScheduledStart)
	visit.SetScheduledEnd(e.ScheduledEnd)
	visit.SetDurationMinutes(e.DurationMinutes)
	visit.SetStatus(listingmodel.VisitStatus(e.Status))
	visit.SetType(listingmodel.VisitMode(e.Type))

	if e.Source.Valid {
		visit.SetSource(e.Source.String)
	}

	// Map optional fields (NULL in schema) - check Valid before accessing
	if e.Source.Valid {
		visit.SetSource(e.Source.String)
	}

	if e.RealtorNotes.Valid {
		visit.SetRealtorNotes(e.RealtorNotes.String)
	}

	if e.OwnerNotes.Valid {
		visit.SetOwnerNotes(e.OwnerNotes.String)
	}

	if e.RejectionReason.Valid {
		visit.SetRejectionReason(e.RejectionReason.String)
	}

	if e.CancelReason.Valid {
		visit.SetCancelReason(e.CancelReason.String)
	}

	if e.FirstOwnerActionAt.Valid {
		visit.SetFirstOwnerActionAt(e.FirstOwnerActionAt.Time)
	}

	return visit
}
