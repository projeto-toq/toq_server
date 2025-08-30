package userservices

import (
	"log/slog"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) CreateAccessToken(secret string, user usermodel.UserInterface, expires int64) (accessToken string, err error) {
	// TODO: Implementar verificação de status quando sistema for migrado
	// Por enquanto, assumir que perfil está ativo se tem role ativo
	profileStatus := user.GetActiveRole() != nil

	// Converter RoleInterface para RoleSlug
	var roleSlug permissionmodel.RoleSlug
	if activeRole := user.GetActiveRole(); activeRole != nil {
		if role := activeRole.GetRole(); role != nil {
			roleSlug = permissionmodel.RoleSlug(role.GetSlug())
		}
	}

	infos := usermodel.UserInfos{
		ID:            user.GetID(),
		ProfileStatus: profileStatus,
		Role:          roleSlug,
	}
	now := time.Now().UTC().Unix()
	claims := jwt.MapClaims{
		string(globalmodel.TokenKey): infos,
		"exp":                        expires,
		"iat":                        now,
		"iss":                        "toq-server",
		"jti":                        uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err = token.SignedString([]byte(secret))
	if err != nil {
		slog.Error("Error trying to generate JWT access token", "err", err)
		return "", utils.ErrInternalServer
	}

	return
}
