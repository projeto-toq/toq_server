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
//	@version		2.0
//	@description	TOQ Server - Real Estate HTTP API Server with hexagonal architecture.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	TOQ Development Team
//	@contact.url	https://github.com/projeto-toq/toq_server
//	@contact.email	projeto.toq@gmail.com

//	@license.name	MIT
//	@license.url	https://github.com/projeto-toq/toq_server/blob/main/LICENSE

//	@host		api.gca.dev.br
//	@BasePath	/api/v2

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	_ "github.com/projeto-toq/toq_server/docs" // This is required for Swagger

	"github.com/projeto-toq/toq_server/internal/core/config"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

type LogConfig struct {
	Level     string
	Format    string
	Output    string
	AddSource bool
}

func (c *LogConfig) slogLevel() slog.Level {
	switch strings.ToLower(c.Level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Handler que separa stdout/stderr por nível
type SplitLevelHandler struct {
	stdoutHandler slog.Handler
	stderrHandler slog.Handler
}

func (h *SplitLevelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.stdoutHandler.Enabled(ctx, level) || h.stderrHandler.Enabled(ctx, level)
}

func (h *SplitLevelHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level >= slog.LevelWarn {
		return h.stderrHandler.Handle(ctx, record)
	}
	return h.stdoutHandler.Handle(ctx, record)
}

func (h *SplitLevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SplitLevelHandler{
		stdoutHandler: h.stdoutHandler.WithAttrs(attrs),
		stderrHandler: h.stderrHandler.WithAttrs(attrs),
	}
}

func (h *SplitLevelHandler) WithGroup(name string) slog.Handler {
	return &SplitLevelHandler{
		stdoutHandler: h.stdoutHandler.WithGroup(name),
		stderrHandler: h.stderrHandler.WithGroup(name),
	}
}

// printUsage writes a friendly CLI help to the configured output.
func printUsage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "TOQ Server - Real Estate HTTP API Server\n\n")
	fmt.Fprintf(w, "Usage:\n  toq_server [flags]\n\n")
	fmt.Fprintf(w, "Flags:\n")
	// List all registered flags with defaults
	flag.PrintDefaults()
	fmt.Fprintf(w, "\nExamples:\n")
	fmt.Fprintf(w, "  toq_server --log-level=debug --log-format=text\n")
	fmt.Fprintf(w, "  toq_server -h\n")
}

func parseFlags() *LogConfig {
	config := &LogConfig{}

	// Custom usage output
	flag.Usage = printUsage

	flag.StringVar(&config.Level, "log-level", "info",
		"Log level (debug, info, warn, error)")
	flag.StringVar(&config.Format, "log-format", "json",
		"Log format (json, text)")
	flag.StringVar(&config.Output, "log-output", "stdout",
		"Log output (stdout, file, local)")
	flag.BoolVar(&config.AddSource, "log-add-source", false,
		"Add source code location to logs")

	// Accept GNU-style --help for convenience (in addition to the default -h)
	for _, a := range os.Args[1:] {
		if a == "--help" {
			flag.Usage()
			os.Exit(0)
		}
	}

	// Parse flags (on invalid flag, flag exits with code 2 and calls Usage)
	flag.Parse()
	return config
}

func setupEarlyLogger(config *LogConfig) *slog.Logger {
	// Parse level
	level := config.slogLevel()

	// Setup handler options
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: config.AddSource,
	}

	var handler slog.Handler

	switch strings.ToLower(config.Output) {
	case "stdout":
		// INFO/DEBUG → stdout, WARN/ERROR → stderr
		var stdoutHandler, stderrHandler slog.Handler

		if strings.ToLower(config.Format) == "text" {
			stdoutHandler = slog.NewTextHandler(os.Stdout, opts)
			stderrHandler = slog.NewTextHandler(os.Stderr, opts)
		} else {
			stdoutHandler = slog.NewJSONHandler(os.Stdout, opts)
			stderrHandler = slog.NewJSONHandler(os.Stderr, opts)
		}

		handler = &SplitLevelHandler{
			stdoutHandler: stdoutHandler,
			stderrHandler: stderrHandler,
		}

	case "file":
		// Todos os logs → ./logs/toq_server.log
		if err := os.MkdirAll("./logs", 0755); err != nil {
			slog.Error("Failed to create logs directory", "error", err)
			os.Exit(1)
		}

		logFile, err := os.OpenFile("./logs/toq_server.log",
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			slog.Error("Failed to open log file", "error", err)
			os.Exit(1)
		}

		if strings.ToLower(config.Format) == "text" {
			handler = slog.NewTextHandler(logFile, opts)
		} else {
			handler = slog.NewJSONHandler(logFile, opts)
		}

	case "local":
		// Todos os logs → ./toq_server.log
		logFile, err := os.OpenFile("./toq_server.log",
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			slog.Error("Failed to open log file", "error", err)
			os.Exit(1)
		}

		if strings.ToLower(config.Format) == "text" {
			handler = slog.NewTextHandler(logFile, opts)
		} else {
			handler = slog.NewJSONHandler(logFile, opts)
		}

	default:
		// Default para stdout com split
		var stdoutHandler, stderrHandler slog.Handler

		if strings.ToLower(config.Format) == "text" {
			stdoutHandler = slog.NewTextHandler(os.Stdout, opts)
			stderrHandler = slog.NewTextHandler(os.Stderr, opts)
		} else {
			stdoutHandler = slog.NewJSONHandler(os.Stdout, opts)
			stderrHandler = slog.NewJSONHandler(os.Stderr, opts)
		}

		handler = &SplitLevelHandler{
			stdoutHandler: stdoutHandler,
			stderrHandler: stderrHandler,
		}
	}

	return slog.New(handler)
}

// main é o ponto de entrada do servidor TOQ HTTP.
// Esta versão usa o novo sistema de bootstrap estruturado e robusto.
func main() {
	// 1. Parse CLI flags
	logConfig := parseFlags()
	globalmodel.SetLoggingRuntimeConfig(globalmodel.LoggingRuntimeConfig{
		Level:     logConfig.slogLevel(),
		Format:    logConfig.Format,
		Output:    logConfig.Output,
		AddSource: logConfig.AddSource,
	})

	// 2. Setup early logger with CLI overrides
	logger := setupEarlyLogger(logConfig)
	slog.SetDefault(logger)

	// 3. Log startup with configuration
	slog.Info("🚀 Iniciando TOQ Server Bootstrap",
		"version", globalmodel.AppVersion,
		"component", "bootstrap",
		"log_level", logConfig.Level,
		"log_format", logConfig.Format,
		"log_output", logConfig.Output)

	// 4. Criar instância do bootstrap
	bootstrap := config.NewBootstrap()

	// 5. Executar bootstrap completo
	if err := bootstrap.Bootstrap(); err != nil {
		slog.Error("❌ Falha crítica durante inicialização",
			"error", err,
			"component", "bootstrap")
		os.Exit(1)
	}

	// 6. Aguardar sinal de shutdown e executar graceful shutdown
	bootstrap.WaitShutdown()

	slog.Info("👋 Servidor TOQ finalizado com sucesso",
		"component", "bootstrap")
}
