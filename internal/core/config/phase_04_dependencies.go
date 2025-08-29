package config

// Phase04_InjectDependencies configura toda a injeção de dependências
// Esta fase configura:
// - Factory Pattern para criação de adapters
// - Validation adapters (CEP, CPF, CNPJ)
// - External service adapters (FCM, Email, SMS)
// - Storage adapters (MySQL, Redis)
// - Repository adapters
func (b *Bootstrap) Phase04_InjectDependencies() error {
	b.logger.Info("🎯 FASE 4: Injeção de Dependências via Factory Pattern")
	b.logger.Debug("Configurando injeção de dependências")

	// 1. Inicializar lifecycle manager para dependências
	if err := b.initializeLifecycleManager(); err != nil {
		return NewBootstrapError("Phase04", "lifecycle_manager", "Failed to initialize lifecycle manager", err)
	}

	// 2. Chamar o método InjectDependencies do config
	if err := b.config.InjectDependencies(b.lifecycleManager); err != nil {
		return NewBootstrapError("Phase04", "inject_dependencies", "Failed to inject dependencies", err)
	}

	b.logger.Info("✅ Injeção de dependências concluída via Factory Pattern")
	return nil
}

// initializeLifecycleManager inicializa o gerenciador de ciclo de vida
func (b *Bootstrap) initializeLifecycleManager() error {
	b.logger.Debug("Inicializando lifecycle manager")

	// O lifecycle manager já foi criado no bootstrap
	// Aqui podemos fazer configurações adicionais se necessário

	b.logger.Debug("✅ Lifecycle manager inicializado")
	return nil
}

// Phase04Rollback executa rollback da Fase 4
func (b *Bootstrap) Phase04Rollback() error {
	b.logger.Info("🔄 Executando rollback da Fase 4")

	// O cleanup será feito automaticamente pelo LifecycleManager
	// Todos os adapters criados têm funções de cleanup registradas

	b.logger.Info("✅ Rollback da Fase 4 concluído")
	return nil
}
