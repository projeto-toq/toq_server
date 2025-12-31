package factory

import (
	"context"

	metricshandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers"
	adminhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/admin_handlers"
	authhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/auth_handlers"
	holidayhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/holiday_handlers"
	listinghandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/listing_handlers"
	mediaprocessinghandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/media_processing_handlers"
	photosessionhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/photo_session_handlers"
	schedulehandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/schedule_handlers"
	userhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/user_handlers"
	visithandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/visit_handlers"
	"github.com/projeto-toq/toq_server/internal/core/cache"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	cepport "github.com/projeto-toq/toq_server/internal/core/port/right/cep"
	cnpjport "github.com/projeto-toq/toq_server/internal/core/port/right/cnpj"
	cpfport "github.com/projeto-toq/toq_server/internal/core/port/right/cpf"
	emailport "github.com/projeto-toq/toq_server/internal/core/port/right/email"
	fcmport "github.com/projeto-toq/toq_server/internal/core/port/right/fcm"
	mediaprocessingcallbackport "github.com/projeto-toq/toq_server/internal/core/port/right/functions/mediaprocessingcallback"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	mediaprocessingqueue "github.com/projeto-toq/toq_server/internal/core/port/right/queue/mediaprocessingqueue"
	smsport "github.com/projeto-toq/toq_server/internal/core/port/right/sms"
	storageport "github.com/projeto-toq/toq_server/internal/core/port/right/storage"
	workflowport "github.com/projeto-toq/toq_server/internal/core/port/right/workflow"

	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	globalrepoport "github.com/projeto-toq/toq_server/internal/core/port/right/repository/global_repository"
	holidayrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/holiday_repository"
	listingrepoport "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	mediaprocessingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/media_processing_repository"
	permissionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/permission_repository"
	photosessionrepo "github.com/projeto-toq/toq_server/internal/core/port/right/repository/photo_session_repository"
	propertycoveragerepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/property_coverage_repository"
	schedulerepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/schedule_repository"
	sessionrepoport "github.com/projeto-toq/toq_server/internal/core/port/right/repository/session_repository"
	userrepoport "github.com/projeto-toq/toq_server/internal/core/port/right/repository/user_repository"
	visitrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/visit_repository"
)

// ValidationAdapters agrupa adapters de validação externa
type ValidationAdapters struct {
	CEP  cepport.CEPPortInterface
	CPF  cpfport.CPFPortInterface
	CNPJ cnpjport.CNPJPortInterface
}

// ExternalServiceAdapters agrupa adapters de serviços externos
type ExternalServiceAdapters struct {
	FCM                     fcmport.FCMPortInterface
	Email                   emailport.EmailPortInterface
	SMS                     smsport.SMSPortInterface
	CloudStorage            storageport.CloudStoragePortInterface
	ListingMediaStorage     storageport.ListingMediaStoragePort
	MediaProcessingQueue    mediaprocessingqueue.QueuePortInterface
	MediaProcessingCallback mediaprocessingcallbackport.CallbackPortInterface
	MediaProcessingWorkflow workflowport.WorkflowPortInterface
	CloseFunc               func() error // Função para cleanup de recursos
}

// StorageAdapters agrupa adapters de armazenamento
type StorageAdapters struct {
	Database  *mysqladapter.Database
	Cache     cache.CacheInterface
	CloseFunc func() error // Função para cleanup de recursos
}

// RepositoryAdapters agrupa todos os repositórios MySQL
type RepositoryAdapters struct {
	User             userrepoport.UserRepoPortInterface
	Global           globalrepoport.GlobalRepoPortInterface
	PropertyCoverage propertycoveragerepository.RepositoryInterface
	Listing          listingrepoport.ListingRepoPortInterface
	MediaProcessing  mediaprocessingrepository.RepositoryInterface
	Holiday          holidayrepository.HolidayRepositoryInterface
	Schedule         schedulerepository.ScheduleRepositoryInterface
	Visit            visitrepository.VisitRepositoryInterface
	PhotoSession     photosessionrepo.PhotoSessionRepositoryInterface
	Session          sessionrepoport.SessionRepoPortInterface
	Permission       permissionrepository.PermissionRepositoryInterface
}

// HTTPHandlers agrupa todos os handlers HTTP
type HTTPHandlers struct {
	UserHandler            *userhandlers.UserHandler
	ListingHandler         *listinghandlers.ListingHandler
	MediaProcessingHandler *mediaprocessinghandlers.MediaProcessingHandler
	AuthHandler            *authhandlers.AuthHandler
	MetricsHandler         *metricshandlers.MetricsHandler
	AdminHandler           *adminhandlers.AdminHandler
	ScheduleHandler        *schedulehandlers.ScheduleHandler
	HolidayHandler         *holidayhandlers.HolidayHandler
	PhotoSessionHandler    *photosessionhandlers.PhotoSessionHandler
	VisitHandler           *visithandlers.VisitHandler
}

// MetricsAdapter contém o adapter de métricas
type MetricsAdapter struct {
	Prometheus metricsport.MetricsPortInterface
}

// AdapterFactoryConfig contém as configurações necessárias para criar adapters
type AdapterFactoryConfig struct {
	Context     context.Context
	Environment *globalmodel.Environment
	Database    *mysqladapter.Database
	Metrics     *MetricsAdapter
}
