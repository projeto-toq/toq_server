package emailadapter

import (
	"gopkg.in/gomail.v2"
)

type EmailAdapter struct {
	dialer *gomail.Dialer
}

func NewEmailAdapter(server string, port int, user string, password string) *EmailAdapter {

	dialer := gomail.NewDialer(server, port, user, password)

	return &EmailAdapter{
		dialer: dialer,
	}
}
