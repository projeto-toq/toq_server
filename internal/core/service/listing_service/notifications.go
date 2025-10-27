package listingservices

import (
	"context"
	"fmt"
	"time"

	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) sendPhotographerReservationSMS(ctx context.Context, phone string, start, end time.Time, listingCode uint32) {
	if phone == "" {
		return
	}

	notifier := ls.gsi.GetUnifiedNotificationService()
	if notifier == nil {
		utils.LoggerFromContext(ctx).Warn("listing.notifications.sms_service_unavailable")
		return
	}

	startLocal := start.In(time.Local)
	endLocal := end.In(time.Local)
	startFormatted := startLocal.Format("02/01 15:04")
	endFormatted := endLocal.Format("15:04")
	body := fmt.Sprintf("Nova sessão de fotos reservada para o anúncio %d em %s-%s. Acesse o app TOQ para aceitar ou recusar.", listingCode, startFormatted, endFormatted)

	req := globalservice.NotificationRequest{
		Type:    globalservice.NotificationTypeSMS,
		To:      phone,
		Subject: "Sessão de fotos reservada",
		Body:    body,
	}

	if err := notifier.SendNotification(ctx, req); err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("listing.notifications.sms_send_error", "err", err, "phone", phone)
	}
}

func (ls *listingService) sendPhotographerCancellationSMS(ctx context.Context, phone string, start time.Time, listingCode uint32) {
	if phone == "" {
		return
	}

	notifier := ls.gsi.GetUnifiedNotificationService()
	if notifier == nil {
		utils.LoggerFromContext(ctx).Warn("listing.notifications.sms_service_unavailable")
		return
	}

	startFormatted := start.In(time.Local).Format("02/01 15:04")
	body := fmt.Sprintf("A sessão de fotos do anúncio %d agendada para %s foi cancelada pelo proprietário.", listingCode, startFormatted)

	req := globalservice.NotificationRequest{
		Type:    globalservice.NotificationTypeSMS,
		To:      phone,
		Subject: "Sessão de fotos cancelada",
		Body:    body,
	}

	if err := notifier.SendNotification(ctx, req); err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("listing.notifications.sms_send_error", "err", err, "phone", phone)
	}
}
