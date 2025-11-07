package userservices

import (
	"context"
	"fmt"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricCreciVerifyTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "creci_documents_verify_total",
		Help: "Total number of verify CRECI documents calls",
	}, []string{"result"})
	metricCreciVerifyMissing = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "creci_documents_verify_missing_total",
		Help: "Total number of missing CRECI documents by type",
	}, []string{"doc"})
)

func init() {
	prometheus.MustRegister(metricCreciVerifyTotal)
	prometheus.MustRegister(metricCreciVerifyMissing)
}

// VerifyCreciDocuments checks S3 for required CRECI images and sets user status to PendingManual inside a DB transaction.
func (us *userService) VerifyCreciDocuments(ctx context.Context) (err error) {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	if us.cloudStorageService == nil {
		return utils.InternalError("Storage service not configured")
	}

	bucket := us.cloudStorageService.GetBucketConfig().Name
	// Documentos exigidos
	docs := []string{"selfie.jpg", "front.jpg", "back.jpg"}
	missing := make([]string, 0, 3)

	for _, d := range docs {
		object := fmt.Sprintf("%d/%s", userID, d)
		exists, e := us.cloudStorageService.ObjectExists(ctx, bucket, object)
		if e != nil {
			if gm := us.globalService.GetMetrics(); gm != nil {
				gm.IncrementErrors("user_service", "creci_verify_error")
			}
			utils.SetSpanError(ctx, e)
			logger.Error("user.verify_creci.object_exists_error", "user_id", userID, "doc", d, "error", e)
			return utils.InternalError("Failed to check document existence")
		}
		if !exists {
			missing = append(missing, d)
			metricCreciVerifyMissing.WithLabelValues(d).Inc()
		}
	}

	if len(missing) > 0 {
		metricCreciVerifyTotal.WithLabelValues("missing").Inc()
		// Keep 422 but via ValidationError-like payload; use BadRequest with details per our helpers
		return utils.NewHTTPErrorWithSource(422, "Missing required documents", map[string]any{
			"missing": missing,
		})
	}

	// Atualiza o status dentro de transação
	tx, e := us.globalService.StartTransaction(ctx)
	if e != nil {
		utils.SetSpanError(ctx, e)
		logger.Error("user.verify_creci.tx_start_error", "error", e)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("user.verify_creci.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	if e = us.repo.UpdateUserRoleStatusByUserID(ctx, userID, int(globalmodel.StatusPendingManual)); e != nil {
		utils.SetSpanError(ctx, e)
		logger.Error("user.verify_creci.update_status_error", "user_id", userID, "error", e)
		err = utils.InternalError("Failed to set user status")
		return
	}
	if e = us.globalService.CommitTransaction(ctx, tx); e != nil {
		utils.SetSpanError(ctx, e)
		logger.Error("user.verify_creci.tx_commit_error", "user_id", userID, "error", e)
		err = utils.InternalError("Failed to commit transaction")
		return
	}

	metricCreciVerifyTotal.WithLabelValues("success").Inc()
	return nil
}
