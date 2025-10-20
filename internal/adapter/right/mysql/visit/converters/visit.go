package converters

import (
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/entity"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// ToVisitModel converts VisitEntity to domain interface.
func ToVisitModel(e entity.VisitEntity) listingmodel.VisitInterface {
	visit := listingmodel.NewVisit()
	visit.SetID(e.ID)
	visit.SetListingID(e.ListingID)
	visit.SetOwnerID(e.OwnerID)
	visit.SetRealtorID(e.RealtorID)
	visit.SetScheduledStart(e.ScheduledStart)
	visit.SetScheduledEnd(e.ScheduledEnd)
	visit.SetStatus(listingmodel.VisitStatus(e.Status))
	if e.CancelReason.Valid {
		visit.SetCancelReason(e.CancelReason.String)
	}
	if e.Notes.Valid {
		visit.SetNotes(e.Notes.String)
	}
	visit.SetCreatedBy(e.CreatedBy)
	if e.UpdatedBy.Valid {
		visit.SetUpdatedBy(e.UpdatedBy.Int64)
	}
	return visit
}

// ToVisitEntity converts domain visit to persistence shape.
func ToVisitEntity(model listingmodel.VisitInterface) entity.VisitEntity {
	var cancelReason sql.NullString
	if value, ok := model.CancelReason(); ok {
		cancelReason = sql.NullString{String: value, Valid: true}
	}

	var notes sql.NullString
	if value, ok := model.Notes(); ok {
		notes = sql.NullString{String: value, Valid: true}
	}

	var updatedBy sql.NullInt64
	if value, ok := model.UpdatedBy(); ok {
		updatedBy = sql.NullInt64{Int64: value, Valid: true}
	}

	return entity.VisitEntity{
		ID:             model.ID(),
		ListingID:      model.ListingID(),
		OwnerID:        model.OwnerID(),
		RealtorID:      model.RealtorID(),
		ScheduledStart: model.ScheduledStart(),
		ScheduledEnd:   model.ScheduledEnd(),
		Status:         string(model.Status()),
		CancelReason:   cancelReason,
		Notes:          notes,
		CreatedBy:      model.CreatedBy(),
		UpdatedBy:      updatedBy,
	}
}
