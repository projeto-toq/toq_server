package visitservice

import (
	"context"
	"time"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	participantPhotoVariant = "small"
	liveStatusOn            = "AO_VIVO"
)

func (s *visitService) decorateParticipantSnapshot(ctx context.Context, snapshot listingmodel.VisitParticipantSnapshot) listingmodel.VisitParticipantSnapshot {
	if snapshot.UserID <= 0 {
		return snapshot
	}
	snapshot.PhotoURL = s.generateParticipantPhotoURL(ctx, snapshot.UserID)
	return snapshot
}

func (s *visitService) generateParticipantPhotoURL(ctx context.Context, userID int64) string {
	if s.userService == nil || userID <= 0 {
		return ""
	}

	impersonatedCtx := utils.SetUserInContext(ctx, usermodel.UserInfos{ID: userID})
	url, err := s.userService.GetPhotoDownloadURL(impersonatedCtx, participantPhotoVariant)
	if err != nil {
		logger := utils.LoggerFromContext(ctx)
		logger.Debug("visit.participant.photo_url_error", "user_id", userID, "err", err)
		return ""
	}
	return url
}

func buildVisitTimeline(visit listingmodel.VisitInterface) VisitTimeline {
	requestedAt := visit.RequestedAt()
	if requestedAt.IsZero() {
		requestedAt = visit.CreatedAt()
	}
	if requestedAt.IsZero() {
		requestedAt = visit.ScheduledStart()
	}

	receivedAt := requestedAt
	var respondedAt *time.Time
	if ts, ok := visit.FirstOwnerActionAt(); ok {
		respondedAt = &ts
	}

	return VisitTimeline{CreatedAt: requestedAt, ReceivedAt: receivedAt, RespondedAt: respondedAt}
}

func computeLiveStatus(visit listingmodel.VisitInterface, now time.Time) string {
	startWindow := visit.ScheduledStart().Add(-2 * time.Hour)
	endWindow := visit.ScheduledEnd().Add(2 * time.Hour)

	if !now.Before(startWindow) && !now.After(endWindow) {
		return liveStatusOn
	}
	return ""
}
