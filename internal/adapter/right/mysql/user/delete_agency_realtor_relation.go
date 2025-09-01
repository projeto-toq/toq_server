package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) DeleteAgencyRealtorRelation(ctx context.Context, tx *sql.Tx, agencyID int64, realtorID int64) (deleted int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `DELETE FROM realtors_agency WHERE realtor_id = ? AND agency_id = ?;`

	deleted, err = ua.Delete(ctx, tx, query, realtorID, agencyID)
	if err != nil {
		slog.Error("mysqluseradapter/DeleteAgencyRealtorRelation: error executing Delete", "error", err)
		return 0, fmt.Errorf("delete realtor-agency relation: %w", err)
	}

	if deleted == 0 {
		return 0, sql.ErrNoRows
	}

	return
}
