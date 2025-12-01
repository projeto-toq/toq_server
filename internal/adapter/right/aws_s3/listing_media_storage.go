package s3adapter

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	storageport "github.com/projeto-toq/toq_server/internal/core/port/right/storage"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	defaultUploadTTLSeconds   = 900
	defaultDownloadTTLSeconds = 3600
)

var (
	filenameSanitizer = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)
	segmentSanitizer  = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
	hexChecksumRegex  = regexp.MustCompile(`^(?i:[0-9a-f]+)$`)
)

// ListingMediaStorageAdapter implements ListingMediaStoragePort backed by S3.
type ListingMediaStorageAdapter struct {
	base         *S3Adapter
	bucket       string
	uploadTTL    time.Duration
	downloadTTL  time.Duration
	putPresigner *s3.PresignClient
	getPresigner *s3.PresignClient
	adminClient  *s3.Client
	readerClient *s3.Client
}

// NewListingMediaStorageAdapter builds an adapter dedicated to listing media prefixes.
func NewListingMediaStorageAdapter(base *S3Adapter, env *globalmodel.Environment) *ListingMediaStorageAdapter {
	if base == nil || base.adminClient == nil || base.readerClient == nil || base.listingBucketName == "" {
		return nil
	}

	uploadTTL := resolveTTL(env.MediaProcessing.Storage.UploadURLTTLSeconds, env.S3.SignedURL.UploadTTLSeconds, defaultUploadTTLSeconds)
	downloadTTL := resolveTTL(env.MediaProcessing.Storage.DownloadURLTTLSeconds, env.S3.SignedURL.DownloadTTLSeconds, defaultDownloadTTLSeconds)

	return &ListingMediaStorageAdapter{
		base:         base,
		bucket:       base.listingBucketName,
		uploadTTL:    time.Duration(uploadTTL) * time.Second,
		downloadTTL:  time.Duration(downloadTTL) * time.Second,
		putPresigner: s3.NewPresignClient(base.adminClient),
		getPresigner: s3.NewPresignClient(base.readerClient),
		adminClient:  base.adminClient,
		readerClient: base.readerClient,
	}
}

func resolveTTL(primary, fallback, def int) int {
	if primary > 0 {
		return primary
	}
	if fallback > 0 {
		return fallback
	}
	return def
}

// GenerateRawUploadURL builds a pre-signed PUT URL for raw uploads.
func (a *ListingMediaStorageAdapter) GenerateRawUploadURL(ctx context.Context, listingID uint64, asset mediaprocessingmodel.MediaAsset, contentType, checksum string) (storageport.SignedURL, error) {
	ctx = utils.ContextWithLogger(ctx)
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "ListingMediaStorage.GenerateRawUploadURL")
	if err != nil {
		return storageport.SignedURL{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if err := a.ensureClients(); err != nil {
		utils.SetSpanError(ctx, err)
		return storageport.SignedURL{}, err
	}

	key := asset.S3KeyRaw()
	if key == "" {
		key = a.buildObjectKey(listingID, "raw", asset)
	}

	checksum, err = normalizeChecksum(checksum)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return storageport.SignedURL{}, derrors.Validation("invalid checksum", map[string]string{"checksum": err.Error()})
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(resolveContentType(contentType)),
	}
	if checksum != "" {
		input.ChecksumSHA256 = aws.String(checksum)
	}

	presignOutput, err := a.putPresigner.PresignPutObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = a.uploadTTL
	})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("adapter.s3.listing.generate_upload_url_failed", "listing_id", listingID, "key", key, "error", err)
		return storageport.SignedURL{}, derrors.Infra("failed to generate upload URL", err)
	}

	headers := map[string]string{
		"Content-Type": aws.ToString(input.ContentType),
	}
	if checksum != "" {
		headers["x-amz-checksum-sha256"] = checksum
	}

	return storageport.SignedURL{
		URL:       presignOutput.URL,
		Method:    httpMethodPut,
		Headers:   headers,
		ExpiresIn: a.uploadTTL,
		ObjectKey: key,
	}, nil
}

// GenerateProcessedDownloadURL builds a pre-signed GET URL for processed assets.
func (a *ListingMediaStorageAdapter) GenerateProcessedDownloadURL(ctx context.Context, listingID uint64, asset mediaprocessingmodel.MediaAsset, resolution string) (storageport.SignedURL, error) {
	ctx = utils.ContextWithLogger(ctx)
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "ListingMediaStorage.GenerateProcessedDownloadURL")
	if err != nil {
		return storageport.SignedURL{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if err := a.ensureClients(); err != nil {
		utils.SetSpanError(ctx, err)
		return storageport.SignedURL{}, err
	}

	key := a.resolveProcessedKey(listingID, asset, resolution)
	presignOutput, err := a.getPresigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = a.downloadTTL
	})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("adapter.s3.listing.generate_download_url_failed", "listing_id", listingID, "key", key, "error", err)
		return storageport.SignedURL{}, derrors.Infra("failed to generate download URL", err)
	}

	return storageport.SignedURL{
		URL:       presignOutput.URL,
		Method:    httpMethodGet,
		Headers:   map[string]string{},
		ExpiresIn: a.downloadTTL,
		ObjectKey: key,
	}, nil
}

// GenerateDownloadURL builds a pre-signed GET URL for any key.
func (a *ListingMediaStorageAdapter) GenerateDownloadURL(ctx context.Context, key string) (storageport.SignedURL, error) {
	ctx = utils.ContextWithLogger(ctx)
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "ListingMediaStorage.GenerateDownloadURL")
	if err != nil {
		return storageport.SignedURL{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if err := a.ensureClients(); err != nil {
		utils.SetSpanError(ctx, err)
		return storageport.SignedURL{}, err
	}

	presignOutput, err := a.getPresigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = a.downloadTTL
	})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("adapter.s3.listing.generate_generic_download_url_failed", "key", key, "error", err)
		return storageport.SignedURL{}, derrors.Infra("failed to generate download URL", err)
	}

	return storageport.SignedURL{
		URL:       presignOutput.URL,
		Method:    httpMethodGet,
		Headers:   map[string]string{},
		ExpiresIn: a.downloadTTL,
		ObjectKey: key,
	}, nil
}

// ValidateObjectChecksum retrieves metadata and compares checksum.
func (a *ListingMediaStorageAdapter) ValidateObjectChecksum(ctx context.Context, bucketKey string, expectedChecksum string) (storageport.StorageObjectMetadata, error) {
	ctx = utils.ContextWithLogger(ctx)
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "ListingMediaStorage.ValidateObjectChecksum")
	if err != nil {
		return storageport.StorageObjectMetadata{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if err := a.ensureClients(); err != nil {
		utils.SetSpanError(ctx, err)
		return storageport.StorageObjectMetadata{}, err
	}

	output, err := a.readerClient.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket:       aws.String(a.bucket),
		Key:          aws.String(bucketKey),
		ChecksumMode: types.ChecksumModeEnabled,
	})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("adapter.s3.listing.head_object_failed", "key", bucketKey, "error", err)
		return storageport.StorageObjectMetadata{}, derrors.Infra("failed to retrieve object metadata", err)
	}

	expected, err := normalizeChecksum(expectedChecksum)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return storageport.StorageObjectMetadata{}, derrors.Validation("invalid checksum", map[string]string{"checksum": err.Error()})
	}
	if expected != "" && output.ChecksumSHA256 != nil && expected != aws.ToString(output.ChecksumSHA256) {
		err := derrors.Conflict("checksum mismatch")
		utils.SetSpanError(ctx, err)
		return storageport.StorageObjectMetadata{}, err
	}

	metadata := storageport.StorageObjectMetadata{
		SizeInBytes: aws.ToInt64(output.ContentLength),
		Checksum:    aws.ToString(output.ChecksumSHA256),
		ETag:        aws.ToString(output.ETag),
		ContentType: aws.ToString(output.ContentType),
	}
	if output.LastModified != nil {
		metadata.LastModified = *output.LastModified
	}

	return metadata, nil
}

// DeleteObject removes a key under the listing bucket.
func (a *ListingMediaStorageAdapter) DeleteObject(ctx context.Context, bucketKey string) error {
	ctx = utils.ContextWithLogger(ctx)
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "ListingMediaStorage.DeleteObject")
	if err != nil {
		return derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if err := a.ensureClients(); err != nil {
		utils.SetSpanError(ctx, err)
		return err
	}

	if err := a.base.DeleteBucketObject(ctx, a.bucket, bucketKey); err != nil {
		utils.SetSpanError(ctx, err)
		return derrors.Infra("failed to delete object", err)
	}

	return nil
}

func (a *ListingMediaStorageAdapter) ensureClients() error {
	if a.base == nil || a.bucket == "" || a.putPresigner == nil || a.getPresigner == nil {
		return derrors.Infra("listing media storage adapter not configured", nil)
	}
	return nil
}

func (a *ListingMediaStorageAdapter) buildObjectKey(listingID uint64, stage string, asset mediaprocessingmodel.MediaAsset) string {
	mediaTypeSegment := mediaTypePathSegment(asset.AssetType())
	if mediaTypeSegment == "" {
		mediaTypeSegment = "misc"
	}

	var metaMap map[string]string
	if asset.Metadata() != "" {
		_ = json.Unmarshal([]byte(asset.Metadata()), &metaMap)
	}

	reference := metaMap["client_id"]
	if reference == "" {
		reference = metaMap["clientId"]
	}
	if reference == "" && asset.Sequence() > 0 {
		prefix := "seq"
		assetTypeStr := string(asset.AssetType())
		if strings.Contains(assetTypeStr, "VERTICAL") {
			prefix = "vertical"
		} else if strings.Contains(assetTypeStr, "HORIZONTAL") {
			prefix = "horizontal"
		}
		reference = fmt.Sprintf("%s-%02d", prefix, asset.Sequence())
	}
	if reference == "" {
		reference = fmt.Sprintf("asset-%s", uuid.NewString())
	}
	reference = sanitizeSegment(reference)
	if reference == "" {
		reference = uuid.NewString()
	}

	filename := metaMap["filename"]
	if filename == "" {
		filename = "file"
	}
	filename = sanitizeFilename(filename, metaMap["content_type"])
	// dateSegment removed to keep path clean
	// dateSegment := time.Now().UTC().Format("2006-01-02")

	return fmt.Sprintf("%d/%s/%s/%s-%s", listingID, stage, mediaTypeSegment, reference, filename)
}

// extractFilename helper to get the filename part from existing keys or metadata
func (a *ListingMediaStorageAdapter) extractFilename(asset mediaprocessingmodel.MediaAsset) string {
	// Tenta extrair do S3KeyProcessed se existir
	if asset.S3KeyProcessed() != "" {
		_, file := path.Split(asset.S3KeyProcessed())
		if file != "" {
			return file
		}
	}
	// Tenta do S3KeyRaw
	if asset.S3KeyRaw() != "" {
		_, file := path.Split(asset.S3KeyRaw())
		if file != "" {
			return file
		}
	}

	// Fallback para metadata
	var metaMap map[string]string
	if asset.Metadata() != "" {
		_ = json.Unmarshal([]byte(asset.Metadata()), &metaMap)
	}

	filename := metaMap["filename"]
	if filename != "" {
		return sanitizeFilename(filename, metaMap["content_type"])
	}

	// Fallback final
	return fmt.Sprintf("asset-%d-%d.bin", asset.ListingIdentityID(), asset.Sequence())
}

func (a *ListingMediaStorageAdapter) resolveProcessedKey(listingID uint64, asset mediaprocessingmodel.MediaAsset, resolution string) string {
	// 1. Se for ZIP, não tem resolução
	if asset.AssetType() == mediaprocessingmodel.MediaAssetTypeZip {
		if asset.S3KeyProcessed() != "" {
			return asset.S3KeyProcessed()
		}
		return a.buildObjectKey(listingID, "processed", asset)
	}

	// 2. Determinar componentes básicos
	stage := "processed"
	if asset.AssetType() == mediaprocessingmodel.MediaAssetTypeThumbnail {
		stage = "thumb"
	}

	mediaTypeSegment := mediaTypePathSegment(asset.AssetType())

	// 3. Obter nome do arquivo
	filename := a.extractFilename(asset)

	// 4. Construir caminho: {listingId}/{stage}/{mediaType}/{resolution}/{filename}
	if resolution == "" {
		resolution = "original"
	}

	// /{listingId}/processed/{mediaType}/{size}/{uuid}.{ext}
	return fmt.Sprintf("%d/%s/%s/%s/%s", listingID, stage, mediaTypeSegment, resolution, filename)
}

func mediaTypePathSegment(assetType mediaprocessingmodel.MediaAssetType) string {
	switch assetType {
	case mediaprocessingmodel.MediaAssetTypePhotoVertical:
		return "photo/vertical"
	case mediaprocessingmodel.MediaAssetTypePhotoHorizontal:
		return "photo/horizontal"
	case mediaprocessingmodel.MediaAssetTypeVideoVertical:
		return "video/vertical"
	case mediaprocessingmodel.MediaAssetTypeVideoHorizontal:
		return "video/horizontal"
	case mediaprocessingmodel.MediaAssetTypeProjectDoc:
		return "project/doc"
	case mediaprocessingmodel.MediaAssetTypeProjectRender:
		return "project/render"
	case mediaprocessingmodel.MediaAssetTypeThumbnail:
		return "thumb"
	case mediaprocessingmodel.MediaAssetTypeZip:
		return "zip"
	default:
		return "misc"
	}
}

func sanitizeFilename(name, contentType string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		trimmed = fmt.Sprintf("asset-%s", uuid.NewString())
	}

	sanitized := filenameSanitizer.ReplaceAllString(trimmed, "-")
	sanitized = strings.Trim(sanitized, "-.")
	if sanitized == "" {
		sanitized = fmt.Sprintf("asset-%s", uuid.NewString())
	}

	ext := path.Ext(trimmed)
	if ext == "" {
		ext = defaultExtension(contentType)
	}
	if ext != "" && !strings.HasSuffix(strings.ToLower(sanitized), strings.ToLower(ext)) {
		sanitized = sanitized + ext
	}

	if len(sanitized) > 96 {
		sanitized = sanitized[:96]
	}

	return sanitized
}

func defaultExtension(contentType string) string {
	switch strings.ToLower(contentType) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/heic":
		return ".heic"
	case "video/mp4":
		return ".mp4"
	case "video/quicktime":
		return ".mov"
	case "application/pdf":
		return ".pdf"
	default:
		return ".bin"
	}
}

func sanitizeSegment(value string) string {
	sanitized := segmentSanitizer.ReplaceAllString(strings.TrimSpace(value), "-")
	sanitized = strings.Trim(sanitized, "-_")
	if sanitized == "" {
		return ""
	}
	if len(sanitized) > 48 {
		sanitized = sanitized[:48]
	}
	return sanitized
}

func resolveContentType(value string) string {
	if strings.TrimSpace(value) == "" {
		return "application/octet-stream"
	}
	return value
}

func normalizeChecksum(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", nil
	}

	if strings.Contains(trimmed, ":") {
		parts := strings.SplitN(trimmed, ":", 2)
		trimmed = parts[1]
	}

	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return "", nil
	}

	if hexChecksumRegex.MatchString(trimmed) && len(trimmed)%2 == 0 {
		decoded, err := hex.DecodeString(trimmed)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(decoded), nil
	}

	if _, err := base64.StdEncoding.DecodeString(trimmed); err == nil {
		return trimmed, nil
	}

	return "", fmt.Errorf("unsupported checksum format")
}

// DownloadFile downloads an object from the listing bucket to memory.
func (a *ListingMediaStorageAdapter) DownloadFile(ctx context.Context, key string) ([]byte, error) {
	ctx = utils.ContextWithLogger(ctx)
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "ListingMediaStorage.DownloadFile")
	if err != nil {
		return nil, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if err := a.ensureClients(); err != nil {
		utils.SetSpanError(ctx, err)
		return nil, err
	}

	buf := manager.NewWriteAtBuffer([]byte{})
	_, err = a.base.downloader.Download(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		utils.SetSpanError(ctx, err)
		return nil, derrors.Infra("failed to download file", err)
	}

	return buf.Bytes(), nil
}

// UploadFile uploads content from memory to the listing bucket.
func (a *ListingMediaStorageAdapter) UploadFile(ctx context.Context, key string, content []byte, contentType string) error {
	ctx = utils.ContextWithLogger(ctx)
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "ListingMediaStorage.UploadFile")
	if err != nil {
		return derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if err := a.ensureClients(); err != nil {
		utils.SetSpanError(ctx, err)
		return err
	}

	_, err = a.base.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		utils.SetSpanError(ctx, err)
		return derrors.Infra("failed to upload file", err)
	}

	return nil
}

const (
	httpMethodPut = "PUT"
	httpMethodGet = "GET"
)

var _ storageport.ListingMediaStoragePort = (*ListingMediaStorageAdapter)(nil)
