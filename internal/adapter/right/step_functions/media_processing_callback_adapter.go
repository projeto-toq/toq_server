package stepfunctionscallbackadapter

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"

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

// ValidateSignature recalculates the HMAC-SHA256 using the shared secret and compares it against the provided value.
func (a *MediaProcessingCallbackAdapter) ValidateSignature(ctx context.Context, providedSignature string, payload []byte) error {
	ctx = utils.ContextWithLogger(ctx)
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "StepFunctionsCallback.ValidateSignature")
	if err != nil {
		return derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	logger := utils.LoggerFromContext(ctx)

	if a.sharedSecret == "" {
		logger.Warn("adapter.stepfunctions.callback.secret_missing")
		return nil
	}

	normalizedSignature := strings.ToLower(strings.TrimSpace(providedSignature))
	if normalizedSignature == "" {
		err := derrors.Forbidden("missing callback signature", derrors.WithPublicMessage("invalid callback signature"))
		utils.SetSpanError(ctx, err)
		return err
	}

	if len(payload) == 0 {
		err := derrors.Forbidden("empty payload for signature validation", derrors.WithPublicMessage("invalid callback signature"))
		utils.SetSpanError(ctx, err)
		return err
	}

	expectedSignature := computeHexHMAC(payload, []byte(a.sharedSecret))
	if subtle.ConstantTimeCompare([]byte(normalizedSignature), []byte(expectedSignature)) != 1 {
		err := derrors.Forbidden("invalid callback secret", derrors.WithPublicMessage("invalid callback signature"))
		utils.SetSpanError(ctx, err)
		logger.Warn("adapter.stepfunctions.callback.signature_mismatch")
		return err
	}

	return nil
}

func computeHexHMAC(payload, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

var _ mediaprocessingcallbackport.CallbackPortInterface = (*MediaProcessingCallbackAdapter)(nil)
