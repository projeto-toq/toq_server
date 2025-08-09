package globalservice

import (
	"context"
	"database/sql"

	cepmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cep_model"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	cepport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cep"
	emailport "github.com/giulio-alfieri/toq_server/internal/core/port/right/email"
	fcmport "github.com/giulio-alfieri/toq_server/internal/core/port/right/fcm"
	gcsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/gcs"
	globalrepository "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/global_repository"
	smsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/sms"
)

type globalService struct {
	globalRepo           globalrepository.GlobalRepoPortInterface
	cep                  cepport.CEPPortInterface
	firebaseCloudMessage fcmport.FCMPortInterface
	email                emailport.EmailPortInterface
	sms                  smsport.SMSPortInterface
	googleCludStorage    gcsport.GCSPortInterface
}

func NewGlobalService(
	globalRepo globalrepository.GlobalRepoPortInterface,
	cep cepport.CEPPortInterface,
	firebaseCloudMessage fcmport.FCMPortInterface,
	email emailport.EmailPortInterface,
	sms smsport.SMSPortInterface,
	googleCloudStorage gcsport.GCSPortInterface,
) GlobalServiceInterface {
	return &globalService{
		globalRepo:           globalRepo,
		cep:                  cep,
		firebaseCloudMessage: firebaseCloudMessage,
		email:                email,
		sms:                  sms,
		googleCludStorage:    googleCloudStorage,
	}
}

type GlobalServiceInterface interface {
	CreateAudit(ctx context.Context, tx *sql.Tx, table globalmodel.TableName, action string, executedBY ...int64) (err error)

	GetConfiguration(ctx context.Context) (configuration map[string]string, err error)

	// Novo sistema de notificação unificado
	GetUnifiedNotificationService() UnifiedNotificationService

	// DEPRECATED: SendNotification será removido em favor do sistema unificado
	SendNotification(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType, code ...string) (err error)

	GetPrivilegeForCache(ctx context.Context, service usermodel.GRPCService, method uint8, role usermodel.UserRole) (privilege usermodel.PrivilegeInterface, err error)
	GetCEP(ctx context.Context, cep string) (address cepmodel.CEPInterface, err error)

	StartTransaction(ctx context.Context) (tx *sql.Tx, err error)
	RollbackTransaction(ctx context.Context, tx *sql.Tx) (err error)
	CommitTransaction(ctx context.Context, tx *sql.Tx) (err error)
}
