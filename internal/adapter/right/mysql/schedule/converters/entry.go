package converters

import (
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entity"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
)

// ToEntryModel converts an EntryEntity to the domain object.
func ToEntryModel(e entity.EntryEntity) schedulemodel.AgendaEntryInterface {
	entry := schedulemodel.NewAgendaEntry()
	entry.SetID(e.ID)
	entry.SetAgendaID(e.AgendaID)
	entry.SetEntryType(schedulemodel.EntryType(e.EntryType))
	entry.SetStartsAt(e.StartsAt)
	entry.SetEndsAt(e.EndsAt)
	entry.SetBlocking(e.Blocking)
	if e.Reason.Valid {
		entry.SetReason(e.Reason.String)
	}
	if e.VisitID.Valid {
		entry.SetVisitID(uint64(e.VisitID.Int64))
	}
	if e.PhotoBookingID.Valid {
		entry.SetPhotoBookingID(uint64(e.PhotoBookingID.Int64))
	}
	return entry
}

// ToEntryEntity converts the domain entry into persistence shape.
func ToEntryEntity(model schedulemodel.AgendaEntryInterface) entity.EntryEntity {
	var reason sql.NullString
	if value, ok := model.Reason(); ok {
		reason = sql.NullString{String: value, Valid: true}
	}

	var visitID sql.NullInt64
	if value, ok := model.VisitID(); ok {
		visitID = sql.NullInt64{Int64: int64(value), Valid: true}
	}

	var photoID sql.NullInt64
	if value, ok := model.PhotoBookingID(); ok {
		photoID = sql.NullInt64{Int64: int64(value), Valid: true}
	}

	return entity.EntryEntity{
		ID:             model.ID(),
		AgendaID:       model.AgendaID(),
		EntryType:      string(model.EntryType()),
		StartsAt:       model.StartsAt(),
		EndsAt:         model.EndsAt(),
		Blocking:       model.Blocking(),
		Reason:         reason,
		VisitID:        visitID,
		PhotoBookingID: photoID,
	}
}
