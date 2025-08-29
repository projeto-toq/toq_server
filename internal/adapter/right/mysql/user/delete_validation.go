package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) DeleteValidation(ctx context.Context, tx *sql.Tx, id int64) (deleted int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `DELETE FROM temp_user_validations WHERE user_id = ?;`

	deleted, err = ua.Delete(ctx, tx, query, id)
	if err != nil {
		slog.Error("mysqluseradapter/DeleteValidation: error executing Delete", "error", err)
		return 0, utils.ErrInternalServer
	}

	if deleted == 0 {
		return 0, utils.ErrInternalServer
	}

	return
}
