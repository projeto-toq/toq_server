package userservices

import (
	"bytes"
	"context"
	"html/template"
	"sync"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const systemUserWelcomeTemplatePath = "internal/core/templates/email_welcome_system_user.html"

// systemUserWelcomeEmailRenderer renders the welcome email template lazily.
type systemUserWelcomeEmailRenderer struct {
	once sync.Once
	tmpl *template.Template
	err  error
}

func newSystemUserWelcomeEmailRenderer() *systemUserWelcomeEmailRenderer {
	return &systemUserWelcomeEmailRenderer{}
}

func (r *systemUserWelcomeEmailRenderer) render(data any) (string, error) {
	r.once.Do(func() {
		r.tmpl, r.err = template.ParseFiles(systemUserWelcomeTemplatePath)
	})
	if r.err != nil {
		return "", r.err
	}
	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// sendSystemUserWelcomeEmail prepares and schedules the welcome notification for system users.
func (us *userService) sendSystemUserWelcomeEmail(ctx context.Context, user usermodel.UserInterface, role permissionmodel.RoleSlug) error {
	ctx = utils.ContextWithLogger(ctx)

	body, err := us.emailRenderer.render(map[string]any{
		"NickName": user.GetNickName(),
		"FullName": user.GetFullName(),
		"Role":     role.String(),
		"ResetURL": us.cfg.SystemUserResetPasswordURL,
	})
	if err != nil {
		return err
	}

	notificationService := us.globalService.GetUnifiedNotificationService()
	emailRequest := globalservice.NotificationRequest{
		Type:    globalservice.NotificationTypeEmail,
		To:      user.GetEmail(),
		Subject: "TOQ - Welcome",
		Body:    body,
	}

	if notifyErr := notificationService.SendNotification(ctx, emailRequest); notifyErr != nil {
		return notifyErr
	}

	return nil
}
