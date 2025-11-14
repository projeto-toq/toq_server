package photosessionservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ReservePhotoSession blocks a slot window for a listing owner with configurable approval mode.
//
// This method orchestrates the complete photo session reservation flow:
//  1. Validates listing ownership and eligibility
//  2. Checks photographer availability (no conflicts)
//  3. Creates agenda entry (blocks the slot)
//  4. Creates booking with status determined by config:
//     - If require_photographer_approval=false: ACCEPTED (automatic)
//     - If require_photographer_approval=true: PENDING_APPROVAL (manual)
//  5. Updates listing status accordingly:
//     - Automatic: StatusPhotosScheduled
//     - Manual: StatusPendingPhotoConfirmation
//  6. Sends FCM notification (automatic mode only)
//
// The operation is transactional: if any step fails, all changes are rolled back.
//
// Configuration:
//   - Approval mode controlled by env.yaml: photo_session.require_photographer_approval
//   - Default: false (automatic approval)
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - input: ReserveSessionInput with userID, listingID, slotID
//
// Returns:
//   - output: ReserveSessionOutput with photoSessionID, slotID, timestamps, photographerID
//   - err: Domain error with appropriate HTTP status code:
//   - 401 (Auth) if user not authorized
//   - 404 (NotFound) if listing not found
//   - 409 (Conflict) if slot unavailable or listing not eligible
//   - 422 (Validation) if input invalid
//   - 500 (Infra) for infrastructure failures
//
// Side Effects:
//   - Creates photographer_agenda_entries record (blocks slot)
//   - Creates photographer_photo_session_bookings record
//   - Updates listings.status
//   - Sends FCM push notification (automatic mode only)
//   - Logs audit entry with approval mode
func (s *photoSessionService) ReservePhotoSession(ctx context.Context, input ReserveSessionInput) (ReserveSessionOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ReserveSessionOutput{}, derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.UserID <= 0 {
		return ReserveSessionOutput{}, derrors.Auth("unauthorized")
	}
	if input.ListingIdentityID <= 0 {
		return ReserveSessionOutput{}, derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "greater_than_zero"})
	}
	if input.SlotID == 0 {
		return ReserveSessionOutput{}, derrors.Validation("slotId must be greater than zero", map[string]any{"slotId": "greater_than_zero"})
	}

	photographerID, slotStartUTC := decodeSlotID(input.SlotID)
	if photographerID == 0 {
		return ReserveSessionOutput{}, derrors.Validation("slotId is invalid", map[string]any{"slotId": "invalid"})
	}

	loc, tzErr := resolveLocation("")
	if tzErr != nil {
		return ReserveSessionOutput{}, tzErr
	}

	slotDuration := time.Duration(s.cfg.SlotDurationMinutes) * time.Minute
	if slotDuration <= 0 {
		slotDuration = 4 * time.Hour
	}

	slotStart := slotStartUTC.In(loc)
	slotEnd := slotStart.Add(slotDuration)
	if !slotEnd.After(slotStart) {
		return ReserveSessionOutput{}, derrors.Validation("slot duration must be positive", map[string]any{"slot": "invalid_duration"})
	}
	if slotEnd.Before(s.now().In(loc)) {
		return ReserveSessionOutput{}, derrors.ErrSlotUnavailable
	}

	// Determine approval mode from configuration
	requireApproval := s.cfg.RequirePhotographerApproval
	approvalMode := "automatic"
	if requireApproval {
		approvalMode = "manual"
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("photo_session.reserve.tx_start_error", "err", txErr)
		return ReserveSessionOutput{}, derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("photo_session.reserve.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ReserveSessionOutput{}, utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.reserve.get_listing_error", "listing_identity_id", input.ListingIdentityID, "err", err)
		return ReserveSessionOutput{}, derrors.Infra("failed to load listing", err)
	}

	if listing.Deleted() {
		return ReserveSessionOutput{}, utils.BadRequest("listing is not available")
	}

	if listing.UserID() != input.UserID {
		return ReserveSessionOutput{}, derrors.Auth("listing does not belong to user")
	}

	if !listingAllowsPhotoSession(listing.Status()) {
		return ReserveSessionOutput{}, derrors.ErrListingNotEligible
	}

	conflicts, err := s.repo.FindBlockingEntries(ctx, tx, photographerID, slotStart.UTC(), slotEnd.UTC())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.reserve.find_blocking_error", "photographer_id", photographerID, "err", err)
		return ReserveSessionOutput{}, derrors.Infra("failed to verify photographer agenda", err)
	}
	if len(conflicts) > 0 {
		return ReserveSessionOutput{}, derrors.ErrSlotUnavailable
	}

	agendaEntry := photosessionmodel.NewAgendaEntry()
	agendaEntry.SetPhotographerUserID(photographerID)
	agendaEntry.SetEntryType(photosessionmodel.AgendaEntryTypePhotoSession)
	agendaEntry.SetSource(photosessionmodel.AgendaEntrySourceBooking)
	agendaEntry.SetStartsAt(slotStart.UTC())
	agendaEntry.SetEndsAt(slotEnd.UTC())
	agendaEntry.SetBlocking(true)
	agendaEntry.SetTimezone(loc.String())
	if input.ListingIdentityID > 0 {
		agendaEntry.SetSourceID(uint64(input.ListingIdentityID))
	}

	entryIDs, err := s.repo.CreateEntries(ctx, tx, []photosessionmodel.AgendaEntryInterface{agendaEntry})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.reserve.create_entry_error", "photographer_id", photographerID, "err", err)
		return ReserveSessionOutput{}, derrors.Infra("failed to create agenda entry", err)
	}
	if len(entryIDs) == 0 {
		return ReserveSessionOutput{}, derrors.Infra("failed to create agenda entry", fmt.Errorf("no entry id returned"))
	}
	entryID := entryIDs[0]

	// Determine booking status based on approval mode configuration
	var bookingStatus photosessionmodel.BookingStatus
	if requireApproval {
		// Manual mode: photographer must approve
		bookingStatus = photosessionmodel.BookingStatusPendingApproval
	} else {
		// Automatic mode: pre-approved
		bookingStatus = photosessionmodel.BookingStatusAccepted
	}

	// Create booking with appropriate status
	booking := photosessionmodel.NewPhotoSessionBooking()
	booking.SetAgendaEntryID(entryID)
	booking.SetPhotographerUserID(photographerID)
	booking.SetListingIdentityID(input.ListingIdentityID)
	booking.SetStartsAt(slotStart.UTC())
	booking.SetEndsAt(slotEnd.UTC())
	booking.SetStatus(bookingStatus)

	bookingID, err := s.repo.CreateBooking(ctx, tx, booking)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.reserve.create_booking_error", "agenda_entry_id", entryID, "err", err)
		return ReserveSessionOutput{}, derrors.Infra("failed to create booking", err)
	}

	// Determine listing status based on approval mode
	var targetListingStatus listingmodel.ListingStatus
	if requireApproval {
		// Manual mode: await photographer approval
		targetListingStatus = listingmodel.StatusPendingPhotoConfirmation
	} else {
		// Automatic mode: directly scheduled
		targetListingStatus = listingmodel.StatusPhotosScheduled
	}

	// Update listing status
	if updateErr := s.listingRepo.UpdateListingStatus(ctx, tx, listing.ID(), targetListingStatus, listing.Status()); updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return ReserveSessionOutput{}, derrors.ErrListingNotEligible
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("photo_session.reserve.update_listing_status_error", "listing_id", listing.ID(), "err", updateErr)
		return ReserveSessionOutput{}, derrors.Infra("failed to update listing status", updateErr)
	}

	// Commit transaction before sending notifications (async operations should not hold transaction)
	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("photo_session.reserve.commit_error", "listing_id", listing.ID(), "err", commitErr)
		return ReserveSessionOutput{}, derrors.Infra("failed to commit reservation", commitErr)
	}
	committed = true

	// Send FCM notification only in automatic mode
	if !requireApproval {
		notificationTitle := "Sessão de Fotos Confirmada"
		notificationBody := "Sua sessão de fotos foi agendada automaticamente e está confirmada!"
		go s.sendOwnerNotifications(context.Background(), input.UserID, notificationTitle, notificationBody, listing.ID(), bookingID)
	}

	// Audit log with approval mode for future analysis
	logger.Info("photo_session.reserve.success",
		"listing_id", listing.ID(),
		"booking_id", bookingID,
		"photographer_id", photographerID,
		"slot_start", slotStart,
		"approval_mode", approvalMode,
		"booking_status", bookingStatus,
		"listing_status", targetListingStatus.String())

	return ReserveSessionOutput{
		PhotoSessionID:    bookingID,
		SlotID:            input.SlotID,
		SlotStart:         slotStart,
		SlotEnd:           slotEnd,
		PhotographerID:    photographerID,
		ListingIdentityID: input.ListingIdentityID,
	}, nil
}

// sendOwnerNotifications sends FCM push notifications to all opted-in devices of the listing owner.
// This function runs asynchronously and logs any errors without propagating them.
//
// Used after automatic photo session approval to notify the owner that the session is confirmed.
//
// Parameters:
//   - ctx: Context for logging (should be background context when called from goroutine)
//   - userID: ID of the listing owner who will receive notifications
//   - title: Push notification title
//   - body: Push notification body text
//   - listingID: ID of the listing associated with the photo session
//   - sessionID: ID of the photo session booking
func (s *photoSessionService) sendOwnerNotifications(ctx context.Context, userID int64, title, body string, listingID int64, sessionID uint64) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Fetch all device tokens for user (only opted-in devices)
	tokens, err := s.globalService.ListDeviceTokensByUserIDIfOptedIn(ctx, userID)
	if err != nil {
		logger.Error("photo_session.notification.list_tokens_error",
			"user_id", userID,
			"listing_id", listingID,
			"err", err)
		return
	}

	if len(tokens) == 0 {
		logger.Info("photo_session.notification.no_tokens",
			"user_id", userID,
			"listing_id", listingID)
		return
	}

	// Get unified notification service
	notifier := s.globalService.GetUnifiedNotificationService()
	if notifier == nil {
		logger.Error("photo_session.notification.service_unavailable",
			"user_id", userID)
		return
	}

	// Send notification to each token (supports multiple devices per user)
	sentCount := 0
	for _, token := range tokens {
		req := globalservice.NotificationRequest{
			Type:    globalservice.NotificationTypeFCM,
			Token:   token,
			Subject: title,
			Body:    body,
		}

		// SendNotification is async by default, but we're already in a goroutine
		if notifErr := notifier.SendNotification(ctx, req); notifErr != nil {
			logger.Warn("photo_session.notification.send_error",
				"user_id", userID,
				"token", token[:min(len(token), 20)]+"...", // Log only token prefix for security
				"err", notifErr)
		} else {
			sentCount++
		}
	}

	logger.Info("photo_session.notification.completed",
		"user_id", userID,
		"listing_id", listingID,
		"session_id", sessionID,
		"tokens_found", len(tokens),
		"notifications_sent", sentCount,
		"title", title)
}
