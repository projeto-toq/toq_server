package mediaprocessingservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/projeto-toq/toq_server/internal/core/templates"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *mediaProcessingService) notifyOwnerMediaReady(ctx context.Context, listing listingmodel.ListingVersionInterface) error {
	notifier := s.globalService.GetUnifiedNotificationService()
	if notifier == nil {
		return derrors.Infra("notification service unavailable", nil)
	}

	ownerID := listing.UserID()
	if ownerID == 0 {
		return derrors.Infra("listing owner undefined", nil)
	}

	tokens, err := s.globalService.ListDeviceTokensByUserIDIfOptedIn(ctx, ownerID)
	if err != nil {
		return derrors.Infra("failed to list owner device tokens", err)
	}
	if len(tokens) == 0 {
		utils.LoggerFromContext(ctx).Warn("service.media.complete.owner_notification.no_tokens",
			"listing_identity_id", listing.ListingIdentityID(),
			"owner_id", ownerID)
		return nil
	}

	listingTitle := listing.Title()
	if strings.TrimSpace(listingTitle) == "" {
		listingTitle = fmt.Sprintf("Anuncio %d", listing.ListingIdentityID())
	}

	rendered, err := templates.RenderMediaOwnerApprovalTemplate(templates.MediaOwnerApprovalTemplateData{
		ListingTitle:   listingTitle,
		ListingID:      listing.ListingIdentityID(),
		ListingVersion: listing.Version(),
	})
	if err != nil {
		return derrors.Infra("failed to render owner notification template", err)
	}

	for _, token := range tokens {
		if token == "" {
			continue
		}
		req := globalservice.NotificationRequest{
			Type:    globalservice.NotificationTypeFCM,
			Subject: rendered.Title,
			Body:    rendered.Body,
			Token:   token,
			Data:    cloneData(rendered.Data),
		}
		if err := notifier.SendNotification(ctx, req); err != nil {
			return derrors.Infra("failed to enqueue owner notification", err)
		}
	}

	utils.LoggerFromContext(ctx).Info("service.media.complete.owner_notification_enqueued",
		"listing_identity_id", listing.ListingIdentityID(),
		"listing_version", listing.Version(),
		"tokens_count", len(tokens))

	return nil
}

func cloneData(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	output := make(map[string]string, len(input))
	for k, v := range input {
		output[k] = v
	}
	return output
}
