package userservices

import (
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) CreateAccessToken(secret string, user usermodel.UserInterface, expires int64) (accessToken string, err error) {
	infos := usermodel.UserInfos{
		ID:            user.GetID(),
		ProfileStatus: user.GetActiveRole().GetStatus() != usermodel.StatusPendingProfile,
		Role:          user.GetActiveRole().GetRole(),
	}

	claims := jwt.MapClaims{
		string(globalmodel.TokenKey): infos,
		"exp":                        expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err = token.SignedString([]byte(secret))
	if err != nil {
		slog.Error("Error trying to generate JWT access token", "err", err)
		return "", status.Errorf(codes.Internal, "Internal server error")
	}

	return
}
