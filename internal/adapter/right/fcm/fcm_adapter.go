package fcmadapter

import (
	"context"
	"log/slog"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"google.golang.org/api/option"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

type FCMAdapter struct {
	client *messaging.Client
}

func NewFCMAdapter(ctx context.Context, env *globalmodel.Environment) (fcm *FCMAdapter, err error) {
	app, err := firebase.NewApp(ctx,
		nil, option.WithCredentialsFile(env.FCM.CredentialsFile))
	if err != nil {
		slog.Error("failed to create fcm app", "error", err)
		err = utils.ErrInternalServer
		return
	}
	fcm = &FCMAdapter{}
	client, err := app.Messaging(ctx)
	if err != nil {
		slog.Error("failed to create fcm client", "error", err)
		err = utils.ErrInternalServer
		return
	}

	fcm.client = client
	return
}
