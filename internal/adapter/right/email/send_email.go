package emailadapter

import (
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"gopkg.in/gomail.v2"
)

func (e *EmailAdapter) SendEmail(notification globalmodel.Notification) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "TOQ_APP@toq.app.br")
	m.SetHeader("To", notification.To)
	m.SetHeader("Subject", notification.Title)
	m.SetBody("text/html", notification.Body)

	return e.dialer.DialAndSend(m)
}
