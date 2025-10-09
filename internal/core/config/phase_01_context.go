package config

import (
	"context"
	"net/http"
	_ "net/http/pprof" // Import pprof para habilitar endpoints de debugging
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// Phase01_InitializeContext inicializa o contexto base e sinais do sistema
// Esta fase configura:
// - Contexto principal com cancelamento
// - Sinais de shutdown graceful
// - Contexto do usuário sistema
// - Diretório de trabalho correto
func (b *Bootstrap) Phase01_InitializeContext() error {
	b.logger.Info("🎯 FASE 1: Inicialização de Contexto e Sinais")
	b.logger.Debug("Configurando contexto base do sistema")

	// 1. Configurar sinais de shutdown
	b.logger.Debug("Configurando sinais de shutdown graceful")
	signal.Notify(b.shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// 2. Ajustar diretório de trabalho se necessário
	if err := b.adjustWorkingDirectory(); err != nil {
		return NewBootstrapError("Phase01", "working_directory", "Failed to adjust working directory", err)
	}

	// 3. Criar contexto do usuário sistema
	if err := b.createSystemUserContext(); err != nil {
		return NewBootstrapError("Phase01", "system_user_context", "Failed to create system user context", err)
	}

	// 4. Inicializar configuração base
	b.config = NewConfig(b.ctx)

	// 5. Configurar pprof server para debugging (desenvolvimento)
	b.startPprofServer()

	b.logger.Info("✅ Contexto e sinais inicializados com sucesso")
	return nil
}

// adjustWorkingDirectory ajusta o diretório de trabalho para a raiz do projeto
func (b *Bootstrap) adjustWorkingDirectory() error {
	wd, err := os.Getwd()
	if err != nil {
		return NewBootstrapError("Phase01", "getwd", "Failed to get current working directory", err)
	}

	// Se estamos no diretório cmd, subir um nível
	if filepath.Base(wd) == "cmd" {
		parentDir := filepath.Dir(wd)
		if err := os.Chdir(parentDir); err != nil {
			return NewBootstrapError("Phase01", "chdir", "Failed to change to project root directory", err)
		}

		b.logger.Info("📁 Diretório de trabalho ajustado",
			"from", wd,
			"to", parentDir)
	}

	return nil
}

// createSystemUserContext cria o contexto com informações do usuário sistema
func (b *Bootstrap) createSystemUserContext() error {
	// Criar informações do usuário sistema
	systemUser := usermodel.UserInfos{
		ID: usermodel.SystemUserID,
	}

	// Adicionar ao contexto
	ctx := context.WithValue(b.ctx, globalmodel.TokenKey, systemUser)
	ctx = context.WithValue(ctx, globalmodel.RequestIDKey, "server_initialization")

	// Atualizar contexto do bootstrap
	b.ctx = ctx

	b.logger.Debug("👤 Contexto do usuário sistema criado",
		"user_id", systemUser.ID)

	return nil
}

// startPprofServer inicia o servidor pprof para debugging em desenvolvimento
func (b *Bootstrap) startPprofServer() {
	// Verificar se pprof deve ser habilitado (apenas em desenvolvimento)
	if os.Getenv("ENABLE_PPROF") != "true" {
		b.logger.Debug("🔍 Servidor pprof desabilitado (definir ENABLE_PPROF=true para habilitar)")
		return
	}

	// Iniciar pprof em goroutine separada
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()

		b.logger.Info("🔍 Iniciando servidor pprof na porta 6060")

		// Import dinâmico do net/http/pprof
		// Isso habilita automaticamente os endpoints /debug/pprof/
		_ = http.DefaultServeMux // Garante que o pprof seja registrado

		// Iniciar servidor pprof
		pprofServer := &http.Server{
			Addr:    "localhost:6060", // Apenas localhost para segurança
			Handler: nil,              // Usa DefaultServeMux com pprof
		}

		b.logger.Info("✅ Servidor pprof iniciado em localhost:6060")

		// Aguardar contexto ser cancelado para shutdown graceful
		<-b.ctx.Done()

		b.logger.Debug("🔍 Parando servidor pprof")

		// Graceful shutdown do pprof
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := pprofServer.Shutdown(shutdownCtx); err != nil {
			b.logger.Error("Erro ao parar servidor pprof", "error", err)
		} else {
			b.logger.Info("✅ Servidor pprof parado gracefully")
		}
	}()
}

// Phase01Rollback executa rollback da Fase 1
func (b *Bootstrap) Phase01Rollback() error {
	b.logger.Info("🔄 Executando rollback da Fase 1")

	// Cancelar contexto
	if b.cancel != nil {
		b.cancel()
	}

	// Limpar sinais
	signal.Stop(b.shutdownChan)

	b.logger.Info("✅ Rollback da Fase 1 concluído")
	return nil
}
