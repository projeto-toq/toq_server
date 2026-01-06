package globalservice

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// NotificationType enumerates the supported notification channels.
type NotificationType string

const (
	// NotificationTypeEmail sends e-mails via the configured adapter.
	NotificationTypeEmail NotificationType = "email"
	// NotificationTypeSMS delivers SMS messages.
	NotificationTypeSMS NotificationType = "sms"
	// NotificationTypeFCM sends push notifications via Firebase Cloud Messaging.
	NotificationTypeFCM NotificationType = "fcm"
)

// NotificationRequest represents the payload handled by the unified notification service.
type NotificationRequest struct {
	// Type defines which adapter must be used (email, sms, fcm).
	Type NotificationType `json:"type"`

	// From optionally overrides the sender when supported by the adapter (email only for now).
	From string `json:"from,omitempty"`

	// To is required for email and SMS notifications (email address or E.164 phone).
	To string `json:"to"`

	// Subject is mandatory for email and FCM (used as email subject or push title).
	Subject string `json:"subject"`

	// Body carries the message content for every channel.
	Body string `json:"body"`

	// ImageURL optionally enriches push notifications with images/icons.
	ImageURL string `json:"imageUrl,omitempty"`

	// Token is mandatory for FCM notifications (device token from clients).
	Token string `json:"token,omitempty"`

	// Data delivers additional key/value metadata to FCM clients.
	Data map[string]string `json:"data,omitempty"`
}

// UnifiedNotificationService centralizes all notification flows for the application.
type UnifiedNotificationService interface {
	// SendNotification triggers an asynchronous delivery preserving trace/log context.
	// It always returns immediately and should be preferred for user-facing flows.
	SendNotification(ctx context.Context, request NotificationRequest) error

	// SendNotificationSync blocks until the underlying adapter finishes the delivery.
	// Use only when the caller must guarantee the delivery result.
	SendNotificationSync(ctx context.Context, request NotificationRequest) error
}

// unifiedNotificationService implementa UnifiedNotificationService
type unifiedNotificationService struct {
	globalService *globalService
}

// NewUnifiedNotificationService creates a new orchestrator around the provided global service.
func NewUnifiedNotificationService(gs *globalService) UnifiedNotificationService {
	return &unifiedNotificationService{
		globalService: gs,
	}
}

// SendNotification dispatches the request asynchronously propagating trace/log metadata
// to the goroutine responsible for the real delivery.
func (ns *unifiedNotificationService) SendNotification(ctx context.Context, request NotificationRequest) error {
	ctx = coreutils.ContextWithLogger(ctx)

	// Executar notificação em goroutine assíncrona
	go func(parentCtx context.Context) {
		notifyCtx := context.Background()

		if sc := trace.SpanFromContext(parentCtx).SpanContext(); sc.IsValid() {
			notifyCtx = trace.ContextWithSpanContext(notifyCtx, sc)
		}

		if requestID := parentCtx.Value(globalmodel.RequestIDKey); requestID != nil {
			notifyCtx = context.WithValue(notifyCtx, globalmodel.RequestIDKey, requestID)
		}

		notifyCtx = coreutils.ContextWithLogger(notifyCtx)
		notifyCtx, spanEnd, _ := coreutils.GenerateBusinessTracer(notifyCtx, "NotificationService.SendNotificationAsync")
		defer spanEnd()

		notifyLogger := coreutils.LoggerFromContext(notifyCtx)

		// Chamar o método interno síncrono
		if err := ns.sendNotificationSync(notifyCtx, request); err != nil {
			coreutils.SetSpanError(notifyCtx, err)
			notifyLogger.Error("notification.async_send_error",
				"type", request.Type,
				"to", request.To,
				"token", request.Token,
				"err", err)
		}
	}(ctx)

	// Retorna imediatamente (sem aguardar o envio)
	return nil
}

// SendNotificationSync executes the delivery inline. Prefer SendNotification unless
// the caller must wait for third-party confirmation.
func (ns *unifiedNotificationService) SendNotificationSync(ctx context.Context, request NotificationRequest) error {
	return ns.sendNotificationSync(ctx, request)
}

// sendNotificationSync contains the synchronous delivery pipeline (tracing + validation + adapter dispatch).
func (ns *unifiedNotificationService) sendNotificationSync(ctx context.Context, request NotificationRequest) error {
	ctx, spanEnd, err := coreutils.GenerateTracer(ctx)
	if err != nil {
		ctx = coreutils.ContextWithLogger(ctx)
		coreutils.LoggerFromContext(ctx).Error("notification.tracer_error", "err", err)
		return coreutils.InternalError("")
	}
	defer spanEnd()

	ctx = coreutils.ContextWithLogger(ctx)
	logger := coreutils.LoggerFromContext(ctx)

	logger.Info("notification.processing",
		"type", request.Type,
		"to", request.To,
		"subject", request.Subject)

	if err := ns.validateRequest(request); err != nil {
		logger.Warn("notification.request_invalid", "err", err, "type", request.Type)
		return coreutils.BadRequest(err.Error())
	}

	switch request.Type {
	case NotificationTypeEmail:
		return ns.sendEmail(ctx, request)
	case NotificationTypeSMS:
		return ns.sendSMS(ctx, request)
	case NotificationTypeFCM:
		return ns.sendFCM(ctx, request)
	default:
		logger.Warn("notification.type_invalid", "type", request.Type)
		return coreutils.BadRequest("unsupported notification type")
	}
}

// validateRequest enforces the required fields per notification type before dispatching.
func (ns *unifiedNotificationService) validateRequest(request NotificationRequest) error {
	if request.Body == "" {
		return fmt.Errorf("body is required for every notification type")
	}

	switch request.Type {
	case NotificationTypeEmail:
		if request.To == "" {
			return fmt.Errorf("email notifications require the recipient address")
		}
		if request.Subject == "" {
			return fmt.Errorf("email notifications require a subject")
		}

	case NotificationTypeSMS:
		if request.To == "" {
			return fmt.Errorf("sms notifications require the destination phone number")
		}

	case NotificationTypeFCM:
		if request.Token == "" {
			return fmt.Errorf("fcm notifications require the device token")
		}
		if request.Subject == "" {
			return fmt.Errorf("fcm notifications require a subject/title")
		}

	default:
		return fmt.Errorf("unsupported notification type %s", request.Type)
	}

	return nil
}

// sendEmail adapts the request to the email adapter contract and propagates infra failures.
func (ns *unifiedNotificationService) sendEmail(ctx context.Context, request NotificationRequest) error {
	ctx = coreutils.ContextWithLogger(ctx)
	logger := coreutils.LoggerFromContext(ctx)
	logger.Debug("notification.email_sending", "to", request.To, "subject", request.Subject)

	notification := globalmodel.Notification{
		Title: request.Subject,
		Body:  request.Body,
		To:    request.To,
		Icon:  "", // Emails do not render icons
	}

	err := ns.globalService.email.SendEmail(ctx, notification)
	if err != nil {
		coreutils.SetSpanError(ctx, err)
		logger.Error("notification.email_send_error", "err", err, "to", request.To)
		return coreutils.InternalError("")
	}

	logger.Info("notification.email_sent", "to", request.To, "subject", request.Subject)
	return nil
}

// sendSMS adapts the request to the SMS adapter contract.
func (ns *unifiedNotificationService) sendSMS(ctx context.Context, request NotificationRequest) error {
	ctx = coreutils.ContextWithLogger(ctx)
	logger := coreutils.LoggerFromContext(ctx)
	logger.Debug("notification.sms_sending", "to", request.To)

	notification := globalmodel.Notification{
		Title: request.Subject, // Subject optional but forwarded when provided
		Body:  request.Body,
		To:    request.To,
		Icon:  "",
	}

	err := ns.globalService.sms.SendSms(notification)
	if err != nil {
		coreutils.SetSpanError(ctx, err)
		logger.Error("notification.sms_send_error", "err", err, "to", request.To)
		return coreutils.InternalError("")
	}

	logger.Info("notification.sms_sent", "to", request.To)
	return nil
}

// sendFCM adapts the request to the FCM adapter contract.
func (ns *unifiedNotificationService) sendFCM(ctx context.Context, request NotificationRequest) error {
	ctx = coreutils.ContextWithLogger(ctx)
	logger := coreutils.LoggerFromContext(ctx)
	logger.Debug("notification.fcm_sending", "token", request.Token, "title", request.Subject)

	notification := globalmodel.Notification{
		Title:       request.Subject,
		Body:        request.Body,
		Icon:        request.ImageURL,
		DeviceToken: request.Token,
		To:          "",
		Data:        cloneStringMap(request.Data),
	}

	err := ns.globalService.firebaseCloudMessage.SendSingleMessage(ctx, notification)
	if err != nil {
		coreutils.SetSpanError(ctx, err)
		logger.Error("notification.fcm_send_error", "err", err, "token", request.Token)
		return coreutils.InternalError("")
	}

	logger.Info("notification.fcm_sent", "token", request.Token, "title", request.Subject)
	return nil
}

// cloneStringMap copies the provided map so downstream adapters can mutate safely.
func cloneStringMap(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(input))
	for key, value := range input {
		cloned[key] = value
	}
	return cloned
}
