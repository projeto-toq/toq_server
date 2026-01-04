package converters

import (
	"database/sql"
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// ToAgendaEntryEntity maps a domain agenda entry to its DB representation, preserving NULL semantics
// for optional columns (source, source_id, starts_at, ends_at, reason, timezone). Empty timezone is
// coerced to the database default America/Sao_Paulo.
func ToAgendaEntryEntity(entry photosessionmodel.AgendaEntryInterface) entity.AgendaEntry {
	var source sql.NullString
	if s := string(entry.Source()); s != "" {
		source = sql.NullString{String: s, Valid: true}
	}

	var sourceID sql.NullInt64
	if id, ok := entry.SourceID(); ok && id != nil {
		sourceID = sql.NullInt64{Int64: int64(*id), Valid: true}
	}

	reasonStr, reasonValid := entry.Reason()
	reason := sql.NullString{String: reasonStr, Valid: reasonValid}

	startsAt := sql.NullTime{}
	if ts := entry.StartsAt(); !ts.IsZero() {
		startsAt = sql.NullTime{Time: ts, Valid: true}
	}

	endsAt := sql.NullTime{}
	if ts := entry.EndsAt(); !ts.IsZero() {
		endsAt = sql.NullTime{Time: ts, Valid: true}
	}

	timezone := entry.Timezone()
	if timezone == "" {
		timezone = "America/Sao_Paulo"
	}

	return entity.AgendaEntry{
		ID:                 entry.ID(),
		PhotographerUserID: entry.PhotographerUserID(),
		EntryType:          string(entry.EntryType()),
		Source:             source,
		SourceID:           sourceID,
		StartsAt:           startsAt,
		EndsAt:             endsAt,
		Blocking:           sql.NullBool{Bool: entry.Blocking(), Valid: true},
		Reason:             reason,
		Timezone:           sql.NullString{String: timezone, Valid: true},
	}
}

// ToAgendaEntryModel converts a database entity into a domain model, applying defaults for NULL values
// (blocking defaults to true when NULL; timezone defaults to America/Sao_Paulo; timestamps remain zero when NULL).
func ToAgendaEntryModel(entity entity.AgendaEntry) photosessionmodel.AgendaEntryInterface {
	model := photosessionmodel.NewAgendaEntry()
	model.SetID(entity.ID)
	model.SetPhotographerUserID(entity.PhotographerUserID)
	model.SetEntryType(photosessionmodel.AgendaEntryType(entity.EntryType))

	if entity.Source.Valid {
		model.SetSource(photosessionmodel.AgendaEntrySource(entity.Source.String))
	} else {
		model.SetSource("")
	}

	if entity.SourceID.Valid {
		model.SetSourceID(uint64(entity.SourceID.Int64))
	} else {
		model.ClearSourceID()
	}

	if entity.StartsAt.Valid {
		model.SetStartsAt(entity.StartsAt.Time)
	} else {
		model.SetStartsAt(time.Time{})
	}

	if entity.EndsAt.Valid {
		model.SetEndsAt(entity.EndsAt.Time)
	} else {
		model.SetEndsAt(time.Time{})
	}

	if entity.Blocking.Valid {
		model.SetBlocking(entity.Blocking.Bool)
	} else {
		model.SetBlocking(true)
	}

	if entity.Reason.Valid {
		model.SetReason(entity.Reason.String)
	} else {
		model.ClearReason()
	}

	if entity.Timezone.Valid && entity.Timezone.String != "" {
		model.SetTimezone(entity.Timezone.String)
	} else {
		model.SetTimezone("America/Sao_Paulo")
	}

	return model
}
