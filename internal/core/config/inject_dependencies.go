package config

import (
	"fmt"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	"github.com/giulio-alfieri/toq_server/internal/core/factory"
	grpclistingport "github.com/giulio-alfieri/toq_server/internal/core/port/left/grpc/listing"
	grpcuserport "github.com/giulio-alfieri/toq_server/internal/core/port/left/grpc/user"
	complexservices "github.com/giulio-alfieri/toq_server/internal/core/service/complex_service"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
)

// InjectDependencies orquestra a criação de todos os adapters usando Factory Pattern
// Aplica princípios SOLID e melhores práticas Go para injeção de dependências
func (c *config) InjectDependencies() (close func() error, err error) {
	slog.Info("Starting dependency injection using Factory Pattern")

	// Criar factory
	adapterFactory := factory.NewAdapterFactory()

	// Configuração para factory
	factoryConfig := factory.AdapterFactoryConfig{
		Context:     c.context,
		Environment: &c.env,
		Database:    nil, // Será definido após criar storage
	}

	// Validar configuração
	if err = factory.ValidateFactoryConfig(factoryConfig); err != nil {
		return nil, fmt.Errorf("invalid factory configuration: %w", err)
	}

	// 1. Criar Storage Adapters (Database + Cache)
	slog.Info("Creating storage adapters")
	storage, err := adapterFactory.CreateStorageAdapters(c.context, &c.env, c.db)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage adapters: %w", err)
	}
	c.assignStorageAdapters(storage)

	// Atualizar factory config com database
	factoryConfig.Database = storage.Database

	// 2. Criar Repository Adapters
	slog.Info("Creating repository adapters")
	repositories, err := adapterFactory.CreateRepositoryAdapters(storage.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository adapters: %w", err)
	}
	c.assignRepositoryAdapters(repositories)

	// 3. Criar Validation Adapters (CEP, CPF, CNPJ, CRECI)
	slog.Info("Creating validation adapters")
	validation, err := adapterFactory.CreateValidationAdapters(&c.env)
	if err != nil {
		return nil, fmt.Errorf("failed to create validation adapters: %w", err)
	}
	c.assignValidationAdapters(validation)

	// 4. Criar External Service Adapters (FCM, Email, SMS)
	slog.Info("Creating external service adapters")
	external, err := adapterFactory.CreateExternalServiceAdapters(c.context, &c.env)
	if err != nil {
		return nil, fmt.Errorf("failed to create external service adapters: %w", err)
	}
	c.assignExternalServiceAdapters(external)

	// 5. Inicializar Services
	c.initializeServices()

	slog.Info("Dependency injection completed successfully using Factory Pattern")

	// Retornar função de cleanup
	return storage.CloseFunc, nil
}

func (c *config) InitGlobalService() {
	slog.Debug("Initializing Global Service")
	c.globalService = globalservice.NewGlobalService(
		c.repositoryAdapters.Global,
		c.cep,
		c.firebaseCloudMessaging,
		c.email,
		c.sms,
		c.googleCloudStorage,
	)

	// Injetar GlobalService no cache Redis para resolver dependência circular
	if c.cache != nil {
		c.cache.SetGlobalService(c.globalService)
		slog.Debug("GlobalService injected into Redis cache")
	}
}

func (c *config) InitUserHandler() {
	slog.Debug("Initializing User Handler")
	c.userService = userservices.NewUserService(
		c.repositoryAdapters.User,
		c.repositoryAdapters.Session,
		c.globalService,
		c.listingService,
		c.cpf,
		c.cnpj,
		c.creci,
		c.googleCloudStorage,
	)
	handler := grpcuserport.NewUserHandler(c.userService)
	pb.RegisterUserServiceServer(c.server, handler)
}

func (c *config) InitComplexHandler() {
	slog.Debug("Initializing Complex Handler")
	c.complexService = complexservices.NewComplexService(
		c.repositoryAdapters.Complex,
		c.globalService,
	)
}

func (c *config) InitListingHandler() {
	slog.Debug("Initializing Listing Handler")
	c.listingService = listingservices.NewListingService(
		c.repositoryAdapters.Listing,
		c.complexService,
		c.globalService,
		c.googleCloudStorage,
	)
	handler := grpclistingport.NewUserHandler(c.listingService)
	pb.RegisterListingServiceServer(c.server, handler)
}
