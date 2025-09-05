package validationservice

import (
	"context"
	"log/slog"

	userrepo "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/user_repository"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type Service interface {
	CleanExpiredValidations(ctx context.Context, limit int) (int64, error)
}

type service struct {
	repo          userrepo.UserRepoPortInterface
	globalService globalservice.GlobalServiceInterface
}

var (
	metricValidationCleanerDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "validation_cleaner_deleted_total",
		Help: "Total number of temp_user_validations rows deleted by the validation cleaner service",
	})
)

func init() {
	prometheus.MustRegister(metricValidationCleanerDeleted)
}

func New(repo userrepo.UserRepoPortInterface, gs globalservice.GlobalServiceInterface) Service {
	return &service{repo: repo, globalService: gs}
}

// CleanExpiredValidations deletes expired validation rows within a transaction boundary
func (s *service) CleanExpiredValidations(ctx context.Context, limit int) (int64, error) {
	// Create tracing span for public entrypoint
	_, end, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return 0, utils.InternalError("")
	}
	defer end()

	// Defensive default if caller passes invalid limit
	if limit <= 0 {
		slog.Warn("validation.cleaner.invalid_limit", "limit", limit)
		limit = 500
	}

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("validation.cleaner.tx_start_error", "err", err)
		return 0, utils.InternalError("")
	}
	defer func() {
		if err != nil { // rollback only when a prior error occurred before commit
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("validation.cleaner.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	n, err := s.repo.DeleteExpiredValidations(ctx, tx, limit)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("validation.cleaner.delete_error", "err", err)
		return 0, utils.InternalError("")
	}
	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		slog.Error("validation.cleaner.tx_commit_error", "err", cmErr)
		return 0, utils.InternalError("")
	}
	if n > 0 {
		slog.Info("validation.cleaner.deleted", "count", n)
		metricValidationCleanerDeleted.Add(float64(n))
	}
	return n, nil
}
