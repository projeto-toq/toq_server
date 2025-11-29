package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/adapter/left/lambda/thumbnails"
	s3adapter "github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/adapter/right/s3"
	imageprocessing "github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/service/image_processing"
)

func main() {
	// 1. Init Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 2. Load Config
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("Failed to load AWS config", "error", err)
		os.Exit(1)
	}

	// 3. Init Adapters
	s3Client := s3.NewFromConfig(cfg)
	storageAdapter := s3adapter.NewS3Adapter(s3Client)

	// 4. Init Service
	svc := imageprocessing.NewThumbnailService(storageAdapter)

	// 5. Init Handler
	h := thumbnails.NewHandler(svc, logger)

	// 6. Start Lambda
	lambda.Start(h.HandleRequest)
}
