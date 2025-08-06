package globalservice

import (
	"context"
	"fmt"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// NotificationHandler interface para handlers de notificação
type NotificationHandler interface {
	HandleEmailChange(ctx context.Context, code string) error
	HandlePhoneChange(ctx context.Context, code string) error
	HandlePasswordChange(ctx context.Context, code string) error
	HandleCreciNotification(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType) error
	HandleRealtorInvite(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType, code string) error
	HandleInviteResponse(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType, code string) error
	HandleUnlink(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType, code string) error
}

// notificationHandler implementa NotificationHandler
type notificationHandler struct {
	sender         NotificationSender
	templateLoader EmailTemplateLoader
}

// NewNotificationHandler cria uma nova instância do handler
func NewNotificationHandler(sender NotificationSender, templateLoader EmailTemplateLoader) NotificationHandler {
	return &notificationHandler{
		sender:         sender,
		templateLoader: templateLoader,
	}
}

// HandleEmailChange processa notificação de mudança de email
func (nh *notificationHandler) HandleEmailChange(ctx context.Context, code string) error {
	htmlBody, err := nh.templateLoader.LoadTemplate(EmailVerificationTemplate, code)
	if err != nil {
		slog.Error("Failed to load email verification template", "error", err)
		return fmt.Errorf("failed to load email verification template: %w", err)
	}

	return nh.sender.SendEmail(ctx,
		"Confirmação de e-mail da TOQ",
		htmlBody,
		"giulio.alfieri@gmail.com") // TODO: Change to user.GetEmail()
}

// HandlePhoneChange processa notificação de mudança de telefone
func (nh *notificationHandler) HandlePhoneChange(ctx context.Context, code string) error {
	body := fmt.Sprintf("Para validar seu telefone cadastrado na TOQ insira o código %s no App:", code)
	return nh.sender.SendSMS(ctx,
		"Confirmação de telefone da TOQ",
		body,
		"+5511999141768") // TODO: Change to user.GetPhoneNumber()
}

// HandlePasswordChange processa notificação de mudança de senha
func (nh *notificationHandler) HandlePasswordChange(ctx context.Context, code string) error {
	htmlBody, err := nh.templateLoader.LoadTemplate(PasswordResetTemplate, code)
	if err != nil {
		slog.Error("Failed to load password reset template", "error", err)
		return fmt.Errorf("failed to load password reset template: %w", err)
	}

	return nh.sender.SendEmail(ctx,
		"Confirmação de troca de senha da TOQ",
		htmlBody,
		"giulio.alfieri@gmail.com") // TODO: Change to user.GetEmail()
}

// HandleCreciNotification processa notificações relacionadas ao CRECI
func (nh *notificationHandler) HandleCreciNotification(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType) error {
	title, body, err := nh.getCreciNotificationContent(notificationType)
	if err != nil {
		return err
	}

	return nh.sender.SendPushOrEmail(ctx, user, title, body)
}

// HandleRealtorInvite processa notificações de convite de corretor
func (nh *notificationHandler) HandleRealtorInvite(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType, code string) error {
	switch notificationType {
	case globalmodel.NotificationRealtorInviteSMS:
		body := fmt.Sprintf("A %s está te convidando a participar da TOQ, vinculado(a) a ela. Baixe a aplicação e aceite o convite.", code)
		return nh.sender.SendSMS(ctx,
			"Convite para participar da TOQ",
			body,
			"+5511999141768") // TODO: Change to user.GetPhoneNumber()

	case globalmodel.NotificationRealtorInvitePush:
		body := fmt.Sprintf("%s, você tem um convite pendente para vínculo a uma imobiliária.", user.GetNickName())
		return nh.sender.SendPushOrEmail(ctx, user, "Convite para vínculo a imobiliária", body)

	default:
		return fmt.Errorf("unsupported realtor invite notification type: %v", notificationType)
	}
}

// HandleInviteResponse processa notificações de resposta de convite
func (nh *notificationHandler) HandleInviteResponse(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType, code string) error {
	title, body, err := nh.getInviteResponseContent(user, notificationType, code)
	if err != nil {
		return err
	}

	return nh.sender.SendPushOrEmail(ctx, user, title, body)
}

// HandleUnlink processa notificações de desvinculação
func (nh *notificationHandler) HandleUnlink(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType, code string) error {
	title, body, err := nh.getUnlinkContent(user, notificationType, code)
	if err != nil {
		return err
	}

	return nh.sender.SendPushOrEmail(ctx, user, title, body)
}

// getCreciNotificationContent retorna o conteúdo para notificações de CRECI
func (nh *notificationHandler) getCreciNotificationContent(notificationType globalmodel.NotificationType) (string, string, error) {
	switch notificationType {
	case globalmodel.NotificationCreciStateUnsupported:
		return "Erro na validação do Creci",
			"O estado informado do seu Creci ainda não é suportado.",
			nil
	case globalmodel.NotificationInvalidCreciState:
		return "Erro na validação do Creci",
			"O estado do creci informado não corresponde ao estado da imagem. Por favor, tente novamente.",
			nil
	case globalmodel.NotificationInvalidCreciNumber:
		return "Erro na validação do Creci",
			"O número do creci informado não corresponde ao número da imagem. Por favor, tente novamente.",
			nil
	case globalmodel.NotificationBadSelfieImage:
		return "Erro na validação do Creci",
			"A imagem da selfie não corresponde a imagem do documento. Por favor, tente novamente.",
			nil
	case globalmodel.NotificationBadCreciImages:
		return "Erro na validação do Creci",
			"As imagens do seu Creci não puderam ser validadas, pois estão com baixa qualidade. Por favor, tente novamente.",
			nil
	case globalmodel.NotificationCreciValidated:
		return "Creci validado",
			"Seu Creci foi validado com sucesso! Agora você pode usar a plataforma.",
			nil
	default:
		return "", "", fmt.Errorf("unsupported CRECI notification type: %v", notificationType)
	}
}

// getInviteResponseContent retorna o conteúdo para notificações de resposta de convite
func (nh *notificationHandler) getInviteResponseContent(user usermodel.UserInterface, notificationType globalmodel.NotificationType, code string) (string, string, error) {
	switch notificationType {
	case globalmodel.NotificationInviteAccepted:
		title := "Corretor aceitou o convite"
		body := fmt.Sprintf("%s, o corretor %s aceitou seu convite e agora está vinculado a esta imobiliária.", user.GetNickName(), code)
		return title, body, nil
	case globalmodel.NotificationInviteRejected:
		title := "Corretor rejeitou o convite"
		body := fmt.Sprintf("%s, o corretor %s rejeitou seu convite para vincular-se a esta imobiliária.", user.GetNickName(), code)
		return title, body, nil
	default:
		return "", "", fmt.Errorf("unsupported invite response notification type: %v", notificationType)
	}
}

// getUnlinkContent retorna o conteúdo para notificações de desvinculação
func (nh *notificationHandler) getUnlinkContent(user usermodel.UserInterface, notificationType globalmodel.NotificationType, code string) (string, string, error) {
	switch notificationType {
	case globalmodel.NotificationAgencyRemovedFromRealtor:
		title := "Corretor cancelou o vínculo"
		body := fmt.Sprintf("O corretor %s cancelou o vínculo com esta imobiliária.", user.GetNickName())
		return title, body, nil
	case globalmodel.NotificationRealtorRemovedFromAgency:
		title := "Imobiliária cancelou o vínculo"
		body := fmt.Sprintf("%s, a imobiliária %s cancelou o vínculo com você.", user.GetNickName(), code)
		return title, body, nil
	default:
		return "", "", fmt.Errorf("unsupported unlink notification type: %v", notificationType)
	}
}
