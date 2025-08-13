// TOQ Server - Real Estate gRPC API Server
//
// This server implements a hexagonal architecture with the following layers:
// - Adapter Layer: gRPC handlers, external service integrations, database adapters
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
// Framework: gRPC with OpenTelemetry observability
// Storage: MySQL with Redis caching
// External Services: FCM, SMS, Email, CEP/CPF/CNPJ validation
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"net/http"
	_ "net/http/pprof" // Enable pprof debugging endpoints

	"github.com/giulio-alfieri/toq_server/internal/core/config"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// main is the entry point of the TOQ gRPC server.
// It orchestrates the complete server initialization following these phases:
// 1. Context and signal setup for graceful shutdown
// 2. Environment and logging configuration
// 3. Core infrastructure (database, telemetry, activity tracking)
// 4. Dependency injection using Factory Pattern
// 5. gRPC server initialization and startup
// 6. Background workers and graceful shutdown handling
func main() {
	// Phase 1: Setup context and graceful shutdown handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup graceful shutdown channel
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Start debug server for pprof (development/debugging)
	// This enables performance profiling endpoints at localhost:6060
	go func() {
		slog.Debug("Starting pprof debug server on :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			slog.Warn("pprof server failed", "error", err)
		}
	}()

	// Phase 2: Initialize system context with request metadata
	// Add system user context for initialization operations
	infos := usermodel.UserInfos{
		ID: usermodel.SystemUserID,
	}
	ctx = context.WithValue(ctx, globalmodel.TokenKey, infos)
	ctx = context.WithValue(ctx, globalmodel.RequestIDKey, "server_initialization")

	slog.Info("🚀 Starting TOQ Server initialization", "version", globalmodel.AppVersion)

	// Phase 3: Configuration and Environment Setup
	config := config.NewConfig(ctx)

	// Load and validate environment configuration
	if err := config.LoadEnv(); err != nil {
		slog.Error("❌ Failed to load environment configuration", "error", err)
		os.Exit(1)
	}
	slog.Info("✅ Environment configuration loaded successfully")

	// Initialize structured logging system
	config.InitializeLog()
	slog.Info("✅ Logging system initialized")

	slog.Info("🔧 TOQ API Server starting", "version", globalmodel.AppVersion)

	// Phase 4: Core Infrastructure Initialization
	// Initialize database connection with proper cleanup
	config.InitializeDatabase()
	defer func() {
		if err := config.GetDatabase().Close(); err != nil {
			slog.Error("❌ Error closing MySQL connection", "error", err)
			os.Exit(1)
		}
		slog.Info("✅ MySQL connection closed successfully")
	}()
	slog.Info("✅ Database connection established")

	// Initialize OpenTelemetry for observability (tracing + metrics)
	shutdownOtel, err := config.InitializeTelemetry()
	if err != nil {
		slog.Error("❌ Failed to initialize OpenTelemetry", "error", err)
		os.Exit(1)
	}
	defer shutdownOtel()
	slog.Info("✅ OpenTelemetry initialized (tracing + metrics)")

	// Initialize activity tracking system for user session management
	if err := config.InitializeActivityTracker(); err != nil {
		slog.Error("❌ Failed to initialize activity tracker", "error", err)
		os.Exit(1)
	}
	slog.Info("✅ Activity tracking system initialized")

	// Phase 5: gRPC Server and Dependency Injection
	// Initialize gRPC server with TLS and interceptors
	config.InitializeGRPC()
	slog.Info("✅ gRPC server configured with TLS and interceptors")

	// Inject all dependencies using Factory Pattern
	// This creates: validation adapters, external services, storage adapters, repositories
	resourceCleanup, err := config.InjectDependencies()
	if err != nil {
		slog.Error("❌ Failed to inject dependencies via Factory Pattern", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := resourceCleanup(); err != nil {
			slog.Error("❌ Error during resource cleanup", "error", err)
			os.Exit(1)
		}
		slog.Info("✅ All resources cleaned up successfully")
	}()
	slog.Info("✅ Dependency injection completed via Factory Pattern")

	// Configure activity tracker with user service (post-dependency injection)
	config.SetActivityTrackerUserService()
	slog.Info("✅ Activity tracker linked with user service")

	// Phase 6: Database Schema and Background Workers
	// Verify and initialize database schema if needed
	config.VerifyDatabase()
	slog.Info("✅ Database schema verified")

	// Initialize background goroutines (workers, cleanup tasks)
	config.InitializeGoRoutines()
	slog.Info("✅ Background workers initialized")

	// Phase 7: Server Startup and Runtime Management
	serviceQty, methodQty := config.GetInfos()
	slog.Info("🌟 TOQ Server ready to serve",
		"services", serviceQty,
		"methods", methodQty,
		"version", globalmodel.AppVersion)

	// Start server in goroutine to allow graceful shutdown handling
	var serverWg sync.WaitGroup
	serverWg.Add(1)

	go func() {
		defer serverWg.Done()
		slog.Info("🚀 Starting gRPC server")
		if err := config.GetGRPCServer().Serve(config.GetListener()); err != nil {
			slog.Error("❌ gRPC server failed", "error", err)
			cancel() // Trigger shutdown
		}
	}()

	// Phase 8: Graceful Shutdown Handling
	// Wait for shutdown signal or server error
	select {
	case <-shutdown:
		slog.Info("🛑 Shutdown signal received, initiating graceful shutdown...")
		cancel() // Cancel context to stop background workers
	case <-ctx.Done():
		slog.Info("🛑 Context cancelled, initiating graceful shutdown...")
	}

	// Graceful shutdown sequence
	slog.Info("⏳ Stopping gRPC server...")

	// Give server time to finish current requests
	shutdownTimeout := time.Second * 30
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	// Stop accepting new connections and wait for current requests
	config.GetGRPCServer().GracefulStop()

	// Wait for server goroutine to finish
	done := make(chan struct{})
	go func() {
		serverWg.Wait()
		close(done)
	}()

	// Wait for graceful stop or timeout
	select {
	case <-done:
		slog.Info("✅ gRPC server stopped gracefully")
	case <-shutdownCtx.Done():
		slog.Warn("⚠️ Graceful shutdown timeout, forcing stop")
		config.GetGRPCServer().Stop()
	}

	// Wait for background workers to complete
	slog.Info("⏳ Waiting for background workers to complete...")

	// Create a channel to signal when workers are done
	workersDone := make(chan struct{})
	go func() {
		config.GetWG().Wait()
		close(workersDone)
	}()

	// Wait for workers to complete or timeout
	workerTimeout := time.Second * 15
	select {
	case <-workersDone:
		slog.Info("✅ Background workers stopped gracefully")
	case <-time.After(workerTimeout):
		slog.Warn("⚠️ Background workers timeout, forcing shutdown")
	}

	slog.Info("👋 TOQ Server shutdown completed successfully")
}
