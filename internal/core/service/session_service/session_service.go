package sessionservice

import (
	"context"
	"log/slog"

	sessionrepoport "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/session_repository"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type Service interface {
	CleanExpiredSessions(ctx context.Context, limit int) (int64, error)
}

type service struct {
	repo          sessionrepoport.SessionRepoPortInterface
	globalService globalservice.GlobalServiceInterface
}

var (
	metricSessionCleanerDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "session_cleaner_deleted_total",
		Help: "Total number of sessions deleted by the session cleaner service",
	})
)

func init() {
	prometheus.MustRegister(metricSessionCleanerDeleted)
}

func New(repo sessionrepoport.SessionRepoPortInterface, gs globalservice.GlobalServiceInterface) Service {
	return &service{repo: repo, globalService: gs}
}

// CleanExpiredSessions deletes expired sessions within a transaction boundary
func (s *service) CleanExpiredSessions(ctx context.Context, limit int) (int64, error) {
	// tracing span for public entrypoint
	_, spanEnd, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return 0, utils.InternalError("")
	}
	defer spanEnd()

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("session.cleaner.tx_start_error", "err", err)
		return 0, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("session.cleaner.tx_rollback_error", "err", rbErr)
			}
		}
	}()
	// Run deletion inside the transaction boundary
	n, err := s.repo.DeleteExpiredSessions(ctx, tx, limit)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return 0, utils.InternalError("")
	}
	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		slog.Error("session.cleaner.tx_commit_error", "err", cmErr)
		return 0, utils.InternalError("")
	}
	if n > 0 {
		slog.Info("session.cleaner.deleted", "count", n)
		metricSessionCleanerDeleted.Add(float64(n))
	}
	return n, nil
}
