package globalservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/events"
	cepmodel "github.com/projeto-toq/toq_server/internal/core/model/cep_model"
	cepport "github.com/projeto-toq/toq_server/internal/core/port/right/cep"
	emailport "github.com/projeto-toq/toq_server/internal/core/port/right/email"
	fcmport "github.com/projeto-toq/toq_server/internal/core/port/right/fcm"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	globalrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/global_repository"
	userrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/user_repository"
	smsport "github.com/projeto-toq/toq_server/internal/core/port/right/sms"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// globalService implements GlobalServiceInterface and orchestrates cross-domain helpers
// such as transactions, configuration lookups, notifications, and shared repositories.
type globalService struct {
	globalRepo           globalrepository.GlobalRepoPortInterface
	userRepo             userrepository.UserRepoPortInterface
	cep                  cepport.CEPPortInterface
	firebaseCloudMessage fcmport.FCMPortInterface
	email                emailport.EmailPortInterface
	sms                  smsport.SMSPortInterface
	eventBus             events.Bus
	metrics              metricsport.MetricsPortInterface
}

// NewGlobalService wires all global dependencies shared by multiple domains.
// Metrics port is optional and may be nil on minimalist setups.
func NewGlobalService(
	globalRepo globalrepository.GlobalRepoPortInterface,
	userRepo userrepository.UserRepoPortInterface,
	cep cepport.CEPPortInterface,
	firebaseCloudMessage fcmport.FCMPortInterface,
	email emailport.EmailPortInterface,
	sms smsport.SMSPortInterface,
	// optional metrics (can be nil in tests or minimal setups)
	metrics metricsport.MetricsPortInterface,
) GlobalServiceInterface {
	return &globalService{
		globalRepo:           globalRepo,
		userRepo:             userRepo,
		cep:                  cep,
		firebaseCloudMessage: firebaseCloudMessage,
		email:                email,
		sms:                  sms,
		eventBus:             events.NewInMemoryBus(),
		metrics:              metrics,
	}
}

// GlobalServiceInterface centralizes helpers required by multiple services (transactions,
// configuration cache, notification fan-out, session events, etc.).
type GlobalServiceInterface interface {
	GetConfiguration(ctx context.Context) (configuration map[string]string, err error)

	// Unified notification service accessor
	GetUnifiedNotificationService() UnifiedNotificationService

	// Event bus accessor (for publishing session events)
	GetEventBus() events.Bus
	GetCEP(ctx context.Context, cep string) (address cepmodel.CEPInterface, err error)
	// StartSessionEventSubscriber starts the subscriber and returns an unsubscribe function
	StartSessionEventSubscriber() func()

	// Optional metrics accessor
	GetMetrics() metricsport.MetricsPortInterface

	StartTransaction(ctx context.Context) (tx *sql.Tx, err error)
	RollbackTransaction(ctx context.Context, tx *sql.Tx) (err error)
	CommitTransaction(ctx context.Context, tx *sql.Tx) (err error)
	// StartReadOnlyTransaction starts a read-only transaction for pure read flows
	StartReadOnlyTransaction(ctx context.Context) (tx *sql.Tx, err error)
	GetUserIDFromContext(ctx context.Context) (int64, error)
	ListDeviceTokensByUserIDIfOptedIn(ctx context.Context, userID int64) ([]string, error)
}

// GetEventBus returns the in-memory event bus instance
func (gs *globalService) GetEventBus() events.Bus {
	return gs.eventBus
}

// GetMetrics returns the metrics port if configured (may be nil)
func (gs *globalService) GetMetrics() metricsport.MetricsPortInterface {
	return gs.metrics
}

// ListDeviceTokensByUserIDIfOptedIn returns all push tokens for a user when promotional opt-in is active.
// It enriches context with tracing/logging metadata before reaching the repository layer.
func (gs *globalService) ListDeviceTokensByUserIDIfOptedIn(ctx context.Context, userID int64) ([]string, error) {
	if gs.userRepo == nil {
		return nil, fmt.Errorf("user repository not configured")
	}

	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("global.list_device_tokens.tracer_error", "err", tracerErr)
		return nil, utils.InternalError("Failed to initialize device token tracing")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tokens, err := gs.userRepo.ListDeviceTokenStringsByUserIDIfOptedIn(ctx, nil, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.list_device_tokens.repo_error", "err", err, "user_id", userID)
		return nil, utils.InternalError("Failed to list user device tokens")
	}

	return tokens, nil
}
