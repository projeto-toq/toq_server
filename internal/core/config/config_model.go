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
	sessionrepository "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/session_repository"
	smsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/sms"
	storageport "github.com/giulio-alfieri/toq_server/internal/core/port/right/storage"
	complexservices "github.com/giulio-alfieri/toq_server/internal/core/service/complex_service"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
	permissionservices "github.com/giulio-alfieri/toq_server/internal/core/service/permission_service"
	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
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
	wg                     *sync.WaitGroup
	readiness              bool
	globalService          globalservice.GlobalServiceInterface
	userService            userservices.UserServiceInterface
	listingService         listingservices.ListingServiceInterface
	complexService         complexservices.ComplexServiceInterface
	permissionService      permissionservices.PermissionServiceInterface
	cep                    cepport.CEPPortInterface
	cpf                    cpfport.CPFPortInterface
	cnpj                   cnpjport.CNPJPortInterface
	email                  emailport.EmailPortInterface
	sms                    smsport.SMSPortInterface
	cloudStorage           storageport.CloudStoragePortInterface
	firebaseCloudMessaging fcmport.FCMPortInterface
	sessionRepo            sessionrepository.SessionRepoPortInterface
	repositoryAdapters     *factory.RepositoryAdapters
	adapterFactory         factory.AdapterFactory
}

type ConfigInterface interface {
	LoadEnv() error
	InitializeLog()
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
	GetDatabase() *sql.DB
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

// InitializeActivityTracker inicializa o activity tracker (delegado para o novo sistema)
func (c *config) InitializeActivityTracker() error {
	// Este método é mantido para compatibilidade com a interface
	// A implementação real está no phase_03_infrastructure.go
	return nil
}

// InitializeGoRoutines inicializa as goroutines (delegado para o novo sistema)
func (c *config) InitializeGoRoutines() {
	// Este método é mantido para compatibilidade com a interface
	// A implementação real está no phase_07_workers.go
}

// SetActivityTrackerUserService conecta o activity tracker ao user service
func (c *config) SetActivityTrackerUserService() {
	// Este método é mantido para compatibilidade com a interface
	// A implementação real está no phase_07_workers.go
	if c.activityTracker != nil && c.userService != nil {
		c.activityTracker.SetUserService(c.userService)
	}
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

// InitializeLog inicializa o sistema de logging (delegado para o novo sistema)
func (c *config) InitializeLog() {
	// Este método é mantido para compatibilidade com a interface
	// A implementação real está no phase_02_config.go
}

// InitializeDatabase inicializa a conexão com o banco (delegado para o novo sistema)
func (c *config) InitializeDatabase() {
	// Este método é mantido para compatibilidade com a interface
	// A implementação real está no phase_03_infrastructure.go
}

// VerifyDatabase verifica a conexão com o banco (delegado para o novo sistema)
func (c *config) VerifyDatabase() {
	// Este método é mantido para compatibilidade com a interface
	// A implementação real está no phase_03_infrastructure.go
}

// InitializeTelemetry inicializa o sistema de telemetria (delegado para o novo sistema)
func (c *config) InitializeTelemetry() (func(), error) {
	// Este método é mantido para compatibilidade com a interface
	// A implementação real está no phase_03_infrastructure.go
	return func() {}, nil
}

// InitializeHTTP inicializa o servidor HTTP (delegado para o novo sistema)
func (c *config) InitializeHTTP() {
	// Este método é mantido para compatibilidade com a interface
	// A implementação real está no phase_06_handlers.go
}

// SetupHTTPHandlersAndRoutes configura handlers e rotas (delegado para o novo sistema)
func (c *config) SetupHTTPHandlersAndRoutes() {
	// Este método é mantido para compatibilidade com a interface
	// A implementação real está no phase_06_handlers.go
}
