package videoprocessing

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/port"
)

// VideoThumbnailService extracts a frame from a video and uploads a resized JPEG thumbnail.
type VideoThumbnailService struct {
	storage    port.StoragePort
	ffmpegPath string
	seekSecond int
	width      int
	quality    int
}

// NewVideoThumbnailService configures the service with sensible defaults and an injected storage adapter.
func NewVideoThumbnailService(storage port.StoragePort, ffmpegPath string, seekSecond, width, quality int) *VideoThumbnailService {
	if ffmpegPath == "" {
		ffmpegPath = "/opt/ffmpeg/ffmpeg"
	}
	if seekSecond <= 0 {
		seekSecond = 1
	}
	if width <= 0 {
		width = 200
	}
	if quality <= 0 || quality > 100 {
		quality = 85
	}

	return &VideoThumbnailService{
		storage:    storage,
		ffmpegPath: ffmpegPath,
		seekSecond: seekSecond,
		width:      width,
		quality:    quality,
	}
}

// GenerateThumbnail downloads the video, extracts a frame, and uploads the JPEG to the processed path.
func (s *VideoThumbnailService) GenerateThumbnail(ctx context.Context, bucket, key string) (string, error) {
	inputPath, outputPath, cleanup, err := s.prepareTempFiles(ctx, bucket, key)
	if err != nil {
		return "", err
	}
	defer cleanup()

	if err := s.runFFmpeg(ctx, inputPath, outputPath); err != nil {
		return "", err
	}

	destKey, err := s.generateKey(key)
	if err != nil {
		return "", err
	}

	file, err := os.Open(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open generated thumbnail: %w", err)
	}
	defer file.Close()

	if err := s.storage.Upload(ctx, bucket, destKey, file, "image/jpeg"); err != nil {
		return "", err
	}

	return destKey, nil
}

func (s *VideoThumbnailService) prepareTempFiles(ctx context.Context, bucket, key string) (string, string, func(), error) {
	reader, err := s.storage.Download(ctx, bucket, key)
	if err != nil {
		return "", "", func() {}, fmt.Errorf("failed to download video: %w", err)
	}
	defer reader.Close()

	inputFile, err := os.CreateTemp("/tmp", "video-input-*.bin")
	if err != nil {
		return "", "", func() {}, fmt.Errorf("failed to create temp input file: %w", err)
	}

	if _, err := io.Copy(inputFile, reader); err != nil {
		cleanup := func() {
			inputFile.Close()
			os.Remove(inputFile.Name())
		}
		return "", "", cleanup, fmt.Errorf("failed to write temp input file: %w", err)
	}

	if err := inputFile.Close(); err != nil {
		cleanup := func() {
			os.Remove(inputFile.Name())
		}
		return "", "", cleanup, fmt.Errorf("failed to close temp input file: %w", err)
	}

	outputFile, err := os.CreateTemp("/tmp", "video-thumb-*.jpg")
	if err != nil {
		cleanup := func() {
			os.Remove(inputFile.Name())
		}
		return "", "", cleanup, fmt.Errorf("failed to create temp output file: %w", err)
	}
	outputPath := outputFile.Name()
	outputFile.Close()

	cleanup := func() {
		os.Remove(inputFile.Name())
		os.Remove(outputPath)
	}

	return inputFile.Name(), outputPath, cleanup, nil
}

func (s *VideoThumbnailService) runFFmpeg(ctx context.Context, inputPath, outputPath string) error {
	cmd := exec.CommandContext(ctx, s.ffmpegPath,
		"-y",
		"-ss", strconv.Itoa(s.seekSecond),
		"-i", inputPath,
		"-vframes", "1",
		"-vf", fmt.Sprintf("scale=%d:-1", s.width),
		"-q:v", strconv.Itoa(s.quality),
		"-f", "image2",
		outputPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w | output: %s", err, string(output))
	}

	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("thumbnail not created: %w", err)
	}

	return nil
}

func (s *VideoThumbnailService) generateKey(originalKey string) (string, error) {
	if !strings.Contains(originalKey, "raw/") {
		return "", fmt.Errorf("invalid key format: must contain 'raw/' segment")
	}

	parts := strings.SplitN(originalKey, "raw/", 2)
	prefix := parts[0]
	suffix := parts[1]

	dateRegex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}/`)
	cleanSuffix := dateRegex.ReplaceAllString(suffix, "")

	dir := path.Dir(cleanSuffix)
	file := path.Base(cleanSuffix)

	newKey := fmt.Sprintf("%sprocessed/%s/thumbnail/%s", prefix, dir, file)
	newKey = strings.ReplaceAll(newKey, "//", "/")
	return newKey, nil
}
