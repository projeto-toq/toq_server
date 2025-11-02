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

	query := `SELECT * FROM temp_wrong_signin WHERE user_id = ?;`

	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, id)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_wrong_signin.query_error", "error", queryErr)
		return nil, queryErr
	}
	defer rows.Close()

	entities, err := rowsToEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_wrong_signin.rows_to_entities_error", "error", err)
		return nil, fmt.Errorf("scan wrong signin rows: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	if len(entities) > 1 {
		errMultiple := errors.New("multiple wrong_signin rows found")
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_wrong_signin.multiple_rows_error", "user_id", id, "error", errMultiple)
		return nil, errMultiple
	}

	wrongSignin, err = userconverters.WrongSignInEntityToDomain(entities[0])
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_wrong_signin.convert_error", "error", err)
		return nil, fmt.Errorf("convert wrong signin entity: %w", err)
	}

	return

}
