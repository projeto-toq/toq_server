package mysqlglobaladapter

import (
	"context"
	"database/sql"
	"fmt"

	globalconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/global/converters"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ga *GlobalAdapter) CreateAudit(ctx context.Context, tx *sql.Tx, audit globalmodel.AuditInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO audit
			(executed_at, executed_by, table_name, table_id, action)
			VALUES (?, ?, ?, ?, ?);`

	entity := globalconverters.AuditDomainToEntity(ctx, audit)

	id, err := ga.Create(ctx, tx, query,
		entity.ExecutedAT,
		entity.ExecutedBY,
		entity.TableName,
		entity.TableID,
		entity.Action)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.global.create_audit.create_error", "error", err)
		return fmt.Errorf("create audit: %w", err)
	}

	audit.SetID(id)
	return
}
