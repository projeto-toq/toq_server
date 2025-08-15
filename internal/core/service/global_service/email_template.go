package globalservice

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"path/filepath"

	"github.com/giulio-alfieri/toq_server/internal/core/utils/paths"
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

	base := paths.BaseDir()
	primary := filepath.Join(base, "internal", "core", "templates", templateFile)
	candidates := []string{primary}

	// fallback subindo diretórios caso executável esteja em /bin
	if _, _, ok := paths.BestFile(filepath.Join("internal", "core", "templates", templateFile)); ok {
		found, all, ok2 := paths.BestFile(filepath.Join("internal", "core", "templates", templateFile))
		candidates = all
		if ok2 {
			primary = found
		}
	}

	slog.Debug("Loading email template", "templateType", templateType, "candidates", candidates, "chosen", primary)

	tmpl, err := template.ParseFiles(primary)
	if err != nil {
		slog.Error("Failed to parse email template", "chosen", primary, "candidates", candidates, "error", err)
		return "", fmt.Errorf("failed to parse email template %s: %w", primary, err)
	}

	var body bytes.Buffer
	templateData := struct{ Code string }{Code: code}
	if err = tmpl.Execute(&body, templateData); err != nil {
		slog.Error("Failed to execute email template", "chosen", primary, "error", err)
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}
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
