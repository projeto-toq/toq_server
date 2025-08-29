package globalservice

import (
	"context"
	"database/sql"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (gs *globalService) CreateAudit(ctx context.Context, tx *sql.Tx, table globalmodel.TableName, action string, executedBY ...int64) (err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	audit := globalmodel.NewAudit()
	if len(executedBY) > 0 {
		audit.SetExecutedBy(executedBY[0])
	} else {
		audit.SetExecutedBy(ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos).ID)
	}
	audit.SetExecutedAt(time.Now().UTC())
	audit.SetTableName(table)
	audit.SetAction(action)

	return gs.globalRepo.CreateAudit(ctx, tx, audit)

}
