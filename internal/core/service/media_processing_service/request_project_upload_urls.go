package mediaprocessingservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// RequestProjectUploadURLs validates a project media manifest (docs/renders) and returns signed PUT URLs for raw uploads.
//
// Flow:
//  1. Guard by feature flag and listing ID validation.
//  2. Validate manifest (size/content/checksum) and whitelist PROJECT_DOC/PROJECT_RENDER.
//  3. TX: load listing, enforce OffPlanHouse + StatusPendingPlanLoading.
//  4. Upsert assets (idempotent), attach metadata, generate raw upload URLs, set status PENDING_UPLOAD.
//  5. Commit and return upload instructions (max TTL across URLs).
func (s *mediaProcessingService) RequestProjectUploadURLs(ctx context.Context, input dto.RequestUploadURLsInput) (dto.RequestUploadURLsOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return dto.RequestUploadURLsOutput{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if !s.cfg.AllowOwnerProjectUpload {
		return dto.RequestUploadURLsOutput{}, derrors.Forbidden("project uploads are disabled", nil)
	}

	if input.ListingIdentityID == 0 {
		return dto.RequestUploadURLsOutput{}, derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "required"})
	}

	validatedFiles, err := s.validateUploadManifest(input)
	if err != nil {
		return dto.RequestUploadURLsOutput{}, err
	}

	for _, f := range validatedFiles {
		if f.AssetType != mediaprocessingmodel.MediaAssetTypeProjectDoc && f.AssetType != mediaprocessingmodel.MediaAssetTypeProjectRender {
			return dto.RequestUploadURLsOutput{}, derrors.Validation("only project asset types are allowed", map[string]any{"assetType": f.AssetType})
		}
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("service.media.project_request_upload.tx_start_error", "err", txErr, "listing_identity_id", input.ListingIdentityID)
		return dto.RequestUploadURLsOutput{}, derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("service.media.project_request_upload.tx_rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
			}
		}
	}()

	listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.RequestUploadURLsOutput{}, derrors.NotFound("listing not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.project_request_upload.get_listing_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return dto.RequestUploadURLsOutput{}, derrors.Infra("failed to load listing", err)
	}

	if err := s.ensureProjectFlowAllowed(listing); err != nil {
		return dto.RequestUploadURLsOutput{}, err
	}

	instructions := make([]dto.UploadInstruction, 0, len(validatedFiles))
	var uploadTTLSeconds int

	for _, file := range validatedFiles {
		existingAsset, err := s.repo.GetAssetBySequence(ctx, tx, uint64(input.ListingIdentityID), file.AssetType, file.Sequence)

		var asset mediaprocessingmodel.MediaAsset
		if err == nil {
			asset = existingAsset
		} else if errors.Is(err, sql.ErrNoRows) {
			asset = mediaprocessingmodel.NewMediaAsset(uint64(input.ListingIdentityID), file.AssetType, file.Sequence)
		} else {
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.project_request_upload.get_asset_error", "err", err, "listing_identity_id", input.ListingIdentityID)
			return dto.RequestUploadURLsOutput{}, derrors.Infra("failed to check existing asset", err)
		}

		metaBytes, metaErr := encodeUploadMetadata(file, input.RequestedBy)
		if metaErr != nil {
			utils.SetSpanError(ctx, metaErr)
			logger.Error("service.media.project_request_upload.metadata_error", "err", metaErr, "listing_identity_id", input.ListingIdentityID)
			return dto.RequestUploadURLsOutput{}, metaErr
		}

		asset.SetTitle(file.Title)
		asset.SetMetadata(string(metaBytes))

		signedURL, genErr := s.storage.GenerateRawUploadURL(ctx, uint64(input.ListingIdentityID), asset, file.ContentType, file.Checksum)
		if genErr != nil {
			utils.SetSpanError(ctx, genErr)
			logger.Error("service.media.project_request_upload.generate_url_error", "err", genErr, "listing_identity_id", input.ListingIdentityID)
			return dto.RequestUploadURLsOutput{}, genErr
		}

		asset.SetS3KeyRaw(signedURL.ObjectKey)
		asset.SetStatus(mediaprocessingmodel.MediaAssetStatusPendingUpload)

		if err := s.repo.UpsertAsset(ctx, tx, asset); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.project_request_upload.upsert_asset_error", "err", err, "listing_identity_id", input.ListingIdentityID)
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
		logger.Error("service.media.project_request_upload.tx_commit_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return dto.RequestUploadURLsOutput{}, derrors.Infra("failed to commit upload request", err)
	}
	committed = true

	return dto.RequestUploadURLsOutput{
		ListingIdentityID:   input.ListingIdentityID,
		UploadURLTTLSeconds: uploadTTLSeconds,
		Files:               instructions,
	}, nil
}

func encodeUploadMetadata(file dto.RequestUploadFile, requestedBy uint64) ([]byte, error) {
	metaMap := make(map[string]string)
	for k, v := range file.Metadata {
		metaMap[k] = v
	}
	metaMap["requested_by"] = fmt.Sprintf("%d", requestedBy)
	metaMap["filename"] = file.Filename
	metaMap["content_type"] = file.ContentType
	metaMap["checksum"] = file.Checksum
	metaMap["size_bytes"] = fmt.Sprintf("%d", file.Bytes)

	return json.Marshal(metaMap)
}
