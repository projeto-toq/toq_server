package config

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	// Marcar como pronto para receber tráfego
	b.config.SetHealthServing(true)

	b.logger.Info("✅ Servidor marcado como ready para receber tráfego")
	return nil
}

// startHTTPServer inicia o servidor HTTP em goroutine
func (b *Bootstrap) startHTTPServer() error {
	b.logger.Debug("Iniciando servidor HTTP")

	// Verificar se o servidor HTTP foi configurado
	httpServer := b.config.GetHTTPServer()
	if httpServer == nil {
		return NewBootstrapError("Phase08", "http_server", "HTTP server not configured", nil)
	}

	// Adicionar ao wait group
	b.wg.Add(1)

	// Canal para sinalizar quando o servidor parar
	serverDone := make(chan error, 1)

	// Iniciar servidor em goroutine
	go func() {
		defer b.wg.Done()

		b.logger.Info("🚀 Iniciando servidor HTTP na porta configurada")

		// Iniciar servidor HTTP
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			b.logger.Error("Servidor HTTP falhou", "error", err)
			serverDone <- err
		} else {
			serverDone <- nil
		}
	}()

	// Aguardar um momento para verificar se o servidor iniciou corretamente
	time.Sleep(100 * time.Millisecond)

	// Verificar se houve erro na inicialização
	select {
	case err := <-serverDone:
		if err != nil {
			return NewBootstrapError("Phase08", "http_server_start", "Failed to start HTTP server", err)
		}
	default:
		// Servidor iniciou corretamente
	}

	b.logger.Info("✅ Servidor HTTP iniciado com sucesso")
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

	// Configurar canal para sinais do sistema
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Aguardar sinal de shutdown ou sinal do sistema
	select {
	case <-b.shutdownChan:
		b.logger.Info("🛑 Sinal de shutdown interno recebido")
	case sig := <-sigChan:
		b.logger.Info("🛑 Sinal do sistema recebido", "signal", sig)
	}

	b.logger.Info("🛑 Iniciando graceful shutdown...")

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
