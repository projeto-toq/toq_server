package globalservice

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetCSPPolicy retrieves the current Content Security Policy persisted in the configuration repository.
func (gs *globalService) GetCSPPolicy(ctx context.Context) (policy globalmodel.ContentSecurityPolicy, err error) {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		return policy, utils.InternalError("")
	}
	defer spanEnd()

	tx, err := gs.globalRepo.StartReadOnlyTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("global.get_csp_policy.tx_start_error", "err", err)
		return policy, utils.InternalError("")
	}
	defer func() {
		if rbErr := gs.RollbackTransaction(ctx, tx); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			utils.SetSpanError(ctx, rbErr)
			slog.Error("global.get_csp_policy.tx_rollback_error", "err", rbErr)
		}
	}()

	policy, err = gs.globalRepo.GetCSPPolicy(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFound := utils.NotFoundError("content security policy")
			utils.SetSpanError(ctx, notFound)
			return globalmodel.ContentSecurityPolicy{}, notFound
		}
		utils.SetSpanError(ctx, err)
		slog.Error("global.get_csp_policy.query_error", "err", err)
		return globalmodel.ContentSecurityPolicy{}, utils.InternalError("")
	}

	return policy.Clone(), nil
}
