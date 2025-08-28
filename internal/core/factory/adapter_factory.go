package factory

import (
	"context"
	"database/sql"

	mysqladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	complexservices "github.com/giulio-alfieri/toq_server/internal/core/service/complex_service"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
)

// LifecycleManager define uma interface para registrar funções de cleanup.
// Isso quebra a dependência de importação circular entre factory e config.
type LifecycleManager interface {
	AddCleanupFunc(f func())
}

// AdapterFactory define a interface principal para criação de adapters
// Implementa o Abstract Factory pattern para organizar a criação de dependências
type AdapterFactory interface {
	// CreateValidationAdapters cria todos os adapters de validação externa (CEP, CPF, CNPJ, CRECI)
	CreateValidationAdapters(env *globalmodel.Environment) (ValidationAdapters, error)

	// CreateExternalServiceAdapters cria adapters de serviços externos (FCM, Email, SMS, GCS)
	CreateExternalServiceAdapters(ctx context.Context, env *globalmodel.Environment) (ExternalServiceAdapters, error)

	// CreateStorageAdapters cria adapters de armazenamento (Database, Cache)
	CreateStorageAdapters(ctx context.Context, env *globalmodel.Environment, db *sql.DB) (StorageAdapters, error)

	// CreateRepositoryAdapters cria todos os repositórios MySQL
	CreateRepositoryAdapters(database *mysqladapter.Database) (RepositoryAdapters, error)

	// CreateHTTPHandlers cria todos os handlers HTTP
	CreateHTTPHandlers(
		userService userservices.UserServiceInterface,
		globalService globalservice.GlobalServiceInterface,
		listingService listingservices.ListingServiceInterface,
		complexService complexservices.ComplexServiceInterface,
	) HTTPHandlers
}

// NewAdapterFactory cria uma nova instância da factory concreta
func NewAdapterFactory(lm LifecycleManager) AdapterFactory {
	return &ConcreteAdapterFactory{lm: lm}
}
