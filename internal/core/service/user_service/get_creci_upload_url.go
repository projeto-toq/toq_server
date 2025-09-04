package userservices

import (
	"context"

	storagemodel "github.com/giulio-alfieri/toq_server/internal/core/model/storage_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricCreciUploadURLTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "creci_upload_url_generated_total",
		Help: "Total number of CRECI document upload URLs generated",
	}, []string{"type"})
	metricCreciUploadURLInvalid = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "creci_upload_url_invalid_total",
		Help: "Total number of invalid CRECI upload URL requests",
	}, []string{"reason"})
)

func init() {
	prometheus.MustRegister(metricCreciUploadURLTotal)
	prometheus.MustRegister(metricCreciUploadURLInvalid)
}

// GetCreciUploadURL generates a signed URL to upload CRECI documents (selfie/front/back)
func (us *userService) GetCreciUploadURL(ctx context.Context, documentType, contentType string) (signedURL string, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil {
		return
	}

	if us.cloudStorageService == nil {
		return "", utils.InternalError("Storage service not configured")
	}

	// Validate documentType against domain constants
	validDocTypes := storagemodel.ValidDocumentTypes()
	if !validDocTypes[documentType] {
		metricCreciUploadURLInvalid.WithLabelValues("invalid_type").Inc()
		return "", utils.ValidationError("documentType", "Unsupported document type")
	}

	// Validate content type (only images)
	allowedContentTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	if !allowedContentTypes[contentType] {
		metricCreciUploadURLInvalid.WithLabelValues("invalid_content_type").Inc()
		return "", utils.ValidationError("contentType", "Only image/jpeg or image/png are allowed")
	}

	// Generate signed URL using storage service
	signedURL, err = us.cloudStorageService.GenerateDocumentUploadURL(userID, storagemodel.DocumentType(documentType), contentType)
	if err != nil {
		// optional metric for failures
		if gsMetrics := us.globalService.GetMetrics(); gsMetrics != nil {
			gsMetrics.IncrementErrors("user_service", "creci_upload_url_error")
		}
		return "", err
	}

	// increment custom metric on success
	metricCreciUploadURLTotal.WithLabelValues(documentType).Inc()

	return signedURL, nil
}
