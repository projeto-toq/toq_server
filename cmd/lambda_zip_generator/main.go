package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	ziphandler "github.com/projeto-toq/toq_server/internal/adapter/left/lambda/zip_handler"
	"github.com/projeto-toq/toq_server/internal/core/config"
	"github.com/projeto-toq/toq_server/internal/core/factory"
	zipprocessingservice "github.com/projeto-toq/toq_server/internal/core/service/zip_processing_service"
)

func main() {
	// Initialize logger
	logger := slog.Default()

	// Load configuration
	bootstrap := config.NewBootstrap()
	if err := bootstrap.Phase02_LoadConfiguration(); err != nil {
		logger.Error("Failed to load configuration", "err", err)
		os.Exit(1)
	}

	env := bootstrap.GetEnvironment()
	if env == nil {
		logger.Error("Environment is nil after loading configuration")
		os.Exit(1)
	}

	lm := bootstrap.GetLifecycleManager()

	// Create Factory
	adapterFactory := factory.NewAdapterFactory(lm)

	// Create External Service Adapters (includes S3)
	externalAdapters, err := adapterFactory.CreateExternalServiceAdapters(context.Background(), env)
	if err != nil {
		logger.Error("Failed to create external adapters", "err", err)
		os.Exit(1)
	}

	// Initialize Service
	zipService := zipprocessingservice.NewZipProcessingService(externalAdapters.ListingMediaStorage)

	// Initialize Handler
	handler := ziphandler.NewZipHandler(zipService, logger)

	// Start Lambda
	lambda.Start(handler.HandleRequest)
}
