package mysqlphotosessionadapter

import (
	"context"
	"database/sql"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// GetEntryByIDForUpdate fetches an agenda entry with FOR UPDATE locking; requires non-nil transaction.
func (a *PhotoSessionAdapter) GetEntryByIDForUpdate(ctx context.Context, tx *sql.Tx, entryID uint64) (photosessionmodel.AgendaEntryInterface, error) {
	return a.getEntry(ctx, tx, entryID, true)
}
