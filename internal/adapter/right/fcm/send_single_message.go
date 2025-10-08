package fcmadapter

import (
	"context"

	"firebase.google.com/go/messaging"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (f *FCMAdapter) SendSingleMessage(ctx context.Context, message globalmodel.Notification) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	iMessage := &messaging.Message{
		Notification: &messaging.Notification{
			Title:    message.Title,
			Body:     message.Body,
			ImageURL: message.Icon,
		},
		Token: string(message.DeviceToken),
	}
	response, err := f.client.Send(ctx, iMessage)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("fcm.send_single.error", "error", err)
		return err
	}
	logger.Info("fcm.send_single.success", "response", response)
	return
}
