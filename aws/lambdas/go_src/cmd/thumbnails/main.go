package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/disintegration/imaging"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

var (
	s3Client *s3.Client
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
}

type ThumbnailOutput struct {
	Status     string                          `json:"status"`
	Thumbnails []mediaprocessingmodel.JobAsset `json:"thumbnails"`
}

func HandleRequest(ctx context.Context, event mediaprocessingmodel.StepFunctionPayload) (mediaprocessingmodel.LambdaResponse, error) {
	logger.Info("Thumbnails Lambda started", "batch_id", event.BatchID, "assets_to_process", len(event.ValidAssets))

	generatedThumbnails := make([]mediaprocessingmodel.JobAsset, 0)

	for i, asset := range event.ValidAssets {
		logger.Info("Inspecting asset", "index", i, "key", asset.Key, "type", asset.Type)

		if !strings.HasPrefix(asset.Type, "PHOTO") {
			logger.Debug("Skipping non-photo asset", "key", asset.Key, "type", asset.Type)
			continue
		}

		logger.Info("Processing photo", "key", asset.Key)

		// Download
		resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &asset.Key,
		})
		if err != nil {
			logger.Error("Failed to download asset", "key", asset.Key, "error", err)
			continue
		}
		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			logger.Error("Failed to decode image", "key", asset.Key, "error", err)
			continue
		}

		// Generate Thumbnails (Small, Medium, Large)
		sizes := map[string]int{
			"small":  320,
			"medium": 640,
			"large":  1280,
		}

		for sizeName, width := range sizes {
			resized := imaging.Resize(img, width, 0, imaging.Lanczos)

			buf := new(bytes.Buffer)
			if err := jpeg.Encode(buf, resized, nil); err != nil {
				logger.Error("Failed to encode thumbnail", "key", asset.Key, "size", sizeName, "error", err)
				continue
			}

			thumbKey := fmt.Sprintf("%d/processed/%s/%s_%s.jpg", event.ListingID, sizeName, asset.Key, sizeName)

			parts := strings.Split(asset.Key, "/")
			if len(parts) > 2 {
				newKey := strings.Replace(asset.Key, "/raw/", "/processed/", 1)
				extIdx := strings.LastIndex(newKey, ".")
				if extIdx != -1 {
					thumbKey = newKey[:extIdx] + "_" + sizeName + ".jpg"
				} else {
					thumbKey = newKey + "_" + sizeName + ".jpg"
				}
			}

			// LOG: Before upload
			logger.Debug("Uploading thumbnail",
				"original_key", asset.Key,
				"size", sizeName,
				"target_key", thumbKey,
			)

			_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
				Bucket:      &bucket,
				Key:         &thumbKey,
				Body:        bytes.NewReader(buf.Bytes()),
				ContentType: aws.String("image/jpeg"),
			})
			if err != nil {
				logger.Error("Failed to upload thumbnail", "key", thumbKey, "error", err)
				continue
			}

			generatedThumbnails = append(generatedThumbnails, mediaprocessingmodel.JobAsset{
				Key:       thumbKey,
				Type:      "THUMBNAIL_" + strings.ToUpper(sizeName),
				SourceKey: asset.Key, // CRUCIAL: Keep link with original
			})
		}
	}

	// LOG: Final output
	logger.Info("Thumbnails generation finished",
		"batch_id", event.BatchID,
		"generated_count", len(generatedThumbnails),
	)

	return mediaprocessingmodel.LambdaResponse{
		Body: ThumbnailOutput{
			Status:     "thumbnails_generated",
			Thumbnails: generatedThumbnails,
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
