package photosessionservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateSessionStatus updates the status of a photo session booking.
func (s *photoSessionService) UpdateSessionStatus(ctx context.Context, input UpdateSessionStatusInput) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.SessionID == 0 {
		return derrors.Validation("sessionId must be greater than zero", map[string]any{"sessionId": "greater_than_zero"})
	}
	if input.PhotographerID == 0 {
		return derrors.Auth("unauthorized")
	}

	statusStr := strings.ToUpper(strings.TrimSpace(input.Status))
	if statusStr == "" {
		return derrors.Validation("status is required", map[string]any{"status": "required"})
	}
	if statusStr != string(photosessionmodel.BookingStatusAccepted) && statusStr != string(photosessionmodel.BookingStatusRejected) {
		return derrors.BadRequest("status must be ACCEPTED or REJECTED")
	}
	status := photosessionmodel.BookingStatus(statusStr)

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.tx_start_error", "err", err)
		return derrors.Infra("failed to start transaction", err)
	}

	committed := false
	defer func() {
		if !committed {
			if rollbackErr := s.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
				utils.SetSpanError(ctx, rollbackErr)
				logger.Error("photo_session.update_status.tx_rollback_error", "err", rollbackErr)
			}
		}
	}()

	booking, err := s.repo.GetBookingByIDForUpdate(ctx, tx, input.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("session not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.get_booking_error", "session_id", input.SessionID, "err", err)
		return derrors.Infra("failed to load session booking", err)
	}

	if booking.PhotographerUserID() != input.PhotographerID {
		return derrors.Forbidden("session does not belong to photographer")
	}
	if booking.Status() != photosessionmodel.BookingStatusPendingApproval {
		return derrors.Conflict("session is not pending approval")
	}

	if _, err = s.repo.GetEntryByIDForUpdate(ctx, tx, booking.AgendaEntryID()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.Infra("agenda entry missing for booking", err)
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.get_entry_error", "agenda_entry_id", booking.AgendaEntryID(), "err", err)
		return derrors.Infra("failed to load agenda entry", err)
	}

	if err := s.repo.UpdateBookingStatus(ctx, tx, booking.ID(), status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("session not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.update_error", "session_id", booking.ID(), "err", err)
		return derrors.Infra("failed to update session status", err)
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.tx_commit_error", "session_id", booking.ID(), "err", err)
		return derrors.Infra("failed to commit transaction", err)
	}
	committed = true

	logger.Info("photo_session.status.updated", "session_id", booking.ID(), "photographer_id", input.PhotographerID, "status", statusStr)
	return nil
}
