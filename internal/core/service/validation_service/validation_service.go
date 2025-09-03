package validationservice

import (
	"context"
	"log/slog"

	userrepo "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/user_repository"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
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
	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		return 0, err
	}
	n, err := s.repo.DeleteExpiredValidations(ctx, tx, limit)
	if err != nil {
		s.globalService.RollbackTransaction(ctx, tx)
		return 0, err
	}
	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		s.globalService.RollbackTransaction(ctx, tx)
		return 0, err
	}
	if n > 0 {
		slog.Info("validation_service.cleaner.deleted", "count", n)
		metricValidationCleanerDeleted.Add(float64(n))
	}
	return n, nil
}
