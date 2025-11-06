package mysqluseradapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetWrongSigninByUserID(ctx context.Context, tx *sql.Tx, id int64) (wrongSignin usermodel.WrongSigninInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query wrong signin tracking by user ID
	// Note: Primary key is user_id, so max 1 row expected
	query := `SELECT user_id, failed_attempts, last_attempt_at 
	          FROM temp_wrong_signin WHERE user_id = ?;`

	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, id)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_wrong_signin.query_error", "error", queryErr)
		return nil, queryErr
	}
	defer rows.Close()

	// Scan rows using type-safe function (replaces rowsToEntities)
	entities, err := scanWrongSigninEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_wrong_signin.scan_error", "error", err)
		return nil, fmt.Errorf("scan wrong signin rows: %w", err)
	}

	// Handle no results
	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	// Safety check: primary key should prevent multiple rows
	if len(entities) > 1 {
		errMultiple := errors.New("multiple wrong_signin rows found")
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_wrong_signin.multiple_rows_error", "user_id", id, "error", errMultiple)
		return nil, errMultiple
	}

	// Convert entity to domain model using type-safe converter
	wrongSignin = userconverters.WrongSignInEntityToDomainTyped(entities[0])

	return

}
