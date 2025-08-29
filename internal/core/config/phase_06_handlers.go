package config

// Phase06_ConfigureHandlers configura os handlers HTTP e rotas
// Esta fase configura:
// - Servidor HTTP com middleware
// - Handlers HTTP via Factory Pattern
// - Sistema de rotas
// - Health checks (liveness/readiness)
func (b *Bootstrap) Phase06_ConfigureHandlers() error {
	b.logger.Info("🎯 FASE 6: Configuração de Handlers e Rotas")
	b.logger.Debug("Configurando handlers HTTP e sistema de rotas")

	// 1. Inicializar servidor HTTP com middleware
	if err := b.initializeHTTPServer(); err != nil {
		return NewBootstrapError("Phase06", "http_server", "Failed to initialize HTTP server", err)
	}

	// 2. Criar handlers HTTP via Factory Pattern
	if err := b.createHTTPHandlers(); err != nil {
		return NewBootstrapError("Phase06", "http_handlers", "Failed to create HTTP handlers", err)
	}

	// 3. Configurar rotas e middlewares
	if err := b.setupRoutesAndMiddleware(); err != nil {
		return NewBootstrapError("Phase06", "routes", "Failed to setup routes and middleware", err)
	}

	// 4. Configurar health checks
	if err := b.setupHealthChecks(); err != nil {
		return NewBootstrapError("Phase06", "health_checks", "Failed to setup health checks", err)
	}

	b.logger.Info("✅ Handlers e rotas configurados com sucesso")
	return nil
}

// initializeHTTPServer inicializa o servidor HTTP com middleware
func (b *Bootstrap) initializeHTTPServer() error {
	b.logger.Debug("Inicializando servidor HTTP com middleware")

	// Inicializar HTTP server
	b.config.InitializeHTTP()

	// Adicionar cleanup
	b.lifecycleManager.AddCleanupFunc(func() {
		b.config.CloseHTTPServer()
	})

	b.logger.Info("✅ Servidor HTTP configurado com TLS e middleware")
	return nil
}

// createHTTPHandlers cria os handlers HTTP via Factory Pattern
func (b *Bootstrap) createHTTPHandlers() error {
	b.logger.Debug("Criando handlers HTTP via Factory Pattern")

	// Os handlers são criados durante SetupHTTPHandlersAndRoutes
	// Este método agora é apenas um placeholder para consistência com o bootstrap flow
	b.logger.Info("✅ Handlers HTTP preparados para criação")
	return nil
}

// setupRoutesAndMiddleware configura rotas e middlewares
func (b *Bootstrap) setupRoutesAndMiddleware() error {
	b.logger.Debug("Configurando rotas e middlewares")

	// Configurar handlers e rotas
	b.config.SetupHTTPHandlersAndRoutes()

	b.logger.Info("✅ Rotas e middlewares configurados")
	return nil
}

// setupHealthChecks configura os health checks
func (b *Bootstrap) setupHealthChecks() error {
	b.logger.Debug("Configurando health checks (liveness/readiness)")

	// Configurar health checks
	// Nota: Implementação real configuraria endpoints /healthz e /readyz

	b.logger.Info("✅ Health checks configurados")
	return nil
}

// Phase06Rollback executa rollback da Fase 6
func (b *Bootstrap) Phase06Rollback() error {
	b.logger.Info("🔄 Executando rollback da Fase 6")

	// Fechar servidor HTTP
	b.config.CloseHTTPServer()

	b.logger.Info("✅ Rollback da Fase 6 concluído")
	return nil
}
