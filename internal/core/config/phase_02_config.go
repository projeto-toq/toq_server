package config

import (
	"fmt"
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
// - Inicialização do sistema de logging
func (b *Bootstrap) Phase02_LoadConfiguration() error {
	b.logger.Info("🎯 FASE 2: Carregamento e Validação de Configuração")
	b.logger.Debug("Carregando configuração do sistema")

	// 1. Carregar variáveis de ambiente
	if err := b.loadEnvironmentConfig(); err != nil {
		return NewBootstrapError("Phase02", "env_load", "Failed to load environment configuration", err)
	}

	// 2. Inicializar sistema de logging (primeira vez com env)
	if err := b.initializeEarlyLogging(); err != nil {
		return NewBootstrapError("Phase02", "early_logging", "Failed to initialize early logging", err)
	}

	// 3. Reconfigurar logging com YAML (ENV ainda tem prioridade)
	if err := b.reconfigureLoggingWithYAML(); err != nil {
		return NewBootstrapError("Phase02", "yaml_logging", "Failed to reconfigure logging with YAML", err)
	}

	// 4. Validar configuração completa
	if err := b.validateConfiguration(); err != nil {
		return NewBootstrapError("Phase02", "validation", "Configuration validation failed", err)
	}

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

// initializeEarlyLogging inicializa o logging baseado apenas em variáveis de ambiente
func (b *Bootstrap) initializeEarlyLogging() error {
	b.logger.Debug("Inicializando logging baseado em variáveis de ambiente")

	// Inicializar logging com configuração de ambiente
	b.config.InitializeLog()

	b.logger.Info("✅ Logging inicial baseado em ENV configurado")
	return nil
}

// reconfigureLoggingWithYAML reconfigura o logging com valores do YAML
func (b *Bootstrap) reconfigureLoggingWithYAML() error {
	b.logger.Debug("Reconfigurando logging com valores YAML")

	// Re-inicializar logging com configuração completa (ENV > YAML > defaults)
	b.config.InitializeLog()

	b.logger.Info("✅ Logging reconfigurado com prioridade ENV > YAML > defaults")
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
		{"logging_config", b.validateLoggingConfig},
		{"telemetry_config", b.validateTelemetryConfig},
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
func (b *Bootstrap) validateLoggingConfig() error {
	b.logger.Debug("Validando configuração de logging")

	// Validar nível de log
	validLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	levelValid := false
	for _, level := range validLevels {
		if strings.ToUpper(b.env.LOG.Level) == level {
			levelValid = true
			break
		}
	}
	if !levelValid {
		return fmt.Errorf("invalid log level: %s (must be one of: %v)", b.env.LOG.Level, validLevels)
	}

	// Validar configuração de arquivo se habilitada
	if b.env.LOG.ToFile {
		if b.env.LOG.Path == "" {
			return fmt.Errorf("log to file enabled but path is empty")
		}

		if b.env.LOG.Filename == "" {
			return fmt.Errorf("log to file enabled but filename is empty")
		}

		// Verificar se o diretório existe ou pode ser criado
		if _, err := os.Stat(b.env.LOG.Path); os.IsNotExist(err) {
			// Tentar criar o diretório
			if err := os.MkdirAll(b.env.LOG.Path, 0755); err != nil {
				return fmt.Errorf("cannot create log directory: %s", b.env.LOG.Path)
			}
		}

		// Verificar se podemos escrever no diretório
		testFile := fmt.Sprintf("%s/.log_test", b.env.LOG.Path)
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			return fmt.Errorf("cannot write to log directory: %s", b.env.LOG.Path)
		}
		os.Remove(testFile) // Limpar arquivo de teste
	}

	b.logger.Debug("✅ Configuração de logging validada com sucesso")
	return nil
}

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
		if !strings.HasPrefix(b.env.TELEMETRY.OTLP.Endpoint, "http://") &&
			!strings.HasPrefix(b.env.TELEMETRY.OTLP.Endpoint, "https://") &&
			!strings.HasPrefix(b.env.TELEMETRY.OTLP.Endpoint, "grpc://") {
			return fmt.Errorf("invalid OTLP endpoint format: %s (must start with http://, https://, or grpc://)", b.env.TELEMETRY.OTLP.Endpoint)
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
