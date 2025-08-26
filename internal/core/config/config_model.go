package config

import (
	"context"
	"database/sql"
	"net"
	"sync"

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
	gcsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/gcs"
	sessionrepository "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/session_repository"
	smsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/sms"
	complexservices "github.com/giulio-alfieri/toq_server/internal/core/service/complex_service"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

type config struct {
	env                    globalmodel.Environment
	db                     *sql.DB
	database               *mysqladapter.Database
	listener               net.Listener
	server                 *grpc.Server
	context                context.Context
	cache                  cache.CacheInterface
	activityTracker        *goroutines.ActivityTracker
	wg                     *sync.WaitGroup
	healthSrv              *health.Server
	readiness              bool
	globalService          globalservice.GlobalServiceInterface
	userService            userservices.UserServiceInterface
	listingService         listingservices.ListingServiceInterface
	complexService         complexservices.ComplexServiceInterface
	cep                    cepport.CEPPortInterface
	cpf                    cpfport.CPFPortInterface
	cnpj                   cnpjport.CNPJPortInterface
	email                  emailport.EmailPortInterface
	sms                    smsport.SMSPortInterface
	googleCloudStorage     gcsport.GCSPortInterface
	firebaseCloudMessaging fcmport.FCMPortInterface
	sessionRepo            sessionrepository.SessionRepoPortInterface
	repositoryAdapters     *factory.RepositoryAdapters
}

type ConfigInterface interface {
	LoadEnv() error
	InitializeLog()
	InitializeDatabase()
	InitializeActivityTracker() error
	VerifyDatabase()
	InitializeTelemetry() (func(), error)
	InitializeGRPC()
	InjectDependencies(*LifecycleManager) error
	InitGlobalService()
	InitUserHandler()
	InitComplexHandler()
	InitListingHandler()
	InitializeGoRoutines()
	SetActivityTrackerUserService()
	GetDatabase() *sql.DB
	GetGRPCServer() *grpc.Server
	GetListener() net.Listener
	GetInfos() (serviceQty int, methodQty int)
	GetWG() *sync.WaitGroup
	GetActivityTracker() *goroutines.ActivityTracker
	SetHealthServing(serving bool)
	StartHTTPHealth()
}

func NewConfig(ctx context.Context) ConfigInterface {
	var wg sync.WaitGroup
	return &config{
		context: ctx,
		wg:      &wg,
	}
}

func (c *config) SetHealthServing(serving bool) {
	if c.healthSrv == nil {
		c.readiness = serving
		return
	}
	if serving {
		c.healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	} else {
		c.healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}
	c.readiness = serving
}
