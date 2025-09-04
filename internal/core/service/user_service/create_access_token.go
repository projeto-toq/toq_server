package userservices

import (
	"log/slog"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// CreateAccessToken generates a signed JWT access token. An active role is
// mandatory for access tokens. If the user has no active role, a domain error
// is returned and no token is issued.
func (us *userService) CreateAccessToken(secret string, user usermodel.UserInterface, expires int64) (accessToken string, err error) {
	// Exigência: access token requer role ativa
	activeRole := user.GetActiveRole()
	if activeRole == nil {
		slog.Warn("cannot issue access token without active role", "user_id", user.GetID())
		return "", utils.ErrUserActiveRoleMissing
	}

	infos := usermodel.UserInfos{
		ID:         user.GetID(),
		UserRoleID: activeRole.GetID(),
		RoleStatus: activeRole.GetStatus(),
	}

	now := time.Now().UTC().Unix()
	claims := jwt.MapClaims{
		string(globalmodel.TokenKey): infos,
		"exp":                        expires,
		"iat":                        now,
		"iss":                        "toq-server",
		"jti":                        uuid.New().String(),
		"typ":                        "access", // tipo explícito para validação
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err = token.SignedString([]byte(secret))
	if err != nil {
		slog.Error("Error trying to generate JWT access token", "err", err)
		return "", utils.InternalError("Failed to sign access token")
	}

	return
}
