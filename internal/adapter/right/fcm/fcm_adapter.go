package fcmadapter

import (
	"context"
	"log/slog"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FCMAdapter struct {
	client *messaging.Client
}

func NewFCMAdapter(ctx context.Context, env *globalmodel.Environment) (fcm *FCMAdapter, err error) {
	app, err := firebase.NewApp(ctx,
		nil, option.WithCredentialsFile(env.FCM.CredentialsFile))
	if err != nil {
		slog.Error("failed to create fcm app", "error", err)
		err = status.Error(codes.Internal, "internal error")
		return
	}
	fcm = &FCMAdapter{}
	client, err := app.Messaging(ctx)
	if err != nil {
		slog.Error("failed to create fcm client", "error", err)
		err = status.Error(codes.Internal, "internal error")
		return
	}

	fcm.client = client
	return
}
