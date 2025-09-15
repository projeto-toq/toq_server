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

func (us *userService) CreateRefreshToken(expired bool, user usermodel.UserInterface, tokens *usermodel.Tokens, jti string) (err error) {

	var expires int64

	secret := globalmodel.GetJWTSecret()

	// Refresh tokens must use refresh TTL (absolute window), not access TTL
	if expired {
		// Mantém comportamento de expiração forçada para cenários de teste
		expires = time.Now().UTC().Add(time.Hour * -1).Unix()
	} else {
		expires = time.Now().UTC().Add(globalmodel.GetRefreshTTL()).Unix()
	}

	// Incluir RoleSlug se usuário tiver uma role ativa (mesmo que refresh não dependa dela, facilita auditoria em clientes)
	var roleSlug permissionmodel.RoleSlug
	if ar := user.GetActiveRole(); ar != nil && ar.GetRole() != nil {
		roleSlug = permissionmodel.RoleSlug(ar.GetRole().GetSlug())
	}
	infos := usermodel.UserInfos{
		ID:         user.GetID(),
		UserRoleID: 0, // Refresh não exige role ID
		RoleSlug:   roleSlug,
	}

	//cria os claims
	now := time.Now().UTC().Unix()
	claims := jwt.MapClaims{
		string(globalmodel.TokenKey): infos,
		"exp":                        expires,
		"iat":                        now,
		"iss":                        "toq-server",
		"jti":                        jti,
		"typ":                        "refresh", // tipo explícito para validação
	}

	// cria o refresh token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//assina o refresh token com a senha
	refreshToken, err := token.SignedString([]byte(secret))
	if err != nil {
		slog.Error("user.create_refresh_token.sign_error", "err", err)
		return utils.InternalError("Failed to sign refresh token")
	}

	//salva na estrutura
	tokens.RefreshToken = refreshToken

	return

}
