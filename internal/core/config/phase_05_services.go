package config

import (
	"fmt"
)

// Phase05_InitializeServices inicializa todos os serviços do sistema
// Esta fase configura:
// - Global Service (primeiro, sem dependências)
// - Permission Service (segundo, usado por outros)
// - User Service (terceiro, depende de permission)
// - Complex Service (quarto)
// - Listing Service (quinto)
// Ordem crítica para evitar dependências circulares
func (b *Bootstrap) Phase05_InitializeServices() error {
	b.logger.Info("🎯 FASE 5: Inicialização de Serviços")
	b.logger.Debug("Inicializando serviços na ordem correta")

	// Ordem crítica de inicialização baseada em dependências
	services := []struct {
		name string
		fn   func() error
	}{
		{"GlobalService", b.initializeGlobalService},
		{"PermissionService", b.initializePermissionService},
		{"UserService", b.initializeUserService},
		{"ComplexService", b.initializeComplexService},
		{"ListingService", b.initializeListingService},
	}

	for _, service := range services {
		b.logger.Debug("Inicializando serviço", "service", service.name)
		if err := service.fn(); err != nil {
			return NewBootstrapError("Phase05", service.name, fmt.Sprintf("Failed to initialize %s", service.name), err)
		}
		b.logger.Info("✅ Serviço inicializado", "service", service.name)
	}

	b.logger.Info("✅ Todos os serviços inicializados com sucesso")
	return nil
}

// initializeGlobalService inicializa o Global Service (primeiro, sem dependências)
func (b *Bootstrap) initializeGlobalService() error {
	b.logger.Debug("Inicializando Global Service")

	// Inicializar Global Service
	b.config.InitGlobalService()

	// Injetar GlobalService no cache Redis para resolver dependência circular
	// Nota: Implementação real faria isso se necessário

	b.logger.Debug("✅ Global Service inicializado")
	return nil
}

// initializePermissionService inicializa o Permission Service (segundo, usado por outros)
func (b *Bootstrap) initializePermissionService() error {
	b.logger.Debug("Inicializando Permission Service")

	// Inicializar Permission Service
	b.config.InitPermissionHandler()

	b.logger.Debug("✅ Permission Service inicializado")
	return nil
}

// initializeUserService inicializa o User Service (terceiro, depende de permission)
func (b *Bootstrap) initializeUserService() error {
	b.logger.Debug("Inicializando User Service")

	// Inicializar User Service (depende do Permission Service)
	b.config.InitUserHandler()

	b.logger.Debug("✅ User Service inicializado")
	return nil
}

// initializeComplexService inicializa o Complex Service (quarto)
func (b *Bootstrap) initializeComplexService() error {
	b.logger.Debug("Inicializando Complex Service")

	// Inicializar Complex Service
	b.config.InitComplexHandler()

	b.logger.Debug("✅ Complex Service inicializado")
	return nil
}

// initializeListingService inicializa o Listing Service (quinto)
func (b *Bootstrap) initializeListingService() error {
	b.logger.Debug("Inicializando Listing Service")

	// Inicializar Listing Service
	b.config.InitListingHandler()

	b.logger.Debug("✅ Listing Service inicializado")
	return nil
}

// Phase05Rollback executa rollback da Fase 5
func (b *Bootstrap) Phase05Rollback() error {
	b.logger.Info("🔄 Executando rollback da Fase 5")

	// Os serviços serão limpos automaticamente quando o contexto for cancelado
	// Não há necessidade de rollback específico nesta fase

	b.logger.Info("✅ Rollback da Fase 5 concluído")
	return nil
}
