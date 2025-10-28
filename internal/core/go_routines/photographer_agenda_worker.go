package goroutines

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"
	"time"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	userrepoport "github.com/projeto-toq/toq_server/internal/core/port/right/repository/user_repository"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// PhotographerAgendaWorkerConfig groups the tunables for the photographer agenda refresher.
type PhotographerAgendaWorkerConfig struct {
	Interval      time.Duration
	HorizonMonths int
	Timezone      string
}

// PhotographerAgendaWorker keeps photographer agendas with the configured horizon.
type PhotographerAgendaWorker struct {
	repo          userrepoport.UserRepoPortInterface
	photoService  photosessionservices.PhotoSessionServiceInterface
	globalService globalservice.GlobalServiceInterface
	cfg           PhotographerAgendaWorkerConfig
	logger        *slog.Logger
}

// NewPhotographerAgendaWorker wires a photographer agenda worker.
func NewPhotographerAgendaWorker(
	repo userrepoport.UserRepoPortInterface,
	photoService photosessionservices.PhotoSessionServiceInterface,
	globalService globalservice.GlobalServiceInterface,
	cfg PhotographerAgendaWorkerConfig,
) *PhotographerAgendaWorker {
	if cfg.Interval <= 0 {
		cfg.Interval = 24 * time.Hour
	}
	if cfg.HorizonMonths <= 0 {
		cfg.HorizonMonths = 3
	}
	if cfg.Timezone == "" {
		cfg.Timezone = "America/Sao_Paulo"
	}
	return &PhotographerAgendaWorker{
		repo:          repo,
		photoService:  photoService,
		globalService: globalService,
		cfg:           cfg,
	}
}

// Start launches the refresher loop until the context is cancelled.
func (w *PhotographerAgendaWorker) Start(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	ctx = utils.ContextWithLogger(ctx)
	w.logger = utils.LoggerFromContext(ctx)
	w.logger.Info("photographer agenda worker started", "interval", w.cfg.Interval)

	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	w.refresh(ctx)

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("photographer agenda worker stopped")
			return
		case <-ticker.C:
			w.refresh(ctx)
		}
	}
}

func (w *PhotographerAgendaWorker) refresh(ctx context.Context) {
	var spanEnd func()
	var err error
	ctx, spanEnd, err = utils.GenerateTracer(ctx)
	if err != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("photographer.agenda.worker.tracer_failed", "error", err)
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := w.globalService.StartReadOnlyTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photographer.agenda.worker.tx_start_failed", "error", err)
		return
	}
	defer func() {
		if rbErr := w.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("photographer.agenda.worker.tx_rollback_failed", "error", rbErr)
		}
	}()

	photographers, err := w.repo.GetUsersByRoleAndStatus(ctx, tx, permissionmodel.RoleSlugPhotographer, permissionmodel.StatusActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photographer.agenda.worker.list_failed", "error", err)
		return
	}

	logger.Debug("photographer.agenda.worker.refreshing", "count", len(photographers))

	for _, photographer := range photographers {
		w.refreshPhotographer(ctx, photographer)
	}
}

func (w *PhotographerAgendaWorker) refreshPhotographer(ctx context.Context, user usermodel.UserInterface) {
	logger := utils.LoggerFromContext(ctx)
	skipCtx := utils.WithSkipTracing(ctx)
	err := w.photoService.RefreshPhotographerAgenda(skipCtx, photosessionservices.EnsureAgendaInput{
		PhotographerID: uint64(user.GetID()),
		Timezone:       w.cfg.Timezone,
		HorizonMonths:  w.cfg.HorizonMonths,
	})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photographer.agenda.worker.refresh_failed", "user_id", user.GetID(), "error", err)
		return
	}
}
