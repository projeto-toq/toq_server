package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

var (
	s3Client *s3.Client
	uploader *manager.Uploader
	logger   *slog.Logger
	bucket   string
)

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	bucket = os.Getenv("MEDIA_BUCKET")
	if bucket == "" {
		bucket = "toq-listing-medias"
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("Failed to load AWS config", "error", err)
		os.Exit(1)
	}

	s3Client = s3.NewFromConfig(cfg)
	// Uploader facilita multipart upload para streams
	uploader = manager.NewUploader(s3Client)
}

type ZipOutput struct {
	ZipKey       string   `json:"zipKey"`
	AssetsZipped int      `json:"assetsZipped"`
	ZipBundles   []string `json:"zipBundles"`
}

func HandleRequest(ctx context.Context, event mediaprocessingmodel.StepFunctionPayload) (mediaprocessingmodel.LambdaResponse, error) {
	logger.Info("Zip Lambda started (Streaming Mode)", "batch_id", event.BatchID, "assets_count", len(event.ValidAssets))

	zipKey := fmt.Sprintf("%d/zip/complete_%d_%d.zip", event.ListingID, event.BatchID, time.Now().Unix())

	// Pipe: writer (zip) -> reader (s3 upload)
	pr, pw := io.Pipe()
	zipWriter := zip.NewWriter(pw)

	// Canal para capturar erro da goroutine de escrita
	errCh := make(chan error, 1)
	zippedCount := 0

	go func() {
		defer pw.Close()
		defer zipWriter.Close()

		for _, asset := range event.ValidAssets {
			// 1. Download do S3 (Stream)
			resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &asset.Key,
			})
			if err != nil {
				logger.Error("Failed to download asset", "key", asset.Key, "error", err)
				errCh <- fmt.Errorf("download failed for %s: %w", asset.Key, err)
				return
			}

			// 2. Criar entrada no ZIP
			// Remove prefixo "listingID/raw/" para limpar a estrutura do zip
			entryName := asset.Key
			prefixToRemove := fmt.Sprintf("%d/raw/", event.ListingID)
			if strings.HasPrefix(asset.Key, prefixToRemove) {
				entryName = strings.TrimPrefix(asset.Key, prefixToRemove)
			}

			f, err := zipWriter.Create(entryName)
			if err != nil {
				resp.Body.Close()
				logger.Error("Failed to create zip entry", "key", asset.Key, "error", err)
				errCh <- fmt.Errorf("zip entry failed for %s: %w", asset.Key, err)
				return
			}

			// 3. Copiar stream (S3 -> Zip)
			if _, err := io.Copy(f, resp.Body); err != nil {
				resp.Body.Close()
				logger.Error("Failed to write zip content", "key", asset.Key, "error", err)
				errCh <- fmt.Errorf("zip write failed for %s: %w", asset.Key, err)
				return
			}
			resp.Body.Close()
			zippedCount++
		}
		close(errCh) // Sucesso
	}()

	// Upload bloqueante lendo do PipeReader
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &zipKey,
		Body:        pr,
		ContentType: aws.String("application/zip"),
	})

	// Verificar se houve erro na goroutine
	if zipErr := <-errCh; zipErr != nil {
		return mediaprocessingmodel.LambdaResponse{}, zipErr
	}

	if err != nil {
		return mediaprocessingmodel.LambdaResponse{}, fmt.Errorf("failed to upload zip: %w", err)
	}

	return mediaprocessingmodel.LambdaResponse{
		Body: ZipOutput{
			ZipKey:       zipKey,
			AssetsZipped: zippedCount,
			ZipBundles:   []string{zipKey},
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
