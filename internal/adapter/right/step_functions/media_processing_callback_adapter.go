package stepfunctionscallbackadapter

import (
	"context"
	"crypto/subtle"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	mediaprocessingcallbackport "github.com/projeto-toq/toq_server/internal/core/port/right/functions/mediaprocessingcallback"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// MediaProcessingCallbackAdapter validates callbacks received from Step Functions or the dispatch Lambda.
type MediaProcessingCallbackAdapter struct {
	sharedSecret string
}

// NewMediaProcessingCallbackAdapter builds the adapter using shared secrets configured in env.yaml.
func NewMediaProcessingCallbackAdapter(env *globalmodel.Environment) *MediaProcessingCallbackAdapter {
	if env == nil {
		return &MediaProcessingCallbackAdapter{}
	}
	return &MediaProcessingCallbackAdapter{sharedSecret: env.MediaProcessing.Callback.SharedSecret}
}

// ValidateSharedSecret compares the provided secret with the configured value.
func (a *MediaProcessingCallbackAdapter) ValidateSharedSecret(ctx context.Context, providedSecret string) error {
	ctx = utils.ContextWithLogger(ctx)
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "StepFunctionsCallback.ValidateSharedSecret")
	if err != nil {
		return derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if a.sharedSecret == "" {
		utils.LoggerFromContext(ctx).Warn("adapter.stepfunctions.callback.secret_missing")
		return nil
	}

	if subtle.ConstantTimeCompare([]byte(a.sharedSecret), []byte(providedSecret)) != 1 {
		err := derrors.Forbidden("invalid callback secret", derrors.WithPublicMessage("invalid callback signature"))
		utils.SetSpanError(ctx, err)
		return err
	}

	return nil
}

var _ mediaprocessingcallbackport.CallbackPortInterface = (*MediaProcessingCallbackAdapter)(nil)
