package config

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/routes"
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	"github.com/projeto-toq/toq_server/internal/core/cache"
	"github.com/projeto-toq/toq_server/internal/core/factory"
	goroutines "github.com/projeto-toq/toq_server/internal/core/go_routines"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	httpport "github.com/projeto-toq/toq_server/internal/core/port/left/http"
	cepport "github.com/projeto-toq/toq_server/internal/core/port/right/cep"
	cnpjport "github.com/projeto-toq/toq_server/internal/core/port/right/cnpj"
	cpfport "github.com/projeto-toq/toq_server/internal/core/port/right/cpf"
	emailport "github.com/projeto-toq/toq_server/internal/core/port/right/email"
	fcmport "github.com/projeto-toq/toq_server/internal/core/port/right/fcm"
	smsport "github.com/projeto-toq/toq_server/internal/core/port/right/sms"
	storageport "github.com/projeto-toq/toq_server/internal/core/port/right/storage"
	complexservices "github.com/projeto-toq/toq_server/internal/core/service/complex_service"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	holidayservices "github.com/projeto-toq/toq_server/internal/core/service/holiday_service"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	permissionservices "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
	sessionservice "github.com/projeto-toq/toq_server/internal/core/service/session_service"
	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
	validationservice "github.com/projeto-toq/toq_server/internal/core/service/validation_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
	"github.com/projeto-toq/toq_server/internal/core/utils/hmacauth"
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
	holidayService         holidayservices.HolidayServiceInterface
	scheduleService        scheduleservices.ScheduleServiceInterface
	photoSessionService    photosessionservices.PhotoSessionServiceInterface
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
	runtimeEnvironment     string
	workersEnabled         bool
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
	InitHolidayService()
	InitScheduleService()
	InitPhotoSessionService()
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
	GetRuntimeEnvironment() string
	AreWorkersEnabled() bool
	SetHealthServing(serving bool)
	GetMaxWrongSigninAttempts() int
	GetTempBlockDuration() time.Duration
	httpport.APIVersionProvider
}

func NewConfig(ctx context.Context) ConfigInterface {
	var wg sync.WaitGroup
	return &config{
		context:            ctx,
		wg:                 &wg,
		workersEnabled:     true,
		runtimeEnvironment: "homo",
	}
}

func (c *config) SetHealthServing(serving bool) {
	c.readiness = serving
}

func (c *config) GetRuntimeEnvironment() string {
	return c.runtimeEnvironment
}

func (c *config) AreWorkersEnabled() bool {
	return c.workersEnabled
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
	if repositories.Holiday == nil {
		slog.Error("repositories.Holiday is nil")
	}
	if repositories.Schedule == nil {
		slog.Error("repositories.Schedule is nil")
	}
	if repositories.Visit == nil {
		slog.Error("repositories.Visit is nil")
	}
	if repositories.PhotoSession == nil {
		slog.Error("repositories.PhotoSession is nil")
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
		User:         repositories.User,
		Global:       repositories.Global,
		Complex:      repositories.Complex,
		Listing:      repositories.Listing,
		Holiday:      repositories.Holiday,
		Schedule:     repositories.Schedule,
		Visit:        repositories.Visit,
		PhotoSession: repositories.PhotoSession,
		Session:      repositories.Session,
		Permission:   repositories.Permission,
		DeviceToken:  repositories.DeviceToken,
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
	c.InitHolidayService()
	c.InitScheduleService()
	c.InitPhotoSessionService()
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
	if !c.workersEnabled {
		slog.Info("Workers desabilitados para este ambiente; ignorando inicialização de goroutines")
		return
	}
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
	if !c.workersEnabled {
		slog.Info("Workers desabilitados; link entre ActivityTracker e UserService não será aplicado")
		return
	}
	if c.activityTracker != nil && c.userService != nil {
		c.activityTracker.SetUserService(c.userService)
		slog.Info("Activity tracker connected to user service")
	} else {
		slog.Warn("Activity tracker or user service not available for connection")
	}
}

// InitializeTempBlockCleaner inicializa o worker de limpeza de bloqueios temporários
func (c *config) InitializeTempBlockCleaner() error {
	if !c.workersEnabled {
		slog.Info("Workers desabilitados; TempBlockCleaner não será inicializado")
		return nil
	}
	if c.userService == nil {
		slog.Error("User service not available for temp block cleaner initialization")
		return fmt.Errorf("user service not initialized")
	}

	if c.globalService == nil {
		slog.Error("Global service not available for temp block cleaner initialization")
		return fmt.Errorf("global service not initialized")
	}

	c.tempBlockCleaner = goroutines.NewTempBlockCleanerWorker(c.userService, c.globalService)
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
		slog.Error("Failed to open MySQL connection", "error", err)
		return
	}

	c.db = db
	c.database = mysqladapter.NewDB(db)
	slog.Info("Database connection opened successfully")
}

// VerifyDatabase verifica a conexão com o banco de dados
func (c *config) VerifyDatabase() {
	if c.db == nil {
		slog.Error("Database connection is nil")
		return
	}
	err := c.db.Ping()
	if err != nil {
		slog.Error("Failed to ping database", "error", err)
	} else {
		slog.Info("Database connection verified")
	}
}

// InitializeHTTP inicializa o servidor HTTP e configura o Gin
func (c *config) InitializeHTTP() {
	// Criar router Gin com modo de debug ou release
	if c.env.APP.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware de recuperação de pânico
	router.Use(gin.Recovery())

	// Configurar logger do Gin para usar o logger do sistema
	router.Use(func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodGet && ctx.Request.URL.Path == "/metrics" {
			ctx.Next()
			return
		}
		coreutils.LoggerFromContext(ctx).Info("Request received",
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
			"remote_addr", ctx.Request.RemoteAddr,
		)
		ctx.Next()
	})

	c.ginRouter = router

	// Parse timeouts from string to time.Duration
	readTimeout, err := time.ParseDuration(c.env.HTTP.ReadTimeout)
	if err != nil {
		slog.Warn("Failed to parse ReadTimeout, using default", "value", c.env.HTTP.ReadTimeout, "error", err)
		readTimeout = 10 * time.Second
	}
	writeTimeout, err := time.ParseDuration(c.env.HTTP.WriteTimeout)
	if err != nil {
		slog.Warn("Failed to parse WriteTimeout, using default", "value", c.env.HTTP.WriteTimeout, "error", err)
		writeTimeout = 10 * time.Second
	}
	idleTimeout, err := time.ParseDuration(c.env.HTTP.IdleTimeout)
	if err != nil {
		slog.Warn("Failed to parse IdleTimeout, using default", "value", c.env.HTTP.IdleTimeout, "error", err)
		idleTimeout = 60 * time.Second
	}

	// Configurar servidor HTTP
	c.httpServer = &http.Server{
		Addr:         c.env.HTTP.Port,
		Handler:      c.ginRouter,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
}

// SetupHTTPHandlersAndRoutes configura os handlers e as rotas HTTP
func (c *config) SetupHTTPHandlersAndRoutes() {
	// Criar handlers HTTP
	c.httpHandlers = c.adapterFactory.CreateHTTPHandlers(
		c.userService,
		c.globalService,
		c.listingService,
		c.complexService,
		c.scheduleService,
		c.holidayService,
		c.permissionService,
		c.photoSessionService,
		c.metricsAdapter,
		c.hmacValidator,
	)

	// Configurar rotas
	routes.SetupRoutes(
		c.ginRouter,
		&c.httpHandlers,
		c.activityTracker,
		c.permissionService,
		c.metricsAdapter,
		c, // Passa o config como APIVersionProvider
	)
}

// BasePath retorna o caminho base da API
func (c *config) BasePath() string {
	return "/api/v2"
}

// Version retorna a versão da API
func (c *config) Version() string {
	return "v2"
}

// InitializeTelemetry inicializa o sistema de telemetria
func (c *config) InitializeTelemetry() (func(), error) {
	tm := NewTelemetryManager(c.env, c.runtimeEnvironment)
	return tm.Initialize(c.context)
}
