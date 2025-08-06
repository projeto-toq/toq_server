package userservices

import (
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) CreateRefreshToken(expired bool, userID int64, tokens *usermodel.Tokens) (err error) {

	var expires int64

	secret := globalmodel.GetJWTSecret()

	expires = us.GetTokenExpiration(expired)

	infos := usermodel.UserInfos{
		ID:            userID,
		ProfileStatus: false,
	}

	//cria os claims
	claims := jwt.MapClaims{
		string(globalmodel.TokenKey): infos,
		"exp":                        expires,
	}

	// cria o refresh token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//assina o refresh token com a senha
	refreshToken, err := token.SignedString([]byte(secret))
	if err != nil {
		slog.Error("error trying to generate jwt refresh token", "error", err)
		return status.Errorf(codes.Internal, "Internal server error")
	}

	//salva na estrutura
	tokens.RefreshToken = refreshToken

	return

}
