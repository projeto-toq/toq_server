package photosessionservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateSessionStatus updates the status of a photo session booking and notifies the listing owner.
func (s *photoSessionService) UpdateSessionStatus(ctx context.Context, input UpdateSessionStatusInput) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

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
	if status == photosessionmodel.BookingStatusAccepted || status == photosessionmodel.BookingStatusRejected {
		// Aceitar/Rejeitar só é permitido em PENDING_APPROVAL
		if booking.Status() != photosessionmodel.BookingStatusPendingApproval {
			return derrors.Conflict("session is not pending approval")
		}
	} else if status == photosessionmodel.BookingStatusDone {
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
	listing, err := s.listingRepo.GetListingByID(ctx, tx, booking.ListingID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("listing not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.get_listing_error", "listing_id", booking.ListingID(), "err", err)
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

// sendOwnerNotifications sends FCM push notifications to all opted-in devices of the listing owner.
// This function runs asynchronously and logs any errors without propagating them.
func (s *photoSessionService) sendOwnerNotifications(ctx context.Context, userID int64, title, body string, listingID int64, sessionID uint64) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Busca todos os device tokens do usuário (apenas dispositivos com opt-in ativo)
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

	// Obtém serviço de notificação unificado
	notifier := s.globalService.GetUnifiedNotificationService()
	if notifier == nil {
		logger.Error("photo_session.notification.service_unavailable",
			"user_id", userID)
		return
	}

	// Envia notificação para cada token (permite múltiplos dispositivos do mesmo usuário)
	sentCount := 0
	for _, token := range tokens {
		req := globalservice.NotificationRequest{
			Type:    globalservice.NotificationTypeFCM,
			Token:   token,
			Subject: title,
			Body:    body,
		}

		// SendNotification é assíncrono por padrão, mas já estamos em goroutine
		if notifErr := notifier.SendNotification(ctx, req); notifErr != nil {
			logger.Warn("photo_session.notification.send_error",
				"user_id", userID,
				"token", token[:min(len(token), 20)]+"...", // Log apenas início do token
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

// min returns the minimum of two integers (helper function)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
