package factory

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/core/cache"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	cepport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cep"
	cnpjport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cnpj"
	cpfport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cpf"
	creciport "github.com/giulio-alfieri/toq_server/internal/core/port/right/creci"
	emailport "github.com/giulio-alfieri/toq_server/internal/core/port/right/email"
	fcmport "github.com/giulio-alfieri/toq_server/internal/core/port/right/fcm"
	gcsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/gcs"
	smsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/sms"

	mysqladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql"
	mysqluseradapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user"
	complexrepoport "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/complex_repository"
	globalrepoport "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/global_repository"
	listingrepoport "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/listing_repository"
	sessionrepoport "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/session_repository"
	userrepoport "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/user_repository"
)

// ValidationAdapters agrupa adapters de validação externa
type ValidationAdapters struct {
	CEP   cepport.CEPPortInterface
	CPF   cpfport.CPFPortInterface
	CNPJ  cnpjport.CNPJPortInterface
	CRECI creciport.CreciPortInterface
}

// ExternalServiceAdapters agrupa adapters de serviços externos
type ExternalServiceAdapters struct {
	FCM       fcmport.FCMPortInterface
	Email     emailport.EmailPortInterface
	SMS       smsport.SMSPortInterface
	GCS       gcsport.GCSPortInterface
	CloseFunc func() error // Função para cleanup de recursos
}

// StorageAdapters agrupa adapters de armazenamento
type StorageAdapters struct {
	Database  *mysqladapter.Database
	Cache     cache.CacheInterface
	CloseFunc func() error // Função para cleanup de recursos
}

// RepositoryAdapters agrupa todos os repositórios MySQL
type RepositoryAdapters struct {
	User        userrepoport.UserRepoPortInterface
	Global      globalrepoport.GlobalRepoPortInterface
	Complex     complexrepoport.ComplexRepoPortInterface
	Listing     listingrepoport.ListingRepoPortInterface
	Session     sessionrepoport.SessionRepoPortInterface
	DeviceToken *mysqluseradapter.DeviceTokenRepository
}

// AdapterFactoryConfig contém as configurações necessárias para criar adapters
type AdapterFactoryConfig struct {
	Context     context.Context
	Environment *globalmodel.Environment
	Database    *mysqladapter.Database
}
