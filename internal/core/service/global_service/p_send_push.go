package globalservice

import (
	"context"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func sendPush(ctx context.Context, gs *globalService, notitificaton globalmodel.Notification) (err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	return gs.firebaseCloudMessage.SendSingleMessage(ctx, notitificaton)
}
