package main

import (
	"context"
	"log/slog"
	"os"

	"zip_generator/internal/adapter"
	"zip_generator/internal/model"
	"zip_generator/internal/service"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	zipService *service.ZipService
	logger     *slog.Logger
)

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	bucket := os.Getenv("MEDIA_BUCKET")
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	s3Adapter, err := adapter.NewS3Adapter(context.Background(), bucket, region)
	if err != nil {
		logger.Error("Failed to initialize S3 adapter", "error", err)
		os.Exit(1)
	}

	zipService = service.NewZipService(s3Adapter, logger)
}

func HandleRequest(ctx context.Context, event model.StepFunctionPayload) (model.ZipOutput, error) {
	return zipService.GenerateZip(ctx, event)
}

func main() {
	lambda.Start(HandleRequest)
}
