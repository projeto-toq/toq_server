package mediaprocessingcallbackport

import "context"

// CallbackPortInterface defines the minimal contract required to validate callbacks from the media pipeline.
type CallbackPortInterface interface {
	ValidateSharedSecret(ctx context.Context, providedSecret string) error
}
