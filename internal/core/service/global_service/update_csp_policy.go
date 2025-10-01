package globalservice

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// UpdateCSPPolicy updates or creates the Content Security Policy ensuring optimistic concurrency through version checks.
func (gs *globalService) UpdateCSPPolicy(ctx context.Context, expectedVersion int64, directives map[string]string) (policy globalmodel.ContentSecurityPolicy, err error) {
	if len(directives) == 0 {
		return policy, utils.BadRequest("directives must not be empty")
	}

	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		return policy, utils.InternalError("")
	}
	defer spanEnd()

	tx, err := gs.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("global.update_csp_policy.tx_start_error", "err", err)
		return policy, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := gs.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("global.update_csp_policy.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	current, repoErr := gs.globalRepo.GetCSPPolicy(ctx, tx)
	switch {
	case repoErr == nil:
		if current.Version != expectedVersion {
			err = utils.ConflictError("content security policy version mismatch")
			utils.SetSpanError(ctx, err)
			return policy, err
		}
		updated := globalmodel.NewContentSecurityPolicy(current.ID, current.Version+1, directives)
		if repoErr = gs.globalRepo.UpdateCSPPolicy(ctx, tx, updated); repoErr != nil {
			utils.SetSpanError(ctx, repoErr)
			slog.Error("global.update_csp_policy.update_error", "err", repoErr)
			err = utils.InternalError("")
			return policy, err
		}
		policy = updated
	case errors.Is(repoErr, sql.ErrNoRows):
		if expectedVersion != 0 {
			err = utils.ConflictError("content security policy version mismatch")
			utils.SetSpanError(ctx, err)
			return policy, err
		}
		candidate := globalmodel.NewContentSecurityPolicy(0, 1, directives)
		policy, repoErr = gs.globalRepo.CreateCSPPolicy(ctx, tx, candidate)
		if repoErr != nil {
			utils.SetSpanError(ctx, repoErr)
			slog.Error("global.update_csp_policy.create_error", "err", repoErr)
			err = utils.InternalError("")
			return policy, err
		}
	default:
		utils.SetSpanError(ctx, repoErr)
		slog.Error("global.update_csp_policy.query_error", "err", repoErr)
		err = utils.InternalError("")
		return policy, err
	}

	if err = gs.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("global.update_csp_policy.tx_commit_error", "err", err)
		return policy, utils.InternalError("")
	}

	return policy.Clone(), nil
}
