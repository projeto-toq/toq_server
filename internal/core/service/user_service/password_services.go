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
		slog.Error("encryptPassword: failed to generate bcrypt hash: %v", "err", err)
		return ""
	}
	return string(hash)
}

func validatePassword(field, password string) (err error) {
	if len(password) < 8 {
		return utils.ValidationError(field, "Password must be at least 8 characters long.")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)
	specialChar := regexp.MustCompile(`[!@#~$%^&*()_+{}":;'?/>.<,]`)

	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z':
			hasLower = true
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case specialChar.MatchString(string(char)):
			hasSpecial = true
		}
	}

	if !(hasUpper && hasLower && hasNumber && hasSpecial) {
		return utils.ValidationError(field, "Password must include uppercase, lowercase, number, and special character.")
	}
	return
}
