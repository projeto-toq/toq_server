package config

// LifecycleManager gerencia o fechamento de recursos
type LifecycleManager struct {
	cleanupFuncs []func()
}

// NewLifecycleManager cria um novo gerenciador
func NewLifecycleManager() *LifecycleManager {
	return &LifecycleManager{
		cleanupFuncs: make([]func(), 0),
	}
}

// AddCleanupFunc adiciona uma função de limpeza
func (lm *LifecycleManager) AddCleanupFunc(f func()) {
	if f != nil {
		lm.cleanupFuncs = append(lm.cleanupFuncs, f)
	}
}

// Cleanup executa todas as funções de limpeza em ordem reversa
func (lm *LifecycleManager) Cleanup() {
	for i := len(lm.cleanupFuncs) - 1; i >= 0; i-- {
		lm.cleanupFuncs[i]()
	}
}
