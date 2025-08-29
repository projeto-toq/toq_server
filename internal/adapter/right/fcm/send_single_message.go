package fcmadapter

import (
	"context"
	"log/slog"

	"firebase.google.com/go/messaging"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (f *FCMAdapter) SendSingleMessage(ctx context.Context, message globalmodel.Notification) (err error) {

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
		slog.Error("failed to send message", "error", err)
		err = utils.ErrInternalServer
		return
	}
	slog.Info("message sent", "response", response)
	return
}
