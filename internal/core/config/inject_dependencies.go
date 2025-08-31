package config

import (
	"fmt"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/factory"
	goroutines "github.com/giulio-alfieri/toq_server/internal/core/go_routines"
	complexservices "github.com/giulio-alfieri/toq_server/internal/core/service/complex_service"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
	permissionservices "github.com/giulio-alfieri/toq_server/internal/core/service/permission_service"
	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
)

// InjectDependencies orquestra a criação de todos os adapters usando Factory Pattern
// Aplica princípios SOLID e melhores práticas Go para injeção de dependências
func (c *config) InjectDependencies(lm *LifecycleManager) (err error) {
	slog.Info("Starting dependency injection using Factory Pattern")
	slog.Debug("InjectDependencies method called on config instance")

	if c == nil {
		slog.Error("DEBUG: config instance is nil")
		return fmt.Errorf("config instance is nil")
	}

	if lm == nil {
		slog.Error("DEBUG: LifecycleManager is nil")
		return fmt.Errorf("lifecycle manager is nil")
	}

	// Criar factory e salvar no config
	c.adapterFactory = factory.NewAdapterFactory(lm)

	// Configuração para factory
	factoryConfig := factory.AdapterFactoryConfig{
		Context:     c.context,
		Environment: &c.env,
		Database:    nil, // Será definido após criar storage
	}

	// Validar configuração
	if err = factory.ValidateFactoryConfig(factoryConfig); err != nil {
		return fmt.Errorf("invalid factory configuration: %w", err)
	}

	// 1. Criar Storage Adapters (Database + Cache)

	storage, err := c.adapterFactory.CreateStorageAdapters(c.context, &c.env, c.db)
	if err != nil {
		return fmt.Errorf("failed to create storage adapters: %w", err)
	}
	c.assignStorageAdapters(storage)
	if storage.CloseFunc != nil {
		lm.AddCleanupFunc(func() { _ = storage.CloseFunc() })
	}

	// NOVO: Criar ActivityTracker após cache estar disponível
	if err := c.createActivityTracker(); err != nil {
		return fmt.Errorf("failed to create activity tracker: %w", err)
	}

	// Atualizar factory config com database
	factoryConfig.Database = storage.Database

	// 2. Criar Repository Adapters

	repositories, err := c.adapterFactory.CreateRepositoryAdapters(storage.Database)
	if err != nil {
		return fmt.Errorf("failed to create repository adapters: %w", err)
	}
	c.assignRepositoryAdapters(repositories)

	// 3. Criar Validation Adapters (CEP, CPF, CNPJ, CRECI)

	validation, err := c.adapterFactory.CreateValidationAdapters(&c.env)
	if err != nil {
		return fmt.Errorf("failed to create validation adapters: %w", err)
	}
	c.assignValidationAdapters(validation)

	// 4. Criar External Service Adapters (FCM, Email, SMS)
	slog.Info("Creating external service adapters")
	external, err := c.adapterFactory.CreateExternalServiceAdapters(c.context, &c.env)
	if err != nil {
		return fmt.Errorf("failed to create external service adapters: %w", err)
	}
	c.assignExternalServiceAdapters(external)
	if external.CloseFunc != nil {
		lm.AddCleanupFunc(func() { _ = external.CloseFunc() })
	}

	// 5. Inicializar Services
	c.initializeServices()

	// 6. Inicializar TempBlockCleanerWorker após permission service estar disponível
	if err := c.InitializeTempBlockCleaner(); err != nil {
		return fmt.Errorf("failed to initialize temp block cleaner: %w", err)
	}

	slog.Info("Dependency injection completed successfully using Factory Pattern")

	return nil
}

func (c *config) InitGlobalService() {
	slog.Debug("Initializing Global Service")

	// Debug: verificar se os adapters estão nil
	if c.repositoryAdapters == nil {
		slog.Error("repositoryAdapters is nil")
		return
	}
	if c.repositoryAdapters.Global == nil {
		slog.Error("repositoryAdapters.Global is nil")
		return
	}
	if c.cep == nil {
		slog.Error("cep adapter is nil")
		return
	}
	if c.firebaseCloudMessaging == nil {
		slog.Error("firebaseCloudMessaging adapter is nil")
		return
	}
	if c.email == nil {
		slog.Error("email adapter is nil")
		return
	}
	if c.sms == nil {
		slog.Error("sms adapter is nil")
		return
	}
	if c.cloudStorage == nil {
		slog.Error("cloudStorage adapter is nil")
		return
	}
	if c.repositoryAdapters.DeviceToken == nil {
		slog.Error("repositoryAdapters.DeviceToken is nil")
		return
	}

	c.globalService = globalservice.NewGlobalService(
		c.repositoryAdapters.Global,
		c.cep,
		c.firebaseCloudMessaging,
		c.email,
		c.sms,
		c.cloudStorage,
		c.repositoryAdapters.DeviceToken,
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
		c.cloudStorage,
		c.permissionService,
	)
	// HTTP handler initialization is done during HTTP server setup
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
		c.cloudStorage,
	)
	// HTTP handler initialization is done during HTTP server setup
}

func (c *config) InitPermissionHandler() {
	slog.Debug("Initializing Permission Handler")
	c.permissionService = permissionservices.NewPermissionService(
		c.repositoryAdapters.Permission,
		c.cache,
		c.globalService,
	)
}

// createActivityTracker inicializa o ActivityTracker com Redis client
func (c *config) createActivityTracker() error {
	if c.cache == nil {
		return fmt.Errorf("cache não inicializado - necessário para ActivityTracker")
	}

	// Obter Redis client do cache
	redisClient := c.cache.GetRedisClient()
	if redisClient == nil {
		return fmt.Errorf("Redis client não disponível no cache")
	}

	// Criar ActivityTracker sem userService (será definido na Phase 07)
	c.activityTracker = goroutines.NewActivityTracker(redisClient, nil)

	slog.Info("✅ ActivityTracker criado com sucesso com Redis client")
	return nil
}
