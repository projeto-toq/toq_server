package photosessionservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateSessionStatus updates the status of a photo session booking and notifies the listing owner.
//
// This endpoint is only available when manual photographer approval is enabled via configuration.
// If automatic approval is enabled (require_photographer_approval=false), this operation returns
// a validation error indicating the feature is disabled.
//
// Supported Status Transitions:
//   - PENDING_APPROVAL → ACCEPTED (photographer accepts)
//   - PENDING_APPROVAL → REJECTED (photographer declines)
//   - ACCEPTED/ACTIVE → DONE (photographer completes session)
//
// Listing Status Changes:
//   - ACCEPTED: StatusPendingPhotoConfirmation → StatusPhotosScheduled
//   - REJECTED: StatusPendingPhotoConfirmation → StatusPendingPhotoScheduling
//   - DONE: StatusPhotosScheduled → StatusPendingPhotoProcessing
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - input: UpdateSessionStatusInput with sessionID, photographerID, status
//
// Returns:
//   - error: Domain error with appropriate HTTP status code:
//   - 400 (BadRequest) if manual approval is disabled in config
//   - 401 (Auth) if photographer not authorized
//   - 403 (Forbidden) if session does not belong to photographer
//   - 404 (NotFound) if session not found
//   - 409 (Conflict) if session not in expected state for transition
//   - 422 (Validation) if input invalid
//   - 500 (Infra) for infrastructure failures
//
// Configuration:
//   - Requires photo_session.require_photographer_approval = true in env.yaml
//
// Side Effects:
//   - Updates photographer_photo_session_bookings.status
//   - Updates listings.status
//   - Sends FCM push notification to listing owner
func (s *photoSessionService) UpdateSessionStatus(ctx context.Context, input UpdateSessionStatusInput) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Check if manual approval mode is enabled
	if !s.cfg.RequirePhotographerApproval {
		logger.Warn("photo_session.update_status.manual_approval_disabled",
			"photographer_id", input.PhotographerID,
			"session_id", input.SessionID)
		return derrors.BadRequest("manual photographer approval is disabled; photo sessions are automatically approved upon reservation")
	}

	// Validações de entrada
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

	// Validar status permitidos: ACCEPTED, REJECTED ou DONE (conclusão)
	if statusStr != string(photosessionmodel.BookingStatusAccepted) &&
		statusStr != string(photosessionmodel.BookingStatusRejected) &&
		statusStr != string(photosessionmodel.BookingStatusDone) {
		return derrors.BadRequest("status must be ACCEPTED, REJECTED or DONE")
	}
	status := photosessionmodel.BookingStatus(statusStr)

	// Inicia transação
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

	// Carrega booking com lock
	booking, err := s.repo.GetBookingByIDForUpdate(ctx, tx, input.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("session not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.get_booking_error", "session_id", input.SessionID, "err", err)
		return derrors.Infra("failed to load session booking", err)
	}

	// Valida que a sessão pertence ao fotógrafo
	if booking.PhotographerUserID() != input.PhotographerID {
		return derrors.Forbidden("session does not belong to photographer")
	}

	// Valida transições de estado permitidas
	switch status {
	case photosessionmodel.BookingStatusAccepted, photosessionmodel.BookingStatusRejected:
		// Aceitar/Rejeitar só é permitido em PENDING_APPROVAL
		if booking.Status() != photosessionmodel.BookingStatusPendingApproval {
			return derrors.Conflict("session is not pending approval")
		}
	case photosessionmodel.BookingStatusDone:
		// Concluir só é permitido em ACCEPTED ou ACTIVE
		if booking.Status() != photosessionmodel.BookingStatusAccepted &&
			booking.Status() != photosessionmodel.BookingStatusActive {
			return derrors.Conflict("session must be accepted or active to be completed")
		}
	}

	// Carrega agenda entry (validação de consistência)
	if _, err = s.repo.GetEntryByIDForUpdate(ctx, tx, booking.AgendaEntryID()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.Infra("agenda entry missing for booking", err)
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.get_entry_error", "agenda_entry_id", booking.AgendaEntryID(), "err", err)
		return derrors.Infra("failed to load agenda entry", err)
	}

	// Atualiza status do booking
	if err := s.repo.UpdateBookingStatus(ctx, tx, booking.ID(), status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("session not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.update_error", "session_id", booking.ID(), "err", err)
		return derrors.Infra("failed to update session status", err)
	}

	// Carrega o listing associado ao booking
	listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, booking.ListingIdentityID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("listing not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.get_listing_error", "listing_identity_id", booking.ListingIdentityID(), "err", err)
		return derrors.Infra("failed to load listing", err)
	}

	// Determina o novo status do listing e mensagem de notificação baseado na decisão do fotógrafo
	var newListingStatus listingmodel.ListingStatus
	var expectedListingStatus listingmodel.ListingStatus
	var notificationTitle string
	var notificationBody string

	switch status {
	case photosessionmodel.BookingStatusAccepted:
		// Fotógrafo aceitou: listing passa de StatusPendingPhotoConfirmation para StatusPhotosScheduled
		expectedListingStatus = listingmodel.StatusPendingPhotoConfirmation
		newListingStatus = listingmodel.StatusPhotosScheduled
		notificationTitle = "Sessão de Fotos Confirmada"
		notificationBody = "O fotógrafo aceitou a sessão de fotos do seu imóvel. A sessão está agendada!"

	case photosessionmodel.BookingStatusRejected:
		// Fotógrafo rejeitou: listing volta de StatusPendingPhotoConfirmation para StatusPendingPhotoScheduling
		expectedListingStatus = listingmodel.StatusPendingPhotoConfirmation
		newListingStatus = listingmodel.StatusPendingPhotoScheduling
		notificationTitle = "Sessão de Fotos Recusada"
		notificationBody = "O fotógrafo não pôde aceitar a sessão. Por favor, reagende para outro horário."

	case photosessionmodel.BookingStatusDone:
		// Fotógrafo concluiu: listing passa de StatusPhotosScheduled para StatusPendingPhotoProcessing
		expectedListingStatus = listingmodel.StatusPhotosScheduled
		newListingStatus = listingmodel.StatusPendingPhotoProcessing
		notificationTitle = "Sessão de Fotos Concluída"
		notificationBody = "A sessão de fotos do seu imóvel foi concluída. As fotos estão sendo processadas."
	}

	// Atualiza o status do listing
	if updateErr := s.listingRepo.UpdateListingStatus(ctx, tx, listing.ID(), newListingStatus, expectedListingStatus); updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			// Listing não está no status esperado - pode ter sido alterado por outro processo
			logger.Warn("photo_session.update_status.listing_status_mismatch",
				"listing_id", listing.ID(),
				"current_status", listing.Status().String(),
				"expected_status", expectedListingStatus.String())
			return derrors.Conflict("listing status has changed")
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("photo_session.update_status.update_listing_status_error", "listing_id", listing.ID(), "err", updateErr)
		return derrors.Infra("failed to update listing status", updateErr)
	}

	// Commit da transação antes de enviar notificações
	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.tx_commit_error", "session_id", booking.ID(), "err", err)
		return derrors.Infra("failed to commit transaction", err)
	}
	committed = true

	// Envia notificações FCM ao proprietário de forma assíncrona
	ownerUserID := listing.UserID()
	go s.sendOwnerNotifications(context.Background(), ownerUserID, notificationTitle, notificationBody, listing.ID(), booking.ID())

	logger.Info("photo_session.status.updated",
		"session_id", booking.ID(),
		"photographer_id", input.PhotographerID,
		"status", statusStr,
		"listing_id", listing.ID(),
		"listing_new_status", newListingStatus.String())

	return nil
}
