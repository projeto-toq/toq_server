package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// Create executa uma query INSERT e retorna o ID gerado
func (pa *PermissionAdapter) Create(ctx context.Context, tx *sql.Tx, query string, args ...any) (id int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.create.prepare_error", "error", err)
		return 0, fmt.Errorf("prepare permission create statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.create.exec_error", "error", err)
		return 0, fmt.Errorf("exec permission create statement: %w", err)
	}

	id, err = result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.create.last_insert_id_error", "error", err)
		return 0, fmt.Errorf("last insert id permission create: %w", err)
	}

	return id, nil
}
