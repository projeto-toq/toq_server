package userservices

import (
	"context"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// CreateAccessToken generates a signed JWT access token. An active role is
// mandatory for access tokens. If the user has no active role, a domain error
// is returned and no token is issued.
func (us *userService) CreateAccessToken(secret string, user usermodel.UserInterface, expires int64) (accessToken string, err error) {
	logger := utils.LoggerFromContext(context.Background())
	// Exigência: access token requer role ativa
	activeRole := user.GetActiveRole()
	if activeRole == nil {
		// Estado inválido de domínio: usuário deve ter sempre uma role ativa
		// Promover a log para Error e retornar 500 para facilitar detecção
		logger.Error("user.create_access_token.active_role_missing", "user_id", user.GetID())
		return "", utils.InternalError("Active role missing unexpectedly")
	}

	infos := usermodel.UserInfos{
		ID:         user.GetID(),
		UserRoleID: activeRole.GetID(),
		RoleSlug:   permissionmodel.RoleSlug(activeRole.GetRole().GetSlug()),
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
		logger.Error("user.create_access_token.sign_error", "error", err)
		return "", utils.InternalError("Failed to sign access token")
	}

	return
}
