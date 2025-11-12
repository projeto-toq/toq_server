package userservices

import (
	"bytes"
	"context"
	"html/template"
	"sync"
	"time"

	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const accountLockedTemplatePath = "internal/core/templates/email_account_locked.html"

// accountLockedEmailRenderer renders the account locked email template lazily.
type accountLockedEmailRenderer struct {
	once sync.Once
	tmpl *template.Template
	err  error
}

func newAccountLockedEmailRenderer() *accountLockedEmailRenderer {
	return &accountLockedEmailRenderer{}
}

func (r *accountLockedEmailRenderer) render(data any) (string, error) {
	r.once.Do(func() {
		r.tmpl, r.err = template.ParseFiles(accountLockedTemplatePath)
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

// sendAccountLockedNotification sends a security alert to user when account is temporarily locked
//
// This method is called ASYNCHRONOUSLY (via goroutine) after user account is blocked due to
// excessive failed signin attempts. It notifies the legitimate user via email/SMS about the
// security event, while the API response remains generic to prevent account enumeration.
//
// Security Strategy:
//   - API returns generic "Invalid credentials" to potential attacker
//   - Legitimate user receives notification via registered email/phone
//   - Notification includes: block reason, duration, and unblock time
//
// Parameters:
//   - ctx: Context for logging and tracing (detached from HTTP request lifecycle)
//   - userID: ID of the blocked user
//
// Notification Details:
//   - Sent to user's registered email address
//   - Subject: Security alert about temporary account lock
//   - Body: HTML template with lock details and password reset link
//   - Includes timestamp when account will be automatically unblocked
//
// Error Handling:
//   - Errors are logged but DO NOT fail the signin flow (fire-and-forget)
//   - User retrieval errors: logs ERROR, aborts notification
//   - Template rendering errors: logs ERROR, aborts notification
//   - Notification send errors: logs WARN (not critical, user can retry signin)
//
// Observability:
//   - Logs ERROR if user data cannot be retrieved
//   - Logs ERROR if template fails to render
//   - Logs WARN if notification fails to send
//   - No logging on success (notification service handles its own logging)
//
// Important Notes:
//   - Method runs in separate goroutine - must handle its own errors
//   - Context may outlive HTTP request - avoid using request-scoped values
//   - Notification is best-effort: if it fails, user can still retry after 15 minutes
//   - Uses HTML email template (internal/core/templates/email_account_locked.html)
//
// Example:
//
//	// Called asynchronously when user is blocked:
//	go us.sendAccountLockedNotification(ctx, userID)
func (us *userService) sendAccountLockedNotification(_ context.Context, userID int64) {
	// Create new context with logger for goroutine (detached from request lifecycle)
	ctx := utils.ContextWithLogger(context.Background())
	logger := utils.LoggerFromContext(ctx)

	// Retrieve user data to get email address (not in transaction, read-only)
	user, err := us.repo.GetUserByID(ctx, nil, userID)
	if err != nil {
		// Log error but don't fail - notification is best-effort
		logger.Error("auth.security_alert.get_user_failed",
			"user_id", userID,
			"error", err)
		return
	}

	// Calculate when account will be automatically unblocked
	now := time.Now().UTC()
	unblockTime := now.Add(us.cfg.TempBlockDuration)

	// Render email body from HTML template
	renderer := newAccountLockedEmailRenderer()
	body, err := renderer.render(map[string]any{
		"NickName":         user.GetNickName(),
		"BlockedAt":        now.Format("02/01/2006 15:04:05 MST"),
		"UnblockAt":        unblockTime.Format("02/01/2006 15:04:05 MST"),
		"FailedAttempts":   us.cfg.MaxWrongSigninAttempts,
		"ResetPasswordURL": us.cfg.SystemUserResetPasswordURL,
	})
	if err != nil {
		logger.Error("auth.security_alert.render_template_failed",
			"user_id", userID,
			"error", err)
		return
	}

	// Prepare notification request
	notificationReq := globalservice.NotificationRequest{
		Type:    globalservice.NotificationTypeEmail,
		To:      user.GetEmail(),
		Subject: "TOQ - Alerta de Segurança",
		Body:    body,
	}

	// Send notification asynchronously (fire-and-forget)
	notificationService := us.globalService.GetUnifiedNotificationService()
	err = notificationService.SendNotification(ctx, notificationReq)
	if err != nil {
		// Log warning but don't fail - notification is non-critical
		logger.Warn("auth.security_alert.notification_failed",
			"user_id", userID,
			"email", user.GetEmail(),
			"error", err)
		return
	}

	// Success - notification service will handle logging
}
