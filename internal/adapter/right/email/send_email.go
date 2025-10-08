package emailadapter

import (
	"context"
	"fmt"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"gopkg.in/gomail.v2"
)

func (e *EmailAdapter) SendEmail(ctx context.Context, notification globalmodel.Notification) error {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	m := gomail.NewMessage()

	// Headers dinâmicos com configuração robusta
	if e.fromName != "" {
		m.SetHeader("From", m.FormatAddress(e.fromEmail, e.fromName))
	} else {
		m.SetHeader("From", e.fromEmail)
	}

	m.SetHeader("To", notification.To)
	m.SetHeader("Subject", notification.Title)
	m.SetBody("text/html", notification.Body)

	// Retry logic para robustez
	var lastErr error
	for attempt := 0; attempt <= e.maxRetries; attempt++ {
		if attempt > 0 {
			// Backoff exponencial: 1s, 2s, 4s...
			waitTime := time.Duration(attempt) * time.Second
			logger.Debug("email.send.retry_wait", "attempt", attempt, "wait", waitTime, "to", notification.To)
			time.Sleep(waitTime)
		}

		logger.Debug("email.send.attempt", "attempt", attempt+1, "to", notification.To, "subject", notification.Title)

		if err := e.dialer.DialAndSend(m); err != nil {
			lastErr = err
			utils.SetSpanError(ctx, err)
			logger.Warn("email.send.failure", "attempt", attempt+1, "error", err, "to", notification.To)
			continue
		}

		logger.Info("email.send.success", "to", notification.To, "subject", notification.Title, "attempts", attempt+1)
		return nil // Sucesso
	}

	return fmt.Errorf("failed to send email after %d attempts: %w", e.maxRetries+1, lastErr)
}
