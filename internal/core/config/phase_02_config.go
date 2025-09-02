package config

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

// Phase02_LoadConfiguration carrega e valida toda a configuração do sistema
// Esta fase configura:
// - Carregamento de variáveis de ambiente
// - Carregamento de configuração YAML
// - Validação de configuração
func (b *Bootstrap) Phase02_LoadConfiguration() error {
	b.logger.Info("🎯 FASE 2: Carregamento e Validação de Configuração")
	b.logger.Debug("Carregando configuração do sistema")

	// 1. Carregar configuração YAML
	if err := b.loadEnvironmentConfig(); err != nil {
		return NewBootstrapError("Phase02", "load_config", "Failed to load configuration", err)
	}

	// 2. Validar configuração completa
	if err := b.validateConfiguration(); err != nil {
		return NewBootstrapError("Phase02", "validation", "Configuration validation failed", err)
	}

	// 3. Aplicar configurações de segurança (JWT secret e TTLs)
	b.applySecurityConfig()

	b.logger.Info("✅ Configuração carregada e validada com sucesso",
		"version", globalmodel.AppVersion)

	return nil
}

// loadEnvironmentConfig carrega as variáveis de ambiente
func (b *Bootstrap) loadEnvironmentConfig() error {
	b.logger.Debug("Carregando variáveis de ambiente")

	if err := b.config.LoadEnv(); err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}

	// Obter o environment carregado do config
	env, err := b.config.GetEnvironment()
	if err != nil {
		return fmt.Errorf("failed to get environment from config: %w", err)
	}

	// Armazenar referência para uso posterior
	b.env = env

	b.logger.Debug("✅ Variáveis de ambiente carregadas")
	return nil
}

// validateConfiguration valida toda a configuração carregada
func (b *Bootstrap) validateConfiguration() error {
	b.logger.Debug("Validando configuração completa")

	// Validar configurações críticas
	validations := []struct {
		name string
		fn   func() error
	}{
		{"database_config", b.validateDatabaseConfig},
		{"http_config", b.validateHTTPConfig},
		{"telemetry_config", b.validateTelemetryConfig},
		{"security_config", b.validateSecurityConfig},
	}

	var validationErrors []error

	for _, validation := range validations {
		if err := validation.fn(); err != nil {
			b.logger.Warn("Validação falhou",
				"validation", validation.name,
				"error", err)
			validationErrors = append(validationErrors, fmt.Errorf("%s: %w", validation.name, err))
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("configuration validation failed with %d errors: %v", len(validationErrors), validationErrors)
	}

	b.logger.Debug("✅ Validação de configuração concluída")
	return nil
}

// validateSecurityConfig valida itens de segurança (JWT e TTLs)
func (b *Bootstrap) validateSecurityConfig() error {
	// Garantir que o ambiente foi carregado
	if b.env == nil {
		return fmt.Errorf("environment not loaded")
	}

	// Validar JWT secret
	if strings.TrimSpace(b.env.JWT.Secret) == "" {
		return fmt.Errorf("jwt.secret is required and must not be empty")
	}

	// Validar TTLs
	if b.env.AUTH.AccessTTLMinutes <= 0 {
		return fmt.Errorf("auth.access_ttl_minutes must be > 0")
	}
	if b.env.AUTH.RefreshTTLDays <= 0 {
		return fmt.Errorf("auth.refresh_ttl_days must be > 0")
	}
	if b.env.AUTH.MaxSessionRotations <= 0 {
		return fmt.Errorf("auth.max_session_rotations must be > 0")
	}

	return nil
}

// applySecurityConfig aplica secret e TTLs no runtime global
func (b *Bootstrap) applySecurityConfig() {
	if b.env == nil {
		b.logger.Warn("Environment not loaded; skipping security config apply")
		return
	}

	// Aplicar JWT secret e TTLs no global model
	globalmodel.SetJWTSecret(b.env.JWT.Secret)
	globalmodel.SetAccessTTL(b.env.AUTH.AccessTTLMinutes)
	globalmodel.SetRefreshTTL(b.env.AUTH.RefreshTTLDays)
	globalmodel.SetMaxSessionRotations(b.env.AUTH.MaxSessionRotations)

	b.logger.Info("🔐 JWT and token TTL configuration applied")
}

// validateDatabaseConfig valida configuração do banco de dados
func (b *Bootstrap) validateDatabaseConfig() error {
	b.logger.Debug("Validando configuração do banco de dados")

	if b.env.DB.URI == "" {
		return fmt.Errorf("database URI is required")
	}

	// Validar formato da URI MySQL
	uri := b.env.DB.URI

	// Verificar se começa com usuário
	if !strings.Contains(uri, ":") {
		return fmt.Errorf("invalid database URI format: missing user credentials")
	}

	// Verificar se tem @tcp(
	if !strings.Contains(uri, "@tcp(") {
		return fmt.Errorf("invalid database URI format: missing @tcp(")
	}

	// Verificar se tem porta
	if !strings.Contains(uri, ":") || !strings.Contains(uri, ")/") {
		return fmt.Errorf("invalid database URI format: missing port or database name")
	}

	// Verificar se tem nome do banco
	// Para MySQL DSN, verificar se contém o nome do banco
	if !strings.Contains(uri, "/toq_db") {
		return fmt.Errorf("invalid database URI format: missing database name '/toq_db'")
	}

	// Verificar parâmetros opcionais
	if strings.Contains(uri, "?") {
		params := strings.Split(uri, "?")[1]
		if params != "" && !strings.Contains(params, "=") {
			return fmt.Errorf("invalid database URI format: malformed parameters")
		}
	}

	b.logger.Debug("✅ Configuração do banco de dados validada com sucesso")
	return nil
}

// validateHTTPConfig valida configuração HTTP
func (b *Bootstrap) validateHTTPConfig() error {
	b.logger.Debug("Validando configuração HTTP")

	// Validar porta HTTP (converter string para int)
	if b.env.HTTP.Port == "" {
		return fmt.Errorf("HTTP port is required")
	}

	portStr := b.env.HTTP.Port
	// Remover ":" incondicionalmente
	portStr = strings.TrimPrefix(portStr, ":")

	portInt := 0
	if _, err := fmt.Sscanf(portStr, "%d", &portInt); err != nil {
		return fmt.Errorf("invalid HTTP port format: %s", b.env.HTTP.Port)
	}

	if portInt <= 0 || portInt > 65535 {
		return fmt.Errorf("invalid HTTP port: %d (must be between 1 and 65535)", portInt)
	}

	// Validar network
	if b.env.HTTP.Network != "" && b.env.HTTP.Network != "tcp" && b.env.HTTP.Network != "tcp4" && b.env.HTTP.Network != "tcp6" {
		return fmt.Errorf("invalid network: %s (must be tcp, tcp4, or tcp6)", b.env.HTTP.Network)
	}

	// Validar timeouts (strings devem ser parseáveis como duration)
	if b.env.HTTP.ReadTimeout != "" {
		if _, err := time.ParseDuration(b.env.HTTP.ReadTimeout); err != nil {
			return fmt.Errorf("invalid read timeout format: %s", b.env.HTTP.ReadTimeout)
		}
	}

	if b.env.HTTP.WriteTimeout != "" {
		if _, err := time.ParseDuration(b.env.HTTP.WriteTimeout); err != nil {
			return fmt.Errorf("invalid write timeout format: %s", b.env.HTTP.WriteTimeout)
		}
	}

	// Validar MaxHeaderBytes
	if b.env.HTTP.MaxHeaderBytes < 0 {
		return fmt.Errorf("invalid max header bytes: %d (must be non-negative)", b.env.HTTP.MaxHeaderBytes)
	}

	// Validar TLS se habilitado
	if b.env.HTTP.TLS.Enabled {
		if b.env.HTTP.TLS.CertPath == "" {
			return fmt.Errorf("TLS enabled but cert path is empty")
		}
		if b.env.HTTP.TLS.KeyPath == "" {
			return fmt.Errorf("TLS enabled but key path is empty")
		}

		// Verificar se os arquivos existem
		if _, err := os.Stat(b.env.HTTP.TLS.CertPath); os.IsNotExist(err) {
			return fmt.Errorf("TLS cert file does not exist: %s", b.env.HTTP.TLS.CertPath)
		}
		if _, err := os.Stat(b.env.HTTP.TLS.KeyPath); os.IsNotExist(err) {
			return fmt.Errorf("TLS key file does not exist: %s", b.env.HTTP.TLS.KeyPath)
		}
	}

	b.logger.Debug("✅ Configuração HTTP validada com sucesso")
	return nil
}

// validateLoggingConfig valida configuração de logging
// validateTelemetryConfig valida configuração de telemetria
func (b *Bootstrap) validateTelemetryConfig() error {
	b.logger.Debug("Validando configuração de telemetria")

	// Se telemetria estiver desabilitada, não há validação adicional necessária
	if !b.env.TELEMETRY.Enabled {
		b.logger.Debug("Telemetria desabilitada, pulando validação")
		return nil
	}

	// Validar OTLP se habilitado
	if b.env.TELEMETRY.OTLP.Enabled {
		if b.env.TELEMETRY.OTLP.Endpoint == "" {
			return fmt.Errorf("OTLP enabled but endpoint is empty")
		}

		// Validar formato do endpoint
		// Accept either a full URL with scheme or a host:port
		if !strings.HasPrefix(b.env.TELEMETRY.OTLP.Endpoint, "http://") &&
			!strings.HasPrefix(b.env.TELEMETRY.OTLP.Endpoint, "https://") &&
			!strings.HasPrefix(b.env.TELEMETRY.OTLP.Endpoint, "grpc://") {
			// Try parsing as host:port
			if _, _, err := net.SplitHostPort(b.env.TELEMETRY.OTLP.Endpoint); err != nil {
				return fmt.Errorf("invalid OTLP endpoint format: %s (must be scheme://host:port or host:port)", b.env.TELEMETRY.OTLP.Endpoint)
			}
		}
	}

	// Validar porta de métricas se especificada
	if b.env.TELEMETRY.METRICS.Port != "" {
		portStr := b.env.TELEMETRY.METRICS.Port
		// Remover ":" incondicionalmente
		portStr = strings.TrimPrefix(portStr, ":")

		portInt := 0
		if _, err := fmt.Sscanf(portStr, "%d", &portInt); err != nil {
			return fmt.Errorf("invalid metrics port format: %s", b.env.TELEMETRY.METRICS.Port)
		}

		if portInt <= 0 || portInt > 65535 {
			return fmt.Errorf("invalid metrics port: %d (must be between 1 and 65535)", portInt)
		}
	}

	b.logger.Debug("✅ Configuração de telemetria validada com sucesso")
	return nil
}

// Phase02Rollback executa rollback da Fase 2
func (b *Bootstrap) Phase02Rollback() error {
	b.logger.Info("🔄 Executando rollback da Fase 2")

	// Não há muito para fazer rollback nesta fase
	// As configurações são apenas carregadas, não criam recursos

	b.logger.Info("✅ Rollback da Fase 2 concluído")
	return nil
}
