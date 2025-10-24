package config

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/projeto-toq/toq_server/internal/core/cache"
	"github.com/projeto-toq/toq_server/internal/core/factory"
)

// Phase03_InitializeInfrastructure inicializa a infraestrutura core do sistema
// Esta fase configura:
// - Conexão com banco de dados
// - Conexão com cache Redis
// - Sistema de telemetria (OpenTelemetry)
// - Activity tracker para sessões
func (b *Bootstrap) Phase03_InitializeInfrastructure() error {
	b.logger.Info("🎯 FASE 3: Inicialização da Infraestrutura Core")
	b.logger.Debug("Configurando infraestrutura fundamental")

	// 1. Inicializar conexão com banco de dados
	if err := b.initializeDatabase(); err != nil {
		return NewBootstrapError("Phase03", "database", "Failed to initialize database connection", err)
	}

	// 2. Inicializar sistema de cache Redis
	if err := b.initializeCache(); err != nil {
		return NewBootstrapError("Phase03", "cache", "Failed to initialize Redis cache", err)
	}

	// 3. Inicializar OpenTelemetry (tracing + metrics)
	if err := b.initializeTelemetry(); err != nil {
		return NewBootstrapError("Phase03", "telemetry", "Failed to initialize OpenTelemetry", err)
	}

	// 4. Inicializar adapter de métricas
	if err := b.initializeMetrics(); err != nil {
		return NewBootstrapError("Phase03", "metrics", "Failed to initialize metrics adapter", err)
	}

	b.logger.Info("✅ Infraestrutura core inicializada com sucesso")
	return nil
}

// initializeDatabase inicializa a conexão com o banco de dados
func (b *Bootstrap) initializeDatabase() error {
	b.logger.Debug("Inicializando conexão com banco de dados")

	// Inicializar database
	b.config.InitializeDatabase()

	// Adicionar cleanup para fechar conexão
	b.lifecycleManager.AddCleanupFunc(func() {
		if db := b.config.GetDatabase(); db != nil {
			if err := db.Close(); err != nil {
				slog.Error("Erro fechando conexão MySQL", "error", err)
			} else {
				slog.Info("Conexão MySQL fechada com sucesso")
			}
		}
	})

	b.logger.Info("✅ Conexão com banco de dados estabelecida")
	return nil
}

// initializeCache inicializa o sistema de cache Redis
func (b *Bootstrap) initializeCache() error {
	b.logger.Debug("Inicializando sistema de cache Redis")

	// Obter configuração do Redis do environment
	redisURL := b.env.REDIS.URL
	if redisURL == "" {
		return NewBootstrapError("Phase03", "redis_config", "Redis URL not configured", nil)
	}

	b.logger.Debug("Conectando ao Redis", "url", redisURL)

	// Testar conexão com Redis
	_, err := cache.NewRedisCache(redisURL, nil)
	if err != nil {
		return NewBootstrapError("Phase03", "redis_connection", "Failed to initialize Redis cache", err)
	}

	// O cache será configurado quando os adapters forem criados na Phase 04
	b.logger.Info("✅ Sistema de cache Redis inicializado com sucesso")
	return nil
}

// initializeTelemetry inicializa o sistema de observabilidade
func (b *Bootstrap) initializeTelemetry() error {
	b.logger.Debug("Inicializando sistema de telemetria OpenTelemetry")

	// Inicializar OpenTelemetry com cleanup
	shutdownOtel, err := b.config.InitializeTelemetry()
	if err != nil {
		return fmt.Errorf("failed to initialize OpenTelemetry: %w", err)
	}

	// Adicionar cleanup
	b.lifecycleManager.AddCleanupFunc(func() {
		slog.Info("Desligando OpenTelemetry...")
		shutdownOtel()
	})

	b.logger.Info("✅ OpenTelemetry inicialização concluída")
	return nil
}

// initializeMetrics inicializa o adapter de métricas Prometheus
func (b *Bootstrap) initializeMetrics() error {
	b.logger.Debug("Inicializando adapter de métricas Prometheus")

	if !b.env.TELEMETRY.METRICS.Enabled {
		b.logger.Info("Métricas Prometheus desabilitadas; adapter não será iniciado")
		return nil
	}

	// Criar factory de adapters
	adapterFactory := factory.NewAdapterFactory(b.lifecycleManager)

	// Criar adapter de métricas
	metricsAdapter := adapterFactory.CreateMetricsAdapter(b.config.GetRuntimeEnvironment())

	// Inicializar métricas
	ctx := context.Background()
	if err := metricsAdapter.Prometheus.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
	}

	// Armazenar no config
	b.config.(*config).metricsAdapter = metricsAdapter

	// Adicionar cleanup
	b.lifecycleManager.AddCleanupFunc(func() {
		ctx := context.Background()
		if err := metricsAdapter.Prometheus.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown metrics", "error", err)
		}
	})

	b.logger.Info("✅ Adapter de métricas Prometheus inicializado")
	return nil
}

// Phase03Rollback executa rollback da Fase 3
func (b *Bootstrap) Phase03Rollback() error {
	b.logger.Info("🔄 Executando rollback da Fase 3")

	// O cleanup será feito automaticamente pelo LifecycleManager
	// Aqui podemos fazer limpeza específica se necessário

	b.logger.Info("✅ Rollback da Fase 3 concluído")
	return nil
}
