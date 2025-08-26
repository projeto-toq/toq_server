package config

import (
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/factory"
)

// assignValidationAdapters atribui os validation adapters às propriedades do config
func (c *config) assignValidationAdapters(adapters factory.ValidationAdapters) {
	slog.Debug("Assigning validation adapters to config")
	c.cep = adapters.CEP
	c.cpf = adapters.CPF
	c.cnpj = adapters.CNPJ
}

// assignExternalServiceAdapters atribui os external service adapters às propriedades do config
func (c *config) assignExternalServiceAdapters(adapters factory.ExternalServiceAdapters) {
	slog.Debug("Assigning external service adapters to config")
	c.firebaseCloudMessaging = adapters.FCM
	c.email = adapters.Email
	c.sms = adapters.SMS
	c.cloudStorage = adapters.CloudStorage
}

// assignStorageAdapters atribui os storage adapters às propriedades do config
func (c *config) assignStorageAdapters(adapters factory.StorageAdapters) {
	slog.Debug("Assigning storage adapters to config")
	c.database = adapters.Database
	c.cache = adapters.Cache
}

// assignRepositoryAdapters atribui os repository adapters para uso nos serviços
func (c *config) assignRepositoryAdapters(adapters factory.RepositoryAdapters) {
	slog.Debug("Assigning repository adapters to config")
	c.sessionRepo = adapters.Session
	// Os outros repositórios serão usados diretamente nos métodos de inicialização
	c.repositoryAdapters = &adapters
}

// initializeServices inicializa todos os serviços usando os adapters criados
func (c *config) initializeServices() {
	slog.Info("Initializing services")

	c.InitGlobalService()
	c.InitComplexHandler()
	c.InitListingHandler()
	c.InitUserHandler()

	slog.Info("All services initialized successfully")
}
