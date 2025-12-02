package mediaprocessingcallbackport

import "context"

// CallbackPortInterface defines the minimal contract required to validate callbacks from the media pipeline.
type CallbackPortInterface interface {
	ValidateSignature(ctx context.Context, providedSignature string, payload []byte) error
}
