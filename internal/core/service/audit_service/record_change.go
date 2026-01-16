package auditservice

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// RecordChange stores a normalized audit event with actor, target and correlation metadata.
func (s *auditService) RecordChange(ctx context.Context, tx *sql.Tx, input auditmodel.RecordInput) (err error) {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("audit.record.tracer_error", "err", tracerErr)
		return derrors.Infra("failed to initialize audit tracer", tracerErr)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err = s.validateInput(input); err != nil {
		return derrors.BadRequest("invalid audit input", derrors.WithDetails(map[string]string{"reason": err.Error()}))
	}

	event := s.buildEvent(ctx, input)

	if err = s.repo.CreateEvent(ctx, tx, event); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("audit.record.persist_error", slog.Any("target", input.Target), slog.String("operation", string(input.Operation)), slog.Any("err", err))
		return derrors.Infra("failed to persist audit event", err)
	}

	return nil
}
