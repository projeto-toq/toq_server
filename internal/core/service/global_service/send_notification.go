package globalservice

import (
	"context"
	"fmt"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// NotificationService interface para o serviço de notificações
type NotificationService interface {
	SendNotification(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType, code ...string) error
}

// notificationService implementa NotificationService
type notificationService struct {
	handler NotificationHandler
}

// NewNotificationService cria uma nova instância do serviço
func NewNotificationService(gs *globalService) NotificationService {
	sender := NewNotificationSender(gs)
	templateLoader := NewEmailTemplateLoader()
	handler := NewNotificationHandler(sender, templateLoader)

	return &notificationService{
		handler: handler,
	}
}

// SendNotification é a função principal para envio de notificações
func (ns *notificationService) SendNotification(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType, code ...string) error {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		slog.Error("Failed to generate tracer", "error", err)
		return err
	}
	defer spanEnd()

	// Validar se código é necessário para tipos específicos
	iCode := ""
	if ns.requiresCode(notificationType) && len(code) == 0 {
		err := utils.ErrInternalServer
		slog.Error("Code required but not provided", "notificationType", notificationType, "error", err)
		return err
	}

	if len(code) > 0 {
		iCode = code[0]
	}

	slog.Info("Processing notification",
		"notificationType", notificationType,
		"userID", user.GetID(),
		"hasCode", len(code) > 0)

	// Processar notificação baseada no tipo
	err = ns.routeNotification(ctx, user, notificationType, iCode)
	if err != nil {
		slog.Error("Failed to send notification",
			"notificationType", notificationType,
			"userID", user.GetID(),
			"error", err)
		return fmt.Errorf("failed to send notification: %w", err)
	}

	slog.Info("Notification sent successfully",
		"notificationType", notificationType,
		"userID", user.GetID())

	return nil
}

// requiresCode verifica se o tipo de notificação requer código
func (ns *notificationService) requiresCode(notificationType globalmodel.NotificationType) bool {
	requiresCodeTypes := []globalmodel.NotificationType{
		globalmodel.NotificationEmailChange,
		globalmodel.NotificationPhoneChange,
		globalmodel.NotificationPasswordChange,
	}

	for _, reqType := range requiresCodeTypes {
		if notificationType == reqType {
			return true
		}
	}
	return false
}

// routeNotification roteia a notificação para o handler apropriado
func (ns *notificationService) routeNotification(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType, code string) error {
	switch notificationType {
	case globalmodel.NotificationEmailChange:
		return ns.handler.HandleEmailChange(ctx, code)

	case globalmodel.NotificationPhoneChange:
		return ns.handler.HandlePhoneChange(ctx, code)

	case globalmodel.NotificationPasswordChange:
		return ns.handler.HandlePasswordChange(ctx, code)

	// CRECI notifications removed - system no longer used
	// case globalmodel.NotificationCreciStateUnsupported,
	// 	globalmodel.NotificationInvalidCreciState,
	// 	globalmodel.NotificationInvalidCreciNumber,
	// 	globalmodel.NotificationBadSelfieImage,
	// 	globalmodel.NotificationBadCreciImages,
	// 	globalmodel.NotificationCreciValidated:
	// 	return ns.handler.HandleCreciNotification(ctx, user, notificationType)

	case globalmodel.NotificationRealtorInviteSMS,
		globalmodel.NotificationRealtorInvitePush:
		return ns.handler.HandleRealtorInvite(ctx, user, notificationType, code)

	case globalmodel.NotificationInviteAccepted,
		globalmodel.NotificationInviteRejected:
		return ns.handler.HandleInviteResponse(ctx, user, notificationType, code)

	case globalmodel.NotificationAgencyRemovedFromRealtor,
		globalmodel.NotificationRealtorRemovedFromAgency:
		return ns.handler.HandleUnlink(ctx, user, notificationType, code)

	default:
		err := fmt.Errorf("unsupported notification type: %v", notificationType)
		slog.Error("Unsupported notification type", "notificationType", notificationType)
		return err
	}
}

// SendNotification mantém a compatibilidade com a interface existente
func (gs *globalService) SendNotification(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType, code ...string) error {
	// Criar o serviço de notificação lazy
	notificationSvc := NewNotificationService(gs)
	return notificationSvc.SendNotification(ctx, user, notificationType, code...)
}
