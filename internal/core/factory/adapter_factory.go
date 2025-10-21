package factory

import (
	"context"
	"database/sql"

	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	complexservices "github.com/projeto-toq/toq_server/internal/core/service/complex_service"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	holidayservices "github.com/projeto-toq/toq_server/internal/core/service/holiday_service"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	permissionservices "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
	"github.com/projeto-toq/toq_server/internal/core/utils/hmacauth"
)

// LifecycleManager define uma interface para registrar funções de cleanup.
// Isso quebra a dependência de importação circular entre factory e config.
type LifecycleManager interface {
	AddCleanupFunc(f func())
}

// AdapterFactory define a interface principal para criação de adapters
// Implementa o Abstract Factory pattern para organizar a criação de dependências
type AdapterFactory interface {
	// CreateValidationAdapters cria todos os adapters de validação externa (CEP, CPF, CNPJ)
	CreateValidationAdapters(env *globalmodel.Environment) (ValidationAdapters, error)

	// CreateExternalServiceAdapters cria adapters de serviços externos (FCM, Email, SMS, GCS)
	CreateExternalServiceAdapters(ctx context.Context, env *globalmodel.Environment) (ExternalServiceAdapters, error)

	// CreateStorageAdapters cria adapters de armazenamento (Database, Cache)
	CreateStorageAdapters(ctx context.Context, env *globalmodel.Environment, db *sql.DB) (StorageAdapters, error)

	// CreateRepositoryAdapters cria todos os repositórios MySQL
	CreateRepositoryAdapters(database *mysqladapter.Database) (RepositoryAdapters, error)

	// CreateMetricsAdapter cria o adapter de métricas
	CreateMetricsAdapter() *MetricsAdapter

	// CreateHTTPHandlers cria todos os handlers HTTP
	CreateHTTPHandlers(
		userService userservices.UserServiceInterface,
		globalService globalservice.GlobalServiceInterface,
		listingService listingservices.ListingServiceInterface,
		complexService complexservices.ComplexServiceInterface,
		scheduleService scheduleservices.ScheduleServiceInterface,
		holidayService holidayservices.HolidayServiceInterface,
		permissionService permissionservices.PermissionServiceInterface,
		metricsAdapter *MetricsAdapter,
		hmacValidator *hmacauth.Validator,
	) HTTPHandlers
}

// NewAdapterFactory creates a new concrete adapter factory instance
func NewAdapterFactory(lm LifecycleManager) AdapterFactory {
	return &ConcreteAdapterFactory{
		lm: lm,
	}
}
