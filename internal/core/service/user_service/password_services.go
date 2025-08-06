package userservices

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Encrypt the user password
func (us *userService) encryptPassword(password string) (hashPassword string) {

	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(password))

	return hex.EncodeToString(hash.Sum(nil))
}

func validatePassword(password string) (err error) {
	if len(password) < 8 {
		return status.Error(codes.InvalidArgument, "Password must be at least 8 characters")
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
		return status.Error(codes.InvalidArgument, "Password must include a mix of uppercase, lowercase, numbers, and special characters")
	}
	return
}
