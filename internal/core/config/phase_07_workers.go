package config

import "fmt"

// Phase07_StartBackgroundWorkers inicializa os workers em background
// Esta fase configura:
// - Session cleaner (limpeza de sessões expiradas)
// - CRECI validator worker (validação de CRECI)
// - Activity tracker workers (rastreamento de atividades)
// - Database maintenance tasks (manutenção do banco)
func (b *Bootstrap) Phase07_StartBackgroundWorkers() error {
	b.logger.Info("🎯 FASE 7: Inicialização de Background Workers")
	b.logger.Debug("Configurando workers em background")

	// 1. Inicializar goroutines do sistema
	if err := b.initializeSystemGoroutines(); err != nil {
		return NewBootstrapError("Phase07", "system_goroutines", "Failed to initialize system goroutines", err)
	}

	// 2. Configurar activity tracker com user service
	if err := b.linkActivityTrackerToUserService(); err != nil {
		return NewBootstrapError("Phase07", "activity_tracker_link", "Failed to link activity tracker to user service", err)
	}

	// 3. Verificar schema do banco de dados
	if err := b.verifyDatabaseSchema(); err != nil {
		return NewBootstrapError("Phase07", "database_schema", "Failed to verify database schema", err)
	}

	b.logger.Info("✅ Background workers inicializados com sucesso")
	return nil
}

// initializeSystemGoroutines inicializa as goroutines do sistema
func (b *Bootstrap) initializeSystemGoroutines() error {
	b.logger.Debug("Inicializando goroutines do sistema")

	// Inicializar background workers
	b.config.InitializeGoRoutines()

	b.logger.Info("✅ Background workers inicializados")
	return nil
}

// linkActivityTrackerToUserService conecta o activity tracker ao user service
func (b *Bootstrap) linkActivityTrackerToUserService() error {
	b.logger.Debug("Verificando e conectando activity tracker ao user service")

	// Verificar se ActivityTracker foi criado na Phase 04
	if b.config.GetActivityTracker() == nil {
		return fmt.Errorf("ActivityTracker não foi criado na Phase 04")
	}

	// Configurar activity tracker com user service
	b.config.SetActivityTrackerUserService()

	b.logger.Info("✅ Activity tracker conectado ao user service")
	return nil
}

// verifyDatabaseSchema verifica e inicializa o schema do banco
func (b *Bootstrap) verifyDatabaseSchema() error {
	b.logger.Debug("Verificando schema do banco de dados")

	// Verificar e inicializar schema se necessário
	b.config.VerifyDatabase()

	b.logger.Info("✅ Schema do banco de dados verificado")
	return nil
}

// Phase07Rollback executa rollback da Fase 7
func (b *Bootstrap) Phase07Rollback() error {
	b.logger.Info("🔄 Executando rollback da Fase 7")

	// Cancelar contexto para parar os workers
	if b.cancel != nil {
		b.cancel()
	}

	// Aguardar workers terminarem (será feito no Phase08)
	b.logger.Info("✅ Rollback da Fase 7 concluído")
	return nil
}
