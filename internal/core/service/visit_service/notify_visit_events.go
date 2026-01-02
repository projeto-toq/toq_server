package visitservice

import (
	"context"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/projeto-toq/toq_server/internal/core/templates"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *visitService) notifyVisitRequested(ctx context.Context, visit listingmodel.VisitInterface) {
	payload, err := templates.RenderVisitOwnerRequest(templates.VisitTemplateData{
		VisitID:           visit.ID(),
		ListingIdentityID: visit.ListingIdentityID(),
		ScheduledStart:    visit.ScheduledStart(),
		ScheduledEnd:      visit.ScheduledEnd(),
		Status:            string(visit.Status()),
	})
	if err != nil {
		utils.LoggerFromContext(ctx).Warn("visit.notify.render_owner_request_error", "visit_id", visit.ID(), "err", err)
		return
	}
	s.dispatchVisitNotification(ctx, visit.OwnerUserID(), payload)
}

func (s *visitService) notifyVisitStatusOwner(ctx context.Context, visit listingmodel.VisitInterface) {
	payload, err := templates.RenderVisitOwnerStatus(templates.VisitTemplateData{
		VisitID:           visit.ID(),
		ListingIdentityID: visit.ListingIdentityID(),
		ScheduledStart:    visit.ScheduledStart(),
		ScheduledEnd:      visit.ScheduledEnd(),
		Status:            string(visit.Status()),
	})
	if err != nil {
		utils.LoggerFromContext(ctx).Warn("visit.notify.render_owner_status_error", "visit_id", visit.ID(), "err", err)
		return
	}
	s.dispatchVisitNotification(ctx, visit.OwnerUserID(), payload)
}

func (s *visitService) notifyVisitStatusRealtor(ctx context.Context, visit listingmodel.VisitInterface) {
	payload, err := templates.RenderVisitRealtorStatus(templates.VisitTemplateData{
		VisitID:           visit.ID(),
		ListingIdentityID: visit.ListingIdentityID(),
		ScheduledStart:    visit.ScheduledStart(),
		ScheduledEnd:      visit.ScheduledEnd(),
		Status:            string(visit.Status()),
	})
	if err != nil {
		utils.LoggerFromContext(ctx).Warn("visit.notify.render_realtor_status_error", "visit_id", visit.ID(), "err", err)
		return
	}
	s.dispatchVisitNotification(ctx, visit.RequesterUserID(), payload)
}

func (s *visitService) dispatchVisitNotification(ctx context.Context, userID int64, payload templates.VisitPayload) {
	if userID == 0 {
		return
	}
	notifier := s.globalService.GetUnifiedNotificationService()
	if notifier == nil {
		utils.LoggerFromContext(ctx).Warn("visit.notify.notifier_unavailable", "user_id", userID)
		return
	}

	tokens, err := s.globalService.ListDeviceTokensByUserIDIfOptedIn(ctx, userID)
	if err != nil {
		utils.LoggerFromContext(ctx).Warn("visit.notify.list_tokens_error", "user_id", userID, "err", err)
		return
	}
	if len(tokens) == 0 {
		return
	}

	for _, token := range tokens {
		if token == "" {
			continue
		}
		req := globalservice.NotificationRequest{
			Type:    globalservice.NotificationTypeFCM,
			Subject: payload.Title,
			Body:    payload.Body,
			Token:   token,
			Data:    cloneVisitData(payload.Data),
		}
		if err := notifier.SendNotification(ctx, req); err != nil {
			utils.LoggerFromContext(ctx).Warn("visit.notify.enqueue_error", "user_id", userID, "token", token, "err", err)
		}
	}
}

func cloneVisitData(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	out := make(map[string]string, len(input))
	for k, v := range input {
		out[k] = v
	}
	return out
}
