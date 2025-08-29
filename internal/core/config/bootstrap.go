package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

// Bootstrap é o orquestrador principal do sistema de inicialização
// Responsável por coordenar todas as fases de inicialização de forma ordenada e robusta
type Bootstrap struct {
	// Contexto e controle de ciclo de vida
	ctx              context.Context
	cancel           context.CancelFunc
	lifecycleManager *LifecycleManager

	// Configuração e estado
	config    ConfigInterface
	env       *globalmodel.Environment
	startTime time.Time

	// Canais de comunicação
	shutdownChan chan os.Signal
	healthStatus HealthStatus

	// Sincronização
	wg         *sync.WaitGroup
	phaseMutex sync.RWMutex

	// Logging estruturado
	logger *slog.Logger
}

// HealthStatus representa o status de saúde do sistema
type HealthStatus struct {
	Overall    bool                   `json:"overall"`
	Phases     map[string]PhaseHealth `json:"phases"`
	StartedAt  time.Time              `json:"started_at"`
	Uptime     time.Duration          `json:"uptime"`
	ErrorCount int                    `json:"error_count"`
}

// PhaseHealth representa o status de uma fase específica
type PhaseHealth struct {
	Name       string        `json:"name"`
	Status     string        `json:"status"` // "pending", "running", "completed", "failed"
	StartedAt  time.Time     `json:"started_at,omitempty"`
	Duration   time.Duration `json:"duration,omitempty"`
	Error      string        `json:"error,omitempty"`
	RetryCount int           `json:"retry_count,omitempty"`
}

// BootstrapConfig contém configurações para o bootstrap
type BootstrapConfig struct {
	ShutdownTimeout time.Duration
	PhaseTimeout    time.Duration
	MaxRetries      int
	RetryDelay      time.Duration
}

// NewBootstrap cria uma nova instância do sistema de bootstrap
func NewBootstrap() *Bootstrap {
	ctx, cancel := context.WithCancel(context.Background())

	return &Bootstrap{
		ctx:              ctx,
		cancel:           cancel,
		lifecycleManager: NewLifecycleManager(),
		startTime:        time.Now(),
		shutdownChan:     make(chan os.Signal, 1),
		wg:               &sync.WaitGroup{},
		healthStatus: HealthStatus{
			Phases: make(map[string]PhaseHealth),
		},
		logger: slog.Default(),
	}
}

// Bootstrap executa todo o processo de inicialização do sistema
// Esta é a função principal que orquestra todas as fases
func (b *Bootstrap) Bootstrap() error {
	b.logger.Info("🚀 Iniciando TOQ Server Bootstrap",
		"version", globalmodel.AppVersion,
		"timestamp", b.startTime.Format(time.RFC3339))

	defer func() {
		if r := recover(); r != nil {
			b.logger.Error("💥 Panic durante inicialização", "panic", r)
			b.Shutdown()
			os.Exit(1)
		}
	}()

	// Executar todas as fases em ordem
	phases := []struct {
		name string
		fn   func() error
	}{
		{"Phase01_InitializeContext", b.Phase01_InitializeContext},
		{"Phase02_LoadConfiguration", b.Phase02_LoadConfiguration},
		{"Phase03_InitializeInfrastructure", b.Phase03_InitializeInfrastructure},
		{"Phase04_InjectDependencies", b.Phase04_InjectDependencies},
		{"Phase05_InitializeServices", b.Phase05_InitializeServices},
		{"Phase06_ConfigureHandlers", b.Phase06_ConfigureHandlers},
		{"Phase07_StartBackgroundWorkers", b.Phase07_StartBackgroundWorkers},
		{"Phase08_StartServer", b.Phase08_StartServer},
	}

	for _, phase := range phases {
		if err := b.executePhase(phase.name, phase.fn); err != nil {
			b.logger.Error("❌ Falha na fase de inicialização",
				"phase", phase.name,
				"error", err)
			return fmt.Errorf("bootstrap failed at phase %s: %w", phase.name, err)
		}
	}

	b.logger.Info("🎉 TOQ Server inicializado com sucesso",
		"total_time", time.Since(b.startTime))

	return nil
}

// executePhase executa uma fase com controle de erro e métricas
func (b *Bootstrap) executePhase(phaseName string, phaseFunc func() error) error {
	b.phaseMutex.Lock()
	b.healthStatus.Phases[phaseName] = PhaseHealth{
		Name:      phaseName,
		Status:    "running",
		StartedAt: time.Now(),
	}
	b.phaseMutex.Unlock()

	start := time.Now()
	b.logger.Info("▶️ Executando fase",
		"phase", phaseName,
		"timestamp", start.Format(time.RFC3339))

	// Executar fase com timeout
	done := make(chan error, 1)
	go func() {
		done <- phaseFunc()
	}()

	select {
	case err := <-done:
		duration := time.Since(start)

		b.phaseMutex.Lock()
		phase := b.healthStatus.Phases[phaseName]
		phase.Duration = duration

		if err != nil {
			phase.Status = "failed"
			phase.Error = err.Error()
			b.healthStatus.ErrorCount++
			b.logger.Error("❌ Fase falhou",
				"phase", phaseName,
				"duration", duration,
				"error", err)
		} else {
			phase.Status = "completed"
			b.logger.Info("✅ Fase concluída",
				"phase", phaseName,
				"duration", duration)
		}

		b.healthStatus.Phases[phaseName] = phase
		b.phaseMutex.Unlock()

		return err

	case <-time.After(2 * time.Minute): // Timeout de 2 minutos por fase
		b.phaseMutex.Lock()
		phase := b.healthStatus.Phases[phaseName]
		phase.Status = "failed"
		phase.Error = "timeout after 2 minutes"
		phase.Duration = time.Since(start)
		b.healthStatus.Phases[phaseName] = phase
		b.healthStatus.ErrorCount++
		b.phaseMutex.Unlock()

		b.logger.Error("⏰ Timeout na fase",
			"phase", phaseName,
			"duration", time.Since(start))

		return fmt.Errorf("phase %s timed out after 2 minutes", phaseName)
	}
}

// Shutdown executa o desligamento graceful do sistema
func (b *Bootstrap) Shutdown() {
	b.logger.Info("🛑 Iniciando shutdown graceful")

	// Cancelar contexto para parar workers
	b.cancel()

	// Aguardar workers terminarem
	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		b.logger.Info("✅ Todos os workers terminaram")
	case <-time.After(30 * time.Second):
		b.logger.Warn("⚠️ Timeout aguardando workers")
	}

	// Executar cleanup
	b.lifecycleManager.Cleanup()

	b.logger.Info("👋 TOQ Server shutdown concluído",
		"uptime", time.Since(b.startTime))
}

// Health retorna o status de saúde atual do sistema
func (b *Bootstrap) Health() HealthStatus {
	b.phaseMutex.RLock()
	defer b.phaseMutex.RUnlock()

	status := b.healthStatus
	status.Uptime = time.Since(b.startTime)

	// Calcular status geral
	status.Overall = true
	for _, phase := range status.Phases {
		if phase.Status == "failed" {
			status.Overall = false
			break
		}
	}

	return status
}

// WaitShutdown aguarda sinal de shutdown
func (b *Bootstrap) WaitShutdown() {
	signal.Notify(b.shutdownChan, syscall.SIGINT, syscall.SIGTERM)
	<-b.shutdownChan
	b.logger.Info("🛑 Sinal de shutdown recebido")
}
