package globalservice

import (
	"context"
	"database/sql"
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateAudit persists an immutable audit trail record for data-changing operations.
// When executedBY is omitted the user ID is extracted from the request context.
func (gs *globalService) CreateAudit(ctx context.Context, tx *sql.Tx, table globalmodel.TableName, action string, executedBY ...int64) error {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("global.audit.tracer_error", "err", tracerErr)
		return utils.InternalError("Failed to initialize audit tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	audit := globalmodel.NewAudit()
	if len(executedBY) > 0 {
		audit.SetExecutedBy(executedBY[0])
	} else if v := ctx.Value(globalmodel.TokenKey); v != nil {
		if infos, ok := v.(usermodel.UserInfos); ok {
			audit.SetExecutedBy(infos.ID)
		}
	}
	audit.SetExecutedAt(time.Now().UTC())
	audit.SetTableName(table)
	audit.SetAction(action)

	if err := gs.globalRepo.CreateAudit(ctx, tx, audit); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.audit.persist_error", "err", err, "table", table, "action", action)
		return utils.InternalError("Failed to persist audit record")
	}
	return nil
}
