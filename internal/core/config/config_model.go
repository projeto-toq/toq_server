package config

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes"
	mysqladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql"
	"github.com/giulio-alfieri/toq_server/internal/core/cache"
	"github.com/giulio-alfieri/toq_server/internal/core/factory"
	goroutines "github.com/giulio-alfieri/toq_server/internal/core/go_routines"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	cepport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cep"
	cnpjport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cnpj"
	cpfport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cpf"
	emailport "github.com/giulio-alfieri/toq_server/internal/core/port/right/email"
	fcmport "github.com/giulio-alfieri/toq_server/internal/core/port/right/fcm"
	smsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/sms"
	storageport "github.com/giulio-alfieri/toq_server/internal/core/port/right/storage"
	complexservices "github.com/giulio-alfieri/toq_server/internal/core/service/complex_service"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
	permissionservices "github.com/giulio-alfieri/toq_server/internal/core/service/permission_service"
	sessionservice "github.com/giulio-alfieri/toq_server/internal/core/service/session_service"
	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
	validationservice "github.com/giulio-alfieri/toq_server/internal/core/service/validation_service"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/hmacauth"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"gopkg.in/yaml.v3"
)

type config struct {
	env                    globalmodel.Environment
	db                     *sql.DB
	database               *mysqladapter.Database
	httpServer             *http.Server
	ginRouter              *gin.Engine
	httpHandlers           factory.HTTPHandlers
	context                context.Context
	cache                  cache.CacheInterface
	activityTracker        *goroutines.ActivityTracker
	tempBlockCleaner       *goroutines.TempBlockCleanerWorker
	sessionService         sessionservice.Service
	wg                     *sync.WaitGroup
	readiness              bool
	globalService          globalservice.GlobalServiceInterface
	userService            userservices.UserServiceInterface
	listingService         listingservices.ListingServiceInterface
	complexService         complexservices.ComplexServiceInterface
	permissionService      permissionservices.PermissionServiceInterface
	metricsAdapter         *factory.MetricsAdapter
	cep                    cepport.CEPPortInterface
	cpf                    cpfport.CPFPortInterface
	cnpj                   cnpjport.CNPJPortInterface
	email                  emailport.EmailPortInterface
	sms                    smsport.SMSPortInterface
	cloudStorage           storageport.CloudStoragePortInterface
	firebaseCloudMessaging fcmport.FCMPortInterface
	repositoryAdapters     *factory.RepositoryAdapters
	adapterFactory         factory.AdapterFactory
	hmacValidator          *hmacauth.Validator
}

type ConfigInterface interface {
	LoadEnv() error
	InitializeDatabase()
	InitializeActivityTracker() error
	VerifyDatabase()
	InitializeTelemetry() (func(), error)
	InitializeHTTP()
	SetupHTTPHandlersAndRoutes()
	InjectDependencies(*LifecycleManager) error
	InitGlobalService()
	InitUserHandler()
	InitComplexHandler()
	InitListingHandler()
	InitPermissionHandler()
	InitializeGoRoutines()
	SetActivityTrackerUserService()
	InitializeTempBlockCleaner() error
	GetDatabase() *sql.DB
	GetEnvironment() (*globalmodel.Environment, error)
	GetHTTPServer() *http.Server
	CloseHTTPServer()
	GetGinRouter() *gin.Engine
	GetHTTPHandlers() *factory.HTTPHandlers
	GetWG() *sync.WaitGroup
	GetActivityTracker() *goroutines.ActivityTracker
	SetHealthServing(serving bool)
}

func NewConfig(ctx context.Context) ConfigInterface {
	var wg sync.WaitGroup
	return &config{
		context: ctx,
		wg:      &wg,
	}
}

func (c *config) SetHealthServing(serving bool) {
	c.readiness = serving
}

func (c *config) GetHTTPServer() *http.Server {
	return c.httpServer
}

func (c *config) CloseHTTPServer() {
	if c.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := c.httpServer.Shutdown(ctx); err != nil {
			slog.Error("Error shutting down HTTP server", "error", err)
		}
	}
}

func (c *config) GetGinRouter() *gin.Engine {
	return c.ginRouter
}

func (c *config) GetHTTPHandlers() *factory.HTTPHandlers {
	return &c.httpHandlers
}

// assignStorageAdapters atribui os adapters de armazenamento ao config
func (c *config) assignStorageAdapters(storage factory.StorageAdapters) {
	c.database = storage.Database
	c.cache = storage.Cache
	c.db = storage.Database.DB
}

// assignRepositoryAdapters atribui os adapters de repositório ao config
func (c *config) assignRepositoryAdapters(repositories factory.RepositoryAdapters) {
	slog.Info("Assigning repository adapters")

	// Log para debug
	if repositories.User == nil {
		slog.Error("repositories.User is nil")
	}
	if repositories.Global == nil {
		slog.Error("repositories.Global is nil")
	}
	if repositories.Complex == nil {
		slog.Error("repositories.Complex is nil")
	}
	if repositories.Listing == nil {
		slog.Error("repositories.Listing is nil")
	}
	if repositories.Session == nil {
		slog.Error("repositories.Session is nil")
	}
	if repositories.Permission == nil {
		slog.Error("repositories.Permission is nil")
	}
	if repositories.DeviceToken == nil {
		slog.Error("repositories.DeviceToken is nil")
	}

	// Criar uma cópia dos repositórios para evitar problemas com ponteiros
	c.repositoryAdapters = &factory.RepositoryAdapters{
		User:        repositories.User,
		Global:      repositories.Global,
		Complex:     repositories.Complex,
		Listing:     repositories.Listing,
		Session:     repositories.Session,
		Permission:  repositories.Permission,
		DeviceToken: repositories.DeviceToken,
	}

	slog.Info("Repository adapters assigned successfully")
}

// assignValidationAdapters atribui os adapters de validação ao config
func (c *config) assignValidationAdapters(validation factory.ValidationAdapters) {
	c.cep = validation.CEP
	c.cpf = validation.CPF
	c.cnpj = validation.CNPJ
}

// assignExternalServiceAdapters atribui os adapters de serviços externos ao config
func (c *config) assignExternalServiceAdapters(external factory.ExternalServiceAdapters) {
	c.firebaseCloudMessaging = external.FCM
	c.email = external.Email
	c.sms = external.SMS
	c.cloudStorage = external.CloudStorage
}

// initializeServices inicializa todos os serviços do sistema
func (c *config) initializeServices() {
	slog.Info("Initializing all services")

	// Inicializar serviços na ordem correta (resolvendo dependências)
	c.InitGlobalService()
	c.InitPermissionHandler()
	c.InitComplexHandler()
	c.InitListingHandler()
	c.InitUserHandler()

	slog.Info("All services initialized successfully")
}

func (c *config) GetDatabase() *sql.DB {
	return c.db
}

func (c *config) GetWG() *sync.WaitGroup {
	return c.wg
}

// InitializeActivityTracker verifica se o activity tracker foi criado
func (c *config) InitializeActivityTracker() error {
	if c.activityTracker == nil {
		slog.Error("ActivityTracker não foi criado na Phase 04 - falha na inicialização")
		return fmt.Errorf("ActivityTracker não foi inicializado")
	}

	slog.Info("✅ ActivityTracker verificado e disponível")
	return nil
}

// InitializeGoRoutines inicializa as goroutines do sistema
func (c *config) InitializeGoRoutines() {
	baseCtx := coreutils.ContextWithLogger(c.context)
	logger := coreutils.LoggerFromContext(baseCtx)

	if c.activityTracker != nil && c.wg != nil {
		// Iniciar worker do activity tracker
		c.wg.Add(1)
		go c.activityTracker.StartBatchWorker(c.wg, coreutils.ContextWithLogger(baseCtx))
		logger.Info("Activity tracker batch worker started")
	} else {
		logger.Warn("Activity tracker or wait group not available for goroutine initialization")
	}

	if c.tempBlockCleaner != nil {
		// Iniciar worker de limpeza de bloqueios temporários
		c.wg.Add(1)
		go func(workerCtx context.Context) {
			defer c.wg.Done()
			c.tempBlockCleaner.Start(workerCtx)
		}(coreutils.ContextWithLogger(baseCtx))
		logger.Info("Temp block cleaner worker started")
	} else {
		logger.Warn("Temp block cleaner not available for goroutine initialization")
	}

	// Start session cleaner using service (if session repo and global service are set)
	if c.repositoryAdapters != nil && c.repositoryAdapters.Session != nil && c.globalService != nil {
		c.sessionService = sessionservice.New(c.repositoryAdapters.Session, c.globalService)
		intervalSecs := c.env.AUTH.SessionCleanerIntervalSeconds
		if intervalSecs <= 0 {
			intervalSecs = 60
		}
		c.wg.Add(1)
		go goroutines.SessionCleaner(c.sessionService, c.wg, coreutils.ContextWithLogger(baseCtx), time.Duration(intervalSecs)*time.Second)
		logger.Info("Session cleaner worker started", "interval_seconds", intervalSecs)
	} else {
		logger.Warn("Session cleaner prerequisites not met; skipping start")
	}

	// Start validation cleaner if user repository and global service are set
	if c.repositoryAdapters != nil && c.repositoryAdapters.User != nil && c.globalService != nil {
		validationSvc := validationservice.New(c.repositoryAdapters.User, c.globalService)
		intervalSecs := c.env.AUTH.ValidationCleanerIntervalSeconds
		if intervalSecs <= 0 {
			intervalSecs = 300 // default 5 minutes
		}
		c.wg.Add(1)
		go func(workerCtx context.Context) {
			defer c.wg.Done()
			goroutines.ValidationCleaner(validationSvc, time.Duration(intervalSecs)*time.Second, workerCtx)
		}(coreutils.ContextWithLogger(baseCtx))
		logger.Info("Validation cleaner worker started", "interval_seconds", intervalSecs)
	} else {
		logger.Warn("Validation cleaner prerequisites not met; skipping start")
	}
}

// SetActivityTrackerUserService conecta o activity tracker ao user service
func (c *config) SetActivityTrackerUserService() {
	if c.activityTracker != nil && c.userService != nil {
		c.activityTracker.SetUserService(c.userService)
		slog.Info("Activity tracker connected to user service")
	} else {
		slog.Warn("Activity tracker or user service not available for connection")
	}
}

// InitializeTempBlockCleaner inicializa o worker de limpeza de bloqueios temporários
func (c *config) InitializeTempBlockCleaner() error {
	if c.permissionService == nil {
		slog.Error("Permission service not available for temp block cleaner initialization")
		return fmt.Errorf("permission service not initialized")
	}

	if c.globalService == nil {
		slog.Error("Global service not available for temp block cleaner initialization")
		return fmt.Errorf("global service not initialized")
	}

	c.tempBlockCleaner = goroutines.NewTempBlockCleanerWorker(c.permissionService, c.globalService)
	slog.Info("✅ TempBlockCleanerWorker initialized")
	return nil
}

// GetTempBlockCleaner retorna o worker de limpeza de bloqueios temporários
func (c *config) GetTempBlockCleaner() *goroutines.TempBlockCleanerWorker {
	return c.tempBlockCleaner
}

// LoadEnv carrega as variáveis de ambiente e configuração YAML
func (c *config) LoadEnv() error {
	// Carregar configuração do arquivo YAML
	configPath := "configs/env.yaml"

	// Ler arquivo de configuração
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Fazer parse do YAML
	env := &globalmodel.Environment{}
	if err := yaml.Unmarshal(data, env); err != nil {
		return fmt.Errorf("failed to parse config YAML: %w", err)
	}

	// Armazenar no config
	c.env = *env

	slog.Info("Configuration loaded successfully from YAML", "path", configPath)
	return nil
}

// InitializeLog inicializa o sistema de logging

// InitializeDatabase inicializa a conexão com o banco de dados
func (c *config) InitializeDatabase() {
	if c.database != nil {
		slog.Info("Database already initialized")
		return
	}

	// Abrir conexão MySQL
	db, err := sql.Open("mysql", c.env.DB.URI)
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		return
	}

	// Testar conexão
	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", "error", err)
		db.Close()
		return
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Criar wrapper Database
	c.database = mysqladapter.NewDB(db)
	c.db = db

	slog.Info("Database connection initialized", "uri", c.env.DB.URI)
}

// VerifyDatabase verifica a conexão com o banco de dados
func (c *config) VerifyDatabase() {
	if c.db == nil {
		slog.Error("Database connection not initialized")
		return
	}

	// Testar conexão
	if err := c.db.Ping(); err != nil {
		slog.Error("Database connection verification failed", "error", err)
		return
	}

	slog.Info("Database connection verified successfully")
}

// InitializeTelemetry inicializa o sistema de telemetria OpenTelemetry
func (c *config) InitializeTelemetry() (func(), error) {
	if !c.env.TELEMETRY.Enabled {
		slog.Info("OpenTelemetry disabled by configuration")
		return func() {}, nil
	}

	// Usar o novo TelemetryManager
	telemetryManager := NewTelemetryManager(c.env)

	shutdownFunc, err := telemetryManager.Initialize(c.context)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize telemetry: %w", err)
	}

	slog.Info("OpenTelemetry initialized successfully")
	return shutdownFunc, nil
}

// InitializeHTTP inicializa o servidor HTTP (implementação real)
func (c *config) InitializeHTTP() {
	if err := c.SetupHTTPServer(); err != nil {
		slog.Error("Failed to setup HTTP server", "error", err)
		return
	}
	slog.Info("HTTP server initialization completed")
}
func (c *config) SetupHTTPServer() error {
	if c.ginRouter == nil {
		c.ginRouter = gin.New() // Usar gin.New() para controle manual dos middlewares
	}

	c.httpServer = &http.Server{
		Addr:           c.env.HTTP.Port,
		Handler:        c.ginRouter,
		ReadTimeout:    parseDuration(c.env.HTTP.ReadTimeout),
		WriteTimeout:   parseDuration(c.env.HTTP.WriteTimeout),
		MaxHeaderBytes: c.env.HTTP.MaxHeaderBytes,
	}

	slog.Info("HTTP server initialized", "port", c.env.HTTP.Port, "read_timeout", c.env.HTTP.ReadTimeout, "write_timeout", c.env.HTTP.WriteTimeout)
	return nil
}

// parseDuration converte string de duração para time.Duration
func parseDuration(durationStr string) time.Duration {
	if durationStr == "" {
		return 30 * time.Second // default
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		slog.Warn("Invalid duration format, using default", "duration", durationStr, "error", err)
		return 30 * time.Second
	}
	return duration
}

// SetupHTTPHandlersAndRoutes configura handlers e rotas
func (c *config) SetupHTTPHandlersAndRoutes() {
	if c.ginRouter == nil {
		slog.Error("Gin router not initialized")
		return
	}

	// 1. Criar handlers via factory pattern
	if err := c.createHTTPHandlers(); err != nil {
		slog.Error("Failed to create HTTP handlers", "error", err)
		return
	}

	// 2. Registrar todas as rotas (auth, user, listing) via routes package com dependências injetadas
	// Isso aplicará os middlewares globais primeiro
	routes.SetupRoutes(c.ginRouter, &c.httpHandlers, c.activityTracker, c.permissionService, c.metricsAdapter, NewStaticAPIVersionProvider())

	// 3. Configurar rotas básicas de health check APÓS middlewares globais serem aplicados
	c.setupBasicRoutes()

	slog.Info("HTTP handlers and routes configured successfully")
}

// createHTTPHandlers creates HTTP handlers using the factory pattern
func (c *config) createHTTPHandlers() error {
	slog.Debug("Creating HTTP handlers via factory pattern")

	if c.userService == nil || c.globalService == nil || c.listingService == nil || c.complexService == nil {
		return fmt.Errorf("required services not initialized")
	}

	if c.hmacValidator == nil {
		validator, err := hmacauth.NewValidator(c.env.GetHMACSecurityConfig())
		if err != nil {
			return fmt.Errorf("failed to initialize HMAC validator: %w", err)
		}
		c.hmacValidator = validator
	}

	// Create handlers using the pre-initialized factory instance
	c.httpHandlers = c.adapterFactory.CreateHTTPHandlers(
		c.userService,
		c.globalService,
		c.listingService,
		c.complexService,
		c.permissionService,
		c.metricsAdapter,
		c.hmacValidator,
	)

	slog.Info("✅ HTTP handlers created successfully via factory")
	return nil
}

// setupBasicRoutes configura rotas básicas com middlewares de métricas
func (c *config) setupBasicRoutes() {
	// Aplicar middlewares de métricas às rotas básicas
	var metricsMiddleware gin.HandlerFunc
	if c.metricsAdapter != nil {
		metricsMiddleware = middlewares.TelemetryMiddleware(c.metricsAdapter.Prometheus)
	} else {
		metricsMiddleware = gin.HandlerFunc(func(ctx *gin.Context) { ctx.Next() })
	}

	// Health check endpoints com middleware de métricas
	c.ginRouter.GET("/healthz", metricsMiddleware, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	c.ginRouter.GET("/readyz", metricsMiddleware, func(ctx *gin.Context) {
		if c.readiness {
			ctx.JSON(200, gin.H{"status": "ready"})
		} else {
			ctx.JSON(503, gin.H{"status": "not ready"})
		}
	})

	// Metrics endpoint (sem middleware adicional para evitar recursão)
	if c.metricsAdapter != nil && c.httpHandlers.MetricsHandler != nil {
		if metricsHandler, ok := c.httpHandlers.MetricsHandler.(interface{ GetMetrics(c *gin.Context) }); ok {
			c.ginRouter.GET("/metrics", metricsHandler.GetMetrics)
		}
	}

	// API base group (v2)
	v1 := c.ginRouter.Group(NewStaticAPIVersionProvider().BasePath())
	{
		v1.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "pong"})
		})
	}
}
