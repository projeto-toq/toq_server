package workflow

import (
	"context"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// WorkflowPortInterface defines the contract for interacting with workflow orchestration engines (e.g. AWS Step Functions).
type WorkflowPortInterface interface {
	StartMediaFinalization(ctx context.Context, input mediaprocessingmodel.MediaFinalizationInput) (executionARN string, err error)
}
