package globalservice

import (
	"context"
	"fmt"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// NotificationType define os tipos de notificação suportados
type NotificationType string

const (
	// NotificationTypeEmail representa notificação por e-mail
	NotificationTypeEmail NotificationType = "email"
	// NotificationTypeSMS representa notificação por SMS
	NotificationTypeSMS NotificationType = "sms"
	// NotificationTypeFCM representa notificação push via Firebase Cloud Messaging
	NotificationTypeFCM NotificationType = "fcm"
)

// NotificationRequest representa uma requisição de notificação
type NotificationRequest struct {
	// Type define o tipo de notificação (email, sms, fcm)
	Type NotificationType `json:"type"`

	// From é opcional e será usado na notificação de email
	From string `json:"from,omitempty"`

	// To é obrigatório para email e sms. Contém o número de telefone para SMS ou endereço de email para email
	To string `json:"to"`

	// Subject é obrigatório para e-mail. Contém o subject do email ou title do FCM
	Subject string `json:"subject"`

	// Body é obrigatório para todos. Contém o corpo da mensagem
	Body string `json:"body"`

	// ImageURL é opcional e conterá imageURL do FCM
	ImageURL string `json:"imageUrl,omitempty"`

	// Token é necessário para FCM, conterá o deviceToken
	Token string `json:"token,omitempty"`
}

// UnifiedNotificationService interface para o novo sistema de notificação unificado
type UnifiedNotificationService interface {
	// SendNotification envia uma notificação de forma ASSÍNCRONA (recomendado)
	// Retorna imediatamente sem aguardar o envio
	SendNotification(ctx context.Context, request NotificationRequest) error

	// SendNotificationSync envia uma notificação de forma SÍNCRONA (usar apenas quando necessário)
	// Aguarda o envio completo antes de retornar
	SendNotificationSync(ctx context.Context, request NotificationRequest) error
}

// unifiedNotificationService implementa UnifiedNotificationService
type unifiedNotificationService struct {
	globalService *globalService
}

// NewUnifiedNotificationService cria uma nova instância do serviço de notificação unificado
func NewUnifiedNotificationService(gs *globalService) UnifiedNotificationService {
	return &unifiedNotificationService{
		globalService: gs,
	}
}

// SendNotification envia uma notificação de forma assíncrona baseada no tipo especificado na requisição
// Esta função é o ponto central do novo sistema de notificação, direcionando
// cada requisição para o adapter apropriado baseado no tipo de notificação.
// IMPORTANTE: Esta função é assíncrona por padrão para garantir que todas as
// notificações não bloqueiem a resposta ao usuário.
func (ns *unifiedNotificationService) SendNotification(ctx context.Context, request NotificationRequest) error {
	// Executar notificação em goroutine assíncrona
	go func() {
		// Criar contexto independente preservando Request ID
		notifyCtx := context.Background()
		if requestID := ctx.Value(globalmodel.RequestIDKey); requestID != nil {
			notifyCtx = context.WithValue(notifyCtx, globalmodel.RequestIDKey, requestID)
		}

		// Chamar o método interno síncrono
		err := ns.sendNotificationSync(notifyCtx, request)
		if err != nil {
			slog.Error("notification.async_send_error",
				"type", request.Type,
				"to", request.To,
				"token", request.Token,
				"err", err)
		}
	}()

	// Retorna imediatamente (sem aguardar o envio)
	return nil
}

// SendNotificationSync envia uma notificação de forma SÍNCRONA
// Use apenas quando realmente precisar aguardar o resultado do envio
func (ns *unifiedNotificationService) SendNotificationSync(ctx context.Context, request NotificationRequest) error {
	return ns.sendNotificationSync(ctx, request)
}

// sendNotificationSync executa o envio síncrono da notificação
// Este método é interno e contém a lógica real de envio
func (ns *unifiedNotificationService) sendNotificationSync(ctx context.Context, request NotificationRequest) error {
	// Gerar tracer para observabilidade
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		slog.Error("notification.tracer_error", "err", err)
		return utils.InternalError("")
	}
	defer spanEnd()

	// Log da requisição de notificação
	slog.Info("notification.processing",
		"type", request.Type,
		"to", request.To,
		"subject", request.Subject)

	// Validar a requisição antes de processar
	if err := ns.validateRequest(request); err != nil {
		slog.Error("notification.request_invalid", "err", err)
		return utils.BadRequest("invalid notification request")
	}

	// Direcionar para o método apropriado baseado no tipo
	switch request.Type {
	case NotificationTypeEmail:
		return ns.sendEmail(ctx, request)
	case NotificationTypeSMS:
		return ns.sendSMS(ctx, request)
	case NotificationTypeFCM:
		return ns.sendFCM(ctx, request)
	default:
		slog.Error("notification.type_invalid", "type", request.Type)
		return utils.BadRequest("unsupported notification type")
	}
}

// validateRequest valida os campos obrigatórios da requisição baseado no tipo
// Cada tipo de notificação tem requisitos específicos que são verificados aqui.
func (ns *unifiedNotificationService) validateRequest(request NotificationRequest) error {
	// Body é obrigatório para todos os tipos
	if request.Body == "" {
		return fmt.Errorf("campo 'body' é obrigatório")
	}

	switch request.Type {
	case NotificationTypeEmail:
		// Para email: to e subject são obrigatórios
		if request.To == "" {
			return fmt.Errorf("campo 'to' é obrigatório para notificações por email")
		}
		if request.Subject == "" {
			return fmt.Errorf("campo 'subject' é obrigatório para notificações por email")
		}

	case NotificationTypeSMS:
		// Para SMS: to é obrigatório (número de telefone)
		if request.To == "" {
			return fmt.Errorf("campo 'to' é obrigatório para notificações por SMS")
		}

	case NotificationTypeFCM:
		// Para FCM: token e subject são obrigatórios
		if request.Token == "" {
			return fmt.Errorf("campo 'token' é obrigatório para notificações FCM")
		}
		if request.Subject == "" {
			return fmt.Errorf("campo 'subject' é obrigatório para notificações FCM (usado como title)")
		}

	default:
		return fmt.Errorf("tipo de notificação não suportado: %s", request.Type)
	}

	return nil
}

// sendEmail processa notificação por email usando o email adapter
// Converte a requisição para o formato esperado pelo adapter de email.
func (ns *unifiedNotificationService) sendEmail(ctx context.Context, request NotificationRequest) error {
	_ = ctx //ignorar temporariament
	slog.Debug("notification.email_sending", "to", request.To, "subject", request.Subject)

	// Construir a notificação no formato esperado pelo adapter
	notification := globalmodel.Notification{
		Title: request.Subject,
		Body:  request.Body,
		To:    request.To,
		Icon:  "", // Email não usa ícone
	}

	// Se From foi especificado, adicionar à notificação (implementação futura)
	// Note: From field support to be implemented in email adapter

	// Enviar através do adapter de email
	err := ns.globalService.email.SendEmail(notification)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("notification.email_send_error", "err", err, "to", request.To)
		return utils.InternalError("")
	}

	slog.Info("notification.email_sent", "to", request.To, "subject", request.Subject)
	return nil
}

// sendSMS processa notificação por SMS usando o SMS adapter
// Converte a requisição para o formato esperado pelo adapter de SMS.
func (ns *unifiedNotificationService) sendSMS(ctx context.Context, request NotificationRequest) error {
	_ = ctx //ignorar temporariamente
	slog.Debug("notification.sms_sending", "to", request.To)

	// Construir a notificação no formato esperado pelo adapter
	notification := globalmodel.Notification{
		Title: request.Subject, // Subject opcional para SMS, usado como título se fornecido
		Body:  request.Body,
		To:    request.To,
		Icon:  "", // SMS não usa ícone
	}

	// Enviar através do adapter de SMS
	err := ns.globalService.sms.SendSms(notification)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("notification.sms_send_error", "err", err, "to", request.To)
		return utils.InternalError("")
	}

	slog.Info("notification.sms_sent", "to", request.To)
	return nil
}

// sendFCM processa notificação push usando o FCM adapter
// Converte a requisição para o formato esperado pelo adapter FCM.
func (ns *unifiedNotificationService) sendFCM(ctx context.Context, request NotificationRequest) error {
	slog.Debug("notification.fcm_sending", "token", request.Token, "title", request.Subject)

	// Construir a notificação no formato esperado pelo adapter
	notification := globalmodel.Notification{
		Title:       request.Subject, // Subject é usado como title do push
		Body:        request.Body,
		Icon:        request.ImageURL, // ImageURL é usado como ícone
		DeviceToken: request.Token,
		To:          "", // FCM não usa To, usa DeviceToken
	}

	// Enviar através do adapter FCM
	err := ns.globalService.firebaseCloudMessage.SendSingleMessage(ctx, notification)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("notification.fcm_send_error", "err", err, "token", request.Token)
		return utils.InternalError("")
	}

	slog.Info("notification.fcm_sent", "token", request.Token, "title", request.Subject)
	return nil
}
