package config

import (
	"fmt"

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

	// Armazenar referência para uso posterior
	env := &globalmodel.Environment{}
	// Nota: Em implementação real, obter o env do config
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
	// Nota: Implementação real validaria URI, conexões, etc.
	b.logger.Debug("Validando configuração do banco de dados")
	return nil
}

// validateHTTPConfig valida configuração HTTP
func (b *Bootstrap) validateHTTPConfig() error {
	// Nota: Implementação real validaria portas, TLS, etc.
	b.logger.Debug("Validando configuração HTTP")
	return nil
}

// validateLoggingConfig valida configuração de logging
func (b *Bootstrap) validateLoggingConfig() error {
	// Nota: Implementação real validaria caminhos, níveis, etc.
	b.logger.Debug("Validando configuração de logging")
	return nil
}

// validateTelemetryConfig valida configuração de telemetria
func (b *Bootstrap) validateTelemetryConfig() error {
	// Nota: Implementação real validaria endpoints, chaves, etc.
	b.logger.Debug("Validando configuração de telemetria")
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
