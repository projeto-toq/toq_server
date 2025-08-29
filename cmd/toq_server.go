// TOQ Server - Real Estate HTTP API Server
//
// This server implements a hexagonal architecture with the following layers:
// - Adapter Layer: HTTP handlers, external service integrations, database adapters
// - Port Layer: Interfaces defining contracts between layers
// - Core Layer: Business logic, domain models, services
//
// The server follows Go best practices:
// - Proper error handling with context
// - Structured logging with slog
// - Resource cleanup with defer statements
// - Dependency injection through factory pattern
// - Clean shutdown with signal handling
//
// Architecture: Hexagonal (Ports & Adapters)
// Framework: HTTP/Gin with OpenTelemetry observability
// Storage: MySQL with Redis caching
// External Services: FCM, SMS, Email, CEP/CPF/CNPJ validation

//	@title			TOQ Server API
//	@version		1.0
//	@description	TOQ Server - Real Estate HTTP API Server with hexagonal architecture
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	TOQ Development Team
//	@contact.url	https://github.com/giulio-alfieri/toq_server
//	@contact.email	support@toq.com

//	@license.name	MIT
//	@license.url	https://github.com/giulio-alfieri/toq_server/blob/main/LICENSE

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
package main

import (
	"log/slog"
	"os"

	_ "github.com/giulio-alfieri/toq_server/docs" // This is required for Swagger

	"github.com/giulio-alfieri/toq_server/internal/core/config"
)

// main é o ponto de entrada do servidor TOQ HTTP.
// Esta versão usa o novo sistema de bootstrap estruturado e robusto.
func main() {
	// 1. Criar instância do bootstrap
	bootstrap := config.NewBootstrap()

	// 2. Executar bootstrap completo
	if err := bootstrap.Bootstrap(); err != nil {
		slog.Error("❌ Falha crítica durante inicialização", "error", err)
		os.Exit(1)
	}

	// 3. Aguardar sinal de shutdown e executar graceful shutdown
	bootstrap.WaitShutdown()

	slog.Info("👋 Servidor TOQ finalizado com sucesso")
}
