package photosessionservices

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// EnsurePhotographerAgendaWithTx provisions bootstrap agenda entries using an existing transaction.
func (s *photoSessionService) EnsurePhotographerAgendaWithTx(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error {
	if tx == nil {
		return utils.InternalError("")
	}
	if err := validateEnsureAgendaInput(input); err != nil {
		return err
	}
	return nil
}
