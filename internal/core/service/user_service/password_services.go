package userservices

import (
	"log/slog"
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// Encrypt the user password
func (us *userService) encryptPassword(password string) (hashPassword string) {
	// Usa bcrypt com custo padrão; erros aqui são improváveis (custo inválido)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("encryptPassword: failed to generate bcrypt hash: %v", err)
		return ""
	}
	return string(hash)
}

func validatePassword(password string) (err error) {
	if len(password) < 8 {
		return utils.ValidationError("password", "Password must be at least 8 characters")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z':
			hasLower = true
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case regexp.MustCompile(`[!@#~$%^&*()_+{}":;'?/>.<,]`).MatchString(string(char)):
			hasSpecial = true
		}
	}

	if !(hasUpper && hasLower && hasNumber && hasSpecial) {
		return utils.ValidationError("password", "Password must include upper, lower, number, and special char")
	}
	return
}
