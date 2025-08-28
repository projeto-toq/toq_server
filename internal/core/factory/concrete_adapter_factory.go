package factory

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/cache"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"

	// HTTP handlers
	listinghandlers "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/listing_handlers"
	userhandlers "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/user_handlers"

	// Validation adapters
	cepadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/cep"
	cnpjadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/cnpj"
	cpfadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/cpf"

	// External service adapters
	emailadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/email"
	fcmadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/fcm"
	smsadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/sms"

	// Storage adapters - AWS S3 (substituindo GCS)
	s3adapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/aws_s3"

	// Storage adapters
	mysqladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql"

	// Repository adapters
	mysqlcomplexadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/complex"
	mysqlglobaladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/global"
	mysqllistingadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing"
	mysqlpermissionadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/permission"
	sessionmysqladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/session"
	mysqluseradapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user"

	// Core services
	complexservice "github.com/giulio-alfieri/toq_server/internal/core/service/complex_service"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	listingservice "github.com/giulio-alfieri/toq_server/internal/core/service/listing_service"
	userservice "github.com/giulio-alfieri/toq_server/internal/core/service/user_service"
)

// ConcreteAdapterFactory implementa a interface AdapterFactory
// Responsável pela criação concreta de todos os adapters do sistema
type ConcreteAdapterFactory struct {
	lm LifecycleManager
}

// CreateValidationAdapters cria e configura todos os adapters de validação externa
// Retorna ValidationAdapters com CEP, CPF, CNPJ e CRECI configurados
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

	slog.Info("Successfully created all external service adapters")

	return ExternalServiceAdapters{
		FCM:          fcm,
		Email:        email,
		SMS:          sms,
		CloudStorage: s3,      // S3 adapter via interface CloudStorage
		CloseFunc:    s3Close, // Função de cleanup do S3
	}, nil
}

// CreateStorageAdapters cria adapters de armazenamento (Database e Cache)
// Inclui MySQL database e Redis cache com função de cleanup
func (f *ConcreteAdapterFactory) CreateStorageAdapters(ctx context.Context, env *globalmodel.Environment, db *sql.DB) (StorageAdapters, error) {
	slog.Info("Creating storage adapters")

	// Database Adapter
	database := mysqladapter.NewDB(db)

	// Redis Cache - criar sem GlobalService inicialmente (será injetado posteriormente)
	redisCache, err := cache.NewRedisCache(env.REDIS.URL, nil)
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
func (f *ConcreteAdapterFactory) CreateRepositoryAdapters(database *mysqladapter.Database) (RepositoryAdapters, error) {
	slog.Info("Creating repository adapters")

	// User Repository
	userRepo := mysqluseradapter.NewUserAdapter(database)

	// Device Token Repository (access through User Repository)
	deviceTokenRepo := userRepo.GetDeviceTokenRepository()

	// Global Repository
	globalRepo := mysqlglobaladapter.NewGlobalAdapter(database)

	// Complex Repository
	complexRepo := mysqlcomplexadapter.NewComplexAdapter(database)

	// Listing Repository
	listingRepo := mysqllistingadapter.NewListingAdapter(database)

	// Session Repository
	sessionRepo := sessionmysqladapter.NewSessionAdapter(database)

	// Permission Repository
	permissionRepo := mysqlpermissionadapter.NewPermissionAdapter(database)

	slog.Info("Successfully created all repository adapters")

	return RepositoryAdapters{
		User:        userRepo,
		Global:      globalRepo,
		Complex:     complexRepo,
		Listing:     listingRepo,
		Session:     sessionRepo,
		Permission:  permissionRepo,
		DeviceToken: deviceTokenRepo,
	}, nil
}

// CreateHTTPHandlers creates and returns all HTTP handlers
func (factory *ConcreteAdapterFactory) CreateHTTPHandlers(
	userService userservice.UserServiceInterface,
	globalService globalservice.GlobalServiceInterface,
	listingService listingservice.ListingServiceInterface,
	complexService complexservice.ComplexServiceInterface,
) HTTPHandlers {
	slog.Info("Creating HTTP handlers")

	// Create user handler using the adapter
	userHandler := userhandlers.NewUserHandlerAdapter(
		userService,
		globalService,
		complexService,
	)

	// Create listing handler using the adapter
	listingHandler := listinghandlers.NewListingHandlerAdapter(
		listingService,
		globalService,
		complexService,
	)

	slog.Info("Successfully created all HTTP handlers")

	return HTTPHandlers{
		UserHandler:    userHandler,
		ListingHandler: listingHandler,
	}
}
