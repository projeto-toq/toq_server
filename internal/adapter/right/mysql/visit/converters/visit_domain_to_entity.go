package converters

import (
	"database/sql"
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/entities"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// ToVisitEntity converts a domain VisitInterface to a database VisitEntity.
//
// Conversion rules:
//   - string -> sql.NullString (Valid=true when non-empty)
//   - VisitStatus -> string (ENUM persisted)
//   - RequestedAt defaults to time.Now().UTC() when zero (aligns with DB default)
//   - Source falls back to "APP" when absent to satisfy NOT NULL/ENUM default
//
// Parameters:
//   - model: VisitInterface from core layer
//
// Returns:
//   - VisitEntity ready for INSERT/UPDATE in listing_visits
func ToVisitEntity(model listingmodel.VisitInterface) entities.VisitEntity {
	// Map mandatory fields directly
	entity := entities.VisitEntity{
		ID:                model.ID(),
		ListingIdentityID: model.ListingIdentityID(),
		ListingVersion:    model.ListingVersion(),
		RequesterUserID:   model.RequesterUserID(),
		OwnerUserID:       model.OwnerUserID(),
		ScheduledStart:    model.ScheduledStart(),
		ScheduledEnd:      model.ScheduledEnd(),
		Status:            string(model.Status()),
	}

	if value, ok := model.Source(); ok {
		entity.Source = sql.NullString{String: value, Valid: true}
	}

	// Ensure source is always populated to satisfy NOT NULL constraint (fallback to APP)
	if !entity.Source.Valid {
		entity.Source = sql.NullString{String: "APP", Valid: true}
	}

	if value, ok := model.Notes(); ok {
		entity.Notes = sql.NullString{String: value, Valid: true}
	}

	// Map optional fields - convert to sql.Null* with Valid based on value presence
	if value, ok := model.RejectionReason(); ok {
		entity.RejectionReason = sql.NullString{String: value, Valid: true}
	}

	if value, ok := model.FirstOwnerActionAt(); ok {
		entity.FirstOwnerActionAt = sql.NullTime{Time: value, Valid: true}
	}

	entity.RequestedAt = model.RequestedAt()
	if entity.RequestedAt.IsZero() {
		entity.RequestedAt = time.Now().UTC()
	}

	return entity
}
