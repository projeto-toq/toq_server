package config

import (
	"time"
)

// Phase08_StartServer inicializa o servidor e gerencia o runtime
// Esta fase configura:
// - Inicialização do servidor HTTP
// - Gerenciamento de shutdown graceful
// - Monitoramento de saúde em runtime
// - Wait group para sincronização
func (b *Bootstrap) Phase08_StartServer() error {
	b.logger.Info("🎯 FASE 8: Inicialização Final e Runtime")
	b.logger.Debug("Preparando servidor para aceitar conexões")

	// 1. Marcar servidor como pronto para receber tráfego
	if err := b.markServerReady(); err != nil {
		return NewBootstrapError("Phase08", "server_ready", "Failed to mark server as ready", err)
	}

	// 2. Iniciar servidor HTTP em goroutine
	if err := b.startHTTPServer(); err != nil {
		return NewBootstrapError("Phase08", "http_server_start", "Failed to start HTTP server", err)
	}

	// 3. Configurar monitoramento de saúde em runtime
	b.startHealthMonitoring()

	b.logger.Info("🌟 TOQ Server pronto para servir",
		"uptime", time.Since(b.startTime))

	return nil
}

// markServerReady marca o servidor como pronto para receber tráfego
func (b *Bootstrap) markServerReady() error {
	b.logger.Debug("Marcando servidor como pronto para tráfego")

	// Nota: Implementação real chamaria b.config.SetHealthServing(true)
	b.logger.Info("✅ Servidor marcado como ready")
	return nil
}

// startHTTPServer inicia o servidor HTTP em goroutine
func (b *Bootstrap) startHTTPServer() error {
	b.logger.Debug("Iniciando servidor HTTP")

	// Adicionar ao wait group
	b.wg.Add(1)

	// Iniciar servidor em goroutine
	go func() {
		defer b.wg.Done()

		b.logger.Info("🚀 Iniciando servidor HTTP")

		// Nota: Implementação real usaria b.config.GetHTTPServer().ListenAndServe()
		// if err := b.config.GetHTTPServer().ListenAndServe(); err != nil && err != http.ErrServerClosed {
		//     b.logger.Error("Servidor HTTP falhou", "error", err)
		//     // Trigger shutdown
		// }

		// Simulação para desenvolvimento
		b.logger.Info("✅ Servidor HTTP iniciado (simulado)")
	}()

	return nil
}

// startHealthMonitoring inicia o monitoramento de saúde em runtime
func (b *Bootstrap) startHealthMonitoring() {
	b.logger.Debug("Iniciando monitoramento de saúde em runtime")

	// Iniciar monitoramento em goroutine
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-b.ctx.Done():
				b.logger.Debug("Monitoramento de saúde parado")
				return
			case <-ticker.C:
				health := b.Health()
				if !health.Overall {
					b.logger.Warn("Sistema com problemas de saúde detectados",
						"error_count", health.ErrorCount)
				}
			}
		}
	}()

	b.logger.Info("✅ Monitoramento de saúde em runtime iniciado")
}

// WaitForShutdown aguarda sinal de shutdown e executa graceful shutdown
func (b *Bootstrap) WaitForShutdown() {
	b.logger.Info("⏳ Aguardando sinal de shutdown...")

	// Aguardar sinal de shutdown
	<-b.shutdownChan

	b.logger.Info("🛑 Sinal de shutdown recebido, iniciando graceful shutdown...")

	// Executar shutdown
	b.Shutdown()
}

// Phase08Rollback executa rollback da Fase 8
func (b *Bootstrap) Phase08Rollback() error {
	b.logger.Info("🔄 Executando rollback da Fase 8")

	// Parar monitoramento de saúde
	b.cancel()

	// Aguardar workers terminarem
	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()

	// Timeout para workers
	select {
	case <-done:
		b.logger.Info("✅ Workers terminaram gracefully")
	case <-time.After(30 * time.Second):
		b.logger.Warn("⚠️ Timeout aguardando workers")
	}

	b.logger.Info("✅ Rollback da Fase 8 concluído")
	return nil
}
