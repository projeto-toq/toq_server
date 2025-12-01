package mediaprocessingservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// RequestUploadURLs validates a manifest and returns signed URLs for raw uploads.
func (s *mediaProcessingService) RequestUploadURLs(ctx context.Context, input dto.RequestUploadURLsInput) (dto.RequestUploadURLsOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return dto.RequestUploadURLsOutput{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ListingIdentityID == 0 {
		return dto.RequestUploadURLsOutput{}, derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "required"})
	}

	validatedFiles, err := s.validateUploadManifest(input)
	if err != nil {
		return dto.RequestUploadURLsOutput{}, err
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("service.media.request_upload.tx_start_error", "err", txErr, "listing_identity_id", input.ListingIdentityID)
		return dto.RequestUploadURLsOutput{}, derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("service.media.request_upload.tx_rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
			}
		}
	}()

	listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.RequestUploadURLsOutput{}, derrors.NotFound("listing not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.request_upload.get_listing_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return dto.RequestUploadURLsOutput{}, derrors.Infra("failed to load listing", err)
	}

	if listing.Status() != listingmodel.StatusPendingPhotoProcessing && listing.Status() != listingmodel.StatusRejectedByOwner {
		return dto.RequestUploadURLsOutput{}, derrors.Conflict("listing is not awaiting media uploads")
	}

	instructions := make([]dto.UploadInstruction, 0, len(validatedFiles))
	var uploadTTLSeconds int

	for _, file := range validatedFiles {
		// IDEMPOTENCY CHECK: Try to find existing asset first to preserve S3 Key
		existingAsset, err := s.repo.GetAssetBySequence(ctx, tx, uint64(input.ListingIdentityID), file.AssetType, file.Sequence)

		var asset mediaprocessingmodel.MediaAsset
		if err == nil {
			// Asset exists: reuse it to keep the same S3 Key (idempotency)
			asset = existingAsset
		} else if errors.Is(err, sql.ErrNoRows) {
			// Asset does not exist: create new
			asset = mediaprocessingmodel.NewMediaAsset(uint64(input.ListingIdentityID), file.AssetType, file.Sequence)
		} else {
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.request_upload.get_asset_error", "err", err, "listing_identity_id", input.ListingIdentityID)
			return dto.RequestUploadURLsOutput{}, derrors.Infra("failed to check existing asset", err)
		}

		asset.SetTitle(file.Title)

		// Prepare metadata
		metaMap := make(map[string]string)
		for k, v := range file.Metadata {
			metaMap[k] = v
		}
		metaMap["requested_by"] = fmt.Sprintf("%d", input.RequestedBy)
		metaMap["filename"] = file.Filename
		metaMap["content_type"] = file.ContentType
		metaMap["checksum"] = file.Checksum
		metaMap["size_bytes"] = fmt.Sprintf("%d", file.Bytes)

		// Set metadata JSON
		metaBytes, _ := json.Marshal(metaMap)
		asset.SetMetadata(string(metaBytes))

		// Generate Signed URL
		// NOTE: If asset.S3KeyRaw is already set (from DB), the adapter will reuse it.
		signedURL, genErr := s.storage.GenerateRawUploadURL(ctx, uint64(input.ListingIdentityID), asset, file.ContentType, file.Checksum)
		if genErr != nil {
			utils.SetSpanError(ctx, genErr)
			logger.Error("service.media.request_upload.generate_url_error", "err", genErr, "listing_identity_id", input.ListingIdentityID)
			return dto.RequestUploadURLsOutput{}, genErr
		}

		asset.SetS3KeyRaw(signedURL.ObjectKey)
		// Reset status to PENDING_UPLOAD in case it was failed/processed, allowing re-upload
		if asset.Status() == mediaprocessingmodel.MediaAssetStatusFailed || asset.Status() == mediaprocessingmodel.MediaAssetStatusProcessed {
			asset.SetStatus(mediaprocessingmodel.MediaAssetStatusPendingUpload)
		}

		// Upsert Asset
		if err := s.repo.UpsertAsset(ctx, tx, asset); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.request_upload.upsert_asset_error", "err", err, "listing_identity_id", input.ListingIdentityID)
			return dto.RequestUploadURLsOutput{}, derrors.Infra("failed to persist asset", err)
		}

		ttl := int(signedURL.ExpiresIn.Seconds())
		if ttl > uploadTTLSeconds {
			uploadTTLSeconds = ttl
		}

		instructions = append(instructions, dto.UploadInstruction{
			AssetType: string(file.AssetType),
			Sequence:  file.Sequence,
			UploadURL: signedURL.URL,
			Method:    signedURL.Method,
			Headers:   signedURL.Headers,
			ObjectKey: signedURL.ObjectKey,
			Title:     file.Title,
		})
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.request_upload.tx_commit_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return dto.RequestUploadURLsOutput{}, derrors.Infra("failed to commit upload request", err)
	}
	committed = true

	return dto.RequestUploadURLsOutput{
		ListingIdentityID:   input.ListingIdentityID,
		UploadURLTTLSeconds: uploadTTLSeconds,
		Files:               instructions,
	}, nil
}

func (s *mediaProcessingService) validateUploadManifest(input dto.RequestUploadURLsInput) ([]dto.RequestUploadFile, error) {
	if len(input.Files) == 0 {
		return nil, derrors.Validation("files are required", map[string]any{"files": "min=1"})
	}

	totalBytes := int64(0)
	validated := make([]dto.RequestUploadFile, 0, len(input.Files))

	// Unique key: AssetType + Sequence
	uniqueSet := make(map[string]struct{}, len(input.Files))

	for idx, file := range input.Files {
		if file.Sequence == 0 {
			return nil, derrors.Validation("sequence must be greater than zero", map[string]any{fmt.Sprintf("files[%d].sequence", idx): "required"})
		}

		key := fmt.Sprintf("%s-%d", file.AssetType, file.Sequence)
		if _, exists := uniqueSet[key]; exists {
			return nil, derrors.Validation("duplicate asset type and sequence", map[string]any{"key": key})
		}
		uniqueSet[key] = struct{}{}

		if err := s.ensureContentTypeAllowed(file.ContentType); err != nil {
			return nil, err
		}
		if file.Bytes <= 0 {
			return nil, derrors.Validation("bytes must be greater than zero", map[string]any{"key": key})
		}
		if s.cfg.MaxFileBytes > 0 && file.Bytes > s.cfg.MaxFileBytes {
			return nil, derrors.Validation("file exceeds allowed size", map[string]any{"key": key, "maxBytes": s.cfg.MaxFileBytes})
		}
		if strings.TrimSpace(file.Checksum) == "" {
			return nil, derrors.Validation("checksum is required", map[string]any{"key": key})
		}

		assetType := strings.ToUpper(strings.TrimSpace(string(file.AssetType)))
		if !isSupportedAssetType(mediaprocessingmodel.MediaAssetType(assetType)) {
			return nil, derrors.Validation("unsupported asset type", map[string]any{"key": key, "assetType": file.AssetType})
		}

		validated = append(validated, file)
		totalBytes += file.Bytes
	}

	if err := s.ensureBatchingLimits(len(validated), totalBytes); err != nil {
		return nil, err
	}

	return validated, nil
}

func isSupportedAssetType(t mediaprocessingmodel.MediaAssetType) bool {
	switch t {
	case mediaprocessingmodel.MediaAssetTypePhotoVertical,
		mediaprocessingmodel.MediaAssetTypePhotoHorizontal,
		mediaprocessingmodel.MediaAssetTypeVideoVertical,
		mediaprocessingmodel.MediaAssetTypeVideoHorizontal,
		mediaprocessingmodel.MediaAssetTypeThumbnail,
		mediaprocessingmodel.MediaAssetTypeZip,
		mediaprocessingmodel.MediaAssetTypeProjectDoc,
		mediaprocessingmodel.MediaAssetTypeProjectRender:
		return true
	}
	return false
}
