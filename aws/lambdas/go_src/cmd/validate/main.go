package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sfn"

	"github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/adapter/left/lambda/validate"
	s3adapter "github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/adapter/right/s3"
	sfnadapter "github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/adapter/right/step_functions"
	validateservice "github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/service/validate"
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

	sfnClient := sfn.NewFromConfig(cfg)
	workflowAdapter := sfnadapter.NewSfnAdapter(sfnClient)

	// 4. Init Service
	bucket := os.Getenv("MEDIA_BUCKET")
	if bucket == "" {
		bucket = "toq-listing-medias"
	}
	stateMachineArn := os.Getenv("STATE_MACHINE_ARN")
	if stateMachineArn == "" {
		stateMachineArn = "arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging"
	}

	svc := validateservice.NewValidateService(storageAdapter, workflowAdapter, bucket, stateMachineArn, logger)

	// 5. Init Handler
	h := validate.NewHandler(svc, logger)

	// 6. Start Lambda
	lambda.Start(h.HandleRequest)
}
