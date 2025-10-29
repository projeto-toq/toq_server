package converters

import (
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// ToAgendaEntryEntity maps a domain agenda entry to its DB representation.
func ToAgendaEntryEntity(entry photosessionmodel.AgendaEntryInterface) entity.AgendaEntry {
	var sourceID sql.NullInt64
	if id, ok := entry.SourceID(); ok && id != nil {
		sourceID = sql.NullInt64{Int64: int64(*id), Valid: true}
	}

	reasonStr, reasonValid := entry.Reason()
	reason := sql.NullString{String: reasonStr, Valid: reasonValid}

	createdAt := sql.NullTime{}
	if t, ok := entry.CreatedAt(); ok {
		createdAt = sql.NullTime{Time: t, Valid: true}
	}

	updatedAt := sql.NullTime{}
	if t, ok := entry.UpdatedAt(); ok {
		updatedAt = sql.NullTime{Time: t, Valid: true}
	}

	tz := entry.Timezone()
	if tz == "" {
		tz = "America/Sao_Paulo"
	}

	return entity.AgendaEntry{
		ID:                 entry.ID(),
		PhotographerUserID: entry.PhotographerUserID(),
		EntryType:          string(entry.EntryType()),
		Source:             string(entry.Source()),
		SourceID:           sourceID,
		StartsAt:           entry.StartsAt(),
		EndsAt:             entry.EndsAt(),
		Blocking:           entry.Blocking(),
		Reason:             reason,
		Timezone:           tz,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}
}

// ToAgendaEntryModel converts a database entity into a domain model instance.
func ToAgendaEntryModel(entity entity.AgendaEntry) photosessionmodel.AgendaEntryInterface {
	model := photosessionmodel.NewAgendaEntry()
	model.SetID(entity.ID)
	model.SetPhotographerUserID(entity.PhotographerUserID)
	model.SetEntryType(photosessionmodel.AgendaEntryType(entity.EntryType))
	model.SetSource(photosessionmodel.AgendaEntrySource(entity.Source))
	if entity.SourceID.Valid {
		model.SetSourceID(uint64(entity.SourceID.Int64))
	} else {
		model.ClearSourceID()
	}
	model.SetStartsAt(entity.StartsAt)
	model.SetEndsAt(entity.EndsAt)
	model.SetBlocking(entity.Blocking)
	if entity.Reason.Valid {
		model.SetReason(entity.Reason.String)
	} else {
		model.ClearReason()
	}
	model.SetTimezone(entity.Timezone)
	if entity.CreatedAt.Valid {
		model.SetCreatedAt(entity.CreatedAt.Time)
	}
	if entity.UpdatedAt.Valid {
		model.SetUpdatedAt(entity.UpdatedAt.Time)
	}
	return model
}
