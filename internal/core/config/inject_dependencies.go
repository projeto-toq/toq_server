package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/factory"
	goroutines "github.com/projeto-toq/toq_server/internal/core/go_routines"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	complexservices "github.com/projeto-toq/toq_server/internal/core/service/complex_service"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	holidayservices "github.com/projeto-toq/toq_server/internal/core/service/holiday_service"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	permissionservices "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
	"github.com/projeto-toq/toq_server/internal/core/utils/hmacauth"
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

	var metrics metricsport.MetricsPortInterface
	if c.metricsAdapter != nil {
		metrics = c.metricsAdapter.Prometheus
	}

	repositories, err := c.adapterFactory.CreateRepositoryAdapters(storage.Database, metrics)
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

	// 5. Inicializar componentes de segurança (HMAC validator)
	if err := c.initializeSecurityComponents(); err != nil {
		return fmt.Errorf("failed to initialize security components: %w", err)
	}

	// 6. Inicializar Services
	c.initializeServices()

	// 7. Inicializar TempBlockCleanerWorker após permission service estar disponível
	if err := c.InitializeTempBlockCleaner(); err != nil {
		return fmt.Errorf("failed to initialize temp block cleaner: %w", err)
	}

	slog.Info("Dependency injection completed successfully using Factory Pattern")

	return nil
}

func (c *config) InitGlobalService() {
	slog.Debug("Initializing Global Service")

	// Optional metrics dependency
	var metrics metricsport.MetricsPortInterface
	if c.metricsAdapter != nil {
		metrics = c.metricsAdapter.Prometheus
	}

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
		metrics,
	)

	// Injetar GlobalService no cache Redis para resolver dependência circular
	if c.cache != nil {
		c.cache.SetGlobalService(c.globalService)
		slog.Debug("GlobalService injected into Redis cache")
	}

	// Start session events subscriber
	if c.globalService != nil {
		_ = c.globalService.StartSessionEventSubscriber() // ignore unsubscribe for now (lifecycle handles full shutdown)
	}
}

func (c *config) InitUserHandler() {
	slog.Debug("Initializing User Handler")
	refreshInterval := time.Duration(c.env.PhotoSession.PhotographerAgendaRefreshIntervalH)
	if refreshInterval <= 0 {
		refreshInterval = 24
	}
	userCfg := userservices.Config{
		SystemUserResetPasswordURL:        c.env.SystemUser.ResetPasswordURL,
		PhotographerTimezone:              c.env.PhotoSession.PhotographerTimezone,
		PhotographerAgendaHorizonMonths:   c.env.PhotoSession.PhotographerHorizonMonths,
		PhotographerAgendaRefreshInterval: refreshInterval * time.Hour,
	}
	c.userService = userservices.NewUserService(
		c.repositoryAdapters.User,
		c.repositoryAdapters.Session,
		c.globalService,
		c.listingService,
		c.photoSessionService,
		c.cpf,
		c.cnpj,
		c.cloudStorage,
		c.permissionService,
		userCfg,
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

func (c *config) InitHolidayService() {
	slog.Debug("Initializing Holiday Service")

	if c.repositoryAdapters == nil {
		slog.Error("repositoryAdapters is nil")
		return
	}

	if c.repositoryAdapters.Holiday == nil {
		slog.Error("repositoryAdapters.Holiday is nil")
		return
	}

	if c.globalService == nil {
		slog.Error("globalService is nil")
		return
	}

	c.holidayService = holidayservices.NewHolidayService(
		c.repositoryAdapters.Holiday,
		c.globalService,
	)
}

func (c *config) InitScheduleService() {
	slog.Debug("Initializing Schedule Service")

	if c.repositoryAdapters == nil {
		slog.Error("repositoryAdapters is nil")
		return
	}

	if c.repositoryAdapters.Schedule == nil {
		slog.Error("repositoryAdapters.Schedule is nil")
		return
	}

	if c.repositoryAdapters.Listing == nil {
		slog.Error("repositoryAdapters.Listing is nil")
		return
	}

	if c.repositoryAdapters.User == nil {
		slog.Error("repositoryAdapters.User is nil")
		return
	}

	if c.globalService == nil {
		slog.Error("globalService is nil")
		return
	}

	serviceConfig, err := scheduleservices.ConfigFromEnvironment(&c.env)
	if err != nil {
		slog.Error("failed to parse schedule configuration", "err", err)
		serviceConfig = scheduleservices.DefaultConfig()
	}

	c.scheduleService = scheduleservices.NewScheduleService(
		c.repositoryAdapters.Schedule,
		c.repositoryAdapters.Listing,
		c.globalService,
		serviceConfig,
	)
}

func (c *config) InitListingHandler() {
	slog.Debug("Initializing Listing Handler")
	c.listingService = listingservices.NewListingService(
		c.repositoryAdapters.Listing,
		c.photoSessionService,
		c.repositoryAdapters.User,
		c.complexService,
		c.globalService,
		c.cloudStorage,
		c.scheduleService,
	)
	// HTTP handler initialization is done during HTTP server setup
}

func (c *config) InitPermissionHandler() {
	slog.Debug("Initializing Permission Handler")
	var metrics metricsport.MetricsPortInterface
	if c.metricsAdapter != nil {
		metrics = c.metricsAdapter.Prometheus
	}
	c.permissionService = permissionservices.NewPermissionService(
		c.repositoryAdapters.Permission,
		c.cache,
		c.globalService,
		metrics,
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
		return fmt.Errorf("redis client não disponível no cache")
	}

	// Criar ActivityTracker sem userService (será definido na Phase 07)
	c.activityTracker = goroutines.NewActivityTracker(redisClient, nil)

	slog.Info("✅ ActivityTracker criado com sucesso com Redis client")
	return nil
}

func (c *config) InitPhotoSessionService() {
	slog.Debug("Initializing Photo Session Service")

	if c.repositoryAdapters == nil || c.repositoryAdapters.PhotoSession == nil {
		slog.Error("repositoryAdapters.PhotoSession is nil")
		return
	}

	if c.holidayService == nil {
		slog.Error("holidayService is nil")
		return
	}

	if c.globalService == nil {
		slog.Error("globalService is nil")
		return
	}

	photoCfg := photosessionservices.Config{
		SlotDurationMinutes: c.env.PhotoSession.SlotDurationMinutes,
		SlotsPerPeriod:      c.env.PhotoSession.SlotsPerPeriod,
		MorningStartHour:    c.env.PhotoSession.MorningStartHour,
		AfternoonStartHour:  c.env.PhotoSession.AfternoonStartHour,
		BusinessStartHour:   c.env.PhotoSession.BusinessStartHour,
		BusinessEndHour:     c.env.PhotoSession.BusinessEndHour,
		AgendaHorizonMonths: c.env.PhotoSession.PhotographerHorizonMonths,
	}

	c.photoSessionService = photosessionservices.NewPhotoSessionService(
		c.repositoryAdapters.PhotoSession,
		c.repositoryAdapters.Listing,
		c.repositoryAdapters.User,
		c.holidayService,
		c.globalService,
		photoCfg,
	)
}

// initializeSecurityComponents ensures security primitives are available before HTTP handler wiring.
func (c *config) initializeSecurityComponents() error {
	if c.hmacValidator != nil {
		slog.Debug("HMAC validator already initialized")
		return nil
	}

	cfg, err := c.GetHMACSecurityConfig()
	if err != nil {
		return err
	}

	validator, err := hmacauth.NewValidator(cfg)
	if err != nil {
		return fmt.Errorf("failed to create HMAC validator: %w", err)
	}

	c.hmacValidator = validator

	// Comentário breve para futuros mantenedores.
	slog.Info("HMAC validator initialized", "encoding", cfg.Encoding, "skew_seconds", cfg.SkewSeconds)
	return nil
}
