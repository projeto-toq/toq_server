package factory

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/cache"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	cepport "github.com/projeto-toq/toq_server/internal/core/port/right/cep"
	cnpjport "github.com/projeto-toq/toq_server/internal/core/port/right/cnpj"
	cpfport "github.com/projeto-toq/toq_server/internal/core/port/right/cpf"
	emailport "github.com/projeto-toq/toq_server/internal/core/port/right/email"
	fcmport "github.com/projeto-toq/toq_server/internal/core/port/right/fcm"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	smsport "github.com/projeto-toq/toq_server/internal/core/port/right/sms"
	storageport "github.com/projeto-toq/toq_server/internal/core/port/right/storage"

	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	mysqluseradapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user"
	complexrepoport "github.com/projeto-toq/toq_server/internal/core/port/right/repository/complex_repository"
	globalrepoport "github.com/projeto-toq/toq_server/internal/core/port/right/repository/global_repository"
	listingrepoport "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	permissionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/permission_repository"
	photosessionrepo "github.com/projeto-toq/toq_server/internal/core/port/right/repository/photo_session_repository"
	sessionrepoport "github.com/projeto-toq/toq_server/internal/core/port/right/repository/session_repository"
	userrepoport "github.com/projeto-toq/toq_server/internal/core/port/right/repository/user_repository"
)

// ValidationAdapters agrupa adapters de validação externa
type ValidationAdapters struct {
	CEP  cepport.CEPPortInterface
	CPF  cpfport.CPFPortInterface
	CNPJ cnpjport.CNPJPortInterface
}

// ExternalServiceAdapters agrupa adapters de serviços externos
type ExternalServiceAdapters struct {
	FCM          fcmport.FCMPortInterface
	Email        emailport.EmailPortInterface
	SMS          smsport.SMSPortInterface
	CloudStorage storageport.CloudStoragePortInterface
	CloseFunc    func() error // Função para cleanup de recursos
}

// StorageAdapters agrupa adapters de armazenamento
type StorageAdapters struct {
	Database  *mysqladapter.Database
	Cache     cache.CacheInterface
	CloseFunc func() error // Função para cleanup de recursos
}

// RepositoryAdapters agrupa todos os repositórios MySQL
type RepositoryAdapters struct {
	User         userrepoport.UserRepoPortInterface
	Global       globalrepoport.GlobalRepoPortInterface
	Complex      complexrepoport.ComplexRepoPortInterface
	Listing      listingrepoport.ListingRepoPortInterface
	PhotoSession photosessionrepo.PhotoSessionRepositoryInterface
	Session      sessionrepoport.SessionRepoPortInterface
	Permission   permissionrepository.PermissionRepositoryInterface
	DeviceToken  *mysqluseradapter.DeviceTokenRepository
}

// HTTPHandlers agrupa todos os handlers HTTP
type HTTPHandlers struct {
	UserHandler    interface{} // User handler interface
	ListingHandler interface{} // Listing handler interface
	AuthHandler    interface{} // Auth handler interface
	MetricsHandler interface{} // Handler para endpoint /metrics
	AdminHandler   interface{} // Admin handler interface
	ComplexHandler interface{} // Complex handler interface
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
