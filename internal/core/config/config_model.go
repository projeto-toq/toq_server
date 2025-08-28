package config

import (
	"context"
	"database/sql"
	"net/http"
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
		c.httpServer.Shutdown(ctx)
	}
}

func (c *config) GetGinRouter() *gin.Engine {
	return c.ginRouter
}

func (c *config) GetHTTPHandlers() *factory.HTTPHandlers {
	return &c.httpHandlers
}
