package sessionservice

import (
	"context"
	"log/slog"

	sessionrepoport "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/session_repository"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
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
	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		return 0, err
	}
	// Run deletion inside the transaction boundary
	n, err := s.repo.DeleteExpiredSessions(ctx, tx, limit)
	if err != nil {
		s.globalService.RollbackTransaction(ctx, tx)
		return 0, err
	}
	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		s.globalService.RollbackTransaction(ctx, tx)
		return 0, err
	}
	if n > 0 {
		slog.Info("session_service.cleaner.deleted", "count", n)
		metricSessionCleanerDeleted.Add(float64(n))
	}
	return n, nil
}
