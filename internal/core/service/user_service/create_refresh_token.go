package userservices

import (
	"log/slog"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/golang-jwt/jwt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) CreateRefreshToken(expired bool, userID int64, tokens *usermodel.Tokens, jti string) (err error) {

	var expires int64

	secret := globalmodel.GetJWTSecret()

	expires = us.GetTokenExpiration(expired)

	infos := usermodel.UserInfos{
		ID:         userID,
		UserRoleID: 0,                            // Refresh token não precisa de UserRoleID específico
		RoleStatus: permissionmodel.StatusActive, // Status padrão para refresh token
	}

	//cria os claims
	now := time.Now().UTC().Unix()
	claims := jwt.MapClaims{
		string(globalmodel.TokenKey): infos,
		"exp":                        expires,
		"iat":                        now,
		"iss":                        "toq-server",
		"jti":                        jti,
	}

	// cria o refresh token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//assina o refresh token com a senha
	refreshToken, err := token.SignedString([]byte(secret))
	if err != nil {
		slog.Error("error trying to generate jwt refresh token", "error", err)
		return utils.ErrInternalServer
	}

	//salva na estrutura
	tokens.RefreshToken = refreshToken

	return

}
