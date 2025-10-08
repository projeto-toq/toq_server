package emailport

import (
	"context"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

type EmailPortInterface interface {
	SendEmail(ctx context.Context, notification globalmodel.Notification) error
}
