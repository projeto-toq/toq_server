package fcmadapter

import (
	"context"
	"log/slog"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FCMAdapter struct {
	client *messaging.Client
}

func NewFCMAdapter(ctx context.Context) (fcm *FCMAdapter, err error) {
	app, err := firebase.NewApp(ctx,
		nil, option.WithCredentialsFile("../configs/fcm_admin.json"))
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
