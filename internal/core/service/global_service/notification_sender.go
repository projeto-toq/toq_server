package globalservice

import (
	"context"
	"fmt"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// NotificationSender interface para envio de notificações
type NotificationSender interface {
	SendEmail(ctx context.Context, title, body, to string) error
	SendSMS(ctx context.Context, title, body, to string) error
	SendPush(ctx context.Context, deviceToken, title, body string) error
	SendPushOrEmail(ctx context.Context, user usermodel.UserInterface, title, body string) error
	SendPushToUserDevices(ctx context.Context, userID int64, title, body string) error
	SendPushToAllOptedInUsers(ctx context.Context, title, body string) error
}

// notificationSender implementa NotificationSender
type notificationSender struct {
	globalService *globalService
}

// NewNotificationSender cria uma nova instância do sender
func NewNotificationSender(gs *globalService) NotificationSender {
	return &notificationSender{
		globalService: gs,
	}
}

// SendEmail envia notificação por email
func (ns *notificationSender) SendEmail(ctx context.Context, title, body, to string) error {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		slog.Error("Failed to generate tracer for email", "error", err)
		return err
	}
	defer spanEnd()

	notification := globalmodel.Notification{
		Title: title,
		Body:  body,
		Icon:  "",
		To:    to,
	}

	slog.Debug("Sending email notification", "to", to, "title", title)

	err = ns.globalService.email.SendEmail(notification)
	if err != nil {
		slog.Error("Failed to send email", "to", to, "title", title, "error", err)
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	slog.Info("Email sent successfully", "to", to, "title", title)
	return nil
}

// SendSMS envia notificação por SMS
func (ns *notificationSender) SendSMS(ctx context.Context, title, body, to string) error {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		slog.Error("Failed to generate tracer for SMS", "error", err)
		return err
	}
	defer spanEnd()

	notification := globalmodel.Notification{
		Title: title,
		Body:  body,
		Icon:  "",
		To:    to,
	}

	slog.Debug("Sending SMS notification", "to", to, "title", title)

	err = ns.globalService.sms.SendSms(notification)
	if err != nil {
		slog.Error("Failed to send SMS", "to", to, "title", title, "error", err)
		return fmt.Errorf("failed to send SMS to %s: %w", to, err)
	}

	slog.Info("SMS sent successfully", "to", to, "title", title)
	return nil
}

// SendPush envia notificação push
func (ns *notificationSender) SendPush(ctx context.Context, deviceToken, title, body string) error {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		slog.Error("Failed to generate tracer for push", "error", err)
		return err
	}
	defer spanEnd()

	notification := globalmodel.Notification{
		Title:       title,
		Body:        body,
		Icon:        "",
		DeviceToken: deviceToken,
	}

	slog.Debug("Sending push notification", "deviceToken", deviceToken, "title", title)

	err = ns.globalService.firebaseCloudMessage.SendSingleMessage(ctx, notification)
	if err != nil {
		slog.Error("Failed to send push notification", "deviceToken", deviceToken, "title", title, "error", err)
		return fmt.Errorf("failed to send push notification: %w", err)
	}

	slog.Info("Push notification sent successfully", "deviceToken", deviceToken, "title", title)
	return nil
}

// SendPushOrEmail envia push notification ou email como fallback
func (ns *notificationSender) SendPushOrEmail(ctx context.Context, user usermodel.UserInterface, title, body string) error {
	if user.GetDeviceToken() == "" {
		slog.Debug("No device token, sending email fallback", "userID", user.GetID())
		// Placeholder - will use user.GetEmail() in production
		return ns.SendEmail(ctx, title, body, "giulio.alfieri@gmail.com")
	}

	slog.Debug("Device token available, sending push notification", "userID", user.GetID())
	return ns.SendPush(ctx, user.GetDeviceToken(), title, body)
}

// SendPushToUserDevices envia push notification para todos os dispositivos de um usuário (se opt_status=1)
func (ns *notificationSender) SendPushToUserDevices(ctx context.Context, userID int64, title, body string) error {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		slog.Error("Failed to generate tracer for multiple push", "error", err)
		return err
	}
	defer spanEnd()

	// Busca os tokens do usuário apenas se ele tem opt_status=1
	tokens, err := ns.globalService.deviceTokenRepo.ListTokensByUserIDIfOptedIn(userID)
	if err != nil {
		slog.Error("Failed to get device tokens for user", "userID", userID, "error", err)
		return fmt.Errorf("failed to get device tokens for user %d: %w", userID, err)
	}

	if len(tokens) == 0 {
		slog.Debug("No device tokens found for opted-in user", "userID", userID)
		return nil
	}

	notification := globalmodel.Notification{
		Title: title,
		Body:  body,
		Icon:  "",
	}

	slog.Debug("Sending push to user devices", "userID", userID, "tokenCount", len(tokens), "title", title)

	err = ns.globalService.firebaseCloudMessage.SendMultipleMessages(ctx, notification, tokens)
	if err != nil {
		slog.Error("Failed to send push to user devices", "userID", userID, "tokenCount", len(tokens), "error", err)
		return fmt.Errorf("failed to send push to user devices: %w", err)
	}

	slog.Info("Push notifications sent to user devices", "userID", userID, "tokenCount", len(tokens), "title", title)
	return nil
}

// SendPushToAllOptedInUsers envia push notification para todos os usuários com opt_status=1
func (ns *notificationSender) SendPushToAllOptedInUsers(ctx context.Context, title, body string) error {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		slog.Error("Failed to generate tracer for broadcast push", "error", err)
		return err
	}
	defer spanEnd()

	// Busca todos os tokens de usuários com opt_status=1
	tokens, err := ns.globalService.deviceTokenRepo.ListTokensByOptedInUsers()
	if err != nil {
		slog.Error("Failed to get device tokens for opted-in users", "error", err)
		return fmt.Errorf("failed to get device tokens for opted-in users: %w", err)
	}

	if len(tokens) == 0 {
		slog.Debug("No device tokens found for opted-in users")
		return nil
	}

	notification := globalmodel.Notification{
		Title: title,
		Body:  body,
		Icon:  "",
	}

	slog.Debug("Sending broadcast push", "tokenCount", len(tokens), "title", title)

	err = ns.globalService.firebaseCloudMessage.SendMultipleMessages(ctx, notification, tokens)
	if err != nil {
		slog.Error("Failed to send broadcast push", "tokenCount", len(tokens), "error", err)
		return fmt.Errorf("failed to send broadcast push: %w", err)
	}

	slog.Info("Broadcast push notifications sent", "tokenCount", len(tokens), "title", title)
	return nil
}
