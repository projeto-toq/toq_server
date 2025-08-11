package factory

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/cache"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"

	// Validation adapters
	cepadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/cep"
	cnpjadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/cnpj"
	cpfadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/cpf"
	creciadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/creci"

	// External service adapters
	emailadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/email"
	fcmadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/fcm"
	smsadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/sms"

	// Storage adapters
	mysqladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql"

	// Repository adapters
	mysqlcomplexadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/complex"
	mysqlglobaladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/global"
	mysqllistingadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing"
	sessionmysqladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/session"
	mysqluseradapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user"
)

// ConcreteAdapterFactory implementa a interface AdapterFactory
// Responsável pela criação concreta de todos os adapters do sistema
type ConcreteAdapterFactory struct{}

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

	// CRECI Adapter
	creci := creciadapter.NewCreciAdapter(context.Background())

	slog.Info("Successfully created all validation adapters")

	return ValidationAdapters{
		CEP:   cep,
		CPF:   cpf,
		CNPJ:  cnpj,
		CRECI: creci,
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

	// Email Adapter
	email := emailadapter.NewEmailAdapter(
		env.EMAIL.SMTPServer,
		env.EMAIL.SMTPPort,
		env.EMAIL.SMTPUser,
		env.EMAIL.SMTPPassword,
	)

	// SMS Adapter
	sms := smsadapter.NewSmsAdapter(
		env.SMS.AccountSid,
		env.SMS.AuthToken,
		env.SMS.MyNumber,
	)

	slog.Info("Successfully created all external service adapters")

	return ExternalServiceAdapters{
		FCM:   fcm,
		Email: email,
		SMS:   sms,
		GCS:   nil, // TODO: Implementar quando necessário
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
// Agrupa repositórios por domínio (User, Global, Complex, Listing, Session)
func (f *ConcreteAdapterFactory) CreateRepositoryAdapters(database *mysqladapter.Database) (RepositoryAdapters, error) {
	slog.Info("Creating repository adapters")

	// User Repository
	userRepo := mysqluseradapter.NewUserAdapter(database)

	// Global Repository
	globalRepo := mysqlglobaladapter.NewGlobalAdapter(database)

	// Complex Repository
	complexRepo := mysqlcomplexadapter.NewComplexAdapter(database)

	// Listing Repository
	listingRepo := mysqllistingadapter.NewListingAdapter(database)

	// Session Repository
	sessionRepo := sessionmysqladapter.NewMySQLSessionAdapter(database.DB)

	slog.Info("Successfully created all repository adapters")

	return RepositoryAdapters{
		User:    userRepo,
		Global:  globalRepo,
		Complex: complexRepo,
		Listing: listingRepo,
		Session: sessionRepo,
	}, nil
}
