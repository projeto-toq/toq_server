package globalservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/events"
	cepmodel "github.com/projeto-toq/toq_server/internal/core/model/cep_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	cepport "github.com/projeto-toq/toq_server/internal/core/port/right/cep"
	emailport "github.com/projeto-toq/toq_server/internal/core/port/right/email"
	fcmport "github.com/projeto-toq/toq_server/internal/core/port/right/fcm"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	globalrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/global_repository"
	userrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/user_repository"
	smsport "github.com/projeto-toq/toq_server/internal/core/port/right/sms"
	storageport "github.com/projeto-toq/toq_server/internal/core/port/right/storage"
)

type globalService struct {
	globalRepo           globalrepository.GlobalRepoPortInterface
	userRepo             userrepository.UserRepoPortInterface
	cep                  cepport.CEPPortInterface
	firebaseCloudMessage fcmport.FCMPortInterface
	email                emailport.EmailPortInterface
	sms                  smsport.SMSPortInterface
	googleCludStorage    storageport.CloudStoragePortInterface
	eventBus             events.Bus
	metrics              metricsport.MetricsPortInterface
}

func NewGlobalService(
	globalRepo globalrepository.GlobalRepoPortInterface,
	userRepo userrepository.UserRepoPortInterface,
	cep cepport.CEPPortInterface,
	firebaseCloudMessage fcmport.FCMPortInterface,
	email emailport.EmailPortInterface,
	sms smsport.SMSPortInterface,
	googleCloudStorage storageport.CloudStoragePortInterface,
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
		googleCludStorage:    googleCloudStorage,
		eventBus:             events.NewInMemoryBus(),
		metrics:              metrics,
	}
}

type GlobalServiceInterface interface {
	CreateAudit(ctx context.Context, tx *sql.Tx, table globalmodel.TableName, action string, executedBY ...int64) (err error)

	GetConfiguration(ctx context.Context) (configuration map[string]string, err error)

	// Novo sistema de notificação unificado
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
func (gs *globalService) ListDeviceTokensByUserIDIfOptedIn(ctx context.Context, userID int64) ([]string, error) {
	if gs.userRepo == nil {
		return nil, fmt.Errorf("user repository not configured")
	}
	return gs.userRepo.ListDeviceTokenStringsByUserIDIfOptedIn(ctx, nil, userID)
}
