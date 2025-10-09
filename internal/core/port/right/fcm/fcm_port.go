package fcmport

import (
	"context"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

type FCMPortInterface interface {
	SendSingleMessage(ctx context.Context, message globalmodel.Notification) (err error)
	SendMultipleMessages(ctx context.Context, message globalmodel.Notification, deviceTokens []string) error
}
