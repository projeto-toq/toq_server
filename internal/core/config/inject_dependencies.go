package config

import (
	"log/slog"
	"os"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	cepadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/cep"
	cnpjadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/cnpj"
	cpfadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/cpf"
	creciadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/creci"
	emailadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/email"
	fcmadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/fcm"

	// gcsadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/google_cloud_storage"
	mysqladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql"
	mysqlcomplexadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/complex"
	mysqlglobaladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/global"
	mysqllistingadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing"
	mysqluseradapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user"
	smsadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/sms"

	"github.com/giulio-alfieri/toq_server/internal/core/cache"
	grpclistingport "github.com/giulio-alfieri/toq_server/internal/core/port/left/grpc/listing"
	grpcuserport "github.com/giulio-alfieri/toq_server/internal/core/port/left/grpc/user"
	complexservices "github.com/giulio-alfieri/toq_server/internal/core/service/complex_service.go"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	listingservices "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
	userservices "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
)

func (c *config) InjectDependencies() (close func() error) {
	c.database = mysqladapter.NewDB(c.db)

	var err error
	c.cep, err = cepadapter.NewCEPAdapter(&c.env)
	if err != nil {
		slog.Error("failed to create cep adapter", "error", err)
		os.Exit(1)
	}

	// c.googleCloudStorage, close = gcsadapter.NewGCSAdapter(c.context)

	c.cpf, err = cpfadapter.NewCPFAdapter(&c.env)
	if err != nil {
		slog.Error("failed to create cpf adapter", "error", err)
		os.Exit(1)
	}

	c.cnpj, err = cnpjadapter.NewCNPJAdapater(&c.env)
	if err != nil {
		slog.Error("failed to create cnpj adapter", "error", err)
		os.Exit(1)
	}

	c.creci = creciadapter.NewCreciAdapter(c.context)

	fcm, err := fcmadapter.NewFCMAdapter(c.context, &c.env)
	if err != nil {
		slog.Error("failed to create fcm adapter", "error", err)
		os.Exit(1)
	}
	c.firebaseCloudMessaging = fcm

	c.email = emailadapter.NewEmailAdapter(c.env.EMAIL.SMTPServer, c.env.EMAIL.SMTPPort, c.env.EMAIL.SMTPUser, c.env.EMAIL.SMTPPassword)
	c.sms = smsadapter.NewSmsAdapter(c.env.SMS.AccountSid, c.env.SMS.AuthToken, c.env.SMS.MyNumber)
	c.InitGlobalService()

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(c.env.REDIS.URL, c.globalService)
	if err != nil {
		slog.Error("failed to create redis cache", "error", err)
		os.Exit(1)
	}
	c.cache = redisCache

	c.InitComplexHandler()
	c.InitListingHandler()
	c.InitUserHandler()

	// Return cleanup function to close cache connection
	return func() error {
		return c.cache.Close()
	}
}

func (c *config) InitGlobalService() {
	repo := mysqlglobaladapter.NewGlobalAdapter(c.database)
	c.globalService = globalservice.NewGlobalService(repo, c.cep, c.firebaseCloudMessaging, c.email, c.sms, c.googleCloudStorage)
}

func (c *config) InitUserHandler() {
	repo := mysqluseradapter.NewUserAdapter(c.database)
	c.userService = userservices.NewUserService(repo, c.globalService, c.listingService, c.cpf, c.cnpj, c.creci, c.googleCloudStorage)
	handler := grpcuserport.NewUserHandler(c.userService)
	pb.RegisterUserServiceServer(c.server, handler)
}

func (c *config) InitComplexHandler() {
	repo := mysqlcomplexadapter.NewComplexAdapter(c.database)
	c.complexService = complexservices.NewComplexService(repo, c.globalService)

}

func (c *config) InitListingHandler() {
	repo := mysqllistingadapter.NewListingAdapter(c.database)
	c.listingService = listingservices.NewListingService(repo, c.complexService, c.globalService, c.googleCloudStorage)
	handler := grpclistingport.NewUserHandler(c.listingService)
	pb.RegisterListingServiceServer(c.server, handler)
}
