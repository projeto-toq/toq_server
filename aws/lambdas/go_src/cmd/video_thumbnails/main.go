package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	videothumbnails "github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/adapter/left/lambda/video_thumbnails"
	s3adapter "github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/adapter/right/s3"
	videoprocessing "github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/service/video_processing"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("Failed to load AWS config", "error", err)
		os.Exit(1)
	}

	s3Client := s3.NewFromConfig(cfg)
	storageAdapter := s3adapter.NewS3Adapter(s3Client)

	ffmpegPath := os.Getenv("FFMPEG_PATH")
	seekSecond := resolveEnvInt("VIDEO_THUMBNAIL_SECOND", 1)
	width := resolveEnvInt("VIDEO_THUMBNAIL_WIDTH", 200)
	quality := resolveEnvInt("VIDEO_THUMBNAIL_QUALITY", 85)

	svc := videoprocessing.NewVideoThumbnailService(storageAdapter, ffmpegPath, seekSecond, width, quality)
	handler := videothumbnails.NewHandler(svc, logger)

	lambda.Start(handler.HandleRequest)
}

func resolveEnvInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	if v, err := strconv.Atoi(raw); err == nil && v > 0 {
		return v
	}
	return fallback
}
