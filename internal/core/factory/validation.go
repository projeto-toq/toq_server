package factory

import (
	"fmt"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

// ValidateFactoryConfig valida se todas as configurações necessárias estão presentes
// Retorna erro se alguma configuração crítica estiver ausente
func ValidateFactoryConfig(config AdapterFactoryConfig) error {
	slog.Debug("Validating factory configuration")

	if config.Context == nil {
		return fmt.Errorf("context is required")
	}

	if config.Environment == nil {
		return fmt.Errorf("environment configuration is required")
	}

	// Validate environment configurations
	if err := validateEnvironment(config.Environment); err != nil {
		return fmt.Errorf("invalid environment configuration: %w", err)
	}

	slog.Debug("Factory configuration validation successful")
	return nil
}

// validateEnvironment valida as configurações de ambiente necessárias
func validateEnvironment(env *globalmodel.Environment) error {
	// Validate CEP configuration
	if env.CEP.Token == "" {
		return fmt.Errorf("CEP token is required")
	}

	// Validate Email configuration
	if env.EMAIL.SMTPServer == "" {
		return fmt.Errorf("SMTP server is required")
	}
	if env.EMAIL.SMTPUser == "" {
		return fmt.Errorf("SMTP user is required")
	}

	// Validate SMS configuration
	if env.SMS.AccountSid == "" {
		return fmt.Errorf("SMS Account SID is required")
	}
	if env.SMS.AuthToken == "" {
		return fmt.Errorf("SMS Auth Token is required")
	}

	// Validate FCM configuration
	if env.FCM.CredentialsFile == "" {
		return fmt.Errorf("FCM credentials file is required")
	}

	// Validate Redis configuration
	if env.REDIS.URL == "" {
		return fmt.Errorf("Redis URL is required")
	}

	return nil
}

// LogAdapterCreation registra informações sobre a criação de adapters
func LogAdapterCreation(adapterType string, count int) {
	slog.Info("Adapter creation completed",
		"type", adapterType,
		"count", count,
	)
}
