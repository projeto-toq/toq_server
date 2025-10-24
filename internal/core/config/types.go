package config

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/core/cache"
	"github.com/projeto-toq/toq_server/internal/core/factory"
	goroutines "github.com/projeto-toq/toq_server/internal/core/go_routines"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	complexservices "github.com/projeto-toq/toq_server/internal/core/service/complex_service"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	holidayservices "github.com/projeto-toq/toq_server/internal/core/service/holiday_service"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	permissionservices "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
)

// Bootstrapper define a interface para o sistema de bootstrap
type Bootstrapper interface {
	Bootstrap() error
	Shutdown()
	Health() HealthStatus
	WaitShutdown()
}

// PhaseExecutor define a interface para execução de fases
type PhaseExecutor interface {
	Execute(ctx context.Context) error
	Rollback(ctx context.Context) error
	Name() string
}

// BootstrapContext contém o contexto compartilhado entre todas as fases
type BootstrapContext struct {
	Context          context.Context
	Config           *globalmodel.Environment
	Database         *sql.DB
	Cache            cache.CacheInterface
	ActivityTracker  *goroutines.ActivityTracker
	GinRouter        *gin.Engine
	HTTPServer       *http.Server
	Factory          factory.AdapterFactory
	Repositories     *factory.RepositoryAdapters
	Services         *ServiceContainer
	Handlers         *factory.HTTPHandlers
	LifecycleManager *LifecycleManager
}

// ServiceContainer agrupa todos os serviços do sistema
type ServiceContainer struct {
	GlobalService       globalservice.GlobalServiceInterface
	UserService         userservices.UserServiceInterface
	ComplexService      complexservices.ComplexServiceInterface
	ListingService      listingservices.ListingServiceInterface
	PermissionService   permissionservices.PermissionServiceInterface
	HolidayService      holidayservices.HolidayServiceInterface
	ScheduleService     scheduleservices.ScheduleServiceInterface
	PhotoSessionService photosessionservices.PhotoSessionServiceInterface
}

// ValidationResult representa o resultado da validação de uma fase
type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

// ErrorHandler define como erros são tratados durante o bootstrap
type ErrorHandler interface {
	HandleError(phase string, err error) error
	ShouldRetry(phase string, attempt int, err error) bool
	GetRetryDelay(phase string, attempt int) time.Duration
}

// Logger define a interface para logging estruturado
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	With(args ...interface{}) Logger
}

// ConfigLoader define a interface para carregamento de configuração
type ConfigLoader interface {
	Load() (*globalmodel.Environment, error)
	Validate(*globalmodel.Environment) error
}

// HealthChecker define a interface para verificação de saúde
type HealthChecker interface {
	CheckHealth() HealthStatus
	IsReady() bool
}

// ResourceManager define a interface para gerenciamento de recursos
type ResourceManager interface {
	Acquire(resource string) error
	Release(resource string)
	Cleanup()
}
