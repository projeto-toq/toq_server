package userservices

import (
	"context"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// SecurityEventLoggerInterface define os métodos para logging de eventos de segurança
type SecurityEventLoggerInterface interface {
	// LogSecurityEvent registra um evento de segurança genérico
	LogSecurityEvent(ctx context.Context, event *usermodel.SecurityEvent) error

	// LogSigninAttempt registra uma tentativa de login
	LogSigninAttempt(ctx context.Context, nationalID string, userID *int64, success bool, errorType *usermodel.SigninErrorType, ipAddress, userAgent string) error

	// LogUserBlocked registra o bloqueio de um usuário
	LogUserBlocked(ctx context.Context, userID int64, reason string, ipAddress, userAgent string) error

	// LogUserUnblocked registra o desbloqueio de um usuário
	LogUserUnblocked(ctx context.Context, userID int64, reason string) error

	// LogInvalidCredentials registra tentativa com credenciais inválidas
	LogInvalidCredentials(ctx context.Context, nationalID string, ipAddress, userAgent string) error

	// LogNoActiveRoles registra usuário sem roles ativos
	LogNoActiveRoles(ctx context.Context, userID int64, nationalID string, ipAddress, userAgent string) error
}
