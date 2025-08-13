package goroutines

import (
	"context"
	"log/slog"
	"sync"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
	"github.com/google/uuid"
)

// CreciValidationWorker is a goroutine that periodically validates CRECI data.
// It uses the provided UserServiceInterface to fetch and validate CRECI data and face data.
// The function runs until the context is done or an error occurs.
// Parameters:
// - service: the user service interface used to fetch and validate CRECI data.
// - wg: a wait group to signal when the goroutine is done.
// - ctx: the context to control the lifecycle of the goroutine.
func CreciValidationWorker(service userservices.UserServiceInterface, wg *sync.WaitGroup, ctx context.Context) {
	slog.Info("CRECI validation routine started")
	ticker := time.NewTicker(globalmodel.ElapseTime)
	defer ticker.Stop()
	defer wg.Done()

	infos := usermodel.UserInfos{}
	infos.ID = 0
	ctx = context.WithValue(ctx, globalmodel.TokenKey, infos)
	ctx = context.WithValue(ctx, globalmodel.RequestIDKey, uuid.New().String())

	for {
		select {
		case <-ctx.Done():
			slog.Info("Creci validation routine stopped")
			return
		case <-ticker.C:
			slog.Info("Creci validation routine ticked")
			realtorsToValidateData, err := service.GetCrecisToValidateByStatus(ctx, usermodel.StatusPendingOCR)
			if err != nil {
				slog.Error("error getting crecis to validate", "error", err)
				return
			}
			service.ValidateCreciData(ctx, realtorsToValidateData)

			realtorsToValidateFace, err := service.GetCrecisToValidateByStatus(ctx, usermodel.StatusPendingFace)
			if err != nil {
				slog.Error("error getting crecis to validate", "error", err)
				return
			}
			service.ValidateCreciFace(ctx, realtorsToValidateFace)

		}
	}
}
