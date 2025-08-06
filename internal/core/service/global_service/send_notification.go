package globalservice

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Adicione esta função auxiliar no início do arquivo
func loadEmailTemplate(code string, emailType int) (string, error) {
	templFile := ""
	if emailType == 1 {
		templFile = "../internal/core/templates/email_verification.html"
	} else if emailType == 2 {
		templFile = "../internal/core/templates/email_reset_password.html"
	}
	tmpl, err := template.ParseFiles(templFile)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, struct{ Code string }{Code: code})
	if err != nil {
		return "", err
	}

	return body.String(), nil
}

func (gs *globalService) SendNotification(ctx context.Context, user usermodel.UserInterface, notificationType globalmodel.NotificationType, code ...string) (err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	iCode := ""
	if len(code) == 0 &&
		(notificationType == globalmodel.NotificationEmailChange ||
			notificationType == globalmodel.NotificationPhoneChange ||
			notificationType == globalmodel.NotificationPasswordChange) {
		err = status.Error(codes.Internal, "code is required for this notification type")
		return
	} else if len(code) > 0 {
		iCode = code[0]
	}

	switch notificationType {
	case globalmodel.NotificationEmailChange:
		htmlBody, err1 := loadEmailTemplate(iCode, 1)
		if err1 != nil {
			return err1
		}
		notification := globalmodel.Notification{
			Title: "Confirmação de e-mail da TOQ",
			Body:  htmlBody,
			Icon:  "",
			To:    "giulio.alfieri@gmail.com", // TODO: Change to user.GetEmail(),
		}
		err = sendEmail(ctx, gs, notification)
	case globalmodel.NotificationPhoneChange:
		notification := globalmodel.Notification{
			Title: "Confirmação de telefone da TOQ",
			Body:  fmt.Sprintf("Para validar seu telefone cadastrado na TOQ insira o código %s no App:", iCode),
			Icon:  "",
			To:    "+5511999141768", //TODO: Change to user.GetPhoneNumber()
		}
		err = sendSMS(ctx, gs, notification)
	case globalmodel.NotificationPasswordChange:
		htmlBody, err1 := loadEmailTemplate(iCode, 2)
		if err1 != nil {
			return err1
		}
		notification := globalmodel.Notification{
			Title: "Confirmação de troca de senha da TOQ",
			Body:  htmlBody,
			Icon:  "",
			To:    "giulio.alfieri@gmail.com", // TODO: Change to user.GetEmail(),
		}
		err = sendEmail(ctx, gs, notification)
	case globalmodel.NotificationCreciStateUnsupported:
		notification := globalmodel.Notification{
			Title: "Erro na validação do Creci",
			Body:  "O estado informado do seu Creci ainda não é suportado.",
			Icon:  "",
		}
		if user.GetDeviceToken() == "" {
			notification.To = "giulio.alfieri@gmail.com" // TODO: Change to user.GetEmail(),
			err = sendEmail(ctx, gs, notification)
		} else {
			notification.DeviceToken = user.GetDeviceToken()
			err = sendPush(ctx, gs, notification)
		}
	case globalmodel.NotificationInvalidCreciState:
		notification := globalmodel.Notification{
			Title: "Erro na validação do Creci",
			Body:  "O estado do creci informado não corresponde ao estado da imagem. Por favor, tente novamente.",
			Icon:  "",
		}
		if user.GetDeviceToken() == "" {
			notification.To = "giulio.alfieri@gmail.com" // TODO: Change to user.GetEmail(),
			err = sendEmail(ctx, gs, notification)
		} else {
			notification.DeviceToken = user.GetDeviceToken()
			err = sendPush(ctx, gs, notification)
		}
	case globalmodel.NotificationInvalidCreciNumber:
		notification := globalmodel.Notification{
			Title: "Erro na validação do Creci",
			Body:  "O número do creci informado não corresponde ao número da imagem. Por favor, tente novamente.",
			Icon:  "",
		}
		if user.GetDeviceToken() == "" {
			notification.To = "giulio.alfieri@gmail.com" // TODO: Change to user.GetEmail(),
			err = sendEmail(ctx, gs, notification)
		} else {
			notification.DeviceToken = user.GetDeviceToken()
			err = sendPush(ctx, gs, notification)
		}
	case globalmodel.NotificationBadSelfieImage:
		notification := globalmodel.Notification{
			Title: "Erro na validação do Creci",
			Body:  "A imagem da selfie não corresponde a imagem do documento. Por favor, tente novamente.",
			Icon:  "",
		}
		if user.GetDeviceToken() == "" {
			notification.To = "giulio.alfieri@gmail.com" // TODO: Change to user.GetEmail(),
			err = sendEmail(ctx, gs, notification)
		} else {
			notification.DeviceToken = user.GetDeviceToken()
			err = sendPush(ctx, gs, notification)
		}
	case globalmodel.NotificationBadCreciImages:
		notification := globalmodel.Notification{
			Title: "Erro na validação do Creci",
			Body:  "As imagens do seu Creci não puderam ser validadas, pois estão com baixa qualidade. Por favor, tente novamente.",
			Icon:  "",
		}
		if user.GetDeviceToken() == "" {
			notification.To = "giulio.alfieri@gmail.com" // TODO: Change to user.GetEmail(),
			err = sendEmail(ctx, gs, notification)
		} else {
			notification.DeviceToken = user.GetDeviceToken()
			err = sendPush(ctx, gs, notification)
		}
	case globalmodel.NotificationCreciValidated:
		notification := globalmodel.Notification{
			Title: "Creci validado",
			Body:  "Seu Creci foi validado com sucesso! Agora você pode usar a ",
			Icon:  "",
		}
		if user.GetDeviceToken() == "" {
			notification.To = "giulio.alfieri@gmail.com" // TODO: Change to user.GetEmail(),
			err = sendEmail(ctx, gs, notification)
		} else {
			notification.DeviceToken = user.GetDeviceToken()
			err = sendPush(ctx, gs, notification)
		}
	case globalmodel.NotificationRealtorInviteSMS:
		notification := globalmodel.Notification{
			Title: "Convite da participar da TOQ",
			Body:  fmt.Sprintf("A %s está te convidando a participar da TOQ, vinculado(a) a ela. Baixe a aplicação e aeite o convite.", iCode),
			Icon:  "",
			To:    "+5511999141768", //TODO: Change to user.GetPhoneNumber()
		}
		err = sendSMS(ctx, gs, notification)
	case globalmodel.NotificationRealtorInvitePush:
		notification := globalmodel.Notification{
			Title: "Convite para vinculo a imobiliária.",
			Body:  fmt.Sprintf("%s, voce tem um convite pendente para vínculo a uma imobiliária.", user.GetNickName()),
			Icon:  "",
		}
		if user.GetDeviceToken() == "" {
			notification.To = "giulio.alfieri@gmail.com" // TODO: Change to user.GetEmail(),
			err = sendEmail(ctx, gs, notification)
		} else {
			notification.DeviceToken = user.GetDeviceToken()
			err = sendPush(ctx, gs, notification)
		}

	case globalmodel.NotificationInviteAccepted:
		notification := globalmodel.Notification{
			Title: "Corretor aceitou o convite",
			Body:  fmt.Sprintf("%s, o corretor %s aceitou seu convite e agora está vinvulado a esta imobiliária.", user.GetNickName(), iCode),
			Icon:  "",
		}
		if user.GetDeviceToken() == "" {
			notification.To = "giulio.alfieri@gmail.com" // TODO: Change to user.GetEmail(),
			err = sendEmail(ctx, gs, notification)
		} else {
			notification.DeviceToken = user.GetDeviceToken()
			err = sendPush(ctx, gs, notification)
		}
	case globalmodel.NotificationInviteRejected:
		notification := globalmodel.Notification{
			Title: "Corretor rejeitou o convite",
			Body:  fmt.Sprintf("%s, o corretor %s rejeitou seu convite para vincular-se a esta imobiliária.", user.GetNickName(), iCode),
			Icon:  "",
		}
		if user.GetDeviceToken() == "" {
			notification.To = "giulio.alfieri@gmail.com" // TODO: Change to user.GetEmail(),
			err = sendEmail(ctx, gs, notification)
		} else {
			notification.DeviceToken = user.GetDeviceToken()
			err = sendPush(ctx, gs, notification)
		}

	case globalmodel.NotificationAgencyRemovedFromRealtor:
		notification := globalmodel.Notification{
			Title: "Corretor cancelou o vínculo",
			Body:  fmt.Sprintf("O corretor %s cancelou o vínculado com esta imobiliária.", user.GetNickName()),
			Icon:  "",
		}
		if user.GetDeviceToken() == "" {
			notification.To = "giulio.alfieri@gmail.com" // TODO: Change to user.GetEmail(),
			err = sendEmail(ctx, gs, notification)
		} else {
			notification.DeviceToken = user.GetDeviceToken()
			err = sendPush(ctx, gs, notification)
		}

	case globalmodel.NotificationRealtorRemovedFromAgency:
		notification := globalmodel.Notification{
			Title: "Imobiliária cacelou o vínculo",
			Body:  fmt.Sprintf("%s, a imobiliária %s cancelou o vínculado com voce.", user.GetNickName(), iCode),
			Icon:  "",
		}
		if user.GetDeviceToken() == "" {
			notification.To = "giulio.alfieri@gmail.com" // TODO: Change to user.GetEmail(),
			err = sendEmail(ctx, gs, notification)
		} else {
			notification.DeviceToken = user.GetDeviceToken()
			err = sendPush(ctx, gs, notification)
		}

	}

	return
}
