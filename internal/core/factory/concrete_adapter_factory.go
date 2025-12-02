package factory

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/core/cache"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"

	// HTTP handlers
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers"
	adminhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/admin_handlers"
	authhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/auth_handlers"
	holidayhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/holiday_handlers"
	listinghandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/listing_handlers"
	mediaprocessinghandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/media_processing_handlers"
	photosessionhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/photo_session_handlers"
	schedulehandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/schedule_handlers"
	userhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/user_handlers"

	// Metrics adapter
	prometheusadapter "github.com/projeto-toq/toq_server/internal/adapter/right/prometheus"

	// Validation adapters
	cepadapter "github.com/projeto-toq/toq_server/internal/adapter/right/cep"
	cnpjadapter "github.com/projeto-toq/toq_server/internal/adapter/right/cnpj"
	cpfadapter "github.com/projeto-toq/toq_server/internal/adapter/right/cpf"

	// External service adapters
	emailadapter "github.com/projeto-toq/toq_server/internal/adapter/right/email"
	fcmadapter "github.com/projeto-toq/toq_server/internal/adapter/right/fcm"
	smsadapter "github.com/projeto-toq/toq_server/internal/adapter/right/sms"

	// Storage adapters - AWS S3 (substituindo GCS)
	s3adapter "github.com/projeto-toq/toq_server/internal/adapter/right/aws_s3"
	sqsmediaprocessingadapter "github.com/projeto-toq/toq_server/internal/adapter/right/aws_sqs/media_processing"
	stepfunctionscallbackadapter "github.com/projeto-toq/toq_server/internal/adapter/right/step_functions"

	// Storage adapters
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"

	// Repository adapters
	mysqlglobaladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/global"
	mysqlholidayadapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/holiday"
	mysqllistingadapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing"
	mysqlmediaprocessingadapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing"
	mysqlpermissionadapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission"
	mysqlphotosessionadapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session"
	mysqlpropertycoverageadapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage"
	mysqlscheduleadapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule"
	sessionmysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/session"
	mysqluseradapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user"
	mysqlvisitadapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit"

	// Core services
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	holidayservice "github.com/projeto-toq/toq_server/internal/core/service/holiday_service"
	listingservice "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	mediaprocessingservice "github.com/projeto-toq/toq_server/internal/core/service/media_processing_service"
	permissionservice "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	photosessionservice "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
	propertycoverageservice "github.com/projeto-toq/toq_server/internal/core/service/property_coverage_service"
	scheduleservice "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
	userservice "github.com/projeto-toq/toq_server/internal/core/service/user_service"
	"github.com/projeto-toq/toq_server/internal/core/utils/hmacauth"

	mediaprocessingcallbackport "github.com/projeto-toq/toq_server/internal/core/port/right/functions/mediaprocessingcallback"
	mediaprocessingqueue "github.com/projeto-toq/toq_server/internal/core/port/right/queue/mediaprocessingqueue"
)

// ConcreteAdapterFactory implementa a interface AdapterFactory
// Responsável pela criação concreta de todos os adapters do sistema
type ConcreteAdapterFactory struct {
	lm LifecycleManager
}

// CreateMetricsAdapter cria o adapter de métricas Prometheus
func (f *ConcreteAdapterFactory) CreateMetricsAdapter(runtimeEnv string) *MetricsAdapter {
	slog.Info("Creating metrics adapter")

	prometheusAdapter := prometheusadapter.NewPrometheusAdapter()

	return &MetricsAdapter{
		Prometheus: prometheusAdapter,
	}
}

// CreateValidationAdapters cria e configura todos os adapters de validação externa
// Retorna ValidationAdapters com CEP, CPF e CNPJ configurados
func (f *ConcreteAdapterFactory) CreateValidationAdapters(env *globalmodel.Environment) (ValidationAdapters, error) {
	slog.Info("Creating validation adapters")

	// CEP Adapter
	cep, err := cepadapter.NewCEPAdapter(env)
	if err != nil {
		return ValidationAdapters{}, fmt.Errorf("failed to create CEP adapter: %w", err)
	}

	// CPF Adapter
	cpf, err := cpfadapter.NewCPFAdapter(env)
	if err != nil {
		return ValidationAdapters{}, fmt.Errorf("failed to create CPF adapter: %w", err)
	}

	// CNPJ Adapter
	cnpj, err := cnpjadapter.NewCNPJAdapter(env)
	if err != nil {
		return ValidationAdapters{}, fmt.Errorf("failed to create CNPJ adapter: %w", err)
	}

	slog.Info("Successfully created all validation adapters")

	return ValidationAdapters{
		CEP:  cep,
		CPF:  cpf,
		CNPJ: cnpj,
	}, nil
}

// CreateExternalServiceAdapters cria adapters para serviços externos
// Inclui FCM (push notifications), Email e SMS
func (f *ConcreteAdapterFactory) CreateExternalServiceAdapters(ctx context.Context, env *globalmodel.Environment) (ExternalServiceAdapters, error) {
	slog.Info("Creating external service adapters")

	// FCM Adapter
	fcm, err := fcmadapter.NewFCMAdapter(ctx, env)
	if err != nil {
		return ExternalServiceAdapters{}, fmt.Errorf("failed to create FCM adapter: %w", err)
	}

	// Email Adapter com configuração robusta
	email := emailadapter.NewEmailAdapter(*env)

	// SMS Adapter
	sms := smsadapter.NewSmsAdapter(
		env.SMS.AccountSid,
		env.SMS.AuthToken,
		env.SMS.MyNumber,
	)

	// S3 Adapter (substituindo GCS)
	s3, s3Close, err := s3adapter.NewS3Adapter(ctx, env)
	if err != nil {
		slog.Warn("failed to create S3 adapter, proceeding without it", "error", err)
		// Não retorna erro, permite que a aplicação continue sem S3
	}

	listingMediaStorage := s3adapter.NewListingMediaStorageAdapter(s3, env)

	var mediaQueue mediaprocessingqueue.QueuePortInterface
	queueAdapter, err := sqsmediaprocessingadapter.NewMediaProcessingQueueAdapter(ctx, env)
	if err != nil {
		slog.Warn("failed to create media processing queue adapter", "error", err)
	} else if queueAdapter != nil {
		mediaQueue = queueAdapter
	} else {
		slog.Warn("media processing queue adapter returned nil (configuration missing?)")
	}

	callbackAdapter := stepfunctionscallbackadapter.NewMediaProcessingCallbackAdapter(env)

	slog.Info("Successfully created all external service adapters")

	return ExternalServiceAdapters{
		FCM:                     fcm,
		Email:                   email,
		SMS:                     sms,
		CloudStorage:            s3, // S3 adapter via interface CloudStorage
		ListingMediaStorage:     listingMediaStorage,
		MediaProcessingQueue:    mediaQueue,
		MediaProcessingCallback: callbackAdapter,
		CloseFunc:               s3Close, // Função de cleanup do S3
	}, nil
}

// CreateStorageAdapters cria adapters de armazenamento (Database e Cache)
// Inclui MySQL database e Redis cache com função de cleanup
func (f *ConcreteAdapterFactory) CreateStorageAdapters(ctx context.Context, env *globalmodel.Environment, db *sql.DB, metrics metricsport.MetricsPortInterface) (StorageAdapters, error) {
	slog.Info("Creating storage adapters")

	// Database Adapter
	database := mysqladapter.NewDB(db)

	// Redis Cache - criar sem GlobalService inicialmente (será injetado posteriormente)
	redisCache, err := cache.NewRedisCache(env.REDIS.URL, nil, metrics)
	if err != nil {
		return StorageAdapters{}, fmt.Errorf("failed to create Redis cache: %w", err)
	}

	// Função de cleanup para fechar recursos
	closeFunc := func() error {
		if redisCache != nil {
			return redisCache.Close()
		}
		return nil
	}

	slog.Info("Successfully created all storage adapters")

	return StorageAdapters{
		Database:  database,
		Cache:     redisCache,
		CloseFunc: closeFunc,
	}, nil
}

// CreateRepositoryAdapters cria todos os repositórios MySQL
// Agrupa repositórios por domínio (User, Global, Complex, Listing, Session, Permission)
func (f *ConcreteAdapterFactory) CreateRepositoryAdapters(database *mysqladapter.Database, metrics metricsport.MetricsPortInterface) (RepositoryAdapters, error) {
	slog.Info("Creating repository adapters")

	// User Repository (with device token management integrated)
	userRepo := mysqluseradapter.NewUserAdapter(database, metrics)

	// Global Repository
	globalRepo := mysqlglobaladapter.NewGlobalAdapter(database, metrics)

	// Listing Repository
	listingRepo := mysqllistingadapter.NewListingAdapter(database, metrics)

	// Property Coverage Repository
	propertyCoverageRepo := mysqlpropertycoverageadapter.NewPropertyCoverageAdapter(database, metrics)

	// Holiday Repository
	holidayRepo := mysqlholidayadapter.NewHolidayAdapter(database, metrics)

	// Schedule Repository
	scheduleRepo := mysqlscheduleadapter.NewScheduleAdapter(database, metrics)

	// Visit Repository
	visitRepo := mysqlvisitadapter.NewVisitAdapter(database, metrics)

	// Photo Session Repository
	photoSessionRepo := mysqlphotosessionadapter.NewPhotoSessionAdapter(database, metrics)

	// Session Repository
	sessionRepo := sessionmysqladapter.NewSessionAdapter(database, metrics)

	// Permission Repository
	permissionRepo := mysqlpermissionadapter.NewPermissionAdapter(database, metrics)

	// Media Processing Repository
	mediaProcessingRepo := mysqlmediaprocessingadapter.NewMediaProcessingAdapter(database, metrics)

	slog.Info("Successfully created all repository adapters")

	return RepositoryAdapters{
		User:             userRepo,
		Global:           globalRepo,
		PropertyCoverage: propertyCoverageRepo,
		Listing:          listingRepo,
		MediaProcessing:  mediaProcessingRepo,
		Holiday:          holidayRepo,
		Schedule:         scheduleRepo,
		Visit:            visitRepo,
		PhotoSession:     photoSessionRepo,
		Session:          sessionRepo,
		Permission:       permissionRepo,
	}, nil
}

// CreateHTTPHandlers creates and returns all HTTP handlers
func (factory *ConcreteAdapterFactory) CreateHTTPHandlers(
	router *gin.Engine,
	userService userservice.UserServiceInterface,
	globalService globalservice.GlobalServiceInterface,
	listingService listingservice.ListingServiceInterface,
	propertyCoverageService propertycoverageservice.PropertyCoverageServiceInterface,
	scheduleService scheduleservice.ScheduleServiceInterface,
	holidayService holidayservice.HolidayServiceInterface,
	permissionService permissionservice.PermissionServiceInterface,
	photoSessionService photosessionservice.PhotoSessionServiceInterface,
	mediaProcessingService mediaprocessingservice.MediaProcessingServiceInterface,
	metricsAdapter *MetricsAdapter,
	callbackValidator mediaprocessingcallbackport.CallbackPortInterface,
	hmacValidator *hmacauth.Validator,
) HTTPHandlers {
	slog.Info("Creating HTTP handlers")

	// Create user handler using the adapter
	userHandlerPort := userhandlers.NewUserHandlerAdapter(
		userService,
		globalService,
		permissionService,
	)

	userHandler, ok := userHandlerPort.(*userhandlers.UserHandler)
	if !ok {
		slog.Error("factory.http_handlers.user_cast_failed")
		return HTTPHandlers{}
	}

	// Create auth handler using the adapter
	authHandler := authhandlers.NewAuthHandlerAdapter(
		userService,
		globalService,
		hmacValidator,
	)

	// Create listing handler using the adapter
	// Note: MediaProcessingService may be nil if dependencies not available
	listingHandlerPort := listinghandlers.NewListingHandlerAdapter(
		listingService,
		globalService,
	)

	listingHandler, ok := listingHandlerPort.(*listinghandlers.ListingHandler)
	if !ok {
		slog.Error("factory.http_handlers.listing_cast_failed")
		return HTTPHandlers{}
	}

	// Create media processing handler
	mediaProcessingHandlerPort := mediaprocessinghandlers.NewMediaProcessingHandler(
		mediaProcessingService,
		slog.Default(),
		callbackValidator,
	)

	mediaProcessingHandler, ok := mediaProcessingHandlerPort.(*mediaprocessinghandlers.MediaProcessingHandler)
	if !ok {
		slog.Error("factory.http_handlers.media_processing_cast_failed")
		return HTTPHandlers{}
	}

	// Create schedule handler using the adapter
	scheduleHandlerPort := schedulehandlers.NewScheduleHandlerAdapter(
		scheduleService,
	)

	scheduleHandler, ok := scheduleHandlerPort.(*schedulehandlers.ScheduleHandler)
	if !ok {
		slog.Error("factory.http_handlers.schedule_cast_failed")
		return HTTPHandlers{}
	}

	// Create holiday handler using the adapter
	holidayHandlerPort := holidayhandlers.NewHolidayHandlerAdapter(
		holidayService,
	)

	holidayHandler, ok := holidayHandlerPort.(*holidayhandlers.HolidayHandler)
	if !ok {
		slog.Error("factory.http_handlers.holiday_cast_failed")
		return HTTPHandlers{}
	}

	// Create admin handler using the adapter
	adminHandler := adminhandlers.NewAdminHandlerAdapter(
		userService,
		listingService,
		permissionService,
		propertyCoverageService,
		router,
	)

	photoSessionHandler := photosessionhandlers.NewPhotoSessionHandler(
		photoSessionService,
		globalService,
	)

	// Create metrics handler (optional)
	var metricsHandler *handlers.MetricsHandler
	if metricsAdapter != nil {
		metricsHandler = handlers.NewMetricsHandler(metricsAdapter.Prometheus)
	}

	slog.Info("Successfully created all HTTP handlers")

	return HTTPHandlers{
		UserHandler:            userHandler,
		ListingHandler:         listingHandler,
		MediaProcessingHandler: mediaProcessingHandler,
		AuthHandler:            authHandler,
		MetricsHandler:         metricsHandler,
		AdminHandler:           adminHandler,
		ScheduleHandler:        scheduleHandler,
		HolidayHandler:         holidayHandler,
		PhotoSessionHandler:    photoSessionHandler,
	}
}
