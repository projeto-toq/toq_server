package mediaprocessingservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	mediaprocessingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/mediaprocessingrepository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	maxBatchReferenceLength = 120
)

type validatedUploadFile struct {
	clientID    string
	assetType   mediaprocessingmodel.MediaAssetType
	orientation mediaprocessingmodel.MediaAssetOrientation
	filename    string
	contentType string
	bytes       int64
	checksum    string
	title       string
	sequence    uint8
	metadata    map[string]string
}

// CreateUploadBatch validates a manifest and returns signed URLs for raw uploads.
func (s *mediaProcessingService) CreateUploadBatch(ctx context.Context, input CreateUploadBatchInput) (CreateUploadBatchOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return CreateUploadBatchOutput{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ListingIdentityID == 0 {
		return CreateUploadBatchOutput{}, derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "required"})
	}

	input.BatchReference = strings.TrimSpace(input.BatchReference)
	if input.BatchReference == "" {
		return CreateUploadBatchOutput{}, derrors.Validation("batchReference is required", map[string]any{"batchReference": "required"})
	}
	if len(input.BatchReference) > maxBatchReferenceLength {
		return CreateUploadBatchOutput{}, derrors.Validation("batchReference is too long", map[string]any{"batchReference": fmt.Sprintf("max %d characters", maxBatchReferenceLength)})
	}

	requestedBy, err := s.resolveRequestedBy(ctx, input.RequestedBy)
	if err != nil {
		return CreateUploadBatchOutput{}, err
	}
	input.RequestedBy = requestedBy

	validatedFiles, totalBytes, err := s.validateUploadManifest(input)
	if err != nil {
		return CreateUploadBatchOutput{}, err
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("service.media.create_batch.tx_start_error", "err", txErr, "listing_identity_id", input.ListingIdentityID)
		return CreateUploadBatchOutput{}, derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("service.media.create_batch.tx_rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
			}
		}
	}()

	listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, input.ListingIdentityID.Int64())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CreateUploadBatchOutput{}, derrors.NotFound("listing not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.create_batch.get_listing_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return CreateUploadBatchOutput{}, derrors.Infra("failed to load listing", err)
	}

	if listing.Status() != listingmodel.StatusPendingPhotoProcessing {
		return CreateUploadBatchOutput{}, derrors.Conflict("listing is not awaiting media uploads")
	}

	if err := s.ensureNoOpenBatch(ctx, tx, input.ListingIdentityID); err != nil {
		return CreateUploadBatchOutput{}, err
	}

	batch := mediaprocessingmodel.NewMediaBatch(input.ListingIdentityID.Uint64(), input.BatchReference, input.RequestedBy)

	batchID, err := s.repo.CreateBatch(ctx, tx, batch)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.create_batch.create_batch_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return CreateUploadBatchOutput{}, derrors.Infra("failed to persist batch", err)
	}

	assets := make([]mediaprocessingmodel.MediaAsset, 0, len(validatedFiles))
	instructions := make([]UploadInstruction, 0, len(validatedFiles))
	var uploadTTLSeconds int

	for _, file := range validatedFiles {
		asset := mediaprocessingmodel.NewMediaAsset(batchID, input.ListingIdentityID.Uint64(), file.assetType, file.sequence)
		asset.SetFilename(file.filename)
		asset.SetContentType(file.contentType)
		if file.orientation != "" {
			asset.UpdateOrientation(file.orientation)
		}
		for key, value := range file.metadata {
			asset.SetMetadata(key, value)
		}
		asset.SetMetadata("client_id", file.clientID)
		if file.title != "" {
			asset.SetMetadata("title", file.title)
		}
		asset.SetMetadata("batch_reference", input.BatchReference)
		asset.SetMetadata("requested_by", fmt.Sprintf("%d", input.RequestedBy))

		signedURL, genErr := s.storage.GenerateRawUploadURL(ctx, input.ListingIdentityID.Uint64(), asset)
		if genErr != nil {
			utils.SetSpanError(ctx, genErr)
			logger.Error("service.media.create_batch.generate_url_error", "err", genErr, "listing_identity_id", input.ListingIdentityID, "client_id", file.clientID)
			return CreateUploadBatchOutput{}, genErr
		}

		asset.UpdateRawObject(signedURL.ObjectKey, file.checksum, file.bytes)
		assets = append(assets, asset)

		ttl := int(signedURL.ExpiresIn.Seconds())
		if ttl > uploadTTLSeconds {
			uploadTTLSeconds = ttl
		}

		instructions = append(instructions, UploadInstruction{
			ClientID:  file.clientID,
			UploadURL: signedURL.URL,
			Method:    signedURL.Method,
			Headers:   signedURL.Headers,
			ObjectKey: signedURL.ObjectKey,
			Sequence:  file.sequence,
			Title:     file.title,
		})
	}

	if err := s.repo.UpsertAssets(ctx, tx, assets); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.create_batch.upsert_assets_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return CreateUploadBatchOutput{}, derrors.Infra("failed to persist assets", err)
	}

	metadata := mediaprocessingmodel.BatchStatusMetadata{
		Message:   "batch_created",
		Details:   map[string]string{"files": fmt.Sprintf("%d", len(assets)), "bytes": fmt.Sprintf("%d", totalBytes)},
		UpdatedBy: input.RequestedBy,
		UpdatedAt: s.nowUTC(),
	}

	if err := s.repo.UpdateBatchStatus(ctx, tx, batchID, mediaprocessingmodel.BatchStatusPendingUpload, metadata); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.create_batch.update_status_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return CreateUploadBatchOutput{}, derrors.Infra("failed to update batch status", err)
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.create_batch.tx_commit_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return CreateUploadBatchOutput{}, derrors.Infra("failed to commit batch creation", err)
	}
	committed = true

	logger.Info("service.media.create_batch.success",
		"listing_identity_id", input.ListingIdentityID,
		"batch_id", batchID,
		"files", len(assets),
		"bytes", totalBytes,
	)

	return CreateUploadBatchOutput{
		ListingIdentityID:   input.ListingIdentityID,
		BatchID:             batchID,
		UploadURLTTLSeconds: uploadTTLSeconds,
		Files:               instructions,
	}, nil
}

func (s *mediaProcessingService) validateUploadManifest(input CreateUploadBatchInput) ([]validatedUploadFile, int64, error) {
	if len(input.Files) == 0 {
		return nil, 0, derrors.Validation("files are required", map[string]any{"files": "min=1"})
	}

	totalBytes := int64(0)
	validated := make([]validatedUploadFile, 0, len(input.Files))
	sequenceSet := make(map[uint8]struct{}, len(input.Files))
	clientSet := make(map[string]struct{}, len(input.Files))

	for idx, file := range input.Files {
		clientID := strings.TrimSpace(file.ClientID)
		if clientID == "" {
			return nil, 0, derrors.Validation("clientId is required", map[string]any{fmt.Sprintf("files[%d].clientId", idx): "required"})
		}
		if _, exists := clientSet[clientID]; exists {
			return nil, 0, derrors.Validation("clientId must be unique", map[string]any{"clientId": clientID})
		}
		clientSet[clientID] = struct{}{}

		if file.Sequence == 0 {
			return nil, 0, derrors.Validation("sequence must be greater than zero", map[string]any{"clientId": clientID})
		}
		if _, exists := sequenceSet[file.Sequence]; exists {
			return nil, 0, derrors.Validation("sequence must be unique", map[string]any{"sequence": file.Sequence})
		}
		sequenceSet[file.Sequence] = struct{}{}

		if err := s.ensureContentTypeAllowed(file.ContentType); err != nil {
			return nil, 0, err
		}
		if file.Bytes <= 0 {
			return nil, 0, derrors.Validation("bytes must be greater than zero", map[string]any{"clientId": clientID})
		}
		if s.cfg.MaxFileBytes > 0 && file.Bytes > s.cfg.MaxFileBytes {
			return nil, 0, derrors.Validation("file exceeds allowed size", map[string]any{"clientId": clientID, "maxBytes": s.cfg.MaxFileBytes})
		}
		if strings.TrimSpace(file.Checksum) == "" {
			return nil, 0, derrors.Validation("checksum is required", map[string]any{"clientId": clientID})
		}

		assetType := strings.ToUpper(strings.TrimSpace(string(file.AssetType)))
		if !isSupportedAssetType(mediaprocessingmodel.MediaAssetType(assetType)) {
			return nil, 0, derrors.Validation("unsupported asset type", map[string]any{"clientId": clientID, "assetType": file.AssetType})
		}

		orientation := mediaprocessingmodel.MediaAssetOrientation(strings.ToUpper(strings.TrimSpace(string(file.Orientation))))
		if requiresOrientation(mediaprocessingmodel.MediaAssetType(assetType)) && !isValidOrientation(orientation) {
			return nil, 0, derrors.Validation("orientation is required for this asset type", map[string]any{"clientId": clientID})
		}

		metadata := cloneStringMap(file.Metadata)
		validated = append(validated, validatedUploadFile{
			clientID:    clientID,
			assetType:   mediaprocessingmodel.MediaAssetType(assetType),
			orientation: orientation,
			filename:    strings.TrimSpace(file.Filename),
			contentType: strings.TrimSpace(file.ContentType),
			bytes:       file.Bytes,
			checksum:    strings.TrimSpace(file.Checksum),
			title:       strings.TrimSpace(file.Title),
			sequence:    file.Sequence,
			metadata:    metadata,
		})

		totalBytes += file.Bytes
	}

	if err := s.ensureBatchingLimits(len(validated), totalBytes); err != nil {
		return nil, 0, err
	}

	return validated, totalBytes, nil
}

func (s *mediaProcessingService) ensureNoOpenBatch(ctx context.Context, tx *sql.Tx, listingIdentityID listingmodel.ListingIdentityID) error {
	filter := mediaprocessingrepository.BatchQueryFilter{
		ListingID: listingIdentityID.Uint64(),
		Statuses: []mediaprocessingmodel.BatchStatus{
			mediaprocessingmodel.BatchStatusPendingUpload,
			mediaprocessingmodel.BatchStatusReceived,
			mediaprocessingmodel.BatchStatusProcessing,
		},
		Limit: 1,
	}

	batches, err := s.repo.ListBatchesByListing(ctx, tx, filter)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return derrors.Infra("failed to check existing batches", err)
	}
	if len(batches) > 0 {
		return derrors.Conflict("listing already has an active media batch")
	}
	return nil
}

func cloneStringMap(input map[string]string) map[string]string {
	if len(input) == 0 {
		return map[string]string{}
	}
	clone := make(map[string]string, len(input))
	for k, v := range input {
		clone[k] = v
	}
	return clone
}

func isSupportedAssetType(assetType mediaprocessingmodel.MediaAssetType) bool {
	switch assetType {
	case mediaprocessingmodel.MediaAssetTypePhotoVertical,
		mediaprocessingmodel.MediaAssetTypePhotoHorizontal,
		mediaprocessingmodel.MediaAssetTypeVideoVertical,
		mediaprocessingmodel.MediaAssetTypeVideoHorizontal,
		mediaprocessingmodel.MediaAssetTypeThumbnail,
		mediaprocessingmodel.MediaAssetTypeZip,
		mediaprocessingmodel.MediaAssetTypeProjectDoc,
		mediaprocessingmodel.MediaAssetTypeProjectRender:
		return true
	default:
		return false
	}
}

func isValidOrientation(orientation mediaprocessingmodel.MediaAssetOrientation) bool {
	switch orientation {
	case mediaprocessingmodel.MediaAssetOrientationHorizontal, mediaprocessingmodel.MediaAssetOrientationVertical:
		return true
	default:
		return false
	}
}

func requiresOrientation(assetType mediaprocessingmodel.MediaAssetType) bool {
	switch assetType {
	case mediaprocessingmodel.MediaAssetTypePhotoVertical,
		mediaprocessingmodel.MediaAssetTypePhotoHorizontal,
		mediaprocessingmodel.MediaAssetTypeVideoVertical,
		mediaprocessingmodel.MediaAssetTypeVideoHorizontal:
		return true
	default:
		return false
	}
}
