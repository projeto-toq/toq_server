package fcmadapter

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/api/option"
)

type FCMAdapter struct {
	client *messaging.Client
}

func NewFCMAdapter(ctx context.Context, env *globalmodel.Environment) (fcm *FCMAdapter, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	app, err := firebase.NewApp(ctx,
		nil, option.WithCredentialsFile(env.FCM.CredentialsFile))
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("fcm.adapter.create_app_error", "error", err)
		return nil, err
	}
	fcm = &FCMAdapter{}
	client, err := app.Messaging(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("fcm.adapter.create_client_error", "error", err)
		return nil, err
	}

	fcm.client = client
	logger.Info("fcm.adapter.initialized")
	return
}
