package config

import (
	"context"
	"database/sql"
	"net"
	"sync"

	mysqladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql"
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
	complexservices "github.com/giulio-alfieri/toq_server/internal/core/service/complex_service.go"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
	"google.golang.org/grpc"
)

type config struct {
	env                    globalmodel.Environment
	db                     *sql.DB
	database               *mysqladapter.Database
	listener               net.Listener
	server                 *grpc.Server
	context                context.Context
	cache                  cache.CacheInterface
	wg                     *sync.WaitGroup
	activity               chan int64
	globalService          globalservice.GlobalServiceInterface
	userService            userservices.UserServiceInterface
	listingService         listingservices.ListingServiceInterface
	complexService         complexservices.ComplexServiceInterface
	cep                    cepport.CEPPortInterface
	cpf                    cpfport.CPFPortInterface
	cnpj                   cnpjport.CNPJPortInterface
	creci                  creciport.CreciPortInterface
	email                  emailport.EmailPortInterface
	sms                    smsport.SMSPortInterface
	googleCloudStorage     gcsport.GCSPortInterface
	firebaseCloudMessaging fcmport.FCMPortInterface
}

type ConfigInterface interface {
	LoadEnv()
	InitializeLog()
	InitializeDatabase()
	VerifyDatabase()
	InitializeTelemetry() func()
	InitializeGRPC()
	InjectDependencies() func() error
	InitGlobalService()
	InitUserHandler()
	InitComplexHandler()
	InitListingHandler()
	InitializeGoRoutines()
	GetDatabase() *sql.DB
	GetGRPCServer() *grpc.Server
	GetListener() net.Listener
	GetInfos() (serviceQty int, methodQty int)
	GetWG() *sync.WaitGroup
}

func NewConfig(ctx context.Context, activity chan int64) ConfigInterface {
	var wg sync.WaitGroup
	return &config{
		context:  ctx,
		activity: activity,
		wg:       &wg,
	}
}
