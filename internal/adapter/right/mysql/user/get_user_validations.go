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

func (ua *UserAdapter) GetUserValidations(ctx context.Context, tx *sql.Tx, id int64) (validation usermodel.ValidationInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT * FROM temp_user_validations WHERE user_id = ?;`

	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, id)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_user_validations.query_error", "error", queryErr)
		return nil, queryErr
	}
	defer rows.Close()

	entities, err := rowsToEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_validations.rows_to_entities_error", "error", err)
		return nil, fmt.Errorf("scan user validations rows: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	if len(entities) > 1 {
		errMultiple := errors.New("multiple validations found for user")
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_user_validations.multiple_validations_error", "user_id", id, "error", errMultiple)
		return nil, errMultiple
	}

	validation, err = userconverters.UserValidationEntityToDomain(entities[0])
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_user_validations.convert_error", "error", err)
		return nil, fmt.Errorf("convert user validation entity: %w", err)
	}

	return

}
