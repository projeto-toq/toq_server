package userservices

import (
	"context"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// securityEventLogger implementa SecurityEventLoggerInterface
type securityEventLogger struct{}

// NewSecurityEventLogger cria uma nova instância de SecurityEventLogger
func NewSecurityEventLogger() SecurityEventLoggerInterface {
	return &securityEventLogger{}
}

// LogSecurityEvent registra um evento de segurança genérico
func (sel *securityEventLogger) LogSecurityEvent(ctx context.Context, event *usermodel.SecurityEvent) error {
	// Prepara os campos de log estruturado
	logFields := []any{
		"eventType", event.EventType,
		"result", event.Result,
		"timestamp", event.Timestamp,
	}

	if event.UserID != nil {
		logFields = append(logFields, "userID", *event.UserID)
	}

	if event.NationalID != "" {
		logFields = append(logFields, "nationalID", event.NationalID)
	}

	if event.IPAddress != "" {
		logFields = append(logFields, "ipAddress", event.IPAddress)
	}

	if event.UserAgent != "" {
		logFields = append(logFields, "userAgent", event.UserAgent)
	}

	if event.ErrorType != nil {
		logFields = append(logFields, "errorType", *event.ErrorType)
	}

	if event.Reason != "" {
		logFields = append(logFields, "reason", event.Reason)
	}

	// Adiciona detalhes se existirem
	if len(event.Details) > 0 {
		logFields = append(logFields, "details", event.Details)
	}

	// Define o nível de log baseado no resultado
	switch event.Result {
	case usermodel.SecurityEventResultSuccess:
		slog.InfoContext(ctx, "Security event logged", logFields...)
	case usermodel.SecurityEventResultFailure:
		slog.WarnContext(ctx, "Security event logged - failure", logFields...)
	case usermodel.SecurityEventResultBlocked:
		slog.ErrorContext(ctx, "Security event logged - user blocked", logFields...)
	default:
		slog.InfoContext(ctx, "Security event logged", logFields...)
	}

	return nil
}

// LogSigninAttempt registra uma tentativa de login
func (sel *securityEventLogger) LogSigninAttempt(ctx context.Context, nationalID string, userID *int64, success bool, errorType *usermodel.SigninErrorType, ipAddress, userAgent string) error {
	var eventType usermodel.SecurityEventType
	var result usermodel.SecurityEventResult

	if success {
		eventType = usermodel.SecurityEventSigninSuccess
		result = usermodel.SecurityEventResultSuccess
	} else {
		eventType = usermodel.SecurityEventSigninFailure
		result = usermodel.SecurityEventResultFailure
	}

	event := usermodel.NewSecurityEvent(eventType, result).
		WithNationalID(nationalID).
		WithIPAddress(ipAddress).
		WithUserAgent(userAgent)

	if userID != nil {
		event = event.WithUserID(*userID)
	}

	if errorType != nil {
		event = event.WithErrorType(*errorType)
	}

	return sel.LogSecurityEvent(ctx, event)
}

// LogUserBlocked registra o bloqueio de um usuário
func (sel *securityEventLogger) LogUserBlocked(ctx context.Context, userID int64, reason string, ipAddress, userAgent string) error {
	event := usermodel.NewSecurityEvent(usermodel.SecurityEventUserBlocked, usermodel.SecurityEventResultBlocked).
		WithUserID(userID).
		WithReason(reason).
		WithIPAddress(ipAddress).
		WithUserAgent(userAgent)

	return sel.LogSecurityEvent(ctx, event)
}

// LogUserUnblocked registra o desbloqueio de um usuário
func (sel *securityEventLogger) LogUserUnblocked(ctx context.Context, userID int64, reason string) error {
	event := usermodel.NewSecurityEvent(usermodel.SecurityEventUserUnblocked, usermodel.SecurityEventResultSuccess).
		WithUserID(userID).
		WithReason(reason)

	return sel.LogSecurityEvent(ctx, event)
}

// LogInvalidCredentials registra tentativa com credenciais inválidas
func (sel *securityEventLogger) LogInvalidCredentials(ctx context.Context, nationalID string, ipAddress, userAgent string) error {
	errorType := usermodel.SigninErrorInvalidCredentials

	event := usermodel.NewSecurityEvent(usermodel.SecurityEventInvalidCredentials, usermodel.SecurityEventResultFailure).
		WithNationalID(nationalID).
		WithErrorType(errorType).
		WithIPAddress(ipAddress).
		WithUserAgent(userAgent)

	return sel.LogSecurityEvent(ctx, event)
}

// LogNoActiveRoles registra usuário sem roles ativos
func (sel *securityEventLogger) LogNoActiveRoles(ctx context.Context, userID int64, nationalID string, ipAddress, userAgent string) error {
	errorType := usermodel.SigninErrorNoActiveRoles

	event := usermodel.NewSecurityEvent(usermodel.SecurityEventNoActiveRoles, usermodel.SecurityEventResultFailure).
		WithUserID(userID).
		WithNationalID(nationalID).
		WithErrorType(errorType).
		WithIPAddress(ipAddress).
		WithUserAgent(userAgent)

	return sel.LogSecurityEvent(ctx, event)
}
