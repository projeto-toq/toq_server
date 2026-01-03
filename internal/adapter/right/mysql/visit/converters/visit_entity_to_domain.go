package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/entities"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// ToVisitModel converts a database VisitEntity to a domain VisitInterface.
//
// Conversion rules:
//   - sql.NullString -> optional setters only when Valid
//   - status string -> VisitStatus enum
//   - RequestedAt propagated only when non-zero (db default already set)
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
	visit.SetStatus(listingmodel.VisitStatus(e.Status))

	if e.Source.Valid {
		visit.SetSource(e.Source.String)
	}

	// Map optional fields (NULL in schema) - check Valid before accessing
	if e.Notes.Valid {
		visit.SetNotes(e.Notes.String)
	}

	if e.RejectionReason.Valid {
		visit.SetRejectionReason(e.RejectionReason.String)
	}

	if e.FirstOwnerActionAt.Valid {
		visit.SetFirstOwnerActionAt(e.FirstOwnerActionAt.Time)
	}

	if !e.RequestedAt.IsZero() {
		visit.SetRequestedAt(e.RequestedAt)
	}

	return visit
}
