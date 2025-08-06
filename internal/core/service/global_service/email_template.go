package globalservice

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"path/filepath"
)

// EmailTemplateType define os tipos de template de email
type EmailTemplateType int

const (
	EmailVerificationTemplate EmailTemplateType = iota + 1
	PasswordResetTemplate
)

// EmailTemplateLoader interface para carregamento de templates
type EmailTemplateLoader interface {
	LoadTemplate(templateType EmailTemplateType, code string) (string, error)
}

// emailTemplateLoader implementa EmailTemplateLoader
type emailTemplateLoader struct{}

// NewEmailTemplateLoader cria uma nova instância do loader
func NewEmailTemplateLoader() EmailTemplateLoader {
	return &emailTemplateLoader{}
}

// LoadTemplate carrega e processa templates de email
func (etl *emailTemplateLoader) LoadTemplate(templateType EmailTemplateType, code string) (string, error) {
	templateFile, err := etl.getTemplateFile(templateType)
	if err != nil {
		return "", err
	}

	templatePath := filepath.Join("internal", "core", "templates", templateFile)

	slog.Debug("Loading email template",
		"templateType", templateType,
		"templatePath", templatePath,
		"code", code)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		slog.Error("Failed to parse email template",
			"templatePath", templatePath,
			"templateType", templateType,
			"error", err)
		return "", fmt.Errorf("failed to parse email template %s: %w", templatePath, err)
	}

	var body bytes.Buffer
	templateData := struct{ Code string }{Code: code}

	err = tmpl.Execute(&body, templateData)
	if err != nil {
		slog.Error("Failed to execute email template",
			"templatePath", templatePath,
			"templateType", templateType,
			"code", code,
			"error", err)
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}

	slog.Debug("Email template loaded successfully",
		"templateType", templateType,
		"templatePath", templatePath)

	return body.String(), nil
}

// getTemplateFile retorna o nome do arquivo do template baseado no tipo
func (etl *emailTemplateLoader) getTemplateFile(templateType EmailTemplateType) (string, error) {
	switch templateType {
	case EmailVerificationTemplate:
		return "email_verification.html", nil
	case PasswordResetTemplate:
		return "email_reset_password.html", nil
	default:
		err := fmt.Errorf("invalid email template type: %d", templateType)
		slog.Error("Invalid email template type", "templateType", templateType, "error", err)
		return "", err
	}
}
