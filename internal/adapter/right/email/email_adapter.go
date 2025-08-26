package emailadapter

import (
	"crypto/tls"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"gopkg.in/gomail.v2"
)

type EmailAdapter struct {
	dialer     *gomail.Dialer
	fromEmail  string
	fromName   string
	maxRetries int
	timeout    time.Duration
}

func NewEmailAdapter(config globalmodel.Environment) *EmailAdapter {
	dialer := gomail.NewDialer(
		config.EMAIL.SMTPServer,
		config.EMAIL.SMTPPort,
		config.EMAIL.SMTPUser,
		config.EMAIL.SMTPPassword,
	)

	// Configuração TLS robusta
	dialer.TLSConfig = &tls.Config{
		ServerName:         config.EMAIL.SMTPServer,
		InsecureSkipVerify: config.EMAIL.SkipVerify,
	}

	if config.EMAIL.UseSSL {
		dialer.SSL = true
	}

	// Configurar timeout
	timeout := 30 * time.Second // default
	if config.EMAIL.TimeoutSecs > 0 {
		timeout = time.Duration(config.EMAIL.TimeoutSecs) * time.Second
	}
	// Note: gomail.Dialer não suporta timeout configurável diretamente
	// O timeout será controlado no contexto de envio

	// Configurar email de origem
	fromEmail := config.EMAIL.FromEmail
	if fromEmail == "" {
		fromEmail = config.EMAIL.SMTPUser
	}

	// Configurar máximo de tentativas
	maxRetries := 3 // default
	if config.EMAIL.MaxRetries > 0 {
		maxRetries = config.EMAIL.MaxRetries
	}

	return &EmailAdapter{
		dialer:     dialer,
		fromEmail:  fromEmail,
		fromName:   config.EMAIL.FromName,
		maxRetries: maxRetries,
		timeout:    timeout,
	}
}
