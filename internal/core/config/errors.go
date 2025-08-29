package config

import (
	"fmt"
	"log/slog"
	"time"
)

// BootstrapError representa um erro estruturado durante o bootstrap
type BootstrapError struct {
	Phase     string                 `json:"phase"`
	Step      string                 `json:"step"`
	Message   string                 `json:"message"`
	Cause     error                  `json:"cause,omitempty"`
	Code      string                 `json:"code"`
	Severity  string                 `json:"severity"` // "low", "medium", "high", "critical"
	Context   map[string]interface{} `json:"context,omitempty"`
	Timestamp string                 `json:"timestamp"`
}

func (e *BootstrapError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s (%s)", e.Phase, e.Step, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("[%s] %s: %s", e.Phase, e.Step, e.Message)
}

func (e *BootstrapError) Unwrap() error {
	return e.Cause
}

// NewBootstrapError cria um novo erro estruturado
func NewBootstrapError(phase, step, message string, cause error) *BootstrapError {
	return &BootstrapError{
		Phase:     phase,
		Step:      step,
		Message:   message,
		Cause:     cause,
		Code:      generateErrorCode(phase, step),
		Severity:  determineSeverity(step),
		Context:   make(map[string]interface{}),
		Timestamp: fmt.Sprintf("%d", now().Unix()),
	}
}

// WithContext adiciona contexto ao erro
func (e *BootstrapError) WithContext(key string, value interface{}) *BootstrapError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithSeverity define a severidade do erro
func (e *BootstrapError) WithSeverity(severity string) *BootstrapError {
	e.Severity = severity
	return e
}

// ErrorHandler implementa tratamento estruturado de erros
type DefaultErrorHandler struct {
	logger     *slog.Logger
	maxRetries int
	retryDelay time.Duration
}

// NewDefaultErrorHandler cria um novo error handler
func NewDefaultErrorHandler(logger *slog.Logger) *DefaultErrorHandler {
	return &DefaultErrorHandler{
		logger:     logger,
		maxRetries: 3,
		retryDelay: time.Second * 2,
	}
}

// HandleError trata um erro ocorrido durante o bootstrap
func (h *DefaultErrorHandler) HandleError(phase string, err error) error {
	if bootstrapErr, ok := err.(*BootstrapError); ok {
		h.logger.Error("Bootstrap error occurred",
			"phase", bootstrapErr.Phase,
			"step", bootstrapErr.Step,
			"message", bootstrapErr.Message,
			"code", bootstrapErr.Code,
			"severity", bootstrapErr.Severity,
			"cause", bootstrapErr.Cause)

		// Log contexto adicional se existir
		for key, value := range bootstrapErr.Context {
			h.logger.Error("Error context",
				"key", key,
				"value", value)
		}
	} else {
		h.logger.Error("Unexpected error during bootstrap",
			"phase", phase,
			"error", err)
	}

	return err
}

// ShouldRetry determina se uma operação deve ser retentada
func (h *DefaultErrorHandler) ShouldRetry(phase string, attempt int, err error) bool {
	if attempt >= h.maxRetries {
		return false
	}

	// Não retentar erros críticos
	if bootstrapErr, ok := err.(*BootstrapError); ok {
		if bootstrapErr.Severity == "critical" {
			return false
		}
	}

	h.logger.Warn("Retrying failed operation",
		"phase", phase,
		"attempt", attempt,
		"max_attempts", h.maxRetries,
		"error", err)

	return true
}

// GetRetryDelay retorna o delay para retry
func (h *DefaultErrorHandler) GetRetryDelay(phase string, attempt int) time.Duration {
	// Exponential backoff
	delay := h.retryDelay * time.Duration(1<<uint(attempt-1))
	if delay > time.Minute {
		delay = time.Minute
	}
	return delay
}

// ValidationError representa erros de validação
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// NewValidationError cria um erro de validação
func NewValidationError(field, message string, value interface{}) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// AggregateError agrupa múltiplos erros
type AggregateError struct {
	Errors []error `json:"errors"`
}

func (e *AggregateError) Error() string {
	if len(e.Errors) == 0 {
		return "no errors"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	return fmt.Sprintf("%d errors occurred", len(e.Errors))
}

// Add adiciona um erro à lista
func (e *AggregateError) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

// HasErrors verifica se há erros
func (e *AggregateError) HasErrors() bool {
	return len(e.Errors) > 0
}

// ErrorCodes define códigos de erro padronizados
const (
	ErrCodeConfigLoad      = "CONFIG_LOAD_FAILED"
	ErrCodeDatabaseConn    = "DATABASE_CONNECTION_FAILED"
	ErrCodeServiceInit     = "SERVICE_INIT_FAILED"
	ErrCodeHandlerConfig   = "HANDLER_CONFIG_FAILED"
	ErrCodeResourceExhaust = "RESOURCE_EXHAUSTED"
	ErrCodeTimeout         = "OPERATION_TIMEOUT"
	ErrCodeValidation      = "VALIDATION_FAILED"
)

// generateErrorCode gera um código de erro baseado na fase e step
func generateErrorCode(phase, step string) string {
	return fmt.Sprintf("%s_%s", phase, step)
}

// determineSeverity determina a severidade baseada no step
func determineSeverity(step string) string {
	switch step {
	case "database_connection", "critical_service_init":
		return "critical"
	case "external_service_connection", "cache_connection":
		return "high"
	case "handler_configuration", "route_setup":
		return "medium"
	default:
		return "low"
	}
}

// now retorna o tempo atual (para facilitar testes)
var now = time.Now
