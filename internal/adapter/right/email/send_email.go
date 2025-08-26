package emailadapter

import (
	"fmt"
	"log/slog"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"gopkg.in/gomail.v2"
)

func (e *EmailAdapter) SendEmail(notification globalmodel.Notification) error {
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
			slog.Debug("Tentando reenvio de email", "attempt", attempt, "wait", waitTime, "to", notification.To)
			time.Sleep(waitTime)
		}

		slog.Debug("Enviando email", "attempt", attempt+1, "to", notification.To, "subject", notification.Title)

		if err := e.dialer.DialAndSend(m); err != nil {
			lastErr = err
			slog.Warn("Falha no envio de email", "attempt", attempt+1, "error", err, "to", notification.To)
			continue
		}

		slog.Info("Email enviado com sucesso", "to", notification.To, "subject", notification.Title, "attempts", attempt+1)
		return nil // Sucesso
	}

	return fmt.Errorf("failed to send email after %d attempts: %w", e.maxRetries+1, lastErr)
}
