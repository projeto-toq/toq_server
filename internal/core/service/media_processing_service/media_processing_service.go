package mediaprocessingservice

import (
	"context"
	"os"
	"strings"
	"time"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	mediaprocessingqueue "github.com/projeto-toq/toq_server/internal/core/port/right/queue/mediaprocessingqueue"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	mediaprocessingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/media_processing_repository"
	storageport "github.com/projeto-toq/toq_server/internal/core/port/right/storage"
	workflowport "github.com/projeto-toq/toq_server/internal/core/port/right/workflow"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
)

// MediaProcessingServiceInterface exposes orchestration helpers for listing media batches.
type MediaProcessingServiceInterface interface {
	RequestUploadURLs(ctx context.Context, input dto.RequestUploadURLsInput) (dto.RequestUploadURLsOutput, error)
	ProcessMedia(ctx context.Context, input dto.ProcessMediaInput) error
	ReconcileStuckJobs(ctx context.Context, timeout time.Duration) error

	// New methods
	GenerateDownloadURLs(ctx context.Context, input dto.GenerateDownloadURLsInput) (dto.GenerateDownloadURLsOutput, error)
	ListMedia(ctx context.Context, input dto.ListMediaInput) (dto.ListMediaOutput, error)

	// Management
	UpdateMedia(ctx context.Context, input dto.UpdateMediaInput) error
	DeleteMedia(ctx context.Context, input dto.DeleteMediaInput) error
	CompleteMedia(ctx context.Context, input dto.CompleteMediaInput) error
	HandleOwnerMediaApproval(ctx context.Context, input dto.ListingMediaApprovalInput) (dto.ListingMediaApprovalOutput, error)

	// Legacy/Internal
	HandleProcessingCallback(ctx context.Context, input dto.HandleProcessingCallbackInput) (dto.HandleProcessingCallbackOutput, error)
}

// Config centralizes tunable parameters leveraged by the service.
type Config struct {
	MaxFilesPerBatch        int
	MaxTotalBytes           int64
	MaxFileBytes            int64
	AllowedContentTypes     []string
	AllowOwnerProjectUpload bool
	RequireAdminReview      bool
}

type mediaProcessingService struct {
	repo                 mediaprocessingrepository.RepositoryInterface
	listingRepo          listingrepository.ListingRepoPortInterface
	globalService        globalservice.GlobalServiceInterface
	storage              storageport.ListingMediaStoragePort
	queue                mediaprocessingqueue.QueuePortInterface
	workflow             workflowport.WorkflowPortInterface
	cfg                  Config
	now                  func() time.Time
	allowedContentLookup map[string]struct{}
}

// NewConfigFromEnvironment translates env.yaml settings into a Config structure.
func NewConfigFromEnvironment(env *globalmodel.Environment) Config {
	cfg := Config{}
	if env == nil {
		return cfg
	}

	cfg.MaxFilesPerBatch = env.MediaProcessing.Limits.MaxFilesPerBatch
	cfg.MaxTotalBytes = env.MediaProcessing.Limits.MaxTotalBytes
	cfg.MaxFileBytes = env.MediaProcessing.Limits.MaxFileBytes
	cfg.AllowedContentTypes = append(cfg.AllowedContentTypes, env.MediaProcessing.Limits.AllowedContentTypes...)
	cfg.AllowOwnerProjectUpload = env.MediaProcessing.Features.AllowOwnerProjectUploads
	cfg.RequireAdminReview = env.MediaProcessing.Features.ListingApprovalAdminReview

	if raw := strings.TrimSpace(os.Getenv("LISTING_APPROVAL_ADMIN_REVIEW")); raw != "" {
		switch strings.ToLower(raw) {
		case "true", "1", "yes":
			cfg.RequireAdminReview = true
		case "false", "0", "no":
			cfg.RequireAdminReview = false
		}
	}

	return cfg
}

// NewMediaProcessingService wires all dependencies required to orchestrate listing media batches.
func NewMediaProcessingService(
	repo mediaprocessingrepository.RepositoryInterface,
	listingRepo listingrepository.ListingRepoPortInterface,
	globalService globalservice.GlobalServiceInterface,
	storage storageport.ListingMediaStoragePort,
	queue mediaprocessingqueue.QueuePortInterface,
	workflow workflowport.WorkflowPortInterface,
	cfg Config,
) (MediaProcessingServiceInterface, error) {
	if repo == nil {
		return nil, derrors.Infra("media processing repository not configured", nil)
	}
	if listingRepo == nil {
		return nil, derrors.Infra("listing repository not configured", nil)
	}
	if globalService == nil {
		return nil, derrors.Infra("global service not configured", nil)
	}
	if storage == nil {
		return nil, derrors.Infra("listing media storage adapter not configured", nil)
	}
	if queue == nil {
		return nil, derrors.Infra("media processing queue adapter not configured", nil)
	}
	if workflow == nil {
		return nil, derrors.Infra("media processing workflow adapter not configured", nil)
	}

	if cfg.MaxFilesPerBatch <= 0 {
		cfg.MaxFilesPerBatch = 60
	}
	if cfg.MaxTotalBytes <= 0 {
		cfg.MaxTotalBytes = 1 << 30 // 1 GiB
	}
	if cfg.MaxFileBytes <= 0 {
		cfg.MaxFileBytes = 256 << 20 // 256 MiB
	}
	if len(cfg.AllowedContentTypes) == 0 {
		cfg.AllowedContentTypes = []string{"image/jpeg", "image/png", "image/heic", "video/mp4", "video/quicktime", "application/pdf"}
	}

	lookup := make(map[string]struct{}, len(cfg.AllowedContentTypes))
	for _, ct := range cfg.AllowedContentTypes {
		normalized := strings.ToLower(strings.TrimSpace(ct))
		if normalized == "" {
			continue
		}
		lookup[normalized] = struct{}{}
	}

	return &mediaProcessingService{
		repo:                 repo,
		listingRepo:          listingRepo,
		globalService:        globalService,
		storage:              storage,
		queue:                queue,
		workflow:             workflow,
		cfg:                  cfg,
		now:                  time.Now,
		allowedContentLookup: lookup,
	}, nil
}

func (s *mediaProcessingService) ensureContentTypeAllowed(contentType string) error {
	if len(s.allowedContentLookup) == 0 {
		return nil
	}
	normalized := strings.ToLower(strings.TrimSpace(contentType))
	if _, ok := s.allowedContentLookup[normalized]; ok {
		return nil
	}
	return derrors.Validation("unsupported content type", map[string]any{"contentType": contentType})
}

func (s *mediaProcessingService) ensureBatchingLimits(count int, totalBytes int64) error {
	if count == 0 {
		return derrors.Validation("at least one file must be provided", map[string]any{"files": "required"})
	}
	if count > s.cfg.MaxFilesPerBatch {
		return derrors.Validation("too many files in batch", map[string]any{"maxFiles": s.cfg.MaxFilesPerBatch})
	}
	if totalBytes <= 0 {
		return derrors.Validation("total size must be positive", map[string]any{"bytes": "positive"})
	}
	if totalBytes > s.cfg.MaxTotalBytes {
		return derrors.Validation("batch exceeds maximum allowed payload", map[string]any{"maxBytes": s.cfg.MaxTotalBytes})
	}
	return nil
}

func (s *mediaProcessingService) nowUTC() time.Time {
	if s.now == nil {
		return time.Now().UTC()
	}
	return s.now().UTC()
}

var _ MediaProcessingServiceInterface = (*mediaProcessingService)(nil)
