package config

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/core/cache"
	"github.com/giulio-alfieri/toq_server/internal/core/factory"
	goroutines "github.com/giulio-alfieri/toq_server/internal/core/go_routines"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
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
	Database         interface{} // *sql.DB
	Cache            cache.CacheInterface
	ActivityTracker  *goroutines.ActivityTracker
	GinRouter        *gin.Engine
	HTTPServer       interface{} // *http.Server
	Factory          factory.AdapterFactory
	Repositories     *factory.RepositoryAdapters
	Services         *ServiceContainer
	Handlers         *factory.HTTPHandlers
	LifecycleManager *LifecycleManager
}

// ServiceContainer agrupa todos os serviços do sistema
type ServiceContainer struct {
	GlobalService     interface{} // GlobalServiceInterface
	UserService       interface{} // UserServiceInterface
	ComplexService    interface{} // ComplexServiceInterface
	ListingService    interface{} // ListingServiceInterface
	PermissionService interface{} // PermissionServiceInterface
}

// PhaseResult representa o resultado da execução de uma fase
type PhaseResult struct {
	Name      string
	Status    string // "success", "failed", "skipped"
	Duration  time.Duration
	Error     error
	StartedAt time.Time
	EndedAt   time.Time
	Metadata  map[string]interface{}
}

// PhaseConfig contém configurações específicas para uma fase
type PhaseConfig struct {
	Name       string
	Timeout    time.Duration
	MaxRetries int
	RetryDelay time.Duration
	Required   bool     // Se true, falha da fase para todo o bootstrap
	DependsOn  []string // Fases das quais esta depende
}

// BootstrapMetrics contém métricas de performance do bootstrap
type BootstrapMetrics struct {
	TotalDuration    time.Duration
	PhaseCount       int
	SuccessfulPhases int
	FailedPhases     int
	SkippedPhases    int
	StartTime        time.Time
	EndTime          time.Time
}

// ValidationResult representa o resultado da validação de uma fase
type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

// BootstrapEvent representa um evento durante o bootstrap
type BootstrapEvent struct {
	Type      string // "phase_start", "phase_end", "error", "warning"
	Phase     string
	Message   string
	Timestamp time.Time
	Data      map[string]interface{}
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
